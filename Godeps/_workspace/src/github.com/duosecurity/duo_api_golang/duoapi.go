package duoapi

import (
    "strings"
    "net/url"
    "sort"
    "crypto/hmac"
    "crypto/sha1"
    "crypto/tls"
    "crypto/x509"
    "encoding/hex"
    "encoding/base64"
    "net/http"
    "time"
    "io/ioutil"
)
var spaceReplacer *strings.Replacer = strings.NewReplacer("+", "%20")

func canonParams(params url.Values) string {
    // Values must be in sorted order
    for key, val := range params {
        sort.Strings(val)
        params[key] = val
    }
    // Encode will place Keys in sorted order
    ordered_params := params.Encode()
    // Encoder turns spaces into +, but we need %XX escaping
    return spaceReplacer.Replace(ordered_params)
}

func canonicalize(method string,
                  host string,
                  uri string,
                  params url.Values,
                  date string) string {
    var canon [5]string
    canon[0] = date
    canon[1] = strings.ToUpper(method)
    canon[2] = strings.ToLower(host)
    canon[3] = uri
    canon[4] = canonParams(params)
    return strings.Join(canon[:], "\n")
}

func sign(ikey string,
          skey string,
          method string,
          host string,
          uri string,
          date string,
          params url.Values) string {
    canon := canonicalize(method, host, uri, params, date)
    mac := hmac.New(sha1.New, []byte(skey))
    mac.Write([]byte(canon))
    sig := hex.EncodeToString(mac.Sum(nil))
    auth := ikey + ":" + sig
    return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

type DuoApi struct {
    ikey string
    skey string
    host string
    userAgent string
    apiClient *http.Client
    authClient *http.Client
}

type apiOptions struct {
    timeout time.Duration
    insecure bool
}

// Optional parameter for NewDuoApi, used to configure timeouts on API calls.
func SetTimeout(timeout time.Duration) func(*apiOptions) {
    return func(opts *apiOptions) {
        opts.timeout = timeout
        return
    }
}

// Optional parameter for testing only.  Bypasses all TLS certificate validation.
func SetInsecure() func(*apiOptions) {
    return func(opts *apiOptions) {
        opts.insecure = true
    }
}

// Build an return a DuoApi struct.
// ikey is your Duo integration key
// skey is your Duo integration secret key
// host is your Duo host
// userAgent allows you to specify the user agent string used when making
//           the web request to Duo.
// options are optional parameters.  Use SetTimeout() to specify a timeout value
//         for Rest API calls.
//
// Example: duoapi.NewDuoApi(ikey,skey,host,userAgent,duoapi.SetTimeout(10*time.Second))
func NewDuoApi(ikey string,
        skey string,
        host string,
        userAgent string,
        options ...func(*apiOptions)) (*DuoApi) {
    opts := apiOptions{}
    for _, o := range options {
        o(&opts)
    }

    // Certificate pinning
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM([]byte(duoPinnedCert))

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
          RootCAs: certPool,
          InsecureSkipVerify: opts.insecure,
      },
    }
    return &DuoApi{
        ikey: ikey,
        skey: skey,
        host: host,
        userAgent: userAgent,
        apiClient: &http.Client{
            Timeout: opts.timeout,
            Transport: tr,
        },
        authClient: &http.Client{
            Transport: tr,
        },
    }
}

type requestOptions struct {
    timeout bool
}

type DuoApiOption func(*requestOptions)

// Pass to Request or SignedRequest to configure a timeout on the request
func UseTimeout(opts *requestOptions) {
    opts.timeout = true
}

func (duoapi *DuoApi) buildOptions(options ...DuoApiOption) (*requestOptions) {
    opts := &requestOptions{}
    for _, o := range options {
        o(opts)
    }
    return opts
}

// Make an unsigned Duo Rest API call.  See Duo's online documentation
// for the available REST API's.
// method is POST or GET
// uri is the URI of the Duo Rest call
// params HTTP query parameters to include in the call.
// options Optional parameters.  Use UseTimeout to toggle whether the
//         Duo Rest API call should timeout or not.
//
// Example: duo.Call("GET", "/auth/v2/ping", nil, duoapi.UseTimeout)
func (duoapi *DuoApi) Call(method string,
        uri string,
        params url.Values,
        options ...DuoApiOption) (*http.Response, []byte, error) {
    opts := duoapi.buildOptions(options...)

    client := duoapi.authClient
    if opts.timeout {
        client = duoapi.apiClient
    }

    url := url.URL{
        Scheme: "https",
        Host: duoapi.host,
        Path: uri,
        RawQuery: params.Encode(),
    }
    request, err := http.NewRequest(method, url.String(), nil)
    if err != nil {
        return nil, nil, err
    }
    resp, err := client.Do(request)
    var body []byte
    if err == nil {
        body, err = ioutil.ReadAll(resp.Body)
        resp.Body.Close()
    }
    return resp, body, err
}

