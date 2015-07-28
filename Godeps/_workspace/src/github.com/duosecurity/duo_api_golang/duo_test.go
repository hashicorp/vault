package duoapi

import (
    "testing"
    "net/url"
    "strings"
)

func TestCanonicalize(t *testing.T) {
    values := url.Values{}
    values.Set("username", "H ell?o")
    values.Set("password", "H-._~i")
    values.Add("password", "A(!'*)")
    params_str := canonicalize("post",
                           "API-XXX.duosecurity.COM",
                           "/auth/v2/ping",
                           values,
                           "5")
    params := strings.Split(params_str, "\n")
    if len(params) != 5 {
        t.Error("Expected 5 parameters, but got " + string(len(params)))
    }
    if params[1] != string("POST") {
        t.Error("Expected POST, but got " + params[1])
    }
    if params[2] != string("api-xxx.duosecurity.com") {
        t.Error("Expected api-xxx.duosecurity.com, but got " + params[2])
    }
    if params[3] != string("/auth/v2/ping") {
        t.Error("Expected /auth/v2/ping, but got " + params[3])
    }
    if params[4] != string("password=A%28%21%27%2A%29&password=H-._~i&username=H%20ell%3Fo") {
        t.Error("Expected sorted escaped params, but got " + params[4])
    }
}

func encodeAndValidate(t *testing.T, input url.Values, output string) {
  values := url.Values{}
  for key, val := range input {
      values.Set(key, val[0])
  }
  params_str := canonicalize("post",
      "API-XXX.duosecurity.com",
      "/auth/v2/ping",
      values,
      "5")
  params := strings.Split(params_str, "\n")
  if params[4] != output {
      t.Error("Mismatch\n" + output + "\n" + params[4])
  }

}

func TestSimple(t *testing.T) {
    values := url.Values{}
    values.Set("realname", "First Last")
    values.Set("username", "root")

    encodeAndValidate(t, values, "realname=First%20Last&username=root")
}

func TestZero(t *testing.T) {
    values := url.Values{}
    encodeAndValidate(t, values, "")
}

func TestOne(t *testing.T) {
    values := url.Values{}
    values.Set("realname", "First Last")
    encodeAndValidate(t, values, "realname=First%20Last")
}

