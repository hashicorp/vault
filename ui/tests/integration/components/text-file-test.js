/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { componentPemBundle } from 'vault/tests/helpers/pki/values';
import { setRunOptions } from 'ember-a11y-testing/test-support';

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
    setRunOptions({
      rules: {
        // TODO: fix textFile / replace with HDS
        label: { enabled: false },
        'label-title-only': { enabled: false },
      },
    });

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
    const file = new Blob([['some content for a file']], { type: 'text/plain' });
    file.name = 'filename.txt';
    await render(hbs`<TextFile @onChange={{this.onChange}} />`);
    await triggerEvent(SELECTORS.fileUpload, 'change', { files: [file] });
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
