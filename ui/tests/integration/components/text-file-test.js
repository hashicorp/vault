import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
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
    this.set('label', 'Some label');
    this.set('file', {});
    this.set('onChange', sinon.spy());
  });

  test('it renders with or without label', async function (assert) {
    this.set('inputOnly', false);

    await render(
      hbs`<TextFile @file={{this.file}} @onChange={{this.onChange}} @uploadOnly={{this.inputOnly}} />`
    );

    assert.dom(SELECTORS.label).hasText('File', 'renders default label');
    assert.dom(SELECTORS.toggle).exists({ count: 1 }, 'toggle exists');
    assert.dom(SELECTORS.fileUpload).exists({ count: 1 }, 'File input shown');

    this.set('inputOnly', true);

    assert.dom(SELECTORS.label).doesNotExist('Label no longer rendered');
    assert.dom(SELECTORS.toggle).doesNotExist('toggle no longer rendered');
    assert.dom(SELECTORS.fileUpload).exists({ count: 1 }, 'File input shown');
  });

  test('it toggles between upload and textarea', async function (assert) {
    await render(hbs`<TextFile @file={{this.file}} @onChange={{this.onChange}} />`);

    assert.dom(SELECTORS.fileUpload).exists({ count: 1 }, 'File input shown');
    assert.dom(SELECTORS.textarea).doesNotExist('Texarea hidden');
    await click(SELECTORS.toggle);
    assert.dom(SELECTORS.textarea).exists({ count: 1 }, 'Textarea shown');
    assert.dom(SELECTORS.fileUpload).doesNotExist('File upload hidden');
  });
});
