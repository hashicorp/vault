package pcf

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/hashicorp/go-hclog"
)

const (
	AuthUsername = "username"
	AuthPassword = "password"

	FoundServiceGUID = "1bf2e7f6-2d1d-41ec-501c-c70"
	FoundAppGUID     = "2d3e834a-3a25-4591-974c-fa5626d5d0a1"
	FoundOrgGUID     = "34a878d0-c2f9-4521-ba73-a9f664e82c7bf"
	FoundSpaceGUID   = "3d2eba6b-ef19-44d5-91dd-1975b0db5cc9"

	UnfoundServiceGUID = "service-id-unfound"
	UnfoundAppGUID     = "app-id-unfound"
	UnfoundOrgID       = "org-id-unfound"
	UnfoundSpaceGUID   = "space-id-unfound"
)

var (
	testServerUrl = ""
	logger        = hclog.Default()
)

func MockServer(loud bool) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if loud {
			logger.Info(fmt.Sprintf("%+v", r))
		}

		// Below, 200's are returned by default, but are included anyways for explicitness.
		pathFields := strings.Split(r.URL.EscapedPath(), "/")
		lastPathField := pathFields[len(pathFields)-1]
		switch lastPathField {
		case "token":
			w.Header().Add("Content-Type", "application/json;charset=UTF-8")
			w.WriteHeader(200)
			w.Write([]byte(tokenResponse))

		case "info":
			w.WriteHeader(200)
			w.Write([]byte(strings.Replace(infoResponse, "{{TEST_URL}}", testServerUrl, -1)))

		case FoundServiceGUID:
			w.WriteHeader(200)
			w.Write([]byte(serviceInstanceResponse))

		case UnfoundServiceGUID:
			w.WriteHeader(404)
			w.Write([]byte(unfoundServiceInstanceResponse))

		case FoundAppGUID:
			w.WriteHeader(200)
			w.Write([]byte(appResponse))

		case UnfoundAppGUID:
			w.WriteHeader(404)
			w.Write([]byte(unfoundAppResponse))

		case FoundOrgGUID:
			w.WriteHeader(200)
			w.Write([]byte(orgResponse))

		case UnfoundOrgID:
			w.WriteHeader(404)
			w.Write([]byte(unfoundOrgResponse))

		case FoundSpaceGUID:
			w.WriteHeader(200)
			w.Write([]byte(spaceResponse))

		case UnfoundSpaceGUID:
			w.WriteHeader(404)
			w.Write([]byte(unfoundSpaceResponse))

		default:
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unexpected object identifier: %s", lastPathField)))
		}
	}))
	testServerUrl = testServer.URL
	return testServer
}

