// Copyright 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tags

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	CategoryURL      = "/com/vmware/cis/tagging/category"
	ErrAlreadyExists = "already_exists"
)

type CategoryCreateSpec struct {
	CreateSpec CategoryCreate `json:"create_spec"`
}

type CategoryUpdateSpec struct {
	UpdateSpec CategoryUpdate `json:"update_spec"`
}

type CategoryCreate struct {
	AssociableTypes []string `json:"associable_types"`
	Cardinality     string   `json:"cardinality"`
	Description     string   `json:"description"`
	Name            string   `json:"name"`
}

type CategoryUpdate struct {
	AssociableTypes []string `json:"associable_types"`
	Cardinality     string   `json:"cardinality"`
	Description     string   `json:"description"`
	Name            string   `json:"name"`
}

type Category struct {
	ID              string   `json:"id"`
	Description     string   `json:"description"`
	Name            string   `json:"name"`
	Cardinality     string   `json:"cardinality"`
	AssociableTypes []string `json:"associable_types"`
	UsedBy          []string `json:"used_by"`
}

func (c *RestClient) CreateCategoryIfNotExist(ctx context.Context, name string, description string, categoryType string, multiValue bool) (*string, error) {
	categories, err := c.GetCategoriesByName(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query category for %s", name)
	}

	if categories == nil {
		var multiValueStr string
		if multiValue {
			multiValueStr = "MULTIPLE"
		} else {
			multiValueStr = "SINGLE"
		}
		categoryCreate := CategoryCreate{[]string{categoryType}, multiValueStr, description, name}
		spec := CategoryCreateSpec{categoryCreate}
		id, err := c.CreateCategory(ctx, &spec)
		if err != nil {
			// in case there are two docker daemon try to create inventory category, query the category once again
			if strings.Contains(err.Error(), "ErrAlreadyExists") {
				if categories, err = c.GetCategoriesByName(ctx, name); err != nil {
					Logger.Debugf("Failed to get inventory category for %s", errors.WithStack(err))
					return nil, errors.Wrap(err, "create inventory category failed")
				}
			} else {
				Logger.Debugf("Failed to create inventory category for %s", errors.WithStack(err))
				return nil, errors.Wrap(err, "create inventory category failed")
			}
		} else {
			return id, nil
		}
	}
	if categories != nil {
		return &categories[0].ID, nil
	}
	// should not happen
	Logger.Debugf("Failed to create inventory for it's existed, but could not query back. Please check system")
	return nil, errors.Errorf("Failed to create inventory for it's existed, but could not query back. Please check system")
}

func (c *RestClient) CreateCategory(ctx context.Context, spec *CategoryCreateSpec) (*string, error) {
	Logger.Debugf("Create category %v", spec)
	stream, _, status, err := c.call(ctx, "POST", CategoryURL, spec, nil)

	Logger.Debugf("Get status code: %d", status)
	if status != http.StatusOK || err != nil {
		Logger.Debugf("Create category failed with status code: %d, error message: %s", status, errors.WithStack(err))
		return nil, errors.Wrapf(err, "Status code: %d", status)
	}

	type RespValue struct {
		Value string
	}

	var pID RespValue
	if err := json.NewDecoder(stream).Decode(&pID); err != nil {
		Logger.Debugf("Decode response body failed for: %s", errors.WithStack(err))
		return nil, errors.Wrap(err, "create category failed")
	}
	return &(pID.Value), nil
}

func (c *RestClient) GetCategory(ctx context.Context, id string) (*Category, error) {
	Logger.Debugf("Get category %s", id)

	stream, _, status, err := c.call(ctx, "GET", fmt.Sprintf("%s/id:%s", CategoryURL, id), nil, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("Get category failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return nil, errors.Errorf("Status code: %d, error: %s", status, err)
	}

	type RespValue struct {
		Value Category
	}

	var pCategory RespValue
	if err := json.NewDecoder(stream).Decode(&pCategory); err != nil {
		Logger.Debugf("Decode response body failed for: %s", errors.WithStack(err))
		return nil, errors.Wrapf(err, "get category %s failed", id)
	}
	return &(pCategory.Value), nil
}

func (c *RestClient) UpdateCategory(ctx context.Context, id string, spec *CategoryUpdateSpec) error {
	Logger.Debugf("Update category %v", spec)
	_, _, status, err := c.call(ctx, "PATCH", fmt.Sprintf("%s/id:%s", CategoryURL, id), spec, nil)

	Logger.Debugf("Get status code: %d", status)
	if status != http.StatusOK || err != nil {
		Logger.Debugf("Update category failed with status code: %d, error message: %s", status, errors.WithStack(err))
		return errors.Wrapf(err, "Status code: %d", status)
	}

	return nil
}

func (c *RestClient) DeleteCategory(ctx context.Context, id string) error {
	Logger.Debugf("Delete category %s", id)

	_, _, status, err := c.call(ctx, "DELETE", fmt.Sprintf("%s/id:%s", CategoryURL, id), nil, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("Delete category failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return errors.Errorf("Status code: %d, error: %s", status, err)
	}
	return nil
}

func (c *RestClient) ListCategories(ctx context.Context) ([]string, error) {
	Logger.Debugf("List all categories")

	stream, _, status, err := c.call(ctx, "GET", CategoryURL, nil, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("Get categories failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return nil, errors.Errorf("Status code: %d, error: %s", status, err)
	}

	type Categories struct {
		Value []string
	}

	var pCategories Categories
	if err := json.NewDecoder(stream).Decode(&pCategories); err != nil {
		Logger.Debugf("Decode response body failed for: %s", errors.WithStack(err))
		return nil, errors.Wrap(err, "list categories failed")
	}
	return pCategories.Value, nil
}

func (c *RestClient) GetCategoriesByName(ctx context.Context, name string) ([]Category, error) {
	Logger.Debugf("Get category %s", name)
	categoryIds, err := c.ListCategories(ctx)
	if err != nil {
		Logger.Debugf("Get category failed for: %s", errors.WithStack(err))
		return nil, errors.Wrapf(err, "get categories by name %s failed", name)
	}

	var categories []Category
	for _, cID := range categoryIds {
		category, err := c.GetCategory(ctx, cID)
		if err != nil {
			Logger.Debugf("Get category %s failed for %s", cID, errors.WithStack(err))
		}
		if category.Name == name {
			categories = append(categories, *category)
		}
	}
	return categories, nil
}
