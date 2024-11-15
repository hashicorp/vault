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
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	TagURL = "/com/vmware/cis/tagging/tag"
)

type TagCreateSpec struct {
	CreateSpec TagCreate `json:"create_spec"`
}

type TagCreate struct {
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type TagUpdateSpec struct {
	UpdateSpec TagUpdate `json:"update_spec"`
}

type TagUpdate struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Tag struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	CategoryID  string   `json:"category_id"`
	UsedBy      []string `json:"used_by"`
}

var Logger = logrus.New()

func (c *RestClient) CreateTagIfNotExist(ctx context.Context, name string, description string, categoryID string) (*string, error) {
	tagCreate := TagCreate{categoryID, description, name}
	spec := TagCreateSpec{tagCreate}
	id, err := c.CreateTag(ctx, &spec)
	if err == nil {
		return id, nil
	}
	Logger.Debugf("Created tag %s failed for %s", errors.WithStack(err))
	// if already exists, query back
	if strings.Contains(err.Error(), ErrAlreadyExists) {
		tagObjs, err := c.GetTagByNameForCategory(ctx, name, categoryID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query tag %s for category %s", name, categoryID)
		}
		if tagObjs != nil {
			return &tagObjs[0].ID, nil
		}

		// should not happen
		return nil, errors.New("Failed to create tag for it's existed, but could not query back. Please check system")
	}

	return nil, errors.Wrap(err, "failed to create tag")
}

func (c *RestClient) DeleteTagIfNoObjectAttached(ctx context.Context, id string) error {
	objs, err := c.ListAttachedObjects(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete tag")
	}
	if objs != nil && len(objs) > 0 {
		Logger.Debugf("tag %s related objects is not empty, do not delete it.", id)
		return nil
	}
	return c.DeleteTag(ctx, id)
}

func (c *RestClient) CreateTag(ctx context.Context, spec *TagCreateSpec) (*string, error) {
	Logger.Debugf("Create Tag %v", spec)
	stream, _, status, err := c.call(ctx, "POST", TagURL, spec, nil)

	Logger.Debugf("Get status code: %d", status)
	if status != http.StatusOK || err != nil {
		Logger.Debugf("Create tag failed with status code: %d, error message: %s", status, errors.WithStack(err))
		return nil, errors.Wrapf(err, "Status code: %d", status)
	}

	type RespValue struct {
		Value string
	}

	var pID RespValue
	if err := json.NewDecoder(stream).Decode(&pID); err != nil {
		Logger.Debugf("Decode response body failed for: %s", errors.WithStack(err))
		return nil, errors.Wrap(err, "create tag failed")
	}
	return &pID.Value, nil
}

func (c *RestClient) GetTag(ctx context.Context, id string) (*Tag, error) {
	Logger.Debugf("Get tag %s", id)

	stream, _, status, err := c.call(ctx, "GET", fmt.Sprintf("%s/id:%s", TagURL, id), nil, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("Get tag failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return nil, errors.Wrapf(err, "Status code: %d", status)
	}

	type RespValue struct {
		Value Tag
	}

	var pTag RespValue
	if err := json.NewDecoder(stream).Decode(&pTag); err != nil {
		Logger.Debugf("Decode response body failed for: %s", errors.WithStack(err))
		return nil, errors.Wrapf(err, "failed to get tag %s", id)
	}
	return &(pTag.Value), nil
}

func (c *RestClient) UpdateTag(ctx context.Context, id string, spec *TagUpdateSpec) error {
	Logger.Debugf("Update tag %v", spec)
	_, _, status, err := c.call(ctx, "PATCH", fmt.Sprintf("%s/id:%s", TagURL, id), spec, nil)

	Logger.Debugf("Get status code: %d", status)
	if status != http.StatusOK || err != nil {
		Logger.Debugf("Update tag failed with status code: %d, error message: %s", status, errors.WithStack(err))
		return errors.Wrapf(err, "Status code: %d", status)
	}

	return nil
}

func (c *RestClient) DeleteTag(ctx context.Context, id string) error {
	Logger.Debugf("Delete tag %s", id)

	_, _, status, err := c.call(ctx, "DELETE", fmt.Sprintf("%s/id:%s", TagURL, id), nil, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("Delete tag failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return errors.Wrapf(err, "Status code: %d", status)
	}
	return nil
}

func (c *RestClient) ListTags(ctx context.Context) ([]string, error) {
	Logger.Debugf("List all tags")

	stream, _, status, err := c.call(ctx, "GET", TagURL, nil, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("Get tags failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return nil, errors.Wrapf(err, "Status code: %d", status)
	}

	return c.handleTagIDList(stream)
}

func (c *RestClient) ListTagsForCategory(ctx context.Context, id string) ([]string, error) {
	Logger.Debugf("List tags for category: %s", id)

	type PostCategory struct {
		CId string `json:"category_id"`
	}
	spec := PostCategory{id}
	stream, _, status, err := c.call(ctx, "POST", fmt.Sprintf("%s/id:%s?~action=list-tags-for-category", TagURL, id), spec, nil)

	if status != http.StatusOK || err != nil {
		Logger.Debugf("List tags for category failed with status code: %s, error message: %s", status, errors.WithStack(err))
		return nil, errors.Wrapf(err, "Status code: %d", status)
	}

	return c.handleTagIDList(stream)
}

func (c *RestClient) handleTagIDList(stream io.ReadCloser) ([]string, error) {
	type Tags struct {
		Value []string
	}

	var pTags Tags
	if err := json.NewDecoder(stream).Decode(&pTags); err != nil {
		Logger.Debugf("Decode response body failed for: %s", errors.WithStack(err))
		return nil, errors.Wrap(err, "failed to decode json")
	}
	return pTags.Value, nil
}

// Get tag through tag name and category id
func (c *RestClient) GetTagByNameForCategory(ctx context.Context, name string, id string) ([]Tag, error) {
	Logger.Debugf("Get tag %s for category %s", name, id)
	tagIds, err := c.ListTagsForCategory(ctx, id)
	if err != nil {
		Logger.Debugf("Get tag failed for %s", errors.WithStack(err))
		return nil, errors.Wrapf(err, "get tag failed for name %s category %s", name, id)
	}

	var tags []Tag
	for _, tID := range tagIds {
		tag, err := c.GetTag(ctx, tID)
		if err != nil {
			Logger.Debugf("Get tag %s failed for %s", tID, errors.WithStack(err))
			return nil, errors.Wrapf(err, "get tag failed for name %s category %s", name, id)
		}
		if tag.Name == name {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}

// Get attached tags through tag name pattern
func (c *RestClient) GetAttachedTagsByNamePattern(ctx context.Context, namePattern string, objID string, objType string) ([]Tag, error) {
	tagIds, err := c.ListAttachedTags(ctx, objID, objType)
	if err != nil {
		Logger.Debugf("Get attached tags failed for %s", errors.WithStack(err))
		return nil, errors.Wrap(err, "get attached tags failed")
	}

	var validName = regexp.MustCompile(namePattern)
	var tags []Tag
	for _, tID := range tagIds {
		tag, err := c.GetTag(ctx, tID)
		if err != nil {
			Logger.Debugf("Get tag %s failed for %s", tID, errors.WithStack(err))
		}
		if validName.MatchString(tag.Name) {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}
