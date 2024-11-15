// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package couchbase

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/go-version"
)

func CheckForOldCouchbaseVersion(hostname, username, password string) (is_old bool, err error) {

	//[TODO] handle list of hostnames

	resp, err := http.Get(fmt.Sprintf("http://%s:%s@%s:8091/pools", username, password, hostname))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type Pools struct {
		ImplementationVersion string `json:"implementationVersion"`
	}
	data := Pools{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	v, err := version.NewVersion(data.ImplementationVersion)

	v650, err := version.NewVersion("6.5.0-0000")
	if err != nil {
		return false, err
	}

	if v.LessThan(v650) {
		return true, nil
	}
	return false, nil

}

func getRootCAfromCouchbase(url string) (Base64pemCA string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(body), nil
}

func createUser(hostname string, port int, adminuser, adminpassword, username, password, rbacName, roles string) (err error) {
	v := url.Values{}

	v.Set("password", password)
	v.Add("roles", roles)
	v.Add("name", rbacName)

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("http://%s:%s@%s:%d/settings/rbac/users/local/%s",
			adminuser, adminpassword, hostname, port, username),
		strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		return fmt.Errorf("createUser returned %s", resp.Status)
	}
	return nil
}

func createGroup(hostname string, port int, adminuser, adminpassword, group, roles string) (err error) {
	v := url.Values{}

	v.Set("roles", roles)

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("http://%s:%s@%s:%d/settings/rbac/groups/%s",
			adminuser, adminpassword, hostname, port, group),
		strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		return fmt.Errorf("createGroup returned %s", resp.Status)
	}
	return nil
}

func waitForBucket(t *testing.T, address, username, password, bucketName string) {
	t.Logf("Waiting for bucket %s...", bucketName)
	f := func() error {
		return checkBucketReady(address, username, password, bucketName)
	}
	bo := backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 10)
	err := backoff.Retry(f, bo)
	if err != nil {
		t.Fatalf("bucket %s installed check failed: %s", bucketName, err)
	}
}

func checkBucketReady(address, username, password, bucket string) (err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s@%s:8091/sampleBuckets", username, password, address))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type installed []struct {
		Name        string `json:"name"`
		Installed   bool   `json:"installed"`
		QuotaNeeded int64  `json:"quotaNeeded"`
	}

	var iresult installed

	err = json.Unmarshal(body, &iresult)
	if err != nil {
		err := &backoff.PermanentError{
			Err: fmt.Errorf("error unmarshaling JSON %s", err),
		}
		return err
	}

	bucketFound := false
	for _, s := range iresult {
		if s.Name == bucket {
			bucketFound = true
			if s.Installed == true {
				return nil // Found & installed
			}
		}
	}

	err = fmt.Errorf("bucket not found")

	if !bucketFound {
		return backoff.Permanent(err)
	}
	return err
}
