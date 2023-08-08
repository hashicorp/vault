/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { rootPem } from 'vault/tests/helpers/pki/values';

const SELECTORS = {
  label: '[data-test-certificate-label]',
  value: '[data-test-certificate-value]',
  icon: '[data-test-certificate-icon]',
  copyButton: '[data-test-copy-button]',
};

module('Integration | Component | certificate-card', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<CertificateCard />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is correct');
    assert.dom(SELECTORS.value).hasNoText('There is no value for the certificate card');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
  });

  test('it renders with a small example value for certificate ', async function (assert) {
    await render(hbs`<CertificateCard @certificateValue="test"/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is correct');
    assert.dom(SELECTORS.value).hasText('test', 'The value for the certificate is correct');
  });

  test('it renders with an example PEM Certificate', async function (assert) {
    const certificate = rootPem;
    this.set('certificate', certificate);
    await render(hbs`<CertificateCard @certificateValue={{this.certificate}}/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is correct');
    assert.dom(SELECTORS.value).hasText(certificate, 'The value for the certificate is correct');
  });
});
