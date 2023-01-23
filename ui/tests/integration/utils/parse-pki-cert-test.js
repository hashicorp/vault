import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { parseCertificate, parseExtensions, parseSubject, formatValues } from 'vault/utils/parse-pki-cert';
import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate } from 'pkijs';
import { addHours, fromUnixTime, isSameDay } from 'date-fns';
import errorMessage from 'vault/utils/error-message';

module('Integration | Util | parse pki certificate', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.pkiCert = `-----BEGIN CERTIFICATE-----\nMIIFFTCCA/2gAwIBAgIULIZoZjgoLLQeYd/I0EQgdUegragwDQYJKoZIhvcNAQEN\nBQAwgdoxDzANBgNVBAYTBkZyYW5jZTESMBAGA1UECBMJQ2hhbXBhZ25lMQ4wDAYD\nVQQHEwVQYXJpczETMBEGA1UECRMKMjM0IHNlc2FtZTEPMA0GA1UEERMGMTIzNDU2\nMSQwDQYDVQQKEwZXaWRnZXQwEwYDVQQKEwxJbmNvcnBvcmF0ZWQxKDAOBgNVBAsT\nB0ZpbmFuY2UwFgYDVQQLEw9IdW1hbiBSZXNvdXJjZXMxGDAWBgNVBAMTD2NvbW1v\nbi1uYW1lLmNvbTETMBEGA1UEBRMKY2VyZWFsMTI5MjAeFw0yMzAxMjEwMDUyMzBa\nFw0zMzAxMTgwMDUzMDBaMIHaMQ8wDQYDVQQGEwZGcmFuY2UxEjAQBgNVBAgTCUNo\nYW1wYWduZTEOMAwGA1UEBxMFUGFyaXMxEzARBgNVBAkTCjIzNCBzZXNhbWUxDzAN\nBgNVBBETBjEyMzQ1NjEkMA0GA1UEChMGV2lkZ2V0MBMGA1UEChMMSW5jb3Jwb3Jh\ndGVkMSgwDgYDVQQLEwdGaW5hbmNlMBYGA1UECxMPSHVtYW4gUmVzb3VyY2VzMRgw\nFgYDVQQDEw9jb21tb24tbmFtZS5jb20xEzARBgNVBAUTCmNlcmVhbDEyOTIwggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDZRug7meAek7/LvKyPqVL0L9hO\n3RQrvotWAGxCUp7gEPVxVBuVH97hwfABazikQQGhXQVeISrwaX7zI945fd3dGx3R\n3iDPrGp3A8KXsaS70luMg6WyIQJ5GM21GIGchACXiIKv+Ln0++0wivFyMw8sA4V2\nbQyZHOsN5puoYqhEFyypw0E3yiyvBW7KuDrkOLzuVSCa1WdYCnpg7O1v/ViM6dIk\no83CH1p1MtQ6ZPgBfB4V6JPAm4R3zhoG0Geg3FziCXm+F2qyfbICyTQLoXXB0YD9\nE5D4jnsGwRvSLIdadxfqZCN740JOHIIZopQLhJDHNjQjTcuqtW8EhC1UJzIjAgMB\nAAGjgdAwgc0wDgYDVR0PAQH/BAQDAgEGMBIGA1UdEwEB/wQIMAYBAf8CAREwHQYD\nVR0OBBYEFAsrMoFu6tt1pybxx9ln6w5QK/2tMB8GA1UdIwQYMBaAFAsrMoFu6tt1\npybxx9ln6w5QK/2tMDcGA1UdEQQwMC6CCGFsdG5hbWUxgghhbHRuYW1lMocEwJ4B\nJoYIdGVzdHVyaTGGCHRlc3R1cmkyMC4GA1UdHgEB/wQkMCKgIDAOggxkbnNuYW1l\nMS5jb20wDoIMZHNubmFtZTIuY29tMA0GCSqGSIb3DQEBDQUAA4IBAQCLIQ/AEVME\n5F9N5kqT0PdJ7PgjCHraWnEa25TH7RxH5mh6BakuUkJr5TFnytDU6TwkVfixgT9j\nT6O+BdB6ILv1u3ECGBQNObq1HtO0NM/Q1IZewEUNIjDVfdXFIxHLLlyxoGiCV/PS\nm/QHHX6K7EezAIdw4OvvO5lfjOzPZ6vaWEab1BCCPgxaWOqQ4U6MX3NzLiP5VqTs\npMFoLJ0yG1yMkW0pr8d1NkqDoZI1JW/DGrQEdYg182ckHogjmjydVE0B00yCzGHh\nOYqj7AHqjkpa9DMZMH22reuiSGNun7o2jEQ9iRt79UEpqkIap3aohsypeqgYCMGf\n6V/JEhjKPzap\n-----END CERTIFICATE-----`;
    this.pssCert = `-----BEGIN CERTIFICATE-----\nMIIDqTCCAl2gAwIBAgIUVY2PTRZl1t/fjfyEwrG4HvGjYekwQQYJKoZIhvcNAQEK\nMDSgDzANBglghkgBZQMEAgEFAKEcMBoGCSqGSIb3DQEBCDANBglghkgBZQMEAgEF\nAKIDAgEgMBoxGDAWBgNVBAMTD2NvbW1vbi1uYW1lLmNvbTAeFw0yMzAxMjEwMTA3\nNDBaFw0yMzAyMjIwMTA4MTBaMBoxGDAWBgNVBAMTD2NvbW1vbi1uYW1lLmNvbTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANlG6DuZ4B6Tv8u8rI+pUvQv\n2E7dFCu+i1YAbEJSnuAQ9XFUG5Uf3uHB8AFrOKRBAaFdBV4hKvBpfvMj3jl93d0b\nHdHeIM+sancDwpexpLvSW4yDpbIhAnkYzbUYgZyEAJeIgq/4ufT77TCK8XIzDywD\nhXZtDJkc6w3mm6hiqEQXLKnDQTfKLK8Fbsq4OuQ4vO5VIJrVZ1gKemDs7W/9WIzp\n0iSjzcIfWnUy1Dpk+AF8HhXok8CbhHfOGgbQZ6DcXOIJeb4XarJ9sgLJNAuhdcHR\ngP0TkPiOewbBG9Ish1p3F+pkI3vjQk4cghmilAuEkMc2NCNNy6q1bwSELVQnMiMC\nAwEAAaN/MH0wDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0O\nBBYEFAsrMoFu6tt1pybxx9ln6w5QK/2tMB8GA1UdIwQYMBaAFAsrMoFu6tt1pybx\nx9ln6w5QK/2tMBoGA1UdEQQTMBGCD2NvbW1vbi1uYW1lLmNvbTBBBgkqhkiG9w0B\nAQowNKAPMA0GCWCGSAFlAwQCAQUAoRwwGgYJKoZIhvcNAQEIMA0GCWCGSAFlAwQC\nAQUAogMCASADggEBAFh+PMwEmxaZR6OtfB0Uvw2vA7Oodmm3W0bYjQlEz8U+Q+JZ\ncIPa4VnRy1QALmKbPCbRApA/gcWzIwtzo1JhLtcDINg2Tl0nj4WvgpIvj0/lQNMq\nmwP7G/K4PyJTv3+y5XwVfepZAZITB0w5Sg5dLC6HP8AGVIaeb3hGNHYvPlE+pbT+\njL0xxzFjOorWoy5fxbWoVyVv9iZ4j0zRnbkYHIi3d8g56VV6Rbyw4WJt6p87lmQ8\n0wbiJTtuew/0Rpuc3PEcR9XfB5ct8bvaGGTSTwh6JQ33ohKKAKjbBNmhBDSP1thQ\n2mTkms/mbDRaTiQKHZx25TmOlLN5Ea1TSS0K6yw=\n-----END CERTIFICATE-----`;
    this.onlyCn = `-----BEGIN CERTIFICATE-----\nMIIDQTCCAimgAwIBAgIUVQy58VgdVpAK9c8SfS31idSv6FUwDQYJKoZIhvcNAQEL\nBQAwGjEYMBYGA1UEAxMPY29tbW9uLW5hbWUuY29tMB4XDTIzMDEyMTAxMjAyOVoX\nDTIzMDIyMjAxMjA1OVowGjEYMBYGA1UEAxMPY29tbW9uLW5hbWUuY29tMIIBIjAN\nBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2UboO5ngHpO/y7ysj6lS9C/YTt0U\nK76LVgBsQlKe4BD1cVQblR/e4cHwAWs4pEEBoV0FXiEq8Gl+8yPeOX3d3Rsd0d4g\nz6xqdwPCl7Gku9JbjIOlsiECeRjNtRiBnIQAl4iCr/i59PvtMIrxcjMPLAOFdm0M\nmRzrDeabqGKoRBcsqcNBN8osrwVuyrg65Di87lUgmtVnWAp6YOztb/1YjOnSJKPN\nwh9adTLUOmT4AXweFeiTwJuEd84aBtBnoNxc4gl5vhdqsn2yAsk0C6F1wdGA/ROQ\n+I57BsEb0iyHWncX6mQje+NCThyCGaKUC4SQxzY0I03LqrVvBIQtVCcyIwIDAQAB\no38wfTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU\nCysygW7q23WnJvHH2WfrDlAr/a0wHwYDVR0jBBgwFoAUCysygW7q23WnJvHH2Wfr\nDlAr/a0wGgYDVR0RBBMwEYIPY29tbW9uLW5hbWUuY29tMA0GCSqGSIb3DQEBCwUA\nA4IBAQDPco+FIHXczf0HTwFAmIVu4HKaeIwDsVPxoUqqWEix8AyCsB5uqpKZasby\nedlrdBohM4dnoV+VmV0de04y95sdo3Ot60hm/czLog3tHg4o7AmfA7saS+5hCL1M\nCJWqoJHRFo0hOWJHpLJRWz5DqRZWspASoVozLOYyjRD+tNBjO5hK4FtaG6eri38t\nOpTt7sdInVODlntpNuuCVprPpHGj4kPOcViQULoFQq5fwyadpdjqSXmEGlt0to5Y\nMbTb4Jhj0HywgO53BUUmMzzY9idXh/8A7ThrM5LtqhxaYHLVhyeo+5e0mgiXKp+n\nQ8Uh4TNNTCvOUlAHycZNaxYTlEPn\n-----END CERTIFICATE-----`;
    // contains unsupported subject and extension oids
    this.unsupportedSANS = `-----BEGIN CERTIFICATE-----\nMIIEjDCCA3SgAwIBAgIUD4EeORgh/i+ZZFOk8KsGKQPWsoIwDQYJKoZIhvcNAQEL\nBQAwgZIxMTAvBgNVBAMMKGZhbmN5LWNlcnQtdW5zdXBwb3J0ZWQtc3Viai1hbmQt\nZXh0LW9pZHMxCzAJBgNVBAYTAlVTMQ8wDQYDVQQIDAZLYW5zYXMxDzANBgNVBAcM\nBlRvcGVrYTESMBAGA1UECgwJQWNtZSwgSW5jMRowGAYJKoZIhvcNAQkBFgtmb29A\nYmFyLmNvbTAeFw0yMzAxMjMxODQ3MjNaFw0zMzAxMjAxODQ3MjNaMIGSMTEwLwYD\nVQQDDChmYW5jeS1jZXJ0LXVuc3VwcG9ydGVkLXN1YmotYW5kLWV4dC1vaWRzMQsw\nCQYDVQQGEwJVUzEPMA0GA1UECAwGS2Fuc2FzMQ8wDQYDVQQHDAZUb3Bla2ExEjAQ\nBgNVBAoMCUFjbWUsIEluYzEaMBgGCSqGSIb3DQEJARYLZm9vQGJhci5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDyYH5qS7krfZ2tA5uZsY2qXbTb\ntGNG1BsyDhZ/qqVlQybjDsHJZwNUbpfhBcCLaKyAwH1R9n54NOOOn6bYgfKWTgy3\nL7224YDAqYe7Y/GPjgI2MRvRfn6t2xzQxtJ0l0k8LeyNcwhiqYLQyOOfDdc127fm\nW40r2nmhLpH0i9e2I/YP1HQ+ldVgVBqeUTntgVSBfrQF56v9mAcvvHEa5sdHqmX4\nJ2lhWTnx9jqb7NZxCem76BlX1Gt5TpP3Ym2ZFVQI9fuPK4O8JVhk1KBCmIgR3Ft+\nPpFUs/c41EMunKJNzveYrInSDScaC6voIJpK23nMAiM1HckLfUUc/4UojD+VAgMB\nAAGjgdcwgdQwHQYDVR0OBBYEFH7tt4enejKTZtYjUKUUx6PXyzlgMB8GA1UdIwQY\nMBaAFH7tt4enejKTZtYjUKUUx6PXyzlgMA4GA1UdDwEB/wQEAwIFoDAgBgNVHSUB\nAf8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBCjBM\nBgNVHREERTBDhwTAngEmhgx1cmlTdXBwb3J0ZWSCEWRucy1OYW1lU3VwcG9ydGVk\noBoGAyoDBKATDBFleGFtcGxlIG90aGVybmFtZTANBgkqhkiG9w0BAQsFAAOCAQEA\nP6ckVJgbcJue+MK3RVDuG+Mh7dl89ynC7NwpQFRjLVZQuoMHZT/dcLlVeFejVXu5\nR+IPLmQU6NV7JAmy4zGap8awf12QTy3g410ecrSF94WWlu8bPoekfUnnP+kfzLPH\nCUAkRKxWDSRKX5C8cMMxacVBBaBIayuusLcHkHmxLLDw34PFzyz61gtZOJq7JYnD\nhU9YsNh6bCDmnBDBsDMOI7h8lBRQwTiWVoSD9YNVvFiY29YvFbJQGdh+pmBtf7E+\n1B/0t5NbvqlQSbhMM0QgYFhuCxr3BGNob7kRjgW4i+oh+Nc5ptA5q70QMaYudqRS\nd8SYWhRdxmH3qcHNPcR1iw==\n-----END CERTIFICATE-----`;
    this.getErrorMessages = (certErrors) => certErrors.map((error) => errorMessage(error));
  });

  test('it parses a certificate with supported values', async function (assert) {
    assert.expect(2);
    const parsedCert = parseCertificate(this.pkiCert);
    assert.propEqual(
      parsedCert,
      {
        alt_names: 'altname1, altname2',
        can_parse: true,
        common_name: 'common-name.com',
        country: 'France',
        exclude_cn_from_sans: true,
        expiry_date: {},
        ip_sans: 'OCTET STRING : C09E0126', // when parsed, should be 192.158.1.38
        issue_date: {},
        locality: 'Paris',
        max_path_length: 17,
        not_valid_after: 1989622380,
        not_valid_before: 1674262350,
        organization: 'Widget',
        ou: 'Finance',
        parsing_errors: [],
        permitted_dns_domains: 'dnsname1.com, dsnname2.com',
        postal_code: '123456',
        province: 'Champagne',
        serial_number: 'cereal1292',
        signature_bits: '512',
        street_address: '234 sesame',
        ttl: '87600h',
        uri_sans: 'testuri1, testuri2',
        use_pss: false,
      },
      'it contains expected attrs, cn is excluded from alt_names (exclude_cn_from_sans: true)'
    );
    assert.ok(
      isSameDay(
        addHours(fromUnixTime(parsedCert.not_valid_before), Number(parsedCert.ttl.split('h')[0])),
        fromUnixTime(parsedCert.not_valid_after),
        'ttl value is correct'
      )
    );
  });

  test('it parses a certificate with use_pass=true and exclude_cn_from_sans=false', async function (assert) {
    assert.expect(2);
    const parsedCert = parseCertificate(this.pssCert);
    assert.propContains(
      parsedCert,
      { signature_bits: '256', ttl: '768h', use_pss: true },
      'returns signature_bits value and use_pss is true'
    );
    assert.propContains(
      parsedCert,
      {
        alt_names: 'common-name.com',
        can_parse: true,
        common_name: 'common-name.com',
        exclude_cn_from_sans: false,
      },
      'common name is included in alt_names'
    );
  });

  test('it returns parsing_errors when certificate has unsupported values', async function (assert) {
    assert.expect(1);
    const parsedCert = parseCertificate(this.unsupportedSANS);
    const parsingErrors = this.getErrorMessages(parsedCert.parsing_errors);
    assert.propEqual(
      parsingErrors,
      [
        'certificate contains unsupported subject OIDs',
        'certificate contains unsupported extension OIDs',
        'subjectAltName contains unsupported types',
      ],
      'it contains expected error messages'
    );
  });

  test('it returns attr as null if nonexistent', async function (assert) {
    assert.expect(1);
    const parsedCert = parseCertificate(this.onlyCn);
    assert.propEqual(
      parsedCert,
      {
        alt_names: 'common-name.com',
        can_parse: true,
        common_name: 'common-name.com',
        country: null,
        exclude_cn_from_sans: false,
        expiry_date: {},
        ip_sans: null,
        issue_date: {},
        locality: null,
        max_path_length: undefined,
        not_valid_after: 1677028859,
        not_valid_before: 1674264029,
        organization: null,
        ou: null,
        parsing_errors: [],
        postal_code: null,
        province: null,
        serial_number: null,
        signature_bits: '256',
        street_address: null,
        ttl: '768h',
        uri_sans: null,
        use_pss: false,
      },
      'it contains expected attrs'
    );
  });

  test('the helper parseSubject returns expected object', async function (assert) {
    const certSchema = (cert) => {
      const cert_base64 = cert.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, '');
      const cert_der = fromBase64(cert_base64);
      const cert_asn1 = asn1js.fromBER(stringToArrayBuffer(cert_der));
      return new Certificate({ schema: cert_asn1.result });
    };

    const supportedCert = certSchema(this.pkiCert);
    assert.propEqual(
      parseSubject(supportedCert.subject.typesAndValues),
      {
        subjErrors: [],
        subjValues: {
          common_name: 'common-name.com',
          country: 'France',
          locality: 'Paris',
          organization: 'Widget',
          ou: 'Finance',
          postal_code: '123456',
          province: 'Champagne',
          serial_number: 'cereal1292',
          street_address: '234 sesame',
        },
      },
      'it returns supported subject values'
    );
    const unsupportedCert = certSchema(this.unsupportedSANS);
    const parsedSubject = parseSubject(unsupportedCert.subject.typesAndValues);
    assert.propEqual(
      this.getErrorMessages(parsedSubject.subjErrors),
      ['certificate contains unsupported subject OIDs'],
      'it returns subject errors'
    );
    assert.ok(parsedSubject.subjErrors[0] instanceof Error, 'subjErrors contain error objects');
  });
});
