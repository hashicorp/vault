package cfclient

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"

	"fmt"
	"net/http"

	"code.cloudfoundry.org/gofileutils/fileutils"
	"github.com/pkg/errors"
)

type BuildpackResponse struct {
	Count     int                 `json:"total_results"`
	Pages     int                 `json:"total_pages"`
	NextUrl   string              `json:"next_url"`
	Resources []BuildpackResource `json:"resources"`
}

type BuildpackResource struct {
	Meta   Meta      `json:"metadata"`
	Entity Buildpack `json:"entity"`
}

type Buildpack struct {
	Guid      string `json:"guid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Locked    bool   `json:"locked"`
	Position  int    `json:"position"`
	Filename  string `json:"filename"`
	Stack     string `json:"stack"`
	c         *Client
}

type BuildpackRequest struct {
	// These are all pointers to the values so that we can tell
	// whether people wanted position 0, or enable/unlock values,
	// vs whether they didn't specify them and want them unchanged/default.
	Name     *string `json:"name,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
	Locked   *bool   `json:"locked,omitempty"`
	Position *int    `json:"position,omitempty"`
	Stack    *string `json:"stack,omitempty"`
}

func (c *Client) CreateBuildpack(bpr *BuildpackRequest) (*Buildpack, error) {
	if bpr.Name == nil || *bpr.Name == "" {
		return nil, errors.New("Unable to create a buidlpack with no name")
	}
	requestUrl := "/v2/buildpacks"
	req := c.NewRequest("POST", requestUrl)
	req.obj = bpr
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating buildpack:")
	}
	bp, err := c.handleBuildpackResp(resp)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating buildpack:")
	}
	return &bp, nil
}

func (c *Client) ListBuildpacks() ([]Buildpack, error) {
	var buildpacks []Buildpack
	requestUrl := "/v2/buildpacks"
	for {
		buildpackResp, err := c.getBuildpackResponse(requestUrl)
		if err != nil {
			return []Buildpack{}, err
		}
		for _, buildpack := range buildpackResp.Resources {
			buildpacks = append(buildpacks, c.mergeBuildpackResource(buildpack))
		}
		requestUrl = buildpackResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return buildpacks, nil
}

func (c *Client) DeleteBuildpack(guid string, async bool) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/buildpacks/%s?async=%t", guid, async)))
	if err != nil {
		return err
	}
	if (async && (resp.StatusCode != http.StatusAccepted)) || (!async && (resp.StatusCode != http.StatusNoContent)) {
		return errors.Wrapf(err, "Error deleting buildpack %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) getBuildpackResponse(requestUrl string) (BuildpackResponse, error) {
	var buildpackResp BuildpackResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return BuildpackResponse{}, errors.Wrap(err, "Error requesting buildpacks")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return BuildpackResponse{}, errors.Wrap(err, "Error reading buildpack request")
	}
	err = json.Unmarshal(resBody, &buildpackResp)
	if err != nil {
		return BuildpackResponse{}, errors.Wrap(err, "Error unmarshalling buildpack")
	}
	return buildpackResp, nil
}

func (c *Client) mergeBuildpackResource(buildpack BuildpackResource) Buildpack {
	buildpack.Entity.Guid = buildpack.Meta.Guid
	buildpack.Entity.CreatedAt = buildpack.Meta.CreatedAt
	buildpack.Entity.UpdatedAt = buildpack.Meta.UpdatedAt
	buildpack.Entity.c = c
	return buildpack.Entity
}

func (c *Client) GetBuildpackByGuid(buildpackGUID string) (Buildpack, error) {
	requestUrl := fmt.Sprintf("/v2/buildpacks/%s", buildpackGUID)
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Buildpack{}, errors.Wrap(err, "Error requesting buildpack info")
	}
	return c.handleBuildpackResp(resp)
}

func (c *Client) handleBuildpackResp(resp *http.Response) (Buildpack, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Buildpack{}, err
	}
	var buildpackResource BuildpackResource
	if err := json.Unmarshal(body, &buildpackResource); err != nil {
		return Buildpack{}, err
	}
	return c.mergeBuildpackResource(buildpackResource), nil
}

func (b *Buildpack) Upload(file io.Reader, fileName string) error {
	var capturedErr error
	fileutils.TempFile("requests", func(requestFile *os.File, err error) {
		if err != nil {
			capturedErr = err
			return
		}
		writer := multipart.NewWriter(requestFile)
		part, err := writer.CreateFormFile("buildpack", fileName)

		if err != nil {
			_ = writer.Close()
			capturedErr = err
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			capturedErr = fmt.Errorf("Error creating upload: %s", err.Error())
			return
		}

		err = writer.Close()
		if err != nil {
			capturedErr = err
			return
		}

		requestFile.Seek(0, 0)
		fileStats, err := requestFile.Stat()
		if err != nil {
			capturedErr = fmt.Errorf("Error getting file info: %s", err)
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/buildpacks/%s/bits", b.c.Config.ApiAddress, b.Guid), requestFile)
		if err != nil {
			capturedErr = err
			return
		}

		req.ContentLength = fileStats.Size()
		contentType := fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())
		req.Header.Set("Content-Type", contentType)
		resp, err := b.c.Do(req) //client.Do() handles the HTTP status code checking for us
		if err != nil {
			capturedErr = err
			return
		}
		defer resp.Body.Close()
	})

	return errors.Wrap(capturedErr, "Error uploading buildpack:")
}

func (b *Buildpack) Update(bpr *BuildpackRequest) error {
	requestUrl := fmt.Sprintf("/v2/buildpacks/%s", b.Guid)
	req := b.c.NewRequest("PUT", requestUrl)
	req.obj = bpr
	resp, err := b.c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error updating buildpack:")
	}
	newBp, err := b.c.handleBuildpackResp(resp)
	if err != nil {
		return errors.Wrap(err, "Error updating buildpack:")
	}
	b.Name = newBp.Name
	b.Locked = newBp.Locked
	b.Enabled = newBp.Enabled
	return nil
}

func (bpr *BuildpackRequest) Lock() {
	b := true
	bpr.Locked = &b
}
func (bpr *BuildpackRequest) Unlock() {
	b := false
	bpr.Locked = &b
}
func (bpr *BuildpackRequest) Enable() {
	b := true
	bpr.Enabled = &b
}
func (bpr *BuildpackRequest) Disable() {
	b := false
	bpr.Enabled = &b
}
func (bpr *BuildpackRequest) SetPosition(i int) {
	bpr.Position = &i
}
func (bpr *BuildpackRequest) SetName(s string) {
	bpr.Name = &s
}
func (bpr *BuildpackRequest) SetStack(s string) {
	bpr.Stack = &s
}
