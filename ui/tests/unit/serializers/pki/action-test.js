import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Serializer | pki/action', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const store = this.owner.lookup('service:store');
    const serializer = store.serializerFor('pki/action');

    assert.ok(serializer);
  });

  module('actionType import', function (hooks) {
    hooks.beforeEach(function () {
      this.actionType = 'import';
      this.pemBundle = `
    -----BEGIN CERTIFICATE-----
    MIIDRTCCAi2gAwIBAgIUdKagCL6TnN5xLkwhPbNY8JEcY0YwDQYJKoZIhvcNAQEL
    BQAwGzEZMBcGA1UEAxMQd3d3LnRlc3QtaW50LmNvbTAeFw0yMzAxMDkxOTA1NTBa
    Fw0yMzAyMTAxOTA2MjBaMBsxGTAXBgNVBAMTEHd3dy50ZXN0LWludC5jb20wggEi
    MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfd5o9JfyRAXH+E1vE2U0xjSqs
    A/cxDqsDXRHBnNJvzAa+7gPKXCDQZbr6chjxLXpP6Bv2/O+dZHq1fo/f6q9PDDGW
    JYIluwbACpe7W1UB7q9xFkZg85yQsNYokGZlwv/AMGpFBxDwVlNGL+4fxvFTv7uF
    mIlDzSIPrzByyCrqAFMNNqNwlAerDt/C6DMZae/rTGXIXsTfUpxPy21bzkeA+70I
    YCV1ffK8UnAeBYNUJ+v8+XgTQ5KhRyQ+fscUkO3T2s6f3O9Q2sWxswkf2YmZB+V1
    cTZ5w6hqiuFdBXz7GRnACi1/gbWbaExQTJRplArFwIHka7dqJh8tYkXDjai3AgMB
    AAGjgYAwfjAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4E
    FgQU68/xXIgvsleKkuA8clK/6YslB/IwHwYDVR0jBBgwFoAU68/xXIgvsleKkuA8
    clK/6YslB/IwGwYDVR0RBBQwEoIQd3d3LnRlc3QtaW50LmNvbTANBgkqhkiG9w0B
    AQsFAAOCAQEAWSff0BH3SJv/XqwN/flqc1CVzOios72/IJ+KBBv0AzFCZ8wJPi+c
    hH1bw7tqi01Bgh595TctogDFN1b6pjN+jrlIP4N+FF9Moj79Q+jHQMnuJomyPuI7
    i07vqUcxgSmvEBBWOWS+/vxe6TfWDg18nyPf127CWQN8IHTo1f/GavX+XmRve6XT
    EWoqcQshEk9i87oqCbaT7B40jgjTAd1r4Cc6P4s1fAGPt9e9eqMj13kTyVDNuCoD
    FSZYalrlkASpg+c9oDQIh2MikGQINXHv/zIEHOW93siKMWeA4ni6phHtMg/p5eJt
    SxnVZsSzj8QLy2uwX1AADR0QUvJzMxptyA==
    -----END CERTIFICATE-----
    -----BEGIN RSA PRIVATE KEY-----
    MIIEowIBAAKCAQEAn3eaPSX8kQFx/hNbxNlNMY0qrAP3MQ6rA10RwZzSb8wGvu4D
    ylwg0GW6+nIY8S16T+gb9vzvnWR6tX6P3+qvTwwxliWCJbsGwAqXu1tVAe6vcRZG
    YPOckLDWKJBmZcL/wDBqRQcQ8FZTRi/uH8bxU7+7hZiJQ80iD68wcsgq6gBTDTaj
    cJQHqw7fwugzGWnv60xlyF7E31KcT8ttW85HgPu9CGAldX3yvFJwHgWDVCfr/Pl4
    E0OSoUckPn7HFJDt09rOn9zvUNrFsbMJH9mJmQfldXE2ecOoaorhXQV8+xkZwAot
    f4G1m2hMUEyUaZQKxcCB5Gu3aiYfLWJFw42otwIDAQABAoIBADC+vZ4Ne4vTtkWl
    Izsj9Y29Chs0xx3uzuWjUGcvib/0zOcWGICF8t3hCuu9btRiQ24jlFDGdnRVH5FV
    E6OtuFLgdlPgOU1RQzn2wvTZcT26+VQHLBI8xVIRTBVwNmzK06Sq6AEbrNjaenAM
    /KwoAuLHzAmFXAgmr0++DIA5oayPWyi5IoyFO7EoRv79Xz5LWfu5j8CKOFXmI5MT
    vEVYM6Gb2xHRa2Ng0SJ4VzwC09GcXlHKRAz+CubJuncvjbcM/EryvexozKkUq4XA
    KqGr9xxdZ4XDlo3Rj9S9P9JaOin0I1mwwz6p+iwMF0zr+/ldjE4oPBdB1PUgSJ7j
    2CZcS1kCgYEAwIZ3UsMIXqkMlkMz/7nu2sqzV3EgQjY5QRoz98ligKg4fhYKz+K4
    yXvJrRyLkwEBaPdLppCZbs4xsuuv3jiqUHV5n7sfpUA5HVKkKh6XY7jnszbqV732
    iB1mQVEjzM92/amew2hDKLGQDW0nglrg6uV+bx0Lnp6Glahr8NOAyk0CgYEA1Ar3
    jTqTkU+NQX7utlxx0HPVL//JH/erp/Gnq9fN8dZhK/yjwX5savUlNHpgePoXf1pE
    lgi21/INQsvp7O2AUKuj96k+jBHQ0SS58AQGFv8iNDkLE57N74vCO6+Xdi1rHj/Y
    7jglr00box/7SOmvb4SZz2o0jm0Ejsg2M0aBuRMCgYEAgTB6F34qOqMDgD1eQka5
    QfXs/Es8E1Ihf08e+jIXuC+poOoXnUINL56ySUizXBS7pnzzNbUoUFNqxB4laF/r
    4YvC7m15ocED0mpnIKBghBlK2VaLUA93xAS+XiwdcszwkuzkTUnEbyUfffL2JSHo
    dZdEDTmXV3wW4Ywfyn2Sma0CgYAeNNG/FLEg6iw9QE/ROqob/+RGyjFklGunqQ0x
    tbRo1xlQotTRI6leMz3xk91aXoYqZjmPBf7GFH0/Hr1cOxkkZM8e4MVAPul4Ybr7
    LheP/xhoSBgD24OKtGYfCoyRETdJP98vUGBN8LYXLt8lK+UKBeHDYmXKRE156ZuP
    AmRIcQKBgFvp+xMoyAsBeOlTjVDZ0mTnFh1yp8f7N3yXdHPpFShwjXjlqLmLO5RH
    mZAvaH0Ux/wCfvwHhdC46jBrs9S4zLBvj3+44NYOzvz2dBWP/5MuXgzFe30h9Yd0
    zUlyEaWm0jY2Ylzax8ECKRL0td2bv36vxOYtTax8MSB15szsnPJ+
    -----END RSA PRIVATE KEY-----
    `;
    });

    test('it serializes only valid params', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        pemBundle: this.pemBundle,
        issuerName: 'do-not-send',
        keyType: 'do-not-send',
      });
      const expectedResult = {
        pem_bundle: this.pemBundle,
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(
        serializedRecord,
        expectedResult,
        'Serializes only parameters valid for import action'
      );
    });
  });

  module('actionType generate-root', function (hooks) {
    hooks.beforeEach(function () {
      this.actionType = 'generate-root';
      this.allKeyFields = {
        keyName: 'key name',
        keyType: 'rsa',
        keyBits: '0',
        keyRef: 'key ref',
        managedKeyName: 'managed name',
        managedKeyId: 'managed id',
      };
      this.withDefaults = {
        exclude_cn_from_sans: false,
        format: 'pem',
        max_path_length: -1,
        not_before_duration: '30s',
        private_key_format: 'der',
      };
    });

    test('it serializes only params with values', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        type: 'external', // only used for endpoint in adapter
        customTtl: '40m', // UI-only value
        issuerName: 'my issuer',
        commonName: undefined,
        foo: 'bar',
      });
      const expectedResult = {
        issuer_name: 'my issuer',
        foo: 'bar',
      };

      // without passing `actionType` it will not compare against an allowlist
      const serializedRecord = record.serialize();
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = external', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'external',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_name: 'key name',
        key_type: 'rsa',
        key_bits: '0',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = internal', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'internal',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_name: 'key name',
        key_type: 'rsa',
        key_bits: '0',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = existing', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'existing',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_ref: 'key ref',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = kms', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'kms',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_name: 'key name',
        managed_key_name: 'managed name',
        managed_key_id: 'managed id',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });
  });
});
