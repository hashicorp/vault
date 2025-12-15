/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const { componentPemBundle } = CERTIFICATES;

module('Integration | Component | text-file', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.label = 'Some label';
    this.onChange = sinon.spy();
    this.owner.lookup('service:flash-messages').registerTypes(['danger']);
  });

  test('it renders with label and toggle by default', async function (assert) {
    await render(hbs`<TextFile @onChange={{this.onChange}} @helpText="this is my help text"/>`);
    const [firstLabel, secondLabel] = findAll('label');
    assert.dom(firstLabel).hasText('File', 'renders default label');
    assert.dom(secondLabel).hasText('Enter as text', 'it renders toggle label');
    assert.dom(GENERAL.textToggle).exists({ count: 1 }, 'toggle exists');
    assert.dom(GENERAL.fileInput).exists({ count: 1 }, 'File input shown');

    assert.dom(GENERAL.tooltipText).hasNoText();
    await click(GENERAL.tooltip('text-file'));
    assert.dom(GENERAL.tooltipText).hasText('this is my help text', 'Tooltip text renders');
  });

  test('it renders without toggle and option for text input when uploadOnly=true', async function (assert) {
    await render(hbs`<TextFile @onChange={{this.onChange}} @uploadOnly={{true}} />`);

    assert.dom('label').exists({ count: 1 }).hasText('File', 'it renders a label just for the file');
    assert.dom(GENERAL.textToggle).doesNotExist('toggle no longer rendered');
    assert.dom(GENERAL.tooltip('text-file')).doesNotExist('tooltip icon no longer rendered');
    assert.dom(GENERAL.fileInput).exists({ count: 1 }, 'File input shown');
  });

  test('it toggles between upload and textarea', async function (assert) {
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);

    assert.dom(GENERAL.fileInput).exists({ count: 1 }, 'File input shown');
    assert.dom(GENERAL.maskedInput).doesNotExist('Texarea hidden');
    await click(GENERAL.textToggle);
    assert.dom(GENERAL.maskedInput).exists({ count: 1 }, 'Textarea shown');
    assert.dom(GENERAL.fileInput).doesNotExist('File upload hidden');
  });

  test('it correctly parses uploaded files', async function (assert) {
    const file = new Blob([['some content for a file']], { type: 'text/plain' });
    file.name = 'filename.txt';
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);
    await triggerEvent(GENERAL.fileInput, 'change', { files: [file] });
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
    const PEM_BUNDLE = componentPemBundle;

    await render(hbs`<TextFile @onChange={{this.onChange}} />`);
    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, PEM_BUNDLE);
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

    await triggerEvent(GENERAL.fileInput, 'change', { files: [this.file] });
    assert.dom(GENERAL.inlineAlert).hasText('There was a problem uploading. Please try again.');
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
