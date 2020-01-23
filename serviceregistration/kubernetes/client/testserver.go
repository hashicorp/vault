package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	TestNamespace = "default"
	TestPodname   = "shell-demo"
)

// TestServer returns an http test server that can be used to test
// Kubernetes client code. It returns its current patches as a map
// so the caller can check current state. Calling the closeFunc
// at the end closes the test server. Responses are provided using
// real responses that have been captured from the Kube API.
func TestServer(t *testing.T) (currentPatches map[string]*Patch, closeFunc func()) {
	currentPatches = make(map[string]*Patch)

	// We're going to have multiple close funcs to call.
	var closers []func()
	closeFunc = func() {
		for _, closer := range closers {
			closer()
		}
	}

	// Edit the url scheme for our test server, and use our
	// fixtures to supply the token and ca.crt.
	scheme = "http://"

	// Plant our token in a place where it can be read for the config.
	tmpToken, err := ioutil.TempFile("", "token")
	if err != nil {
		t.Fatal(err)
	}
	closers = append(closers, func() {
		os.Remove(tmpToken.Name())
	})
	if _, err = tmpToken.WriteString(token); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	if err := tmpToken.Close(); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	tokenFile = tmpToken.Name()

	tmpCACrt, err := ioutil.TempFile("", "ca.crt")
	if err != nil {
		closeFunc()
		t.Fatal(err)
	}
	closers = append(closers, func() {
		os.Remove(tmpCACrt.Name())
	})
	if _, err = tmpCACrt.WriteString(caCrt); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	if err := tmpCACrt.Close(); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	rootCAFile = tmpCACrt.Name()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		namespace, podName, err := parsePath(r.URL.Path)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unable to parse %s: %s", r.URL.Path, err.Error())))
			return
		}

		switch {
		case namespace != TestNamespace, podName != TestPodname:
			w.WriteHeader(404)
			w.Write([]byte(notFoundResponse))
			return
		case r.Method == http.MethodGet:
			w.WriteHeader(200)
			w.Write([]byte(getPodResponse))
			return
		case r.Method == http.MethodPatch:
			var patches []interface{}
			if err := json.NewDecoder(r.Body).Decode(&patches); err != nil {
				w.WriteHeader(400)
				w.Write([]byte(fmt.Sprintf("unable to decode patches %s: %s", r.URL.Path, err.Error())))
				return
			}
			for _, patch := range patches {
				patchMap := patch.(map[string]interface{})
				p := patchMap["path"].(string)
				currentPatches[p] = &Patch{
					Operation: Parse(patchMap["op"].(string)),
					Path:      p,
					Value:     patchMap["value"],
				}
			}
			w.WriteHeader(200)
			w.Write([]byte(updatePodTagsResponse))
			return
		default:
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unexpected request method: %s", r.Method)))
		}
	}))
	closers = append(closers, ts.Close)

	// ts.URL example: http://127.0.0.1:35681
	urlFields := strings.Split(ts.URL, "://")
	if len(urlFields) != 2 {
		closeFunc()
		t.Fatal("received unexpected test url: " + ts.URL)
	}
	urlFields = strings.Split(urlFields[1], ":")
	if len(urlFields) != 2 {
		closeFunc()
		t.Fatal("received unexpected test url: " + ts.URL)
	}
	if err := os.Setenv(EnvVarKubernetesServiceHost, urlFields[0]); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	if err := os.Setenv(EnvVarKubernetesServicePort, urlFields[1]); err != nil {
		closeFunc()
		t.Fatal(err)
	}
	return currentPatches, closeFunc
}

// The path should be formatted like this:
// fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName)
func parsePath(urlPath string) (namespace, podName string, err error) {
	original := urlPath
	podName = path.Base(urlPath)
	urlPath = strings.TrimSuffix(urlPath, "/pods/"+podName)
	namespace = path.Base(urlPath)
	if original != fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName) {
		return "", "", fmt.Errorf("received unexpected path: %s", original)
	}
	return namespace, podName, nil
}