// Make a signed Duo Rest API call.  See Duo's online documentation
// for the available REST API's.
// method is POST or GET
// uri is the URI of the Duo Rest call
// params HTTP query parameters to include in the call.
// options Optional parameters.  Use UseTimeout to toggle whether the
//         Duo Rest API call should timeout or not.
//
// Example: duo.SignedCall("GET", "/auth/v2/check", nil, duoapi.UseTimeout)
func (duoapi *DuoApi) SignedCall(method string,
        uri string,
        params url.Values,
        options ...DuoApiOption) (*http.Response, []byte, error) {
    opts := duoapi.buildOptions(options...)

    now := time.Now().UTC().Format(time.RFC1123Z)
    auth_sig := sign(duoapi.ikey, duoapi.skey, method, duoapi.host, uri, now, params)

    url := url.URL{
        Scheme: "https",
        Host: duoapi.host,
        Path: uri,
        RawQuery: params.Encode(),
    }
    request, err := http.NewRequest(method, url.String(), nil)
    if err != nil {
        return nil, nil, err
    }
    request.Header.Set("Authorization", auth_sig)
    request.Header.Set("Date", now)

    client := duoapi.authClient
    if opts.timeout {
        client = duoapi.apiClient
    }
    resp, err := client.Do(request)
    var body []byte
    if err == nil {
        body, err = ioutil.ReadAll(resp.Body)
        resp.Body.Close()
    }
    return resp, body, err
}

