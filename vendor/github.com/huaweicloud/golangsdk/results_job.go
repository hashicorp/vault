package golangsdk

import (
	"fmt"
	"strings"
)

type JobResponse struct {
	URI   string `json:"uri"`
	JobID string `json:"job_id"`
}

type JobStatus struct {
	Status     string                 `json:"status"`
	Entities   map[string]interface{} `json:"entities"`
	JobID      string                 `json:"job_id"`
	JobType    string                 `json:"job_type"`
	ErrorCode  string                 `json:"error_code"`
	FailReason string                 `json:"fail_reason"`
}

func (r Result) ExtractJobResponse() (*JobResponse, error) {
	job := new(JobResponse)
	err := r.ExtractInto(job)
	return job, err
}

func (r Result) ExtractJobStatus() (*JobStatus, error) {
	job := new(JobStatus)
	err := r.ExtractInto(job)
	return job, err
}

func GetJobEndpoint(endpoint string) string {
	n := strings.Index(endpoint[8:len(endpoint)], "/")
	if n == -1 {
		return endpoint
	}
	return endpoint[0 : n+8]
}

func WaitForJobSuccess(client *ServiceClient, uri string, secs int) error {
	uri = strings.Replace(uri, "v1", "v1.0", 1)

	return WaitFor(secs, func() (bool, error) {
		job := new(JobStatus)
		_, err := client.Get(GetJobEndpoint(client.Endpoint)+uri, &job, nil)
		if err != nil {
			return false, err
		}
		fmt.Printf("JobStatus: %+v.\n", job)

		if job.Status == "SUCCESS" {
			return true, nil
		}
		if job.Status == "FAIL" {
			err = fmt.Errorf("Job failed with code %s: %s.\n", job.ErrorCode, job.FailReason)
			return false, err
		}

		return false, nil
	})
}

func GetJobEntity(client *ServiceClient, uri string, label string) (interface{}, error) {
	uri = strings.Replace(uri, "v1", "v1.0", 1)

	job := new(JobStatus)
	_, err := client.Get(GetJobEndpoint(client.Endpoint)+uri, &job, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JobStatus: %+v.\n", job)

	if job.Status == "SUCCESS" {
		if e := job.Entities[label]; e != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("Unexpected conversion error in GetJobEntity.")
}
