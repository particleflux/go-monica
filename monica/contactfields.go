package monica

import (
	"context"
)

// ContactFieldTypeService handles communication with the contactfieldType related methods of the API.
// API docs: https://www.monicahq.com/api/contactfieldtypes
type ContactFieldTypeService service

type ContactFieldType struct {
	// Id is the unique id of the contact field type
	Id int `json:"id"`
	// Object is always "contactfieldtype"
	Object string `json:"object"`
	// Name is the displayed name
	Name string `json:"name"`
	// FontawesomeIcon is a reference to a fontawesome icon
	// example: "fa fa-envelope-open-o"
	FontawesomeIcon string `json:"fontawesome_icon"`
	Protocol        string `json:"protocol"`
	// Delible defines whether the contactFieldType is deletable
	Delible bool   `json:"delible"`
	Type    string `json:"type"`
	Account struct {
		Id int `json:"id"`
	} `json:"account"`

	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`
}

type ContactField struct {
	Id               int              `json:"id"`
	Object           string           `json:"object"`
	Data             string           `json:"data"`
	ContactFieldType ContactFieldType `json:"contact_field_type"`
	Account          struct {
		Id int `json:"id"`
	} `json:"account"`
	Contact Contact `json:"contact"`

	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`
}

type ContactFieldTypeListOptions struct {
	ListOptions
}

type listContactFieldTypeResponse struct {
	Data *[]*ContactFieldType ` json:"data"`
	Meta ListMeta             `json:"meta"`
}

func (s *ContactFieldTypeService) ListContactFieldTypes(ctx context.Context, opts *ContactFieldTypeListOptions) (*[]*ContactFieldType, *ListMeta, error) {
	url, err := addOptions("contactfieldtypes", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	response := new(listContactFieldTypeResponse)
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}
