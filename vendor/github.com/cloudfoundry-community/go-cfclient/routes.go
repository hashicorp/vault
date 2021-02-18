package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type RoutesResponse struct {
	Count     int              `json:"total_results"`
	Pages     int              `json:"total_pages"`
	NextUrl   string           `json:"next_url"`
	Resources []RoutesResource `json:"resources"`
}

type RoutesResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Route `json:"entity"`
}

type RouteRequest struct {
	DomainGuid string `json:"domain_guid"`
	SpaceGuid  string `json:"space_guid"`
	Host       string `json:"host"` // required for http routes
	Path       string `json:"path"`
	Port       int    `json:"port"`
}

type Route struct {
	Guid                string `json:"guid"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	Host                string `json:"host"`
	Path                string `json:"path"`
	DomainGuid          string `json:"domain_guid"`
	DomainURL           string `json:"domain_url"`
	SpaceGuid           string `json:"space_guid"`
	ServiceInstanceGuid string `json:"service_instance_guid"`
	Port                int    `json:"port"`
	c                   *Client
}

// CreateRoute creates a regular http route
func (c *Client) CreateRoute(routeRequest RouteRequest) (Route, error) {
	routesResource, err := c.createRoute("/v2/routes", routeRequest)
	if nil != err {
		return Route{}, err
	}
	return c.mergeRouteResource(routesResource), nil
}

// CreateTcpRoute creates a TCP route
func (c *Client) CreateTcpRoute(routeRequest RouteRequest) (Route, error) {
	routesResource, err := c.createRoute("/v2/routes?generate_port=true", routeRequest)
	if nil != err {
		return Route{}, err
	}
	return c.mergeRouteResource(routesResource), nil
}

// BindRoute associates the specified route with the application
func (c *Client) BindRoute(routeGUID, appGUID string) error {
	resp, err := c.DoRequest(c.NewRequest("PUT", fmt.Sprintf("/v2/routes/%s/apps/%s", routeGUID, appGUID)))
	if err != nil {
		return errors.Wrapf(err, "Error binding route %s to app %s", routeGUID, appGUID)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error binding route %s to app %s, response code: %d", routeGUID, appGUID, resp.StatusCode)
	}
	return nil
}

func (c *Client) GetRouteByGuid(guid string) (Route, error) {
	var route RoutesResource

	r := c.NewRequest("GET", fmt.Sprintf("/v2/routes/%s", guid))
	resp, err := c.DoRequest(r)
	if err != nil {
		return route.Entity, errors.Wrap(err, "Error requesting route")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&route)
	if err != nil {
		return route.Entity, errors.Wrap(err, "Error unmarshalling route response body")
	}

	route.Entity.Guid = route.Meta.Guid
	route.Entity.CreatedAt = route.Meta.CreatedAt
	route.Entity.UpdatedAt = route.Meta.UpdatedAt
	route.Entity.c = c
	return route.Entity, nil
}

func (c *Client) ListRoutesByQuery(query url.Values) ([]Route, error) {
	return c.fetchRoutes("/v2/routes?" + query.Encode())
}

func (c *Client) fetchRoutes(requestUrl string) ([]Route, error) {
	var routes []Route
	for {
		routesResp, err := c.getRoutesResponse(requestUrl)
		if err != nil {
			return []Route{}, err
		}
		for _, route := range routesResp.Resources {
			route.Entity.Guid = route.Meta.Guid
			route.Entity.CreatedAt = route.Meta.CreatedAt
			route.Entity.UpdatedAt = route.Meta.UpdatedAt
			route.Entity.c = c
			routes = append(routes, route.Entity)
		}
		requestUrl = routesResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return routes, nil
}

func (c *Client) ListRoutes() ([]Route, error) {
	return c.ListRoutesByQuery(nil)
}

func (r *Route) Domain() (*Domain, error) {
	req := r.c.NewRequest("GET", r.DomainURL)
	resp, err := r.c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "requesting domain for route "+r.DomainURL)
	}

	defer resp.Body.Close()
	var domain DomainResource
	if err = json.NewDecoder(resp.Body).Decode(&domain); err != nil {
		return nil, errors.Wrap(err, "unmarshalling domain")
	}

	return r.c.mergeDomainResource(domain), nil
}

func (c *Client) getRoutesResponse(requestUrl string) (RoutesResponse, error) {
	var routesResp RoutesResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return RoutesResponse{}, errors.Wrap(err, "Error requesting routes")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return RoutesResponse{}, errors.Wrap(err, "Error reading routes body")
	}
	err = json.Unmarshal(resBody, &routesResp)
	if err != nil {
		return RoutesResponse{}, errors.Wrap(err, "Error unmarshalling routes")
	}
	return routesResp, nil
}

func (c *Client) createRoute(requestUrl string, routeRequest RouteRequest) (RoutesResource, error) {
	var routeResp RoutesResource
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(routeRequest)
	if err != nil {
		return RoutesResource{}, errors.Wrap(err, "Error creating route - failed to serialize request body")
	}
	r := c.NewRequestWithBody("POST", requestUrl, buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return RoutesResource{}, errors.Wrap(err, "Error creating route")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return RoutesResource{}, errors.Wrap(err, "Error creating route")
	}
	err = json.Unmarshal(resBody, &routeResp)
	if err != nil {
		return RoutesResource{}, errors.Wrap(err, "Error unmarshalling routes")
	}
	routeResp.Entity.c = c
	return routeResp, nil
}

func (c *Client) DeleteRoute(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/routes/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting route %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) mergeRouteResource(rr RoutesResource) Route {
	rr.Entity.Guid = rr.Meta.Guid
	rr.Entity.CreatedAt = rr.Meta.CreatedAt
	rr.Entity.UpdatedAt = rr.Meta.UpdatedAt
	rr.Entity.c = c
	return rr.Entity
}
