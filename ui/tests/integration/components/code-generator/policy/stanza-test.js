/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { ACL_CAPABILITIES, PolicyStanza } from 'core/utils/code-generators/policy';

module('Integration | Component | code-generator/policy/stanza', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.stanza = new PolicyStanza();
    this.onDelete = Sinon.spy();
    this.onChange = Sinon.spy();

    this.renderComponent = () => {
      return render(hbs`
        <CodeGenerator::Policy::Stanza
          @index="0"
          @onChange={{this.onChange}}
          @onDelete={{this.onDelete}}
          @stanza={{this.stanza}}
        />`);
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();
    assert
      .dom(GENERAL.inputByAttr('path'))
      .hasValue('')
      .hasAttribute('placeholder', 'Enter a resource path')
      .hasAttribute('aria-label', 'Resource path')
      .hasAttribute('autocomplete', 'off');
    assert.dom(GENERAL.inputByAttr('path')).hasValue('');
    assert.dom(GENERAL.button('Delete')).exists({ count: 1 });
    // Assert checkboxes
    assert.dom('fieldset input[type="checkbox"]').exists({ count: 7 }, 'it renders 7 checkboxes');
    ACL_CAPABILITIES.forEach((capability) => {
      assert.dom(GENERAL.fieldLabel(capability)).hasText(capability);
      assert.dom(GENERAL.checkboxByAttr(capability)).isNotChecked();
    });
    // Assert preview toggle
    assert.dom(GENERAL.toggleInput('preview')).exists().isNotChecked();
    assert.dom(GENERAL.fieldLabel('preview')).hasText('Show preview');
    // Check empty preview state
    await click(GENERAL.toggleInput('preview'));
    assert.dom(GENERAL.toggleInput('preview')).isChecked();
    assert.dom(GENERAL.fieldLabel('preview')).hasText('Hide preview');
    const expectedPreview = `path "" {
    capabilities = []
  }`;
    assert.dom(GENERAL.fieldByAttr('preview')).hasText(expectedPreview);
  });

  test('it renders policy preview', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'some/api/path');
    await click(GENERAL.checkboxByAttr('update'));
    await click(GENERAL.checkboxByAttr('patch'));
    await click(GENERAL.toggleInput('preview'));
    let expectedPreview = `path "some/api/path" {
    capabilities = ["update", "patch"]
  }`;
    assert.dom(GENERAL.fieldByAttr('preview')).hasText(expectedPreview, 'it renders initial preview');
    // Toggle back to add more capabilities then check preview again
    await click(GENERAL.toggleInput('preview'));
    await typeIn(GENERAL.inputByAttr('path'), '/*');
    await click(GENERAL.checkboxByAttr('patch')); // uncheck
    await click(GENERAL.checkboxByAttr('list')); // check new
    // Confirm policy preview updated
    await click(GENERAL.toggleInput('preview'));
    expectedPreview = `path "some/api/path/*" {
    capabilities = ["update", "list"]
  }`;
    assert.dom(GENERAL.fieldByAttr('preview')).hasText(expectedPreview, 'it updates preview');
  });

  test('it maintains checkbox state when toggling to show and hide preview', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'some/api/path');
    await click(GENERAL.checkboxByAttr('update'));
    assert.dom(GENERAL.checkboxByAttr('update')).isChecked();
    // Toggle to show preview
    await click(GENERAL.toggleInput('preview'));
    assert.dom(GENERAL.toggleInput('preview')).isChecked();
    // Toggle back to checkboxes
    await click(GENERAL.toggleInput('preview'));
    assert.dom(GENERAL.toggleInput('preview')).isNotChecked();
    assert.dom(GENERAL.checkboxByAttr('update')).isChecked('update is still checked after viewing preview');
  });

  test('it selects and unselects capabilities', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.checkboxByAttr('update')); // first onChange call
    assert.dom(GENERAL.checkboxByAttr('update')).isChecked();
    let expectedSet = new Set(['update']);
    assert.deepEqual(
      this.stanza.capabilities,
      expectedSet,
      `has expected capabilities: ${[...expectedSet].join(', ')}`
    );
    // Check "delete"
    await click(GENERAL.checkboxByAttr('delete')); // second onChange call
    assert.dom(GENERAL.checkboxByAttr('delete')).isChecked();
    expectedSet = new Set(['update', 'delete']);
    assert.deepEqual(
      this.stanza.capabilities,
      expectedSet,
      `has expected capabilities: ${[...expectedSet].join(', ')}`
    );
    // Uncheck "delete"
    await click(GENERAL.checkboxByAttr('delete')); // third onChange call
    assert.dom(GENERAL.checkboxByAttr('delete')).isNotChecked();
    expectedSet = new Set(['update']);
    assert.deepEqual(
      this.stanza.capabilities,
      expectedSet,
      `has expected capabilities: ${[...expectedSet].join(', ')}`
    );
    assert.strictEqual(this.onChange.callCount, 3, 'onChange is called every time a capability is selected');
  });

  test('it selects all capabilities and updates @stanza', async function (assert) {
    await this.renderComponent();
    // check in random order to assert generator orders them
    for (const capability of ['list', 'read', 'sudo', 'create', 'delete', 'patch', 'update']) {
      await click(GENERAL.checkboxByAttr(capability));
    }
    await click(GENERAL.toggleInput('preview'));
    const expectedPreview = `path "" {
    capabilities = ["create", "read", "update", "delete", "list", "patch", "sudo"]
  }`;
    assert.dom(GENERAL.fieldByAttr('preview')).hasText(expectedPreview);
    assert.deepEqual(
      this.stanza.capabilities,
      new Set(['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo']),
      'stanza includes every capability, in order'
    );
  });

  test('it updates @stanza when path changes', async function (assert) {
    await this.renderComponent();
    await typeIn(GENERAL.inputByAttr('path'), 'my/super/secret/*');
    assert.strictEqual(this.stanza.path, 'my/super/secret/*', '"path" is updated');
  });

  test('it calls onChange when path changes', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'my/super/secret/*');
    assert.true(this.onChange.calledOnce, 'onChange is called');
  });

  test('it calls onChange when a checkbox is selected', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.checkboxByAttr('update'));
    assert.true(this.onChange.calledOnce, 'onChange is called');
  });

  test('it calls onDelete', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Delete'));
    assert.true(this.onDelete.calledOnce, 'onDelete is called');
  });
});
