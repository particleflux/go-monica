package monica

import "context"

// GenderService handles communication with the gender related methods of the API.
// API docs: https://www.monicahq.com/api/genders
type GenderService service

type Gender struct {
	Id      int    `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Name    string `json:"name"`
	Account struct {
		Id int `json:"id,omitempty"`
	} `json:"account,omitempty"`

	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`
}

type GenderListOptions struct {
	ListOptions
}

type listGenderResponse struct {
	Data *[]*Gender ` json:"data"`
	Meta ListMeta `json:"meta"`
}

func (s *GenderService) ListGenders(ctx context.Context, opts *GenderListOptions) (*[]*Gender, *ListMeta, error) {
	url, err := addOptions("genders", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	response := new(listGenderResponse)
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}