// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Comments = (*comments)(nil)

// Comments describes all the comment related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/comments
type Comments interface {
	// List all comments of the given run.
	List(ctx context.Context, runID string) (*CommentList, error)

	// Read a comment by its ID.
	Read(ctx context.Context, commentID string) (*Comment, error)

	// Create a new comment with the given options.
	Create(ctx context.Context, runID string, options CommentCreateOptions) (*Comment, error)
}

// Comments implements Comments.
type comments struct {
	client *Client
}

// CommentList represents a list of comments.
type CommentList struct {
	*Pagination
	Items []*Comment
}

// Comment represents a Terraform Enterprise comment.
type Comment struct {
	ID   string `jsonapi:"primary,comments"`
	Body string `jsonapi:"attr,body"`
}

type CommentCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,comments"`

	// Required: Body of the comment.
	Body string `jsonapi:"attr,body"`
}

// List all comments of the given run.
func (s *comments) List(ctx context.Context, runID string) (*CommentList, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/comments", url.PathEscape(runID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	cl := &CommentList{}
	err = req.Do(ctx, cl)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// Create a new comment with the given options.
func (s *comments) Create(ctx context.Context, runID string, options CommentCreateOptions) (*Comment, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/comments", url.PathEscape(runID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	comm := &Comment{}
	err = req.Do(ctx, comm)
	if err != nil {
		return nil, err
	}

	return comm, err
}

// Read a comment by its ID.
func (s *comments) Read(ctx context.Context, commentID string) (*Comment, error) {
	if !validStringID(&commentID) {
		return nil, ErrInvalidCommentID
	}

	u := fmt.Sprintf("comments/%s", url.PathEscape(commentID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	comm := &Comment{}
	err = req.Do(ctx, comm)
	if err != nil {
		return nil, err
	}

	return comm, nil
}

func (o CommentCreateOptions) valid() error {
	if !validString(&o.Body) {
		return ErrInvalidCommentBody
	}

	return nil
}
