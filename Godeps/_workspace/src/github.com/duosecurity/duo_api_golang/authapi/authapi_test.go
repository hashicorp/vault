package authapi

import (
    "testing"
    "fmt"
    "strings"
    "net/http"
    "net/http/httptest"
    "time"
    "github.com/duosecurity/duo_api_golang"
  )

func buildAuthApi(url string) *AuthApi {
      ikey := "eyekey"
      skey := "esskey"
      host := strings.Split(url, "//")[1]
      userAgent := "GoTestClient"
      return NewAuthApi(*duoapi.NewDuoApi(ikey,
        skey,
        host,
        userAgent,
        duoapi.SetTimeout(1*time.Second),
        duoapi.SetInsecure()))
}

// Timeouts are set to 1 second.  Take 15 seconds to respond and verify
// that the client times out.
func TestTimeout(t *testing.T) {
    ts := httptest.NewTLSServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        time.Sleep(15*time.Second)
    }))

    duo := buildAuthApi(ts.URL)

    start := time.Now()
    _, err := duo.Ping()
    duration := time.Since(start)
    if duration.Seconds() > 2 {
        t.Error("Timeout took %d seconds", duration.Seconds())
    }
    if err == nil {
        t.Error("Expected timeout error.")
    }
}

// Test a successful ping request / response.
func TestPing(t *testing.T) {
    ts := httptest.NewTLSServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
      fmt.Fprintln(w, `
        {
          "stat": "OK",
          "response": {
            "time": 1357020061,
            "unexpected_parameter" : "blah"
          }
        }`)
    }))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    result, err := duo.Ping()
    if err != nil {
      t.Error("Unexpected error from Ping call" + err.Error())
    }
    if result.Stat != "OK" {
        t.Error("Expected OK, but got " + result.Stat)
    }
    if result.Response.Time != 1357020061 {
        t.Errorf("Expected 1357020061, but got %d", result.Response.Time)
    }
}

// Test a successful Check request / response.
func TestCheck(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
          fmt.Fprintln(w, `
            {
              "stat": "OK",
              "response": {
                "time": 1357020061
              }
            }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    result, err := duo.Check()
    if err != nil {
      t.Error("Failed TestCheck: " + err.Error())
    }
    if result.Stat != "OK" {
        t.Error("Expected OK, but got " + result.Stat)
    }
    if result.Response.Time != 1357020061 {
        t.Errorf("Expected 1357020061, but got %d", result.Response.Time)
    }
}

// Test a successful logo request / response.
func TestLogo(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "image/png")
            w.Write([]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00" +
                            "\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00" +
                            "\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx" +
                            "\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00" +
                            "\x00\x00\x00IEND\xaeB`\x82"))
        }))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    _, err := duo.Logo()
    if err != nil {
      t.Error("Failed TestCheck: " + err.Error())
    }
}

// Test a failure logo reqeust / response.
func TestLogoError(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
            // Return a 400, as if the logo was not found.
            w.WriteHeader(400)
            fmt.Fprintln(w, `
                {
                    "stat": "FAIL",
                    "code": 40002,
                    "message": "Logo not found",
                    "message_detail": "Why u no have logo?"
                  }`)
        }))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    res, err := duo.Logo()
    if err != nil {
      t.Error("Failed TestCheck: " + err.Error())
    }
    if res.Stat != "FAIL" {
        t.Error("Expected FAIL, but got " + res.Stat)
    }
    if res.Code == nil || *res.Code != 40002 {
        t.Error("Unexpected response code.")
    }
    if res.Message == nil || *res.Message != "Logo not found" {
        t.Error("Unexpected message.")
    }
    if res.Message_Detail == nil || *res.Message_Detail != "Why u no have logo?" {
        t.Error("Unexpected message detail.")
    }
}

