package monica

import "context"

// CountriesService handles communication with the country related
// methods of the Monica API.
//
// API docs: https://www.monicahq.com/api/countries
type CountriesService service

type Country struct {
	Id     string `json:"id"`
	Object string `json:"object"`
	Name   string `json:"name"`
	Iso    string `json:"iso"`
}

type listCountryResponse struct {
	Countries map[string]*Country `json:"data"`
}

type CountryListOptions struct {
	// yep, empty, no options currently, not even paging, this returns all at once
}

func (s *CountriesService) ListCountries(ctx context.Context, opts *CountryListOptions) (*map[string]*Country, error) {
	url, err := addOptions("countries", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	countryList := new(listCountryResponse)
	_, err = s.client.Do(ctx, req, countryList)
	if err != nil {
		return nil, err
	}

	return &countryList.Countries, nil
}