func TestPrintableAsciiCharaceters(t *testing.T) {
    values := url.Values{}
    values.Set("digits", "0123456789")
    values.Set("letters", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    values.Set("punctuation", "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
    values.Set("whitespace", "\t\n\x0b\x0c\r ")
    encodeAndValidate(t, values, "digits=0123456789&letters=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&punctuation=%21%22%23%24%25%26%27%28%29%2A%2B%2C-.%2F%3A%3B%3C%3D%3E%3F%40%5B%5C%5D%5E_%60%7B%7C%7D~&whitespace=%09%0A%0B%0C%0D%20")
}

func TestSortOrderWithCommonPrefix(t *testing.T) {
    values := url.Values{}
    values.Set("foo", "1")
    values.Set("foo_bar", "2")
    encodeAndValidate(t, values, "foo=1&foo_bar=2")
}

func TestUnicodeFuzzValues(t *testing.T) {
    values := url.Values{}
    values.Set("bar", "⠕ꪣ㟏䮷㛩찅暎腢슽ꇱ")
    values.Set("baz", "ෳ蒽噩馅뢤갺篧潩鍊뤜")
    values.Set("foo", "퓎훖礸僀訠輕ﴋ耤岳왕")
    values.Set("qux", "讗졆-芎茚쳊ꋔ谾뢲馾")
    encodeAndValidate(t, values, "bar=%E2%A0%95%EA%AA%A3%E3%9F%8F%E4%AE%B7%E3%9B%A9%EC%B0%85%E6%9A%8E%E8%85%A2%EC%8A%BD%EA%87%B1&baz=%E0%B7%B3%E8%92%BD%E5%99%A9%E9%A6%85%EB%A2%A4%EA%B0%BA%E7%AF%A7%E6%BD%A9%E9%8D%8A%EB%A4%9C&foo=%ED%93%8E%ED%9B%96%E7%A4%B8%E5%83%80%E8%A8%A0%E8%BC%95%EF%B4%8B%E8%80%A4%E5%B2%B3%EC%99%95&qux=%E8%AE%97%EC%A1%86-%E8%8A%8E%E8%8C%9A%EC%B3%8A%EA%8B%94%E8%B0%BE%EB%A2%B2%E9%A6%BE")
}

func TestUnicodeFuzzKeysAndValues(t *testing.T) {
    values := url.Values{}
    values.Set("䚚⡻㗐軳朧倪ࠐ킑È셰",
        "ཅ᩶㐚敌숿鬉ꯢ荃ᬧ惐")
    values.Set("瑉繋쳻姿﹟获귌逌쿑砓",
        "趷倢鋓䋯⁽蜰곾嘗ॆ丰")
    values.Set("瑰錔逜麮䃘䈁苘豰ᴱꁂ",
        "៙ந鍘꫟ꐪ䢾ﮖ濩럿㋳")
    values.Set("싅Ⱍ☠㘗隳F蘅⃨갡头",
        "ﮩ䆪붃萋☕㹮攭ꢵ핫U")
    encodeAndValidate(t, values, "%E4%9A%9A%E2%A1%BB%E3%97%90%E8%BB%B3%E6%9C%A7%E5%80%AA%E0%A0%90%ED%82%91%C3%88%EC%85%B0=%E0%BD%85%E1%A9%B6%E3%90%9A%E6%95%8C%EC%88%BF%E9%AC%89%EA%AF%A2%E8%8D%83%E1%AC%A7%E6%83%90&%E7%91%89%E7%B9%8B%EC%B3%BB%E5%A7%BF%EF%B9%9F%E8%8E%B7%EA%B7%8C%E9%80%8C%EC%BF%91%E7%A0%93=%E8%B6%B7%E5%80%A2%E9%8B%93%E4%8B%AF%E2%81%BD%E8%9C%B0%EA%B3%BE%E5%98%97%E0%A5%86%E4%B8%B0&%E7%91%B0%E9%8C%94%E9%80%9C%E9%BA%AE%E4%83%98%E4%88%81%E8%8B%98%E8%B1%B0%E1%B4%B1%EA%81%82=%E1%9F%99%E0%AE%A8%E9%8D%98%EA%AB%9F%EA%90%AA%E4%A2%BE%EF%AE%96%E6%BF%A9%EB%9F%BF%E3%8B%B3&%EC%8B%85%E2%B0%9D%E2%98%A0%E3%98%97%E9%9A%B3F%E8%98%85%E2%83%A8%EA%B0%A1%E5%A4%B4=%EF%AE%A9%E4%86%AA%EB%B6%83%E8%90%8B%E2%98%95%E3%B9%AE%E6%94%AD%EA%A2%B5%ED%95%ABU")
}

func TestSign(t *testing.T) {
    values := url.Values{}
    values.Set("realname", "First Last")
    values.Set("username", "root")
    res := sign("DIWJ8X6AEYOR5OMC6TQ1",
                "Zh5eGmUq9zpfQnyUIu5OL9iWoMMv5ZNmk3zLJ4Ep",
                "POST",
                "api-XXXXXXXX.duosecurity.com",
                "/accounts/v1/account/list",
                "Tue, 21 Aug 2012 17:29:18 -0000",
                values)
    if res != "Basic RElXSjhYNkFFWU9SNU9NQzZUUTE6MmQ5N2Q2MTY2MzE5Nzgx" +
              "YjVhM2EwN2FmMzlkMzY2ZjQ5MTIzNGVkYw==" {
        t.Error("Signature did not produce output documented at " +
                "https://www.duosecurity.com/docs/authapi :(")
    }
}

func TestV2Canonicalize(t *testing.T) {
    values := url.Values{}
    values.Set("䚚⡻㗐軳朧倪ࠐ킑È셰",
        "ཅ᩶㐚敌숿鬉ꯢ荃ᬧ惐")
    values.Set("瑉繋쳻姿﹟获귌逌쿑砓",
        "趷倢鋓䋯⁽蜰곾嘗ॆ丰")
    values.Set("瑰錔逜麮䃘䈁苘豰ᴱꁂ",
        "៙ந鍘꫟ꐪ䢾ﮖ濩럿㋳")
    values.Set("싅Ⱍ☠㘗隳F蘅⃨갡头",
        "ﮩ䆪붃萋☕㹮攭ꢵ핫U")
    canon := canonicalize(
        "PoSt",
        "foO.BAr52.cOm",
        "/Foo/BaR2/qux",
        values,
        "Fri, 07 Dec 2012 17:18:00 -0000")
    expected := "Fri, 07 Dec 2012 17:18:00 -0000\nPOST\nfoo.bar52.com\n/Foo/BaR2/qux\n%E4%9A%9A%E2%A1%BB%E3%97%90%E8%BB%B3%E6%9C%A7%E5%80%AA%E0%A0%90%ED%82%91%C3%88%EC%85%B0=%E0%BD%85%E1%A9%B6%E3%90%9A%E6%95%8C%EC%88%BF%E9%AC%89%EA%AF%A2%E8%8D%83%E1%AC%A7%E6%83%90&%E7%91%89%E7%B9%8B%EC%B3%BB%E5%A7%BF%EF%B9%9F%E8%8E%B7%EA%B7%8C%E9%80%8C%EC%BF%91%E7%A0%93=%E8%B6%B7%E5%80%A2%E9%8B%93%E4%8B%AF%E2%81%BD%E8%9C%B0%EA%B3%BE%E5%98%97%E0%A5%86%E4%B8%B0&%E7%91%B0%E9%8C%94%E9%80%9C%E9%BA%AE%E4%83%98%E4%88%81%E8%8B%98%E8%B1%B0%E1%B4%B1%EA%81%82=%E1%9F%99%E0%AE%A8%E9%8D%98%EA%AB%9F%EA%90%AA%E4%A2%BE%EF%AE%96%E6%BF%A9%EB%9F%BF%E3%8B%B3&%EC%8B%85%E2%B0%9D%E2%98%A0%E3%98%97%E9%9A%B3F%E8%98%85%E2%83%A8%EA%B0%A1%E5%A4%B4=%EF%AE%A9%E4%86%AA%EB%B6%83%E8%90%8B%E2%98%95%E3%B9%AE%E6%94%AD%EA%A2%B5%ED%95%ABU"
    if canon != expected {
        t.Error("Mismatch!\n" + expected + "\n" + canon)
    }
}

func TestNewDuo(t *testing.T) {
    duo := NewDuoApi("ABC", "123", "api-XXXXXXX.duosecurity.com", "go-client")
    if duo == nil {
        t.Fatal("Failed to create a new Duo Api")
    }
}