const (
	tokenResponse = `{
	"access_token": "eyJhbGciOiJSUzI1NiIsImprdSI6Imh0dHBzOi8vdWFhLmRldi5jZmRldi5zaC90b2tlbl9rZXlzIiwia2lkIjoia2V5LTEiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiIxM2NiMzAyYjFjNjY0MDdkOWY3MDM2YzJjMmUxZDEyMCIsInN1YiI6IjYxMWM3ZWVhLWZmZDAtNGU5OC04MmYwLWY0YjU0YWZmNmRjYiIsInNjb3BlIjpbImNsaWVudHMucmVhZCIsIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsInNjaW0ucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJ1YWEudXNlciIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy5yZWFkIiwiY2xvdWRfY29udHJvbGxlci5yZWFkIiwicGFzc3dvcmQud3JpdGUiLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwibmV0d29yay5hZG1pbiIsImRvcHBsZXIuZmlyZWhvc2UiLCJzY2ltLndyaXRlIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI2MTFjN2VlYS1mZmQwLTRlOTgtODJmMC1mNGI1NGFmZjZkY2IiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW4iLCJhdXRoX3RpbWUiOjE1NTgzNzUwODksInJldl9zaWciOiIxOTA1YTEzOSIsImlhdCI6MTU1ODM3NTA4OSwiZXhwIjoxNTU4Mzc1Njg5LCJpc3MiOiJodHRwczovL3VhYS5kZXYuY2ZkZXYuc2gvb2F1dGgvdG9rZW4iLCJ6aWQiOiJ1YWEiLCJhdWQiOlsic2NpbSIsImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCIsImNmIiwiY2xpZW50cyIsInVhYSIsIm9wZW5pZCIsImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJuZXR3b3JrIl19.KSdNhoQSTCh_3zJPLvxeAhEyAfVTvHN1mKprHqfDJJ79WaaEsUM-mLO68QWPvBgON5dx8dOE8GaQw--xpqpqNwncb7MN8jmz_lZxgw-6oOf_O-bYJmGsaxX-ETlMLKvuqUljSC5KvB16zBkRtAP2IhQsMOV-PGdx2Lz4CqBkzALHL4MUlnaaI6Z1O-zMVhFFunpmY-mYZqaHNw_35cNohieehq1TrrqVdHCiNkNVYi7LQPS93Ow8VC6I3GFNzNr6EAjmHu9tEq3sTKAfsBg8zEWjB_25cpiWW5gL-dPhZd4KSgp3wOh1K4kpWw7NKpLnPxf7mcRH4IgNDZPJqkqAjA",
	"token_type": "bearer",
	"id_token": "eyJhbGciOiJSUzI1NiIsImprdSI6Imh0dHBzOi8vdWFhLmRldi5jZmRldi5zaC90b2tlbl9rZXlzIiwia2lkIjoia2V5LTEiLCJ0eXAiOiJKV1QifQ.eyJzdWIiOiI2MTFjN2VlYS1mZmQwLTRlOTgtODJmMC1mNGI1NGFmZjZkY2IiLCJhdWQiOlsiY2YiXSwiaXNzIjoiaHR0cHM6Ly91YWEuZGV2LmNmZGV2LnNoL29hdXRoL3Rva2VuIiwiZXhwIjoxNTU4Mzc1Njg5LCJpYXQiOjE1NTgzNzUwODksImFtciI6WyJwd2QiXSwiYXpwIjoiY2YiLCJzY29wZSI6WyJvcGVuaWQiXSwiZW1haWwiOiJhZG1pbiIsInppZCI6InVhYSIsIm9yaWdpbiI6InVhYSIsImp0aSI6IjEzY2IzMDJiMWM2NjQwN2Q5ZjcwMzZjMmMyZTFkMTIwIiwicHJldmlvdXNfbG9nb25fdGltZSI6MTU1ODM3NDk0NTEyMCwiZW1haWxfdmVyaWZpZWQiOnRydWUsImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX25hbWUiOiJhZG1pbiIsInJldl9zaWciOiIxOTA1YTEzOSIsInVzZXJfaWQiOiI2MTFjN2VlYS1mZmQwLTRlOTgtODJmMC1mNGI1NGFmZjZkY2IiLCJhdXRoX3RpbWUiOjE1NTgzNzUwODl9.eOv9O17i1naYiycCwlXFu2Xh2xjBRNBagq61AX1y2Upb7ek42VFaAi92PAZN9rmcU9i3trvERen0Hv7aIottLM7U-MTKMBnHXjqr1fY5oWyWxGruWsM0T9RBu4g9dbs8hyqIh_be9KdiL4PSybChV7-RspF1kMa58OUvpgQbQhgOMMWKKODYVXeeY8z241octX_ST-5tZv_josk12sworPQbZCwA5QbUjmCNSc_fHg9xe4Ra_Wecq3hmmspHrHW8gTc6ggoWUzxbbCKo1rF2PIVHzJ_61cLaHBepax9DvhCYnSJtDjlG5lPy41dxc01dOAD-JLEaV-CigtrWntFUXQ",
	"refresh_token": "eyJhbGciOiJSUzI1NiIsImprdSI6Imh0dHBzOi8vdWFhLmRldi5jZmRldi5zaC90b2tlbl9rZXlzIiwia2lkIjoia2V5LTEiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiIwOTRlYWQ0ZThiYWM0Nzk1ODJmMDI2ZmMwMjUwNTA2Yy1yIiwic3ViIjoiNjExYzdlZWEtZmZkMC00ZTk4LTgyZjAtZjRiNTRhZmY2ZGNiIiwiaWF0IjoxNTU4Mzc1MDg5LCJleHAiOjE1NjA5NjcwODksImNpZCI6ImNmIiwiY2xpZW50X2lkIjoiY2YiLCJpc3MiOiJodHRwczovL3VhYS5kZXYuY2ZkZXYuc2gvb2F1dGgvdG9rZW4iLCJ6aWQiOiJ1YWEiLCJhdWQiOlsic2NpbSIsImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCIsImNmIiwiY2xpZW50cyIsInVhYSIsIm9wZW5pZCIsImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJuZXR3b3JrIl0sImdyYW50ZWRfc2NvcGVzIjpbImNsaWVudHMucmVhZCIsIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsInNjaW0ucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJ1YWEudXNlciIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy5yZWFkIiwiY2xvdWRfY29udHJvbGxlci5yZWFkIiwicGFzc3dvcmQud3JpdGUiLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwibmV0d29yay5hZG1pbiIsImRvcHBsZXIuZmlyZWhvc2UiLCJzY2ltLndyaXRlIl0sImFtciI6WyJwd2QiXSwiYXV0aF90aW1lIjoxNTU4Mzc1MDg5LCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX25hbWUiOiJhZG1pbiIsIm9yaWdpbiI6InVhYSIsInVzZXJfaWQiOiI2MTFjN2VlYS1mZmQwLTRlOTgtODJmMC1mNGI1NGFmZjZkY2IiLCJyZXZfc2lnIjoiMTkwNWExMzkifQ.LFkoBtAWGL1x1bUo0ak16f-NeWpBS6NZspVwzaVhBv4xg7qxDryUayE5M2BQOGMb4tZLOU2cYyO2uu4li70u0LgJk7k3OZ0-hxKvjX4sJcoiLlJFCEsFzq_yG6iUFnA2w2kA70IQtACvAAHO--Jz0L1QGA8ebt20z7Rup0FufyDJFFevhbppzYb6AfghhnrB-yZbZU9rPq4Q8DWDTN0nMOBn05CA52NRKoj2157JXLRimEG7SZW6dhXUhdjbCvSz1WKiG6fS3fHK5ncqyQtuSqLfI0Naq1v77wfSzbvc0MB-IM4CPYc-ODhWbHFoV1z8kV6dWXm2ng7OyZe3u3A7Fw",
	"expires_in": 599,
	"scope": "clients.read openid routing.router_groups.write scim.read cloud_controller.admin uaa.user routing.router_groups.read cloud_controller.read password.write cloud_controller.write network.admin doppler.firehose scim.write",
	"jti": "13cb302b1c66407d9f7036c2c2e1d120"
}`

	infoResponse = `{
	"name": "",
	"build": "",
	"support": "",
	"version": 0,
	"description": "",
	"authorization_endpoint": "https://login.dev.cfdev.sh",
	"token_endpoint": "{{TEST_URL}}",
	"min_cli_version": null,
	"min_recommended_cli_version": null,
	"app_ssh_endpoint": "ssh.dev.cfdev.sh:2222",
	"app_ssh_host_key_fingerprint": "96:4d:89:2d:39:18:bc:16:e1:d3:d8:44:f8:16:af:85",
	"app_ssh_oauth_client": "ssh-proxy",
	"doppler_logging_endpoint": "wss://doppler.dev.cfdev.sh:443",
	"api_version": "2.133.0",
	"osbapi_version": "2.14",
	"routing_endpoint": "https://api.dev.cfdev.sh/routing",
	"user": "611c7eea-ffd0-4e98-82f0-f4b54aff6dcb"
}`

	serviceInstanceResponse = `{
	"metadata": {
		"guid": "1bf2e7f6-2d1d-41ec-501c-c70",
		"url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70",
		"created_at": "2016-06-08T16:41:29Z",
		"updated_at": "2016-06-08T16:41:26Z"
	},
	"entity": {
		"name": "name-1508",
		"credentials": {
			"creds-key-38": "creds-val-38"
		},
		"service_guid": "a14baddf-1ccc-5299-0152-ab9s49de4422",
		"service_plan_guid": "779d2df0-9cdd-48e8-9781-ea05301cedb1",
		"space_guid": "3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
		"gateway_data": null,
		"dashboard_url": null,
		"type": "managed_service_instance",
		"last_operation": {
			"type": "create",
			"state": "succeeded",
			"description": "service broker-provided description",
			"updated_at": "2016-06-08T16:41:29Z",
			"created_at": "2016-06-08T16:41:29Z"
		},
		"tags": [
			"accounting",
			"mongodb"
		],
		"space_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
		"service_url": "/v2/services/a14baddf-1ccc-5299-0152-ab9s49de4422",
		"service_plan_url": "/v2/service_plans/779d2df0-9cdd-48e8-9781-ea05301cedb1",
		"service_bindings_url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70/service_bindings",
		"service_keys_url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70/service_keys",
		"routes_url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70/routes",
		"shared_from_url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70/shared_from",
		"shared_to_url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70/shared_to",
		"service_instance_parameters_url": "/v2/service_instances/1bf2e7f6-2d1d-41ec-501c-c70/parameters"
	}
}`

	unfoundServiceInstanceResponse = `{
	"description": "The service instance could not be found: service-id-unfound",
	"error_code": "CF-ServiceInstanceNotFound",
	"code": 60004
}`

	appResponse = `{
	"metadata": {
		"guid": "2d3e834a-3a25-4591-974c-fa5626d5d0a1",
		"url": "/v2/apps/2d3e834a-3a25-4591-974c-fa5626d5d0a1",
		"created_at": "2016-06-08T16:41:44Z",
		"updated_at": "2016-06-08T16:41:44Z"
	},
	"entity": {
		"name": "name-2401",
		"production": false,
		"space_guid": "3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
		"stack_guid": "7e03186d-a438-4285-b3b7-c426532e1df2",
		"buildpack": null,
		"detected_buildpack": null,
		"detected_buildpack_guid": null,
		"environment_json": null,
		"memory": 1024,
		"instances": 1,
		"disk_quota": 1024,
		"state": "STOPPED",
		"version": "df19a7ea-2003-4ecb-a909-e630e43f2719",
		"command": null,
		"console": false,
		"debug": null,
		"staging_task_id": null,
		"package_state": "PENDING",
		"health_check_http_endpoint": "",
		"health_check_type": "port",
		"health_check_timeout": null,
		"staging_failed_reason": null,
		"staging_failed_description": null,
		"diego": false,
		"docker_image": null,
		"docker_credentials": {
			"username": null,
			"password": null
		},
		"package_updated_at": "2016-06-08T16:41:45Z",
		"detected_start_command": "",
		"enable_ssh": true,
		"ports": null,
		"space_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
		"stack_url": "/v2/stacks/7e03186d-a438-4285-b3b7-c426532e1df2",
		"routes_url": "/v2/apps/2d3e834a-3a25-4591-974c-fa5626d5d0a1/routes",
		"events_url": "/v2/apps/2d3e834a-3a25-4591-974c-fa5626d5d0a1/events",
		"service_bindings_url": "/v2/apps/2d3e834a-3a25-4591-974c-fa5626d5d0a1/service_bindings",
		"route_mappings_url": "/v2/apps/2d3e834a-3a25-4591-974c-fa5626d5d0a1/route_mappings"
	}
}`

	unfoundAppResponse = `{
	"description": "The app could not be found: app-id-unfound",
	"error_code": "CF-AppNotFound",
	"code": 100004
}`

	orgResponse = `{
	"metadata": {
		"guid": "34a878d0-c2f9-4521-ba73-a9f664e82c7bf",
		"url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf",
		"created_at": "2019-05-17T22:49:40Z",
		"updated_at": "2019-05-17T22:49:40Z"
	},
	"entity": {
		"name": "system",
		"billing_enabled": false,
		"quota_definition_guid": "b172ff20-ae6d-4a13-a554-dc22f3844fb0",
		"status": "active",
		"default_isolation_segment_guid": null,
		"quota_definition_url": "/v2/quota_definitions/b172ff20-ae6d-4a13-a554-dc22f3844fb0",
		"spaces_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/spaces",
		"domains_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/domains",
		"private_domains_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/private_domains",
		"users_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/users",
		"managers_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/managers",
		"billing_managers_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/billing_managers",
		"auditors_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/auditors",
		"app_events_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/app_events",
		"space_quota_definitions_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf/space_quota_definitions"
	}
}`

	unfoundOrgResponse = `{
	"description": "The organization could not be found: org-id-unfound",
	"error_code": "CF-OrganizationNotFound",
	"code": 30003
}`

	spaceResponse = `{
	"metadata": {
		"guid": "3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
		"url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
		"created_at": "2019-05-17T22:53:30Z",
		"updated_at": "2019-05-17T22:53:30Z"
	},
	"entity": {
		"name": "cfdev-space",
		"organization_guid": "34a878d0-c2f9-4521-ba73-a9f664e82c7bf",
		"space_quota_definition_guid": null,
		"isolation_segment_guid": null,
		"allow_ssh": true,
		"organization_url": "/v2/organizations/34a878d0-c2f9-4521-ba73-a9f664e82c7bf",
		"developers_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/developers",
		"managers_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/managers",
		"auditors_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/auditors",
		"apps_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/apps",
		"routes_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/routes",
		"domains_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/domains",
		"service_instances_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/service_instances",
		"app_events_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/app_events",
		"events_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/events",
		"security_groups_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/security_groups",
		"staging_security_groups_url": "/v2/spaces/3d2eba6b-ef19-44d5-91dd-1975b0db5cc9/staging_security_groups"
	}
}`

	unfoundSpaceResponse = `{
	"description": "The app space could not be found: space-id-unfound",
	"error_code": "CF-SpaceNotFound",
	"code": 40004
}`
)
