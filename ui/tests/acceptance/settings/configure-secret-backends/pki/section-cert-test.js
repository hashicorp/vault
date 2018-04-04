import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section-cert';

moduleForAcceptance('Acceptance | settings/configure/secrets/pki/cert', {
  beforeEach() {
    return authLogin();
  },
});

const PEM_BUNDLE = `-----BEGIN CERTIFICATE-----
MIIDGjCCAgKgAwIBAgIUFvnhb2nQ8+KNS3SzjlfYDMHGIRgwDQYJKoZIhvcNAQEL
BQAwDTELMAkGA1UEAxMCZmEwHhcNMTgwMTEwMTg1NDI5WhcNMTgwMjExMTg1NDU5
WjANMQswCQYDVQQDEwJmYTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AN2VtBn6EMlA4aYre/xoKHxlgNDxJnfSQWfs6yF/K201qPnt4QF9AXChatbmcKVn
OaURq+XEJrGVgF/u2lSos3NRZdhWVe8o3/sOetsGxcrd0gXAieOSmkqJjp27bYdl
uY3WsxhyiPvdfS6xz39OehsK/YCB6qCzwB4eEfSKqbkvfDL9sLlAiOlaoHC9pczf
6/FANKp35UDwInSwmq5vxGbnWk9zMkh5Jq6hjOWHZnVc2J8J49PYvkIM8uiHDgOE
w71T2xM5plz6crmZnxPCOcTKIdF7NTEP2lUfiqc9lONV9X1Pi4UclLPHJf5bwTmn
JaWgbKeY+IlF61/mgxzhC7cCAwEAAaNyMHAwDgYDVR0PAQH/BAQDAgEGMA8GA1Ud
EwEB/wQFMAMBAf8wHQYDVR0OBBYEFLDtc6+HZN2lv60JSDAZq3+IHoq7MB8GA1Ud
IwQYMBaAFLDtc6+HZN2lv60JSDAZq3+IHoq7MA0GA1UdEQQGMASCAmZhMA0GCSqG
SIb3DQEBCwUAA4IBAQDVt6OddTV1MB0UvF5v4zL1bEB9bgXvWx35v/FdS+VGn/QP
cC2c4ZNukndyHhysUEPdqVg4+up1aXm4eKXzNmGMY/ottN2pEhVEWQyoIIA1tH0e
8Kv/bysYpHZKZuoGg5+mdlHS2p2Dh2bmYFyBLJ8vaeP83NpTs2cNHcmEvWh/D4UN
UmYDODRN4qh9xYruKJ8i89iMGQfbdcq78dCC4JwBIx3bysC8oF4lqbTYoYNVTnAi
LVqvLdHycEOMlqV0ecq8uMLhPVBalCmIlKdWNQFpXB0TQCsn95rCCdi7ZTsYk5zv
Q4raFvQrZth3Cz/X5yPTtQL78oBYrmHzoQKDFJ2z
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA3ZW0GfoQyUDhpit7/GgofGWA0PEmd9JBZ+zrIX8rbTWo+e3h
AX0BcKFq1uZwpWc5pRGr5cQmsZWAX+7aVKizc1Fl2FZV7yjf+w562wbFyt3SBcCJ
45KaSomOnbtth2W5jdazGHKI+919LrHPf056Gwr9gIHqoLPAHh4R9IqpuS98Mv2w
uUCI6VqgcL2lzN/r8UA0qnflQPAidLCarm/EZudaT3MySHkmrqGM5YdmdVzYnwnj
09i+Qgzy6IcOA4TDvVPbEzmmXPpyuZmfE8I5xMoh0Xs1MQ/aVR+Kpz2U41X1fU+L
hRyUs8cl/lvBOaclpaBsp5j4iUXrX+aDHOELtwIDAQABAoIBACLdk2Ei/9Eq7FaB
MRkeKoCoWASIbU0dQD1iAf1bTTH554Sr8WOSj89xFqaJy9+6xk864Jleq9f1diWi
J6h6gwH6JNRNgWgIPnX6aUpdXnH1RT6ydP/h6XUg/9fBzhIn53Jx/ewy2WsIBtJ6
F/QoHP50VD8MMibnIaubf6fCycHhc97u4BKM2QdnAugn1sWjSiTIoYmFw/3Ej8mB
bItLWZTg9oMASgCtDwPEstlKn7yPqirOJj+G/a+6sIcP2fynd0fISsfLZ0ovN+yW
d3SV3orC0RNj83GVwYykqwCc/3pP0mRfX9fl8DKbXusITqUiGL8LGb+H6YDDpbNU
5Fj7VwECgYEA5P6aIcGfCZayEJtHKlTCA2/KBkGTOP/0iNKWhntBQT/GK+bjmr+D
GO1zR8ZFEIRdlUA5MjC9wU2AQikgFQzzmtz604Wt34fDN2NFrxq8sWN7Hjr65Fjf
ivJ6faT5r5gcNEq3EM/GLF9oJH8M+B5ccFe9iXH8AbmZHOO0FZtYxIcCgYEA97dm
Kj1qyuKlINXKt4KXdYMuIT+Z3G1B92wNN9TY/eJZgCJ7zlNcinUF/OFbiGgsk4t+
P0yVMs8BENQML0TH4Gebf4HfnDFno4J1M9HDt6HSMhsLKyvFYjFvb8hF4SPrY1pF
wW3lM3zMMzAVi8044vRrTvxfxL8QJX+1Hesye1ECgYAT5/H8Fzm8+qWV/fmMu3t2
EwSr0I18uftG3Y+KNzKv+lw+ur50WEuMIjAQQDMGwYrlC4UtUMFeCV+p4KtSSSLw
Bl+jfY5kzQdyTCXll9xpSy2LrjLbIMKl8Hgnbezqj7176jbJtlYSy2RhL84vz2vX
tDjcttTiTYD62uxvqGZqBwKBgFQ3tPM9aDZL8coFBWN4cZfRHnjNT7kCKEA/KwtF
QPSn5LfMgXz3GGo2OO/tihoJGMac0TIiDkN03y7ieLYFU1L2xoYGGIjYvxx2+PPC
KCEhUf4Y9aYavoOQvQsq8p8FgDyJ71dAzoC/uAjbGygpgGKgqG71HHYeYxXsoh3m
3YXRAoGAE7MBnVJWiIN5s63gGz9f9V6k1dPLfcPE1I0KMM/SDOIV0oLMsYQecTTB
ZzkXwRCdcJARkaKulTfjby7+oGpQydP8iZr+CNKFxwf838UbhhsXHnN6rc62qzYD
BXUV2Uwtxf+QCphnlht9muX2fsLIzDJea0JipWj1uf2H8OZsjE8=
-----END RSA PRIVATE KEY-----`;