// Test a successful enroll request / response.
func TestEnroll(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
          if r.FormValue("username") != "49c6c3097adb386048c84354d82ea63d" {
            t.Error("TestEnroll failed to set 'username' query parameter:" +
                    r.RequestURI)
          }
          if r.FormValue("valid_secs") != "10" {
            t.Error("TestEnroll failed to set 'valid_secs' query parameter: " +
                     r.RequestURI)
          }
          fmt.Fprintln(w, `
            {
              "stat": "OK",
              "response": {
                "activation_barcode": "https://api-eval.duosecurity.com/frame/qr?value=8LIRa5danrICkhHtkLxi-cKLu2DWzDYCmBwBHY2YzW5ZYnYaRxA",
                "activation_code": "duo://8LIRa5danrICkhHtkLxi-cKLu2DWzDYCmBwBHY2YzW5ZYnYaRxA",
                "expiration": 1357020061,
                "user_id": "DU94SWSN4ADHHJHF2HXT",
                "username": "49c6c3097adb386048c84354d82ea63d"
              }
            }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    result, err := duo.Enroll(EnrollUsername("49c6c3097adb386048c84354d82ea63d"), EnrollValidSeconds(10))
    if err != nil {
      t.Error("Failed TestEnroll: " + err.Error())
    }
    if result.Stat != "OK" {
        t.Error("Expected OK, but got " + result.Stat)
    }
    if result.Response.Activation_Barcode != "https://api-eval.duosecurity.com/frame/qr?value=8LIRa5danrICkhHtkLxi-cKLu2DWzDYCmBwBHY2YzW5ZYnYaRxA" {
        t.Error("Unexpected activation_barcode: " + result.Response.Activation_Barcode)
    }
    if result.Response.Activation_Code != "duo://8LIRa5danrICkhHtkLxi-cKLu2DWzDYCmBwBHY2YzW5ZYnYaRxA" {
        t.Error("Unexpected activation code: " + result.Response.Activation_Code)
    }
    if result.Response.Expiration != 1357020061 {
        t.Errorf("Unexpected expiration time: %d", result.Response.Expiration)
    }
    if result.Response.User_Id != "DU94SWSN4ADHHJHF2HXT" {
        t.Error("Unexpected user id: " + result.Response.User_Id)
    }
    if result.Response.Username != "49c6c3097adb386048c84354d82ea63d" {
        t.Error("Unexpected username: " + result.Response.Username)
    }
}

// Test a succesful enroll status request / response.
func TestEnrollStatus(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
          if r.FormValue("user_id") != "49c6c3097adb386048c84354d82ea63d" {
            t.Error("TestEnrollStatus failed to set 'user_id' query parameter:" +
                    r.RequestURI)
          }
          if r.FormValue("activation_code") != "10" {
            t.Error("TestEnrollStatus failed to set 'activation_code' query parameter: " +
                     r.RequestURI)
          }
          fmt.Fprintln(w, `
            {
            "stat": "OK",
            "response": "success"
            }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    result, err := duo.EnrollStatus("49c6c3097adb386048c84354d82ea63d", "10")
    if err != nil {
      t.Error("Failed TestEnrollStatus: " + err.Error())
    }
    if result.Stat != "OK" {
        t.Error("Expected OK, but got " + result.Stat)
    }
    if result.Response != "success" {
        t.Error("Unexpected response: " + result.Response)
    }
}

// Test a successful preauth with user id.  The client doesn't enforce api requirements,
// such as requiring only one of user id or username, but we'll cover the username
// in another test anyway.
func TestPreauthUserId(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
          if r.FormValue("ipaddr") != "127.0.0.1" {
            t.Error("TestPreauth failed to set 'ipaddr' query parameter:" +
                    r.RequestURI)
          }
          if r.FormValue("user_id") != "10" {
            t.Error("TestEnrollStatus failed to set 'user_id' query parameter: " +
                     r.RequestURI)
          }
          if r.FormValue("trusted_device_token") != "l33t" {
            t.Error("TestEnrollStatus failed to set 'trusted_device_token' query parameter: " +
                     r.RequestURI)
          }
          fmt.Fprintln(w, `
            {
              "stat": "OK",
              "response": {
                "result": "auth",
                "status_msg": "Account is active",
                "devices": [
                  {
                    "device": "DPFZRS9FB0D46QFTM891",
                    "type": "phone",
                    "number": "XXX-XXX-0100",
                    "name": "",
                    "capabilities": [
                        "push",
                        "sms",
                        "phone"
                    ]
                  },
                  {
                    "device": "DHEKH0JJIYC1LX3AZWO4",
                    "type": "token",
                    "name": "0"
                  }
                ]
              }
            }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    res, err := duo.Preauth(PreauthUserId("10"), PreauthIpAddr("127.0.0.1"), PreauthTrustedToken("l33t"))
    if err != nil {
      t.Error("Failed TestPreauthUserId: " + err.Error())
    }
    if res.Stat != "OK" {
        t.Error("Unexpected stat: " + res.Stat)
    }
    if res.Response.Result != "auth" {
        t.Error("Unexpected response result: " + res.Response.Result)
    }
    if res.Response.Status_Msg != "Account is active" {
        t.Error("Unexpected status message: " + res.Response.Status_Msg)
    }
    if len(res.Response.Devices) != 2 {
        t.Errorf("Unexpected devices length: %d", len(res.Response.Devices))
    }
    if res.Response.Devices[0].Device != "DPFZRS9FB0D46QFTM891" {
        t.Error("Unexpected [0] device name: " + res.Response.Devices[0].Device)
    }
    if res.Response.Devices[0].Type != "phone" {
        t.Error("Unexpected [0] device type: " + res.Response.Devices[0].Type)
    }
    if res.Response.Devices[0].Number != "XXX-XXX-0100" {
        t.Error("Unexpected [0] device number: " + res.Response.Devices[0].Number)
    }
    if res.Response.Devices[0].Name != "" {
        t.Error("Unexpected [0] devices name :" + res.Response.Devices[0].Name)
    }
    if len(res.Response.Devices[0].Capabilities) != 3 {
        t.Errorf("Unexpected [0] device capabilities length: %d", len(res.Response.Devices[0].Capabilities))
    }
    if res.Response.Devices[0].Capabilities[0] != "push" {
        t.Error("Unexpected [0] device capability: " + res.Response.Devices[0].Capabilities[0])
    }
    if res.Response.Devices[0].Capabilities[1] != "sms" {
        t.Error("Unexpected [0] device capability: " + res.Response.Devices[0].Capabilities[1])
    }
    if res.Response.Devices[0].Capabilities[2] != "phone" {
        t.Error("Unexpected [0] device capability: " + res.Response.Devices[0].Capabilities[2])
    }
    if res.Response.Devices[1].Device != "DHEKH0JJIYC1LX3AZWO4" {
        t.Error("Unexpected [1] device name: " + res.Response.Devices[1].Device)
    }
    if res.Response.Devices[1].Type != "token" {
        t.Error("Unexpected [1] device type: " + res.Response.Devices[1].Type)
    }
    if res.Response.Devices[1].Name != "0" {
        t.Error("Unexpected [1] devices name :" + res.Response.Devices[1].Name)
    }
    if len(res.Response.Devices[1].Capabilities) != 0 {
        t.Errorf("Unexpected [1] device capabilities length: %d", len(res.Response.Devices[1].Capabilities))
    }
}

// Test preauth enroll with username, and an enroll response.
func TestPreauthEnroll(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
          if r.FormValue("username") != "10" {
            t.Error("TestEnrollStatus failed to set 'username' query parameter: " +
                     r.RequestURI)
          }
          fmt.Fprintln(w, `
              {
                "stat": "OK",
                "response": {
                  "enroll_portal_url": "https://api-3945ef22.duosecurity.com/portal?48bac5d9393fb2c2",
                  "result": "enroll",
                  "status_msg": "Enroll an authentication device to proceed"
                }
              }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    res, err := duo.Preauth(PreauthUsername("10"))
    if err != nil {
      t.Error("Failed TestPreauthEnroll: " + err.Error())
    }
    if res.Stat != "OK" {
        t.Error("Unexpected stat: " + res.Stat)
    }
    if res.Response.Enroll_Portal_Url != "https://api-3945ef22.duosecurity.com/portal?48bac5d9393fb2c2" {
        t.Error("Unexpected enroll portal URL: " + res.Response.Enroll_Portal_Url)
    }
    if res.Response.Result != "enroll" {
        t.Error("Unexpected response result: " + res.Response.Result)
    }
    if res.Response.Status_Msg != "Enroll an authentication device to proceed" {
        t.Error("Unexpected status msg: " + res.Response.Status_Msg)
    }
}

// Test an authentication request / response.  This won't work against the Duo
// server, because the request parameters included are illegal.  But we can
// verify that the go code sets the query parameters correctly.
func TestAuth(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
            expected := map[string]string {
                "username" : "username value",
                "user_id" : "user_id value",
                "factor" : "auto",
                "ipaddr" : "40.40.40.10",
                "async" : "1",
                "device" : "primary",
                "type" : "request",
                "display_username" : "display username",

            }
            for key, value := range expected {
                if r.FormValue(key) != value {
                    t.Errorf("TestAuth failed to set '%s' query parameter: " +
                         r.RequestURI, key)
                }
            }
            fmt.Fprintln(w, `
                {
                    "stat": "OK",
                    "response": {
                      "result": "allow",
                      "status": "allow",
                      "status_msg": "Success. Logging you in..."
                    }
                  }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    res, err := duo.Auth("auto",
                       AuthUserId("user_id value"),
                       AuthUsername("username value"),
                       AuthIpAddr("40.40.40.10"),
                       AuthAsync(),
                       AuthDevice("primary"),
                       AuthType("request"),
                       AuthDisplayUsername("display username"),
                       )
    if err != nil {
      t.Error("Failed TestAuth: " + err.Error())
    }
    if res.Stat != "OK" {
        t.Error("Unexpected stat: " + res.Stat)
    }
    if res.Response.Result != "allow" {
        t.Error("Unexpected response result: " + res.Response.Result)
    }
    if res.Response.Status != "allow" {
        t.Error("Unexpected response status: " + res.Response.Status)
    }
    if res.Response.Status_Msg != "Success. Logging you in..." {
        t.Error("Unexpected response status msg: " + res.Response.Status_Msg)
    }
}

// Test AuthStatus request / response.
func TestAuthStatus(t *testing.T) {
    ts := httptest.NewTLSServer(
      http.HandlerFunc(
        func (w http.ResponseWriter, r *http.Request) {
            expected := map[string]string {
                "txid" : "4",
            }
            for key, value := range expected {
                if r.FormValue(key) != value {
                    t.Errorf("TestAuthStatus failed to set '%s' query parameter: " +
                         r.RequestURI, key)
                }
            }
            fmt.Fprintln(w, `
            {
                "stat": "OK",
                "response": {
                  "result": "waiting",
                  "status": "pushed",
                  "status_msg": "Pushed a login request to your phone..."
                }
            }`)}))
    defer ts.Close()

    duo := buildAuthApi(ts.URL)

    res, err := duo.AuthStatus("4")
    if err != nil {
      t.Error("Failed TestAuthStatus: " + err.Error())
    }

    if res.Stat != "OK" {
        t.Error("Unexpected stat: " + res.Stat)
    }
    if res.Response.Result != "waiting" {
        t.Error("Unexpected response result: " + res.Response.Result)
    }
    if res.Response.Status != "pushed" {
        t.Error("Unexpected response status: " + res.Response.Status)
    }
    if res.Response.Status_Msg != "Pushed a login request to your phone..." {
        t.Error("Unexpected response status msg: " + res.Response.Status_Msg)
    }
}
