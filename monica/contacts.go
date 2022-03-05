package monica

import (
	"context"
	"fmt"
)

// ContactsService handles communication with the tag related methods of the API.
// API docs: https://www.monicahq.com/api/contacts
type ContactsService service

type Contact struct {
	// output only
	Id     int    `json:"id,omitempty"`
	Object string `json:"object,omitempty"`
	HashId string `json:"hash_id,omitempty"`

	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Gender      string `json:"gender"`
	Description string `json:"description"`

	BirthdateDay        int  `json:"birthdate_day,omitempty"`
	BirthdateMonth      int  `json:"birthdate_month,omitempty"`
	BirthdateYear       int  `json:"birthdate_year,omitempty"`
	IsBirthdateKnown    bool `json:"is_birthdate_known"`
	BirthdateIsAgeBased bool `json:"birthdate_is_age_based,omitempty"`
	BirthdateAge        int  `json:"birthdate_age,omitempty"`

	IsPartial bool `json:"is_partial,omitempty"`

	IsDeceased             bool `json:"is_deceased"`
	DeceasedDateDay        int  `json:"deceased_date_day,omitempty"`
	DeceasedDateMonth      int  `json:"deceased_date_month,omitempty"`
	DeceasedDateYear       int  `json:"deceased_date_year,omitempty"`
	DeceasedDateIsAgeBased bool `json:"deceased_date_is_age_based,omitempty"`
	IsDeceasedDateKnown    bool `json:"is_deceased_date_known"`
}

type ContactInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	GenderId    int    `json:"gender_id"`
	Description string `json:"description,omitempty"`

	BirthdateDay        int  `json:"birthdate_day,omitempty"`
	BirthdateMonth      int  `json:"birthdate_month,omitempty"`
	BirthdateYear       int  `json:"birthdate_year,omitempty"`
	IsBirthdateKnown    bool `json:"is_birthdate_known"`
	BirthdateIsAgeBased bool `json:"birthdate_is_age_based,omitempty"`
	BirthdateAge        int  `json:"birthdate_age,omitempty"`

	IsPartial bool `json:"is_partial,omitempty"`

	IsDeceased             bool `json:"is_deceased"`
	DeceasedDateDay        int  `json:"deceased_date_day,omitempty"`
	DeceasedDateMonth      int  `json:"deceased_date_month,omitempty"`
	DeceasedDateYear       int  `json:"deceased_date_year,omitempty"`
	DeceasedDateIsAgeBased bool `json:"deceased_date_is_age_based,omitempty"`
	IsDeceasedDateKnown    bool `json:"is_deceased_date_known"`
}

type ContactSearchListOptions struct {
	ListOptions
	Query string `url:"query,omitempty""`
}

type listContactsResponse struct {
	Data *[]*Contact `json:"data"`
	Meta ListMeta    `json:"meta"`
}

type updateContactCareerInput struct {
	Job     string `json:"job,omitempty"`
	Company string `json:"company,omitempty"`
}

type CreateContactFieldInput struct {
	ContactFieldTypeId int `json:"contact_field_type_id"`
	ContactId          int `json:"contact_id"`
	// Data value of the contact field, max 255 characters
	Data string `json:"data"`
}

type addTagInput struct {
	Tags []string `json:"tags"`
}

func (s *ContactsService) CreateContact(ctx context.Context, input *ContactInput) (*Contact, error) {
	req, err := s.client.NewRequest("POST", "contacts", *input)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Contact `json:"data"`
	}{}

	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (s *ContactsService) DeleteContact(ctx context.Context, id int) error {
	url := fmt.Sprintf("contacts/%d", id)
	req, err := s.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	response := struct {
		Deleted bool   `json:"deleted"`
		Id      string `json:"id"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return err
	}

	return nil
}

func (s *ContactsService) SearchContacts(ctx context.Context, opts *ContactSearchListOptions) (*[]*Contact, *ListMeta, error) {
	url, err := addOptions("contacts", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	response := new(listContactsResponse)
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}

func (s *ContactsService) UpdateContactCareer(ctx context.Context, contactId int, job string, company string) (*Contact, error) {
	url := fmt.Sprintf("contacts/%d/work", contactId)
	body := updateContactCareerInput{
		Job:     job,
		Company: company,
	}

	req, err := s.client.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Contact `json:"data"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (s *ContactsService) UpdateContact(ctx context.Context, contactId int, contactInput ContactInput) (*Contact, error) {
	url := fmt.Sprintf("contacts/%d", contactId)

	req, err := s.client.NewRequest("PUT", url, contactInput)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Contact `json:"data"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (s *ContactsService) CreateContactField(ctx context.Context, input *CreateContactFieldInput) (*ContactField, error) {
	req, err := s.client.NewRequest("POST", "contactfields", input)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *ContactField `json:"data"`
	}{}

	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (s *ContactsService) AddTags(ctx context.Context, contactId int, tags []string) (*Contact, error) {
	url := fmt.Sprintf("contacts/%d/setTags", contactId)
	req, err := s.client.NewRequest("POST", url, addTagInput{Tags: tags})
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Contact `json:"data"`
	}{}

	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
