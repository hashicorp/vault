/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';

const SELECTORS = {
  label: '[data-test-certificate-label]',
  value: '[data-test-certificate-value]',
  icon: '[data-test-certificate-icon]',
  copyButton: '[data-test-copy-button]',
  copyIcon: '[data-test-icon="clipboard-copy"]',
};
const { rootPem, rootDer } = CERTIFICATES;

module('Integration | Component | certificate-card', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<CertificateCard />`);

    assert.dom(SELECTORS.label).hasNoText('There is no label because there is no value');
    assert.dom(SELECTORS.value).hasNoText('There is no value because none was provided');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyIcon).exists('The copy icon renders');
  });

  test('it renders with an example PEM Certificate', async function (assert) {
    this.certificate = rootPem;
    await render(hbs`<CertificateCard @data={{this.certificate}} />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(this.certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
    assert
      .dom(SELECTORS.copyButton)
      .hasAttribute('data-test-copy-button', this.certificate, 'copy value is the same as data');
  });

  test('it renders with an example DER Certificate', async function (assert) {
    this.certificate = rootDer;
    await render(hbs`<CertificateCard @data={{this.certificate}} />`);

    assert.dom(SELECTORS.label).hasText('DER Format', 'The label text is DER Format');
    assert.dom(SELECTORS.value).hasText(this.certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
    assert
      .dom(SELECTORS.copyButton)
      .hasAttribute('data-test-copy-button', this.certificate, 'copy value is the same as data');
  });

  test('it renders with the PEM Format label regardless of the value provided when @isPem is true', async function (assert) {
    this.certificate = 'example-certificate-text';
    await render(hbs`<CertificateCard @data={{this.certificate}} @isPem={{true}}/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(this.certificate, 'The data rendered is correct');
  });

  test('it renders with an example CA Chain', async function (assert) {
    this.caChain = [
      '-----BEGIN CERTIFICATE-----\nMIIDIDCCA...\n-----END CERTIFICATE-----\n',
      '-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBA...\n-----END RSA PRIVATE KEY-----\n',
    ];

    await render(hbs`<CertificateCard @data={{this.caChain}} />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(this.caChain.join(','), 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
    assert
      .dom(SELECTORS.copyButton)
      .hasAttribute(
        'data-test-copy-button',
        this.caChain.join('\n'),
        'copy value is array converted to a string'
      );
  });
});