const mountAndNav = (assert, prefix) => {
  const path = `${prefix}pki-${new Date().getTime()}`;
  mountSupportedSecretBackend(assert, 'pki', path);
  page.visit({ backend: path });
  return path;
};

test('cert config: generate', function(assert) {
  mountAndNav(assert);
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');
  });

  page.form.generateCA();
  andThen(() => {
    assert.ok(page.form.rows().count > 0, 'shows all of the rows');
    assert.ok(page.form.certificateIsPresent, 'the certificate is included');
  });

  page.form.back();
  page.form.generateCA();
  andThen(() => {
    assert.ok(
      page.flash.latestMessage.includes('You tried to generate a new root CA'),
      'shows warning message'
    );
  });
});

test('cert config: upload', function(assert) {
  mountAndNav(assert);
  andThen(() => {
    assert.equal(page.form.downloadLinks().count, 0, 'there are no download links');
  });

  page.form.uploadCA(PEM_BUNDLE);
  andThen(() => {
    assert.ok(
      page.flash.latestMessage.startsWith('The certificate for this backend has been updated'),
      'flash message displays properly'
    );
  });
});

test('cert config: sign intermediate and set signed intermediate', function(assert) {
  let csrVal, intermediateCert;
  const rootPath = mountAndNav(assert, 'root-');
  page.form.generateCA();

  const intermediatePath = mountAndNav(assert, 'intermediate-');
  page.form.generateCA('Intermediate CA', 'intermediate');
  andThen(() => {
    // cache csr
    csrVal = page.form.csr;
  });
  page.form.back();

  page.visit({ backend: rootPath });
  page.form.signIntermediate('Intermediate CA');
  andThen(() => {
    page.form.csrField(csrVal).submit();
  });
  andThen(() => {
    intermediateCert = page.form.certificate;
  });
  page.form.back();
  page.visit({ backend: intermediatePath });

  andThen(() => {
    page.form.setSignedIntermediateBtn().signedIntermediate(intermediateCert).submit();
  });
  andThen(() => {
    assert.ok(
      page.flash.latestMessage.startsWith('The certificate for this backend has been updated'),
      'flash message displays properly'
    );
    assert.equal(page.form.downloadLinks().count, 3, 'includes the caChain download link');
  });
});
