import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

const SELECTORS = {
  label: '[data-test-text-file-label]',
  toggle: '[data-test-text-toggle]',
  textarea: '[data-test-text-file-textarea]',
  fileUpload: '[data-test-text-file-input]',
};
module('Integration | Component | text-file', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.label = 'Some label';
    this.onChange = sinon.spy();
    this.owner.lookup('service:flash-messages').registerTypes(['danger']);
  });

  test('it renders with label and toggle by default', async function (assert) {
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);

    assert.dom(SELECTORS.label).hasText('File', 'renders default label');
    assert.dom(SELECTORS.toggle).exists({ count: 1 }, 'toggle exists');
    assert.dom(SELECTORS.fileUpload).exists({ count: 1 }, 'File input shown');
  });

  test('it renders without toggle and option for text input when uploadOnly=true', async function (assert) {
    await render(hbs`<TextFile @onChange={{this.onChange}} @uploadOnly={{true}} />`);

    assert.dom(SELECTORS.label).doesNotExist('Label no longer rendered');
    assert.dom(SELECTORS.toggle).doesNotExist('toggle no longer rendered');
    assert.dom(SELECTORS.fileUpload).exists({ count: 1 }, 'File input shown');
  });

  test('it toggles between upload and textarea', async function (assert) {
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);

    assert.dom(SELECTORS.fileUpload).exists({ count: 1 }, 'File input shown');
    assert.dom(SELECTORS.textarea).doesNotExist('Texarea hidden');
    await click(SELECTORS.toggle);
    assert.dom(SELECTORS.textarea).exists({ count: 1 }, 'Textarea shown');
    assert.dom(SELECTORS.fileUpload).doesNotExist('File upload hidden');
  });

  test('it correctly parses uploaded files', async function (assert) {
    this.file = new File(['some content for a file'], 'filename.txt');
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);
    await triggerEvent(SELECTORS.fileUpload, 'change', { files: [this.file] });
    assert.propEqual(
      this.onChange.lastCall.args[0],
      {
        filename: 'filename.txt',
        value: 'some content for a file',
      },
      'parent callback function is called with correct arguments'
    );
  });

  test('it correctly submits text input', async function (assert) {
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
-----END CERTIFICATE-----`;

    await render(hbs`<TextFile @onChange={{this.onChange}} />`);
    await click(SELECTORS.toggle);
    await fillIn(SELECTORS.textarea, PEM_BUNDLE);
    assert.propEqual(
      this.onChange.lastCall.args[0],
      {
        filename: '',
        value: PEM_BUNDLE,
      },
      'parent callback function is called with correct text area input'
    );
  });

  test('it throws an error when it cannot read the file', async function (assert) {
    this.file = { foo: 'bar' };
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);

    await triggerEvent(SELECTORS.fileUpload, 'change', { files: [this.file] });
    assert
      .dom('[data-test-field-validation="text-file"]')
      .hasText('There was a problem uploading. Please try again.');
    assert.propEqual(
      this.onChange.lastCall.args[0],
      {
        filename: '',
        value: '',
      },
      'parent callback function is called with cleared out values'
    );
  });
});
