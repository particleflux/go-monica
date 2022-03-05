package monica

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	userAgent = "go-monica"

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRetryAfter    = "Retry-After"
)

var errNonNilContext = errors.New("context must be non-nil")

type Client struct {
	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	// BaseUrl should include the `/api/` part
	BaseURL    *url.URL
	HTTPClient *http.Client
	UserAgent  string

	rateLimit   Rate
	accessToken string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	//Activity *ActivityService
	Contacts          *ContactsService
	ContactFieldTypes *ContactFieldTypeService
	Countries         *CountriesService
	Genders           *GenderService
	Tags              *TagsService
}

type service struct {
	client *Client
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	Limit int `url:"limit,omitempty"`
}

type successResponse struct {
	Data interface{} `json:"data"`
}

// NewClient instantiates a new client with given baseUrl and oauth accessToken
func NewClient(baseUrl, accessToken string) *Client {
	base, _ := url.Parse(baseUrl)

	client := &Client{
		BaseURL: base,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		UserAgent:   userAgent,
		accessToken: accessToken,
	}

	client.common.client = client

	client.Contacts = (*ContactsService)(&client.common)
	client.ContactFieldTypes = (*ContactFieldTypeService)(&client.common)
	client.Countries = (*CountriesService)(&client.common)
	client.Genders = (*GenderService)(&client.common)
	client.Tags = (*TagsService)(&client.common)

	return client
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	return req, nil
}

// Response is a Monica API response. This wraps the standard http.Response.
type Response struct {
	*http.Response

	// Explicitly specify the Rate type so Rate's String() receiver doesn't
	// propagate to Response.
	Rate Rate

	Meta ListMeta `json:"meta"`

	Data interface{} `json:"data"`
}

type Rate struct {
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset Timestamp `json:"reset"`
}

// ListMeta contains meta data returned by API for lists/collections
type ListMeta struct {
	CurrentPage int    `json:"current_page"`
	From        int    `json:"from"`
	LastPage    int    `json:"last_page"`
	Path        string `json:"path"`
	PerPage     int    `json:"per_page"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

func (r *Response) populatePageValues() {
	//json.NewDecoder(r.Body).Decode(r)
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.Rate = parseRate(r)
	//response.populatePageValues()
	return response
}

// BareDo sends an API request and lets you handle the api response. If an error
// or API Error occurs, the error will contain more information. Otherwise you
// are supposed to read and close the response's Body. If rate limit is exceeded
// and reset time is in the future, BareDo returns *RateLimitError immediately
// without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	response := newResponse(resp)

	c.rateLimit = response.Rate

	err = CheckResponse(resp)
	if err != nil {
		defer resp.Body.Close()
	}
	return response, err
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	decErr := json.NewDecoder(resp.Body).Decode(v)
	if decErr == io.EOF {
		decErr = nil // ignore EOF errors caused by empty response body
	}
	if decErr != nil {
		err = decErr
	}

	return resp, err
}

// parseRate parses the rate related headers.
func parseRate(r *http.Response) Rate {
	var rate Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRetryAfter); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.Reset = Timestamp{time.Now().Add(time.Duration(v) * time.Second)}
		}
	}
	return rate
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have response
// body, and a JSON response body that maps to ErrorResponse.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// *AcceptedError for 202 Accepted status codes,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	response, _ := ioutil.ReadAll(r.Body)

	return errors.New(
		fmt.Sprintf("status code: %d\nResponse:\n%s\n", r.StatusCode, response),
	)
}