// These are examples captured from real life.
const (
	caCrt = `-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTE5MTIxMDIzMDUxOVoXDTI5MTIwODIzMDUxOVowFTETMBEGA1UE
AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANFi
/RIdMHd865X6JygTb9riX01DA3QnR+RoXDXNnj8D3LziLG2n8ItXMJvWbU3sxxyy
nX9HxJ0SIeexj1cYzdQBtJDjO1/PeuKc4CZ7zCukCAtHz8mC7BDPOU7F7pggpcQ0
/t/pa2m22hmCu8aDF9WlUYHtJpYATnI/A5vz/VFLR9daxmkl59Qo3oHITj7vAzSx
/75r9cibpQyJ+FhiHOZHQWYY2JYw2g4v5hm5hg5SFM9yFcZ75ISI9ebyFFIl9iBY
zAk9jqv1mXvLr0Q39AVwMTamvGuap1oocjM9NIhQvaFL/DNqF1ouDQjCf5u2imLc
TraO1/2KO8fqwOZCOrMCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQBtVZCwCPqUUUpIClAlE9nc2fo2bTs9gsjXRmqdQ5oaSomSLE93
aJWYFuAhxPXtlApbLYZfW2m1sM3mTVQN60y0uE4e1jdSN1ErYQ9slJdYDAMaEmOh
iSexj+Nd1scUiMHV9lf3ps5J8sYeCpwZX3sPmw7lqZojTS12pANBDcigsaj5RRyN
9GyP3WkSQUsTpWlDb9Fd+KNdkCVw7nClIpBPA2KW4BQKw/rNSvOFD61mbzc89lo0
Q9IFGQFFF8jO18lbyWqnRBGXcS4/G7jQ3S7C121d14YLUeAYOM7pJykI1g4CLx9y
vitin0L6nprauWkKO38XgM4T75qKZpqtiOcT
-----END CERTIFICATE-----
`
	getPodResponse = `{
  "kind": "Pod",
  "apiVersion": "v1",
  "metadata": {
    "name": "shell-demo",
	"labels": {"fizz": "buzz"},
    "namespace": "default",
    "selfLink": "/api/v1/namespaces/default/pods/shell-demo",
    "uid": "7ecb93ff-aa64-426d-b330-2c0b2c0957a2",
    "resourceVersion": "87798",
    "creationTimestamp": "2020-01-10T19:22:40Z",
    "annotations": {
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"shell-demo\",\"namespace\":\"default\"},\"spec\":{\"containers\":[{\"image\":\"nginx\",\"name\":\"nginx\",\"volumeMounts\":[{\"mountPath\":\"/usr/share/nginx/html\",\"name\":\"shared-data\"}]}],\"dnsPolicy\":\"Default\",\"hostNetwork\":true,\"volumes\":[{\"emptyDir\":{},\"name\":\"shared-data\"}]}}\n"
    }
  },
  "spec": {
    "volumes": [{
      "name": "shared-data",
      "emptyDir": {}
    }, {
      "name": "default-token-5fjt9",
      "secret": {
        "secretName": "default-token-5fjt9",
        "defaultMode": 420
      }
    }],
    "containers": [{
      "name": "nginx",
      "image": "nginx",
      "resources": {},
      "volumeMounts": [{
        "name": "shared-data",
        "mountPath": "/usr/share/nginx/html"
      }, {
        "name": "default-token-5fjt9",
        "readOnly": true,
        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
      }],
      "terminationMessagePath": "/dev/termination-log",
      "terminationMessagePolicy": "File",
      "imagePullPolicy": "Always"
    }],
    "restartPolicy": "Always",
    "terminationGracePeriodSeconds": 30,
    "dnsPolicy": "Default",
    "serviceAccountName": "default",
    "serviceAccount": "default",
    "nodeName": "minikube",
    "hostNetwork": true,
    "securityContext": {},
    "schedulerName": "default-scheduler",
    "tolerations": [{
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "effect": "NoExecute",
      "tolerationSeconds": 300
    }, {
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "effect": "NoExecute",
      "tolerationSeconds": 300
    }],
    "priority": 0,
    "enableServiceLinks": true
  },
  "status": {
    "phase": "Running",
    "conditions": [{
      "type": "Initialized",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T19:22:40Z"
    }, {
      "type": "Ready",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T20:20:55Z"
    }, {
      "type": "ContainersReady",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T20:20:55Z"
    }, {
      "type": "PodScheduled",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T19:22:40Z"
    }],
    "hostIP": "192.168.99.100",
    "podIP": "192.168.99.100",
    "podIPs": [{
      "ip": "192.168.99.100"
    }],
    "startTime": "2020-01-10T19:22:40Z",
    "containerStatuses": [{
      "name": "nginx",
      "state": {
        "running": {
          "startedAt": "2020-01-10T20:20:55Z"
        }
      },
      "lastState": {
        "terminated": {
          "exitCode": 0,
          "reason": "Completed",
          "startedAt": "2020-01-10T19:22:53Z",
          "finishedAt": "2020-01-10T20:12:03Z",
          "containerID": "docker://ed8bc068cd313ea5adb72780e8015ab09ecb61ea077e39304b4a3fe581f471c4"
        }
      },
      "ready": true,
      "restartCount": 1,
      "image": "nginx:latest",
      "imageID": "docker-pullable://nginx@sha256:8aa7f6a9585d908a63e5e418dc5d14ae7467d2e36e1ab4f0d8f9d059a3d071ce",
      "containerID": "docker://a8ee34466791bc6f082f271f40cdfc43625cea81831b1029b1e90b4f6949f6df",
      "started": true
    }],
    "qosClass": "BestEffort"
  }
}
`
	notFoundResponse = `{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods \"shell-dem\" not found",
  "reason": "NotFound",
  "details": {
    "name": "shell-dem",
    "kind": "pods"
  },
  "code": 404
}
`
	token                 = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjZVQU91ckJYcTZKRHQtWHpaOExib2EyUlFZQWZObms2d25mY3ZtVm1NNUUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tNWZqdDkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImY0NGUyMDIxLTU2YWItNDEzNC1hMjMxLTBlMDJmNjhmNzJhNiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.hgMbuT0hlxG04fDvI_Iyxtbwc8M-i3q3K7CqIGC_jYSjVlyezHN_0BeIB3rE0_M2xvbIs6chsWFZVsK_8Pj6ho7VT0x5PWy5n6KsqTBz8LPpjWpsaxpYQos0RzgA3KLnuzZE8Cl-v-PwWQK57jgbS4AdlXujQXdtLXJNwNAKI0pvCASA6UXP55_X845EsJkyT1J-bURSS3Le3g9A4pDoQ_MUv7hqa-p7yQEtFfYCkq1KKrUJZMRjmS4qda1rg-Em-dw9RFvQtPodRYF0DKT7A7qgmLUfIkuky3NnsQtvaUo8ZVtUiwIEfRdqw1oQIY4CSYz-wUl2xZa7n2QQBROE7w`
	updatePodTagsResponse = `{
  "kind": "Pod",
  "apiVersion": "v1",
  "metadata": {
    "name": "shell-demo",
    "namespace": "default",
    "selfLink": "/api/v1/namespaces/default/pods/shell-demo",
    "uid": "7ecb93ff-aa64-426d-b330-2c0b2c0957a2",
    "resourceVersion": "96433",
    "creationTimestamp": "2020-01-10T19:22:40Z",
    "labels": {
      "fizz": "buzz",
      "foo": "bar"
    },
    "annotations": {
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"shell-demo\",\"namespace\":\"default\"},\"spec\":{\"containers\":[{\"image\":\"nginx\",\"name\":\"nginx\",\"volumeMounts\":[{\"mountPath\":\"/usr/share/nginx/html\",\"name\":\"shared-data\"}]}],\"dnsPolicy\":\"Default\",\"hostNetwork\":true,\"volumes\":[{\"emptyDir\":{},\"name\":\"shared-data\"}]}}\n"
    }
  },
  "spec": {
    "volumes": [{
      "name": "shared-data",
      "emptyDir": {}
    }, {
      "name": "default-token-5fjt9",
      "secret": {
        "secretName": "default-token-5fjt9",
        "defaultMode": 420
      }
    }],
    "containers": [{
      "name": "nginx",
      "image": "nginx",
      "resources": {},
      "volumeMounts": [{
        "name": "shared-data",
        "mountPath": "/usr/share/nginx/html"
      }, {
        "name": "default-token-5fjt9",
        "readOnly": true,
        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
      }],
      "terminationMessagePath": "/dev/termination-log",
      "terminationMessagePolicy": "File",
      "imagePullPolicy": "Always"
    }],
    "restartPolicy": "Always",
    "terminationGracePeriodSeconds": 30,
    "dnsPolicy": "Default",
    "serviceAccountName": "default",
    "serviceAccount": "default",
    "nodeName": "minikube",
    "hostNetwork": true,
    "securityContext": {},
    "schedulerName": "default-scheduler",
    "tolerations": [{
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "effect": "NoExecute",
      "tolerationSeconds": 300
    }, {
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "effect": "NoExecute",
      "tolerationSeconds": 300
    }],
    "priority": 0,
    "enableServiceLinks": true
  },
  "status": {
    "phase": "Running",
    "conditions": [{
      "type": "Initialized",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T19:22:40Z"
    }, {
      "type": "Ready",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T20:20:55Z"
    }, {
      "type": "ContainersReady",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T20:20:55Z"
    }, {
      "type": "PodScheduled",
      "status": "True",
      "lastProbeTime": null,
      "lastTransitionTime": "2020-01-10T19:22:40Z"
    }],
    "hostIP": "192.168.99.100",
    "podIP": "192.168.99.100",
    "podIPs": [{
      "ip": "192.168.99.100"
    }],
    "startTime": "2020-01-10T19:22:40Z",
    "containerStatuses": [{
      "name": "nginx",
      "state": {
        "running": {
          "startedAt": "2020-01-10T20:20:55Z"
        }
      },
      "lastState": {
        "terminated": {
          "exitCode": 0,
          "reason": "Completed",
          "startedAt": "2020-01-10T19:22:53Z",
          "finishedAt": "2020-01-10T20:12:03Z",
          "containerID": "docker://ed8bc068cd313ea5adb72780e8015ab09ecb61ea077e39304b4a3fe581f471c4"
        }
      },
      "ready": true,
      "restartCount": 1,
      "image": "nginx:latest",
      "imageID": "docker-pullable://nginx@sha256:8aa7f6a9585d908a63e5e418dc5d14ae7467d2e36e1ab4f0d8f9d059a3d071ce",
      "containerID": "docker://a8ee34466791bc6f082f271f40cdfc43625cea81831b1029b1e90b4f6949f6df",
      "started": true
    }],
    "qosClass": "BestEffort"
  }
}
`
)
