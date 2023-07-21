/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const SELECTORS = {
  label: '[data-test-certificate-label]',
  value: '[data-test-certificate-value]',
};

module('Integration | Component | certificate-card', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders without a certificate value', async function (assert) {
    await render(hbs`<CertificateCard />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is correct');
    assert.dom(SELECTORS.value).hasNoText('The is no value for the certificate');
  });

  test('it renders with a small example value for certificate ', async function (assert) {
    await render(hbs`<CertificateCard @certificateValue="test"/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is correct');
    assert.dom(SELECTORS.value).hasText('test', 'The value for the certificate is correct');
  });

  test('it renders with an example Kubernetes CA Certificate', async function (assert) {
    const certificate = `
      -----BEGIN CERTIFICATE-----
      MIICUTCCAfugAwIBAgIBADANBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL
      MAkGA1UECBMCUE4xCzAJBgNVBAcTAkNOMQswCQYDVQQKEwJPTjELMAkGA1UECxMC
      VU4xFDASBgNVBAMTC0hlcm9uZyBZYW5nMB4XDTA1MDcxNTIxMTk0N1oXDTA1MDgx
      NDIxMTk0N1owVzELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAlBOMQswCQYDVQQHEwJD
      TjELMAkGA1UEChMCT04xCzAJBgNVBAsTAlVOMRQwEgYDVQQDEwtIZXJvbmcgWWFu
      ZzBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQCp5hnG7ogBhtlynpOS21cBewKE/B7j
      V14qeyslnr26xZUsSVko36ZnhiaO/zbMOoRcKK9vEcgMtcLFuQTWDl3RAgMBAAGj
      gbEwga4wHQYDVR0OBBYEFFXI70krXeQDxZgbaCQoR4jUDncEMH8GA1UdIwR4MHaA
      FFXI70krXeQDxZgbaCQoR4jUDncEoVukWTBXMQswCQYDVQQGEwJDTjELMAkGA1UE
      CBMCUE4xCzAJBgNVBAcTAkNOMQswCQYDVQQKEwJPTjELMAkGA1UECxMCVU4xFDAS
      BgNVBAMTC0hlcm9uZyBZYW5nggEAMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEE
      BQADQQA/ugzBrjjK9jcWnDVfGHlk3icNRq0oV7Ri32z/+HQX67aRfgZu7KWdI+Ju
      Wm7DCfrPNGVwFWUQOmsPue9rZBgO
      -----END CERTIFICATE-----
    `;
    this.set('certificate', certificate);
    await render(hbs`<CertificateCard @certificateValue={{this.certificate}}/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is correct');
    assert.dom(SELECTORS.value).hasText(certificate, 'The value for the CA Certificate is correct');
  });
});
