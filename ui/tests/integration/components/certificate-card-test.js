/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { rootPem } from 'vault/tests/helpers/pki/values';
import { rootDer } from 'vault/tests/helpers/pki/values';

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

    assert.dom(SELECTORS.label).hasNoText('There is no label because there is no value');
    assert.dom(SELECTORS.value).hasNoText('There is no value because none was provided');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
  });

  test('it renders with an example PEM Certificate', async function (assert) {
    const certificate = rootPem;
    this.set('certificate', certificate);
    await render(hbs`<CertificateCard @data={{this.certificate}} />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
  });

  test('it renders with an example DER Certificate', async function (assert) {
    const certificate = rootDer;
    this.set('certificate', certificate);
    await render(hbs`<CertificateCard @data={{this.certificate}} />`);

    assert.dom(SELECTORS.label).hasText('DER Format', 'The label text is DER Format');
    assert.dom(SELECTORS.value).hasText(certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
  });

  test('it renders with the PEM Format label regardless of the value provided when @isPem is true', async function (assert) {
    const certificate = 'example-certificate-text';
    this.set('certificate', certificate);
    await render(hbs`<CertificateCard @data={{this.certificate}} @isPem={{true}}/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');
    assert.dom(SELECTORS.copyButton).exists('The copy button exists');
  });
});
