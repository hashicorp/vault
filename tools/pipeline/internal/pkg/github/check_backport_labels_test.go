// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"testing"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/stretchr/testify/require"
)

// fakeLabelsAPI is an in-memory BackportLabelsAPI that serves labels in pages
// of pageSize.
type fakeLabelsAPI struct {
	labels     []*libgithub.Label
	pageSize   int
	created    []*libgithub.Label
	deleted    []string
	failList   bool
	failCreate bool
	failDelete bool
}

func (f *fakeLabelsAPI) ListLabels(_ context.Context, _, _ string, opts *libgithub.ListOptions) ([]*libgithub.Label, *libgithub.Response, error) {
	if f.failList {
		return nil, nil, errors.New("boom")
	}
	size := f.pageSize
	if size == 0 {
		size = len(f.labels) + 1
	}
	page := max(opts.Page, 1)
	start := (page - 1) * size
	end := min(start+size, len(f.labels))
	resp := &libgithub.Response{}
	if end < len(f.labels) {
		resp.NextPage = page + 1
	}
	return f.labels[start:end], resp, nil
}

func (f *fakeLabelsAPI) CreateLabel(_ context.Context, _, _ string, l *libgithub.Label) (*libgithub.Label, *libgithub.Response, error) {
	if f.failCreate {
		return nil, nil, errors.New("create boom")
	}
	f.created = append(f.created, l)
	return l, nil, nil
}

func (f *fakeLabelsAPI) DeleteLabel(_ context.Context, _, _, name string) (*libgithub.Response, error) {
	if f.failDelete {
		return nil, errors.New("delete boom")
	}
	f.deleted = append(f.deleted, name)
	return nil, nil
}

func newLabel(name, color, desc string) *libgithub.Label {
	return &libgithub.Label{Name: &name, Color: &color, Description: &desc}
}

func newBackportTestReq(versions ...string) *CheckBackportLabelsReq {
	cfg := &releases.VersionsConfig{ActiveVersion: &releases.ActiveVersion{Versions: map[string]*releases.Version{}}}
	for _, v := range versions {
		cfg.ActiveVersion.Versions[v] = &releases.Version{}
	}
	return &CheckBackportLabelsReq{
		Owner:             "hashicorp",
		Repo:              "vault-enterprise",
		VersionsDecodeRes: &releases.DecodeRes{Config: cfg},
	}
}

// Test_CheckBackportLabelsReq_Run verifies that missing and stale labels are
// detected, unmanaged labels are ignored, and that pagination is followed.
func Test_CheckBackportLabelsReq_Run(t *testing.T) {
	t.Parallel()

	api := &fakeLabelsAPI{pageSize: 2, labels: []*libgithub.Label{
		newLabel("bug", "ff0000", ""),
		newLabel("backport/1.19.x", "70F5B5", ""),
		newLabel("backport-failed", "ededed", ""),
		newLabel("backport/skip", "ededed", ""),
		newLabel("backport/1.17.x", "70F5B5", ""),
	}}

	res, err := newBackportTestReq("1.19.x", "1.20.x").Run(context.Background(), api)
	require.NoError(t, err)
	require.False(t, res.InSync)
	require.Equal(t, []string{"backport/1.19.x", "backport/1.20.x"}, res.Expected)
	require.Equal(t, []string{"backport/1.20.x"}, res.Missing)
	require.Equal(t, []string{"backport/1.17.x"}, res.Stale)
	require.Empty(t, api.created)
	require.Empty(t, api.deleted)
}

// Test_CheckBackportLabelsReq_Run_InSync verifies that a matching label set
// reports in sync.
func Test_CheckBackportLabelsReq_Run_InSync(t *testing.T) {
	t.Parallel()

	api := &fakeLabelsAPI{labels: []*libgithub.Label{newLabel("backport/1.19.x", "70F5B5", "")}}
	res, err := newBackportTestReq("1.19.x").Run(context.Background(), api)
	require.NoError(t, err)
	require.True(t, res.InSync)
}

// Test_CheckBackportLabelsReq_Run_Fix verifies that missing labels are created
// with the color and description of an existing label and that stale labels
// are deleted.
func Test_CheckBackportLabelsReq_Run_Fix(t *testing.T) {
	t.Parallel()

	api := &fakeLabelsAPI{labels: []*libgithub.Label{
		newLabel("backport/1.17.x", "123456", "desc"),
	}}
	req := newBackportTestReq("1.19.x")
	req.CreateMissing = true
	req.DeleteStale = true

	res, err := req.Run(context.Background(), api)
	require.NoError(t, err)
	require.True(t, res.InSync)
	require.Equal(t, []string{"backport/1.19.x"}, res.Created)
	require.Equal(t, []string{"backport/1.17.x"}, res.Deleted)
	require.Len(t, api.created, 1)
	for _, l := range api.created {
		require.Equal(t, "123456", l.GetColor())
		require.Equal(t, "desc", l.GetDescription())
	}
	require.Equal(t, []string{"backport/1.17.x"}, api.deleted)
}

// Test_CheckBackportLabelsReq_Run_DefaultColor verifies the fallback color when
// no backport labels exist to copy from.
func Test_CheckBackportLabelsReq_Run_DefaultColor(t *testing.T) {
	t.Parallel()

	api := &fakeLabelsAPI{}
	req := newBackportTestReq("1.19.x")
	req.CreateMissing = true

	_, err := req.Run(context.Background(), api)
	require.NoError(t, err)
	require.Len(t, api.created, 1)
	require.Equal(t, defaultBackportLabelColor, api.created[0].GetColor())
}