const duoPinnedCert string = "subject= /C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert Assured ID Root CA\n" +
"-----BEGIN CERTIFICATE-----\n" +
"MIIDtzCCAp+gAwIBAgIQDOfg5RfYRv6P5WD8G/AwOTANBgkqhkiG9w0BAQUFADBl\n" +
"MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\n" +
"d3cuZGlnaWNlcnQuY29tMSQwIgYDVQQDExtEaWdpQ2VydCBBc3N1cmVkIElEIFJv\n" +
"b3QgQ0EwHhcNMDYxMTEwMDAwMDAwWhcNMzExMTEwMDAwMDAwWjBlMQswCQYDVQQG\n" +
"EwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNl\n" +
"cnQuY29tMSQwIgYDVQQDExtEaWdpQ2VydCBBc3N1cmVkIElEIFJvb3QgQ0EwggEi\n" +
"MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtDhXO5EOAXLGH87dg+XESpa7c\n" +
"JpSIqvTO9SA5KFhgDPiA2qkVlTJhPLWxKISKityfCgyDF3qPkKyK53lTXDGEKvYP\n" +
"mDI2dsze3Tyoou9q+yHyUmHfnyDXH+Kx2f4YZNISW1/5WBg1vEfNoTb5a3/UsDg+\n" +
"wRvDjDPZ2C8Y/igPs6eD1sNuRMBhNZYW/lmci3Zt1/GiSw0r/wty2p5g0I6QNcZ4\n" +
"VYcgoc/lbQrISXwxmDNsIumH0DJaoroTghHtORedmTpyoeb6pNnVFzF1roV9Iq4/\n" +
"AUaG9ih5yLHa5FcXxH4cDrC0kqZWs72yl+2qp/C3xag/lRbQ/6GW6whfGHdPAgMB\n" +
"AAGjYzBhMA4GA1UdDwEB/wQEAwIBhjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQW\n" +
"BBRF66Kv9JLLgjEtUYunpyGd823IDzAfBgNVHSMEGDAWgBRF66Kv9JLLgjEtUYun\n" +
"pyGd823IDzANBgkqhkiG9w0BAQUFAAOCAQEAog683+Lt8ONyc3pklL/3cmbYMuRC\n" +
"dWKuh+vy1dneVrOfzM4UKLkNl2BcEkxY5NM9g0lFWJc1aRqoR+pWxnmrEthngYTf\n" +
"fwk8lOa4JiwgvT2zKIn3X/8i4peEH+ll74fg38FnSbNd67IJKusm7Xi+fT8r87cm\n" +
"NW1fiQG2SVufAQWbqz0lwcy2f8Lxb4bG+mRo64EtlOtCt/qMHt1i8b5QZ7dsvfPx\n" +
"H2sMNgcWfzd8qVttevESRmCD1ycEvkvOl77DZypoEd+A5wwzZr8TDRRu838fYxAe\n" +
"+o0bJW1sj6W3YQGx0qMmoRBxna3iw/nDmVG3KwcIzi7mULKn+gpFL6Lw8g==\n" +
"-----END CERTIFICATE-----\n" +
"\n" +
"subject= /C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert Global Root CA\n" +
"-----BEGIN CERTIFICATE-----\n" +
"MIIDrzCCApegAwIBAgIQCDvgVpBCRrGhdWrJWZHHSjANBgkqhkiG9w0BAQUFADBh\n" +
"MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\n" +
"d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD\n" +
"QTAeFw0wNjExMTAwMDAwMDBaFw0zMTExMTAwMDAwMDBaMGExCzAJBgNVBAYTAlVT\n" +
"MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j\n" +
"b20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IENBMIIBIjANBgkqhkiG\n" +
"9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4jvhEXLeqKTTo1eqUKKPC3eQyaKl7hLOllsB\n" +
"CSDMAZOnTjC3U/dDxGkAV53ijSLdhwZAAIEJzs4bg7/fzTtxRuLWZscFs3YnFo97\n" +
"nh6Vfe63SKMI2tavegw5BmV/Sl0fvBf4q77uKNd0f3p4mVmFaG5cIzJLv07A6Fpt\n" +
"43C/dxC//AH2hdmoRBBYMql1GNXRor5H4idq9Joz+EkIYIvUX7Q6hL+hqkpMfT7P\n" +
"T19sdl6gSzeRntwi5m3OFBqOasv+zbMUZBfHWymeMr/y7vrTC0LUq7dBMtoM1O/4\n" +
"gdW7jVg/tRvoSSiicNoxBN33shbyTApOB6jtSj1etX+jkMOvJwIDAQABo2MwYTAO\n" +
"BgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUA95QNVbR\n" +
"TLtm8KPiGxvDl7I90VUwHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUw\n" +
"DQYJKoZIhvcNAQEFBQADggEBAMucN6pIExIK+t1EnE9SsPTfrgT1eXkIoyQY/Esr\n" +
"hMAtudXH/vTBH1jLuG2cenTnmCmrEbXjcKChzUyImZOMkXDiqw8cvpOp/2PV5Adg\n" +
"06O/nVsJ8dWO41P0jmP6P6fbtGbfYmbW0W5BjfIttep3Sp+dWOIrWcBAI+0tKIJF\n" +
"PnlUkiaY4IBIqDfv8NZ5YBberOgOzW6sRBc4L0na4UU+Krk2U886UAb3LujEV0ls\n" +
"YSEY1QSteDwsOoBrp+uvFRTp2InBuThs4pFsiv9kuXclVzDAGySj4dzp30d8tbQk\n" +
"CAUw7C29C79Fv1C5qfPrmAESrciIxpg0X40KPMbp1ZWVbd4=\n" +
"-----END CERTIFICATE-----\n" +
"\n" +
"subject= /C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert High Assurance EV Root CA\n" +
"-----BEGIN CERTIFICATE-----\n" +
"MIIDxTCCAq2gAwIBAgIQAqxcJmoLQJuPC3nyrkYldzANBgkqhkiG9w0BAQUFADBs\n" +
"MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\n" +
"d3cuZGlnaWNlcnQuY29tMSswKQYDVQQDEyJEaWdpQ2VydCBIaWdoIEFzc3VyYW5j\n" +
"ZSBFViBSb290IENBMB4XDTA2MTExMDAwMDAwMFoXDTMxMTExMDAwMDAwMFowbDEL\n" +
"MAkGA1UEBhMCVVMxFTATBgNVBAoTDERpZ2lDZXJ0IEluYzEZMBcGA1UECxMQd3d3\n" +
"LmRpZ2ljZXJ0LmNvbTErMCkGA1UEAxMiRGlnaUNlcnQgSGlnaCBBc3N1cmFuY2Ug\n" +
"RVYgUm9vdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMbM5XPm\n" +
"+9S75S0tMqbf5YE/yc0lSbZxKsPVlDRnogocsF9ppkCxxLeyj9CYpKlBWTrT3JTW\n" +
"PNt0OKRKzE0lgvdKpVMSOO7zSW1xkX5jtqumX8OkhPhPYlG++MXs2ziS4wblCJEM\n" +
"xChBVfvLWokVfnHoNb9Ncgk9vjo4UFt3MRuNs8ckRZqnrG0AFFoEt7oT61EKmEFB\n" +
"Ik5lYYeBQVCmeVyJ3hlKV9Uu5l0cUyx+mM0aBhakaHPQNAQTXKFx01p8VdteZOE3\n" +
"hzBWBOURtCmAEvF5OYiiAhF8J2a3iLd48soKqDirCmTCv2ZdlYTBoSUeh10aUAsg\n" +
"EsxBu24LUTi4S8sCAwEAAaNjMGEwDgYDVR0PAQH/BAQDAgGGMA8GA1UdEwEB/wQF\n" +
"MAMBAf8wHQYDVR0OBBYEFLE+w2kD+L9HAdSYJhoIAu9jZCvDMB8GA1UdIwQYMBaA\n" +
"FLE+w2kD+L9HAdSYJhoIAu9jZCvDMA0GCSqGSIb3DQEBBQUAA4IBAQAcGgaX3Nec\n" +
"nzyIZgYIVyHbIUf4KmeqvxgydkAQV8GK83rZEWWONfqe/EW1ntlMMUu4kehDLI6z\n" +
"eM7b41N5cdblIZQB2lWHmiRk9opmzN6cN82oNLFpmyPInngiK3BD41VHMWEZ71jF\n" +
"hS9OMPagMRYjyOfiZRYzy78aG6A9+MpeizGLYAiJLQwGXFK3xPkKmNEVX58Svnw2\n" +
"Yzi9RKR/5CYrCsSXaQ3pjOLAEFe4yHYSkVXySGnYvCoCWw9E1CAx2/S6cCZdkGCe\n" +
"vEsXCS+0yx5DaMkHJ8HSXPfqIbloEpw8nL+e/IBcm2PN7EeqJSdnoDfzAIJ9VNep\n" +
"+OkuE6N36B9K\n" +
"-----END CERTIFICATE-----\n" +
"\n" +
"subject= /C=US/O=SecureTrust Corporation/CN=SecureTrust CA\n" +
"-----BEGIN CERTIFICATE-----\n" +
"MIIDuDCCAqCgAwIBAgIQDPCOXAgWpa1Cf/DrJxhZ0DANBgkqhkiG9w0BAQUFADBI\n" +
"MQswCQYDVQQGEwJVUzEgMB4GA1UEChMXU2VjdXJlVHJ1c3QgQ29ycG9yYXRpb24x\n" +
"FzAVBgNVBAMTDlNlY3VyZVRydXN0IENBMB4XDTA2MTEwNzE5MzExOFoXDTI5MTIz\n" +
"MTE5NDA1NVowSDELMAkGA1UEBhMCVVMxIDAeBgNVBAoTF1NlY3VyZVRydXN0IENv\n" +
"cnBvcmF0aW9uMRcwFQYDVQQDEw5TZWN1cmVUcnVzdCBDQTCCASIwDQYJKoZIhvcN\n" +
"AQEBBQADggEPADCCAQoCggEBAKukgeWVzfX2FI7CT8rU4niVWJxB4Q2ZQCQXOZEz\n" +
"Zum+4YOvYlyJ0fwkW2Gz4BERQRwdbvC4u/jep4G6pkjGnx29vo6pQT64lO0pGtSO\n" +
"0gMdA+9tDWccV9cGrcrI9f4Or2YlSASWC12juhbDCE/RRvgUXPLIXgGZbf2IzIao\n" +
"wW8xQmxSPmjL8xk037uHGFaAJsTQ3MBv396gwpEWoGQRS0S8Hvbn+mPeZqx2pHGj\n" +
"7DaUaHp3pLHnDi+BeuK1cobvomuL8A/b01k/unK8RCSc43Oz969XL0Imnal0ugBS\n" +
"8kvNU3xHCzaFDmapCJcWNFfBZveA4+1wVMeT4C4oFVmHursCAwEAAaOBnTCBmjAT\n" +
"BgkrBgEEAYI3FAIEBh4EAEMAQTALBgNVHQ8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB\n" +
"/zAdBgNVHQ4EFgQUQjK2FvoE/f5dS3rD/fdMQB1aQ68wNAYDVR0fBC0wKzApoCeg\n" +
"JYYjaHR0cDovL2NybC5zZWN1cmV0cnVzdC5jb20vU1RDQS5jcmwwEAYJKwYBBAGC\n" +
"NxUBBAMCAQAwDQYJKoZIhvcNAQEFBQADggEBADDtT0rhWDpSclu1pqNlGKa7UTt3\n" +
"6Z3q059c4EVlew3KW+JwULKUBRSuSceNQQcSc5R+DCMh/bwQf2AQWnL1mA6s7Ll/\n" +
"3XpvXdMc9P+IBWlCqQVxyLesJugutIxq/3HcuLHfmbx8IVQr5Fiiu1cprp6poxkm\n" +
"D5kuCLDv/WnPmRoJjeOnnyvJNjR7JLN4TJUXpAYmHrZkUjZfYGfZnMUFdAvnZyPS\n" +
"CPyI6a6Lf+Ew9Dd+/cYy2i2eRDAwbO4H3tI0/NL/QPZL9GZGBlSm8jIKYyYwa5vR\n" +
"3ItHuuG51WLQoqD0ZwV4KWMabwTW+MZMo5qxN7SN5ShLHZ4swrhovO0C7jE=\n" +
"-----END CERTIFICATE-----\n" +
"\n" +
"subject= /C=US/O=SecureTrust Corporation/CN=Secure Global CA\n" +
"-----BEGIN CERTIFICATE-----\n" +
"MIIDvDCCAqSgAwIBAgIQB1YipOjUiolN9BPI8PjqpTANBgkqhkiG9w0BAQUFADBK\n" +
"MQswCQYDVQQGEwJVUzEgMB4GA1UEChMXU2VjdXJlVHJ1c3QgQ29ycG9yYXRpb24x\n" +
"GTAXBgNVBAMTEFNlY3VyZSBHbG9iYWwgQ0EwHhcNMDYxMTA3MTk0MjI4WhcNMjkx\n" +
"MjMxMTk1MjA2WjBKMQswCQYDVQQGEwJVUzEgMB4GA1UEChMXU2VjdXJlVHJ1c3Qg\n" +
"Q29ycG9yYXRpb24xGTAXBgNVBAMTEFNlY3VyZSBHbG9iYWwgQ0EwggEiMA0GCSqG\n" +
"SIb3DQEBAQUAA4IBDwAwggEKAoIBAQCvNS7YrGxVaQZx5RNoJLNP2MwhR/jxYDiJ\n" +
"iQPpvepeRlMJ3Fz1Wuj3RSoC6zFh1ykzTM7HfAo3fg+6MpjhHZevj8fcyTiW89sa\n" +
"/FHtaMbQbqR8JNGuQsiWUGMu4P51/pinX0kuleM5M2SOHqRfkNJnPLLZ/kG5VacJ\n" +
"jnIFHovdRIWCQtBJwB1g8NEXLJXr9qXBkqPFwqcIYA1gBBCWeZ4WNOaptvolRTnI\n" +
"HmX5k/Wq8VLcmZg9pYYaDDUz+kulBAYVHDGA76oYa8J719rO+TMg1fW9ajMtgQT7\n" +
"sFzUnKPiXB3jqUJ1XnvUd+85VLrJChgbEplJL4hL/VBi0XPnj3pDAgMBAAGjgZ0w\n" +
"gZowEwYJKwYBBAGCNxQCBAYeBABDAEEwCwYDVR0PBAQDAgGGMA8GA1UdEwEB/wQF\n" +
"MAMBAf8wHQYDVR0OBBYEFK9EBMJBfkiD2045AuzshHrmzsmkMDQGA1UdHwQtMCsw\n" +
"KaAnoCWGI2h0dHA6Ly9jcmwuc2VjdXJldHJ1c3QuY29tL1NHQ0EuY3JsMBAGCSsG\n" +
"AQQBgjcVAQQDAgEAMA0GCSqGSIb3DQEBBQUAA4IBAQBjGghAfaReUw132HquHw0L\n" +
"URYD7xh8yOOvaliTFGCRsoTciE6+OYo68+aCiV0BN7OrJKQVDpI1WkpEXk5X+nXO\n" +
"H0jOZvQ8QCaSmGwb7iRGDBezUqXbpZGRzzfTb+cnCDpOGR86p1hcF895P4vkp9Mm\n" +
"I50mD1hp/Ed+stCNi5O/KU9DaXR2Z0vPB4zmAve14bRDtUstFJ/53CYNv6ZHdAbY\n" +
"iNE6KTCEztI5gGIbqMdXSbxqVVFnFUq+NQfk1XWYN3kwFNspnWzFacxHVaIw98xc\n" +
"f8LDmBxrThaA63p4ZUWiABqvDA1VZDRIuJK58bRQKfJPIx/abKwfROHdI3hRW8cW\n" +
"-----END CERTIFICATE-----\n"
