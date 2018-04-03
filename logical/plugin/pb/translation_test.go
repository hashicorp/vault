package pb

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
)

func TestTranslation_Errors(t *testing.T) {
	errs := []error{
		nil,
		errors.New("test"),
		errutil.UserError{Err: "test"},
		errutil.InternalError{Err: "test"},
		logical.CodedError(403, "test"),
		&logical.StatusBadRequest{Err: "test"},
		logical.ErrUnsupportedOperation,
		logical.ErrUnsupportedPath,
		logical.ErrInvalidRequest,
		logical.ErrPermissionDenied,
		logical.ErrMultiAuthzPending,
	}

	for _, err := range errs {
		pe := ErrToProtoErr(err)
		e := ProtoErrToErr(pe)

		if !reflect.DeepEqual(e, err) {
			t.Fatalf("Errs did not match: %#v, %#v", e, err)
		}
	}
}

func TestTranslation_StorageEntry(t *testing.T) {
	tCases := []*logical.StorageEntry{
		nil,
		&logical.StorageEntry{Key: "key", Value: []byte("value")},
		&logical.StorageEntry{Key: "key1", Value: []byte("value1"), SealWrap: true},
		&logical.StorageEntry{Key: "key1", SealWrap: true},
	}

	for _, c := range tCases {
		p := LogicalStorageEntryToProtoStorageEntry(c)
		e := ProtoStorageEntryToLogicalStorageEntry(p)

		if !reflect.DeepEqual(c, e) {
			t.Fatalf("Entries did not match: %#v, %#v", e, c)
		}
	}
}

func TestTranslation_Request(t *testing.T) {
	tCases := []*logical.Request{
		nil,
		&logical.Request{
			ID:                       "ID",
			ReplicationCluster:       "RID",
			Operation:                logical.CreateOperation,
			Path:                     "test/foo",
			ClientToken:              "token",
			ClientTokenAccessor:      "accessor",
			DisplayName:              "display",
			MountPoint:               "test",
			MountType:                "secret",
			MountAccessor:            "test-231234",
			ClientTokenRemainingUses: 1,
			EntityID:                 "tester",
			PolicyOverride:           true,
			Unauthenticated:          true,
			Connection: &logical.Connection{
				RemoteAddr: "localhost",
			},
		},
		&logical.Request{
			ID:                 "ID",
			ReplicationCluster: "RID",
			Operation:          logical.CreateOperation,
			Path:               "test/foo",
			Data: map[string]interface{}{
				"string": "string",
				"bool":   true,
				"array":  []interface{}{"1", "2"},
				"map": map[string]interface{}{
					"key": "value",
				},
			},
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				LeaseID: "LeaseID",
			},
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				DisplayName: "test",
				Policies:    []string{"test", "Test"},
				Metadata: map[string]string{
					"test": "test",
				},
				ClientToken: "token",
				Accessor:    "accessor",
				Period:      5 * time.Second,
				NumUses:     1,
				EntityID:    "id",
				Alias: &logical.Alias{
					MountType:     "type",
					MountAccessor: "accessor",
					Name:          "name",
				},
				GroupAliases: []*logical.Alias{
					&logical.Alias{
						MountType:     "type",
						MountAccessor: "accessor",
						Name:          "name",
					},
				},
			},
			Headers: map[string][]string{
				"X-Vault-Test": []string{"test"},
			},
			ClientToken:         "token",
			ClientTokenAccessor: "accessor",
			DisplayName:         "display",
			MountPoint:          "test",
			MountType:           "secret",
			MountAccessor:       "test-231234",
			WrapInfo: &logical.RequestWrapInfo{
				TTL:      time.Second,
				Format:   "token",
				SealWrap: true,
			},
			ClientTokenRemainingUses: 1,
			EntityID:                 "tester",
			PolicyOverride:           true,
			Unauthenticated:          true,
		},
	}

	for _, c := range tCases {
		p, err := LogicalRequestToProtoRequest(c)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ProtoRequestToLogicalRequest(p)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c, r) {
			t.Fatalf("Requests did not match: \n%#v, \n%#v", c, r)
		}
	}
}

func TestTranslation_Response(t *testing.T) {
	tCases := []*logical.Response{
		nil,
		&logical.Response{
			Data: map[string]interface{}{
				"data": "blah",
			},
			Warnings: []string{"warning"},
		},
		&logical.Response{
			Data: map[string]interface{}{
				"string": "string",
				"bool":   true,
				"array":  []interface{}{"1", "2"},
				"map": map[string]interface{}{
					"key": "value",
				},
			},
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				LeaseID: "LeaseID",
			},
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				DisplayName: "test",
				Policies:    []string{"test", "Test"},
				Metadata: map[string]string{
					"test": "test",
				},
				ClientToken: "token",
				Accessor:    "accessor",
				Period:      5 * time.Second,
				NumUses:     1,
				EntityID:    "id",
				Alias: &logical.Alias{
					MountType:     "type",
					MountAccessor: "accessor",
					Name:          "name",
				},
				GroupAliases: []*logical.Alias{
					&logical.Alias{
						MountType:     "type",
						MountAccessor: "accessor",
						Name:          "name",
					},
				},
			},
			WrapInfo: &wrapping.ResponseWrapInfo{
				TTL:             time.Second,
				Token:           "token",
				Accessor:        "accessor",
				CreationTime:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				WrappedAccessor: "wrapped-accessor",
				WrappedEntityID: "id",
				Format:          "token",
				CreationPath:    "test/foo",
				SealWrap:        true,
			},
		},
	}

	for _, c := range tCases {
		p, err := LogicalResponseToProtoResponse(c)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ProtoResponseToLogicalResponse(p)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c, r) {
			t.Fatalf("Requests did not match: \n%#v, \n%#v", c, r)
		}
	}
}
