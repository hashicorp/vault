// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
)

// Project represents a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/
type Project struct {
	ID        *int       `json:"id,omitempty"`
	URL       *string    `json:"url,omitempty"`
	OwnerURL  *string    `json:"owner_url,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Body      *string    `json:"body,omitempty"`
	Number    *int       `json:"number,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`

	// The User object that generated the project.
	Creator *User `json:"creator,omitempty"`
}

func (p Project) String() string {
	return Stringify(p)
}

// ListProjects lists the projects for a repo.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-projects
func (s *RepositoriesService) ListProjects(owner, repo string, opt *ListOptions) ([]*Project, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	projects := []*Project{}
	resp, err := s.client.Do(req, &projects)
	if err != nil {
		return nil, resp, err
	}

	return projects, resp, err
}

// GetProject gets a GitHub Project for a repo.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#get-a-project
func (s *RepositoriesService) GetProject(owner, repo string, number int) (*Project, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	project := &Project{}
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// ProjectOptions specifies the parameters to the
// RepositoriesService.CreateProject and
// RepositoriesService.UpdateProject methods.
type ProjectOptions struct {
	// The name of the project. (Required for creation; optional for update.)
	Name string `json:"name,omitempty"`
	// The body of the project. (Optional.)
	Body string `json:"body,omitempty"`
}

// CreateProject creates a GitHub Project for the specified repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#create-a-project
func (s *RepositoriesService) CreateProject(owner, repo string, projectOptions *ProjectOptions) (*Project, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects", owner, repo)
	req, err := s.client.NewRequest("POST", u, projectOptions)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	project := &Project{}
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// UpdateProject updates a repository project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#update-a-project
func (s *RepositoriesService) UpdateProject(owner, repo string, number int, projectOptions *ProjectOptions) (*Project, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v", owner, repo, number)
	req, err := s.client.NewRequest("PATCH", u, projectOptions)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	project := &Project{}
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// DeleteProject deletes a GitHub Project from a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#delete-a-project
func (s *RepositoriesService) DeleteProject(owner, repo string, number int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v", owner, repo, number)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// ProjectColumn represents a column of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/
type ProjectColumn struct {
	ID         *int       `json:"id,omitempty"`
	Name       *string    `json:"name,omitempty"`
	ProjectURL *string    `json:"project_url,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`
	UpdatedAt  *Timestamp `json:"updated_at,omitempty"`
}

// ListProjectColumns lists the columns of a GitHub Project for a repo.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-columns
func (s *RepositoriesService) ListProjectColumns(owner, repo string, number int, opt *ListOptions) ([]*ProjectColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v/columns", owner, repo, number)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	columns := []*ProjectColumn{}
	resp, err := s.client.Do(req, &columns)
	if err != nil {
		return nil, resp, err
	}

	return columns, resp, err
}

// GetProjectColumn gets a column of a GitHub Project for a repo.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#get-a-column
func (s *RepositoriesService) GetProjectColumn(owner, repo string, columnID int) (*ProjectColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v", owner, repo, columnID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	column := &ProjectColumn{}
	resp, err := s.client.Do(req, column)
	if err != nil {
		return nil, resp, err
	}

	return column, resp, err
}

// ProjectColumnOptions specifies the parameters to the
// RepositoriesService.CreateProjectColumn and
// RepositoriesService.UpdateProjectColumn methods.
type ProjectColumnOptions struct {
	// The name of the project column. (Required for creation and update.)
	Name string `json:"name"`
}

// CreateProjectColumn creates a column for the specified (by number) project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#create-a-column
func (s *RepositoriesService) CreateProjectColumn(owner, repo string, number int, columnOptions *ProjectColumnOptions) (*ProjectColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/%v/columns", owner, repo, number)
	req, err := s.client.NewRequest("POST", u, columnOptions)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	column := &ProjectColumn{}
	resp, err := s.client.Do(req, column)
	if err != nil {
		return nil, resp, err
	}

	return column, resp, err
}

// UpdateProjectColumn updates a column of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#update-a-column
func (s *RepositoriesService) UpdateProjectColumn(owner, repo string, columnID int, columnOptions *ProjectColumnOptions) (*ProjectColumn, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v", owner, repo, columnID)
	req, err := s.client.NewRequest("PATCH", u, columnOptions)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	column := &ProjectColumn{}
	resp, err := s.client.Do(req, column)
	if err != nil {
		return nil, resp, err
	}

	return column, resp, err
}

// DeleteProjectColumn deletes a column from a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#delete-a-column
func (s *RepositoriesService) DeleteProjectColumn(owner, repo string, columnID int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v", owner, repo, columnID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// ProjectColumnMoveOptions specifies the parameters to the
// RepositoriesService.MoveProjectColumn method.
type ProjectColumnMoveOptions struct {
	// Position can be one of "first", "last", or "after:<column-id>", where
	// <column-id> is the ID of a column in the same project. (Required.)
	Position string `json:"position"`
}

// MoveProjectColumn moves a column within a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-column
func (s *RepositoriesService) MoveProjectColumn(owner, repo string, columnID int, moveOptions *ProjectColumnMoveOptions) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v/moves", owner, repo, columnID)
	req, err := s.client.NewRequest("POST", u, moveOptions)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// ProjectCard represents a card in a column of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/
type ProjectCard struct {
	ColumnURL  *string    `json:"column_url,omitempty"`
	ContentURL *string    `json:"content_url,omitempty"`
	ID         *int       `json:"id,omitempty"`
	Note       *string    `json:"note,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`
	UpdatedAt  *Timestamp `json:"updated_at,omitempty"`
}

// ListProjectCards lists the cards in a column of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#list-projects-cards
func (s *RepositoriesService) ListProjectCards(owner, repo string, columnID int, opt *ListOptions) ([]*ProjectCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v/cards", owner, repo, columnID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	cards := []*ProjectCard{}
	resp, err := s.client.Do(req, &cards)
	if err != nil {
		return nil, resp, err
	}

	return cards, resp, err
}

// GetProjectCard gets a card in a column of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#get-a-project-card
func (s *RepositoriesService) GetProjectCard(owner, repo string, columnID int) (*ProjectCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v", owner, repo, columnID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	card := &ProjectCard{}
	resp, err := s.client.Do(req, card)
	if err != nil {
		return nil, resp, err
	}

	return card, resp, err
}

// ProjectCardOptions specifies the parameters to the
// RepositoriesService.CreateProjectCard and
// RepositoriesService.UpdateProjectCard methods.
type ProjectCardOptions struct {
	// The note of the card. Note and ContentID are mutually exclusive.
	Note string `json:"note,omitempty"`
	// The ID (not Number) of the Issue or Pull Request to associate with this card.
	// Note and ContentID are mutually exclusive.
	ContentID int `json:"content_id,omitempty"`
	// The type of content to associate with this card. Possible values are: "Issue", "PullRequest".
	ContentType string `json:"content_type,omitempty"`
}

// CreateProjectCard creates a card in the specified column of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#create-a-project-card
func (s *RepositoriesService) CreateProjectCard(owner, repo string, columnID int, cardOptions *ProjectCardOptions) (*ProjectCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/%v/cards", owner, repo, columnID)
	req, err := s.client.NewRequest("POST", u, cardOptions)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	card := &ProjectCard{}
	resp, err := s.client.Do(req, card)
	if err != nil {
		return nil, resp, err
	}

	return card, resp, err
}

// UpdateProjectCard updates a card of a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#update-a-project-card
func (s *RepositoriesService) UpdateProjectCard(owner, repo string, cardID int, cardOptions *ProjectCardOptions) (*ProjectCard, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v", owner, repo, cardID)
	req, err := s.client.NewRequest("PATCH", u, cardOptions)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	card := &ProjectCard{}
	resp, err := s.client.Do(req, card)
	if err != nil {
		return nil, resp, err
	}

	return card, resp, err
}

// DeleteProjectCard deletes a card from a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#delete-a-project-card
func (s *RepositoriesService) DeleteProjectCard(owner, repo string, cardID int) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v", owner, repo, cardID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}

// ProjectCardMoveOptions specifies the parameters to the
// RepositoriesService.MoveProjectCard method.
type ProjectCardMoveOptions struct {
	// Position can be one of "top", "bottom", or "after:<card-id>", where
	// <card-id> is the ID of a card in the same project.
	Position string `json:"position"`
	// ColumnID is the ID of a column in the same project. Note that ColumnID
	// is required when using Position "after:<card-id>" when that card is in
	// another column; otherwise it is optional.
	ColumnID int `json:"column_id,omitempty"`
}

// MoveProjectCard moves a card within a GitHub Project.
//
// GitHub API docs: https://developer.github.com/v3/repos/projects/#move-a-project-card
func (s *RepositoriesService) MoveProjectCard(owner, repo string, cardID int, moveOptions *ProjectCardMoveOptions) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/projects/columns/cards/%v/moves", owner, repo, cardID)
	req, err := s.client.NewRequest("POST", u, moveOptions)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeProjectsPreview)

	return s.client.Do(req, nil)
}
