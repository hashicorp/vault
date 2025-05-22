/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const SELECTORS = {
  label: '[data-test-certificate-label]',
  value: '[data-test-certificate-value]',
  icon: '[data-test-certificate-icon]',
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
    // Sinon spy for clipboard
    const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();

    this.certificate = rootPem;
    await render(hbs`<CertificateCard @data={{this.certificate}} />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(this.certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');

    await click(GENERAL.copyButton);
    assert.true(clipboardSpy.calledOnce, 'Clipboard was called once');
    assert.strictEqual(
      clipboardSpy.firstCall.args[0],
      this.certificate,
      'copy certificate copied the correct text'
    );
    // Restore original clipboard
    clipboardSpy.restore(); // cleanup
  });

  test('it renders with an example DER Certificate', async function (assert) {
    // Sinon spy for clipboard
    const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();

    this.certificate = rootDer;
    await render(hbs`<CertificateCard @data={{this.certificate}} />`);

    assert.dom(SELECTORS.label).hasText('DER Format', 'The label text is DER Format');
    assert.dom(SELECTORS.value).hasText(this.certificate, 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');

    await click(GENERAL.copyButton);
    assert.true(clipboardSpy.calledOnce, 'Clipboard was called once');
    assert.strictEqual(
      clipboardSpy.firstCall.args[0],
      this.certificate,
      'copy certificate copied the correct text'
    );
    // Restore original clipboard
    clipboardSpy.restore(); // cleanup
  });

  test('it renders with the PEM Format label regardless of the value provided when @isPem is true', async function (assert) {
    this.certificate = 'example-certificate-text';
    await render(hbs`<CertificateCard @data={{this.certificate}} @isPem={{true}}/>`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(this.certificate, 'The data rendered is correct');
  });

  test('it renders with an example CA Chain', async function (assert) {
    // Sinon spy for clipboard
    const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();

    this.caChain = [
      '-----BEGIN CERTIFICATE-----\nMIIDIDCCA...\n-----END CERTIFICATE-----\n',
      '-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBA...\n-----END RSA PRIVATE KEY-----\n',
    ];

    await render(hbs`<CertificateCard @data={{this.caChain}} />`);

    assert.dom(SELECTORS.label).hasText('PEM Format', 'The label text is PEM Format');
    assert.dom(SELECTORS.value).hasText(this.caChain.join(','), 'The data rendered is correct');
    assert.dom(SELECTORS.icon).exists('The certificate icon exists');

    await click(GENERAL.copyButton);
    assert.true(clipboardSpy.calledOnce, 'Clipboard was called once');
    assert.strictEqual(
      clipboardSpy.firstCall.args[0],
      this.caChain.join('\n'),
      'copy value is array converted to a string'
    );
    // Restore original clipboard
    clipboardSpy.restore(); // cleanup
  });
});