// Test_CheckBackportLabelsReq_Run_Errors verifies request validation, decode
// failures, and list failures.
func Test_CheckBackportLabelsReq_Run_Errors(t *testing.T) {
	t.Parallel()

	_, err := (*CheckBackportLabelsReq)(nil).Run(context.Background(), &fakeLabelsAPI{})
	require.ErrorContains(t, err, "uninitialized")

	req := newBackportTestReq("1.19.x")
	req.Repo = ""
	_, err = req.Run(context.Background(), &fakeLabelsAPI{})
	require.ErrorContains(t, err, "repo is required")

	req = newBackportTestReq("1.19.x")
	req.VersionsDecodeRes = &releases.DecodeRes{Err: errors.New("decode failed")}
	_, err = req.Run(context.Background(), &fakeLabelsAPI{})
	require.ErrorContains(t, err, "decode failed")

	_, err = newBackportTestReq("1.19.x").Run(context.Background(), &fakeLabelsAPI{failList: true})
	require.ErrorContains(t, err, "listing labels")
}

// Test_CheckBackportLabelsRes_Output verifies the rendered outputs.
func Test_CheckBackportLabelsRes_Output(t *testing.T) {
	t.Parallel()

	res := &CheckBackportLabelsRes{
		Expected: []string{"backport/1.19.x", "backport/1.20.x"},
		Missing:  []string{"backport/1.20.x"},
		Stale:    []string{"backport/1.17.x"},
	}

	table := res.ToTable()
	require.Contains(t, table, "backport/1.20.x")
	require.Contains(t, table, "stale")

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"missing":["backport/1.20.x"]`)

	out, err := res.ToGithubOutput()
	require.NoError(t, err)
	require.Contains(t, string(out), "in_sync=false\n")
	require.Contains(t, string(out), "stale=backport/1.17.x\n")
}

// Test_CheckBackportLabelsReq_Run_EmptyVersions verifies that an empty active
// version list never deletes labels, and that report-only still lists every
// label as stale.
func Test_CheckBackportLabelsReq_Run_EmptyVersions(t *testing.T) {
	t.Parallel()

	labels := []*libgithub.Label{
		newLabel("backport/1.19.x", "70F5B5", ""),
		newLabel("backport/1.21.x", "70F5B5", ""),
	}

	api := &fakeLabelsAPI{labels: labels}
	req := newBackportTestReq()
	req.DeleteStale = true
	_, err := req.Run(context.Background(), api)
	require.ErrorContains(t, err, "no active versions")
	require.Empty(t, api.deleted)

	api = &fakeLabelsAPI{labels: labels}
	res, err := newBackportTestReq().Run(context.Background(), api)
	require.NoError(t, err)
	require.False(t, res.InSync)
	require.Equal(t, []string{"backport/1.19.x", "backport/1.21.x"}, res.Stale)
	require.Empty(t, api.deleted)
}

// Test_CheckBackportLabelsReq_Run_PartialFailure verifies that Run stops at the
// first create or delete error and reports what was already changed.
func Test_CheckBackportLabelsReq_Run_PartialFailure(t *testing.T) {
	t.Parallel()

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		api := &fakeLabelsAPI{failCreate: true, labels: []*libgithub.Label{newLabel("backport/1.17.x", "70F5B5", "")}}
		req := newBackportTestReq("1.19.x", "1.20.x")
		req.CreateMissing = true
		req.DeleteStale = true

		res, err := req.Run(context.Background(), api)
		require.ErrorContains(t, err, `creating label "backport/1.19.x"`)
		require.NotNil(t, res)
		require.Empty(t, res.Created)
		require.Empty(t, api.deleted)
		require.False(t, res.InSync)
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()

		api := &fakeLabelsAPI{failDelete: true, labels: []*libgithub.Label{
			newLabel("backport/1.17.x", "70F5B5", ""),
			newLabel("backport/1.18.x", "70F5B5", ""),
		}}
		req := newBackportTestReq("1.19.x")
		req.CreateMissing = true
		req.DeleteStale = true

		res, err := req.Run(context.Background(), api)
		require.ErrorContains(t, err, `deleting label "backport/1.17.x"`)
		require.NotNil(t, res)
		require.Equal(t, []string{"backport/1.19.x"}, res.Created)
		require.Empty(t, res.Deleted)
		require.False(t, res.InSync)
	})
}

// Test_CheckBackportLabelsReq_Run_SingleFixFlag verifies that each fix flag
// only acts on its own kind of drift.
func Test_CheckBackportLabelsReq_Run_SingleFixFlag(t *testing.T) {
	t.Parallel()

	labels := func() []*libgithub.Label { return []*libgithub.Label{newLabel("backport/1.17.x", "70F5B5", "")} }

	t.Run("delete stale only", func(t *testing.T) {
		t.Parallel()

		api := &fakeLabelsAPI{labels: labels()}
		req := newBackportTestReq("1.19.x")
		req.DeleteStale = true

		res, err := req.Run(context.Background(), api)
		require.NoError(t, err)
		require.Empty(t, api.created)
		require.Equal(t, []string{"backport/1.17.x"}, api.deleted)
		require.Equal(t, []string{"backport/1.19.x"}, res.Missing)
		require.Empty(t, res.Stale)
		require.False(t, res.InSync)
	})

	t.Run("create missing only", func(t *testing.T) {
		t.Parallel()

		api := &fakeLabelsAPI{labels: labels()}
		req := newBackportTestReq("1.19.x")
		req.CreateMissing = true

		res, err := req.Run(context.Background(), api)
		require.NoError(t, err)
		require.Empty(t, api.deleted)
		require.Len(t, api.created, 1)
		require.Empty(t, res.Missing)
		require.Equal(t, []string{"backport/1.17.x"}, res.Stale)
		require.False(t, res.InSync)
	})
}
