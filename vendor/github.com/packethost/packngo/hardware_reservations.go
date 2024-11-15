package packngo

import (
	"path"
)

const hardwareReservationBasePath = "/hardware-reservations"

// HardwareReservationService interface defines available hardware reservation functions
type HardwareReservationService interface {
	Get(hardwareReservationID string, getOpt *GetOptions) (*HardwareReservation, *Response, error)
	List(projectID string, listOpt *ListOptions) ([]HardwareReservation, *Response, error)
	Move(string, string) (*HardwareReservation, *Response, error)
}

// HardwareReservationServiceOp implements HardwareReservationService
type HardwareReservationServiceOp struct {
	client requestDoer
}

// HardwareReservation struct
type HardwareReservation struct {
	ID            string    `json:"id,omitempty"`
	ShortID       string    `json:"short_id,omitempty"`
	Facility      Facility  `json:"facility,omitempty"`
	Plan          Plan      `json:"plan,omitempty"`
	Provisionable bool      `json:"provisionable,omitempty"`
	Spare         bool      `json:"spare,omitempty"`
	SwitchUUID    string    `json:"switch_uuid,omitempty"`
	Intervals     int       `json:"intervals,omitempty"`
	CurrentPeriod int       `json:"current_period,omitempty"`
	Href          string    `json:"href,omitempty"`
	Project       Project   `json:"project,omitempty"`
	Device        *Device   `json:"device,omitempty"`
	CreatedAt     Timestamp `json:"created_at,omitempty"`
}

type hardwareReservationRoot struct {
	HardwareReservations []HardwareReservation `json:"hardware_reservations"`
	Meta                 meta                  `json:"meta"`
}

// List returns all hardware reservations for a given project
func (s *HardwareReservationServiceOp) List(projectID string, opts *ListOptions) (reservations []HardwareReservation, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, hardwareReservationBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(hardwareReservationRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		reservations = append(reservations, subset.HardwareReservations...)
		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Get returns a single hardware reservation
func (s *HardwareReservationServiceOp) Get(hardwareReservationdID string, opts *GetOptions) (*HardwareReservation, *Response, error) {
	if validateErr := ValidateUUID(hardwareReservationdID); validateErr != nil {
		return nil, nil, validateErr
	}
	hardwareReservation := new(HardwareReservation)

	endpointPath := path.Join(hardwareReservationBasePath, hardwareReservationdID)
	apiPathQuery := opts.WithQuery(endpointPath)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, hardwareReservation)
	if err != nil {
		return nil, resp, err
	}

	return hardwareReservation, resp, err
}

// Move a hardware reservation to another project
func (s *HardwareReservationServiceOp) Move(hardwareReservationdID, projectID string) (*HardwareReservation, *Response, error) {
	if validateErr := ValidateUUID(hardwareReservationdID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	hardwareReservation := new(HardwareReservation)
	apiPath := path.Join(hardwareReservationBasePath, hardwareReservationdID, "move")
	body := map[string]string{}
	body["project_id"] = projectID

	resp, err := s.client.DoRequest("POST", apiPath, body, hardwareReservation)
	if err != nil {
		return nil, resp, err
	}

	return hardwareReservation, resp, err
}
