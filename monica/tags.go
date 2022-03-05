package monica

import (
	"context"
	"fmt"
)

// TagsService handles communication with the tag related methods of the API.
// API docs: https://www.monicahq.com/api/tags
type TagsService service

type Tag struct {
	Id       int    `json:"id,omitempty"`
	Object   string `json:"object,omitempty"`
	Name     string `json:"name"`
	NameSlug string `json:"name_slug,omitempty"`

	Account  struct {
		Id int `json:"id,omitempty"`
	} `json:"account,omitempty"`

	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`
}

type TagListOptions struct {
	ListOptions
}

type listTagsResponse struct {
	Data *[]*Tag `json:"data"`
	Meta ListMeta `json:"meta"`
}

type createTagRequest struct {
	Name string `json:"name"`
}

func (s *TagsService) ListTags(ctx context.Context, opts *TagListOptions) (*[]*Tag, *ListMeta, error) {
	url, err := addOptions("tags", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	response := new(listTagsResponse)
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}

// GetTag Retrieves information about a single specified tag
func (s *TagsService) GetTag(ctx context.Context, id int) (*Tag, error) {
	url := fmt.Sprintf("tags/%d", id)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Tag `json:"data"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// CreateTag Creates a tag with given name
func (s *TagsService) CreateTag(ctx context.Context, name string) (*Tag, error) {
	body := createTagRequest{Name: name}
	req, err := s.client.NewRequest("POST", "tags", body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Tag `json:"data"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// UpdateTag Renames a tag
func (s *TagsService) UpdateTag(ctx context.Context, id int, name string) (*Tag, error) {
	url := fmt.Sprintf("tags/%d", id)
	body := createTagRequest{Name: name}
	req, err := s.client.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data *Tag `json:"data"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// DeleteTag Delete a tag by id
func (s *TagsService) DeleteTag(ctx context.Context, id int) error {
	url := fmt.Sprintf("tags/%d", id)
	req, err := s.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	response := struct {
		Deleted bool `json:"deleted"`
		Id string `json:"id"`
	}{}
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return err
	}

	return nil
}