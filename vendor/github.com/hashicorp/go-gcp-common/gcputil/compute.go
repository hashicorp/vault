package gcputil

import (
	"fmt"
	"regexp"
	"time"

	"google.golang.org/api/compute/v1"
)

var gcpLabelRegex = regexp.MustCompile("^(?P<key>[a-z]([\\w-]+)?):(?P<value>[\\w-]*)$")

func ParseGcpLabels(labels []string) (parsed map[string]string, invalid []string) {
	parsed = map[string]string{}
	invalid = []string{}

	for _, labelStr := range labels {
		matches := gcpLabelRegex.FindStringSubmatch(labelStr)
		if len(matches) == 0 {
			invalid = append(invalid, labelStr)
			continue
		}

		captureNames := gcpLabelRegex.SubexpNames()
		var keyPtr, valPtr *string
		for i, name := range captureNames {
			if name == "key" {
				keyPtr = &matches[i]
			} else if name == "value" {
				valPtr = &matches[i]
			}
		}

		if keyPtr == nil || valPtr == nil || len(*keyPtr) < 1 {
			invalid = append(invalid, labelStr)
			continue
		} else {
			parsed[*keyPtr] = *valPtr
		}
	}

	return parsed, invalid
}

type CustomJWTClaims struct {
	Google *GoogleJWTClaims `json:"google,omitempty"`
}

type GoogleJWTClaims struct {
	Compute *GCEIdentityMetadata `json:"compute_engine,omitempty"`
}

type GCEIdentityMetadata struct {
	// ProjectId is the ID for the project where you created the instance.
	ProjectId string `json:"project_id"  structs:"project_id" mapstructure:"project_id"`

	// ProjectNumber is the unique ID for the project where you created the instance.
	ProjectNumber int64 `json:"project_number" structs:"project_number" mapstructure:"project_number"`

	// Zone is the zone where the instance is located.
	Zone string `json:"zone" structs:"zone" mapstructure:"zone"`

	// InstanceId is the unique ID for the instance to which this token belongs. This ID is unique and never reused.
	InstanceId string `json:"instance_id" structs:"instance_id" mapstructure:"instance_id"`

	// InstanceName is the name of the instance to which this token belongs. This name can be reused by several
	// instances over time, so use the instance_id value to identify a unique instance ID.
	InstanceName string `json:"instance_name" structs:"instance_name" mapstructure:"instance_name"`

	// CreatedAt is a unix timestamp indicating when you created the instance.
	CreatedAt int64 `json:"instance_creation_timestamp" structs:"instance_creation_timestamp" mapstructure:"instance_creation_timestamp"`
}

// GetVerifiedInstance returns the Instance as described by the identity metadata or an error.
// If the instance has an invalid status or its creation timestamp does not match the metadata value,
// this  will return nil and an error.
func (meta *GCEIdentityMetadata) GetVerifiedInstance(gceClient *compute.Service) (*compute.Instance, error) {
	instance, err := gceClient.Instances.Get(meta.ProjectId, meta.Zone, meta.InstanceName).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to find instance associated with token: %v", err)
	}

	if !IsValidInstanceStatus(instance.Status) {
		return nil, fmt.Errorf("authenticating instance %s found but has invalid status '%s'", instance.Name, instance.Status)
	}

	// Parse the metadata CreatedAt into time.
	metaTime := time.Unix(meta.CreatedAt, 0)

	// Parse instance creationTimestamp into time.
	actualTime, err := time.Parse(time.RFC3339Nano, instance.CreationTimestamp)
	if err != nil {
		return nil, fmt.Errorf("instance 'creationTimestamp' field could not be parsed into time: %s", instance.CreationTimestamp)
	}

	// Return an error if the metadata creation timestamp is before the instance creation timestamp.
	delta := float64(metaTime.Sub(actualTime)) / float64(time.Second)
	if delta < -1 {
		return nil, fmt.Errorf("metadata instance_creation_timestamp %d is before instance's creation time %d", actualTime.Unix(), metaTime.Unix())
	}
	return instance, nil
}

var validInstanceStates map[string]struct{} = map[string]struct{}{
	"PROVISIONING": struct{}{},
	"RUNNING":      struct{}{},
	"STAGING":      struct{}{},
}

func IsValidInstanceStatus(status string) bool {
	_, ok := validInstanceStates[status]
	return ok
}
