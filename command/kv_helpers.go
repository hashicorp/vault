package command

import (
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

// kvData is a helper struct for `listRecursive'.
type kvData struct {
	path   string      // Name (path) of the secret.
	secret *api.Secret // Data belonging to the secret (sub-paths too).
	err    error       // Errors when fetching the secret (if any).
}

// kvRecursiveListParams contains a list of parameters
// to be passed to `kvListRecursive'.
type kvListRecursiveParams struct {
	v2     *bool       // Flag to store if the API is v2.
	client *api.Client // Client object to use for fetching secrets.

	data   []*kvData      // List of all the secrets fetched recursively.
	depth  int32          // Depth of the recursion.
	track  int32          // Helper variable to track the recursion depth.
	filter *regexp.Regexp // Filter to be applied for matching keys.

	wg  sync.WaitGroup // Keep track of launched goroutines.
	sem chan int32     // Semaphone for throttling goroutines.
	mux sync.Mutex     // Mutex to append entries to `data'.
	tck *time.Ticker   // Avoid busy loops when polling for goroutines.
}

func kvReadRequest(client *api.Client, path string, params map[string]string) (*api.Secret, error) {
	r := client.NewRequest("GET", "/v1/"+path)
	for k, v := range params {
		r.Params.Set(k, v)
	}
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func kvPreflightVersionRequest(client *api.Client, path string) (string, int, error) {
	// We don't want to use a wrapping call here so save any custom value and
	// restore after
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)

	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if secret == nil {
		return "", 0, errors.New("nil response from pre-flight request")
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}

func isKVv2(path string, client *api.Client) (string, bool, error) {
	mountPath, version, err := kvPreflightVersionRequest(client, path)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

func addPrefixToVKVPath(p, mountPath, apiPrefix string) string {
	switch {
	case p == mountPath, p == strings.TrimSuffix(mountPath, "/"):
		return path.Join(mountPath, apiPrefix)
	default:
		p = strings.TrimPrefix(p, mountPath)
		return path.Join(mountPath, apiPrefix, p)
	}
}

func removePrefixFromVKVPath(p, apiPrefix string) string {
	return strings.Replace(p, apiPrefix, "", 1)
}

func getHeaderForMap(header string, data map[string]interface{}) string {
	maxKey := 0
	for k := range data {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}

	// 4 for the column spaces and 5 for the len("value")
	totalLen := maxKey + 4 + 5

	equalSigns := totalLen - (len(header) + 2)

	// If we have zero or fewer equal signs bump it back up to two on either
	// side of the header.
	if equalSigns <= 0 {
		equalSigns = 4
	}

	// If the number of equal signs is not divisible by two add a sign.
	if equalSigns%2 != 0 {
		equalSigns = equalSigns + 1
	}

	return fmt.Sprintf("%s %s %s", strings.Repeat("=", equalSigns/2), header, strings.Repeat("=", equalSigns/2))
}

func kvParseVersionsFlags(versions []string) []string {
	versionsOut := make([]string, 0, len(versions))
	for _, v := range versions {
		versionsOut = append(versionsOut, strutil.ParseStringSlice(v, ",")...)
	}

	return versionsOut
}

// kvListRecursive is a helper function for listing paths
// upto a given depth, recursively.
func kvListRecursive(r *kvListRecursiveParams, base string, sub []interface{}) {
	defer func() {
		// Increment the go-routine tracker.
		atomic.AddInt32(&r.track, 1)

		// Decrement the wait-group count.
		r.wg.Done()
	}()

	// Wait in the queue.
	r.sem <- int32(1)

	// No need to recursively list values that don't have any sub-paths.
	if !strings.HasSuffix(base, "/") {
		return
	}

	// Strip trailing slash, because we want to append
	// it with each of the children's paths.
	base = ensureNoTrailingSlash(base)

	// Calculate the recursion depth of this call.
	d := int32(strings.Count(base, "/"))
	if *r.v2 {
		d--
	}

	// Base case for return.
	if len(sub) <= 0 || (r.depth != -1 && d >= r.depth) {
		return
	}

	var tmp string

	// For each child, construct the new path and call recursively.
	for _, p := range sub {
		path := fmt.Sprintf("%s/%s", base, p)
		secret, err := r.client.Logical().List(path)

		tmp = path
		if *r.v2 {
			tmp = removePrefixFromVKVPath(tmp, "metadata/")
		}

		// Append the entry (that matches the filter) to the list.
		if r.filter.MatchString(tmp) {
			r.mux.Lock()
			r.data = append(r.data, &kvData{tmp, secret, err})
			r.mux.Unlock()
		}

		if err != nil {
			continue
		}

		// Only call if there are children.
		if s, ok := extractListData(secret); ok {
			r.wg.Add(1)
			go kvListRecursive(r, path, s)
		}
	}
}
