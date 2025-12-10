/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { ACL_CAPABILITIES } from 'core/utils/code-generators/policy';

module('Integration | Component | code-generator/policy/builder', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onPolicyChange = Sinon.spy();
    this.policyName = undefined;
    this.renderComponent = () => {
      return render(hbs`
        <CodeGenerator::Policy::Builder @onPolicyChange={{this.onPolicyChange}} @policyName={{this.policyName}} />`);
    };

    this.assertPolicyUpdate = (assert, expected, message) => {
      const [policy] = this.onPolicyChange.lastCall.args;
      assert.strictEqual(policy, expected, `onPolicyChange is called ${message}`);
    };

    this.assertEmptyTemplate = async (assert, { index } = {}) => {
      const container = index ? GENERAL.cardContainer(index) : '';
      assert.dom(`${container} ${GENERAL.inputByAttr('path')}`).hasValue('');
      assert.dom(`${container} ${GENERAL.toggleInput('preview')}`).isNotChecked();
      ACL_CAPABILITIES.forEach((capability) => {
        assert.dom(`${container} ${GENERAL.checkboxByAttr(capability)}`).isNotChecked();
      });
      // check empty preview state
      await click(`${container} ${GENERAL.toggleInput('preview')}`);
      const expectedPreview = `path "" {
    capabilities = []
  }`;
      assert
        .dom(GENERAL.fieldByAttr('preview'))
        .exists('it renders preview')
        .hasText(expectedPreview, 'preview is empty');
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();
    await this.assertEmptyTemplate(assert);
    assert.dom(GENERAL.button('Add rule')).exists({ count: 1 });
    assert.dom(GENERAL.revealButton('Automation snippets')).hasAttribute('aria-expanded', 'false');
    await click(GENERAL.revealButton('Automation snippets'));
    assert.dom(GENERAL.revealButton('Automation snippets')).hasAttribute('aria-expanded', 'true');
    assert.dom(GENERAL.inputByAttr('terraform')).isChecked();
    assert.dom(GENERAL.inputByAttr('cli')).isNotChecked();
  });

  test('it renders default snippets', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.revealButton('Automation snippets'));
    let expectedSnippet = `resource "vault_policy" "<local identifier>" {
  name = "<policy name>"

  policy = <<EOT
  path "" {
    capabilities = []
}
EOT
}`;
    assert
      .dom(GENERAL.fieldByAttr('snippets'))
      .hasText(expectedSnippet, 'it renders empty terraform snippet');

    expectedSnippet = `vault policy write <policy name> - <<EOT
  path "" {
    capabilities = []
}
EOT`;
    await click(GENERAL.inputByAttr('cli'));
    assert.dom(GENERAL.inputByAttr('cli')).isChecked();
    assert.dom(GENERAL.inputByAttr('terraform')).isNotChecked();
    assert.dom(GENERAL.fieldByAttr('snippets')).hasText(expectedSnippet, 'it renders empty cli snippet');
  });

  test('it includes namespace in snippet for non-root namespaces', async function (assert) {
    const namespace = this.owner.lookup('service:namespace');
    namespace.path = 'admin';
    await this.renderComponent();
    await click(GENERAL.revealButton('Automation snippets'));
    const expectedSnippet = `resource "vault_policy" "<local identifier>" {
  namespace = "admin"

  name = "<policy name>"
  
  policy = <<EOT
  path "" {
    capabilities = []
}
EOT
}`;
    assert
      .dom(GENERAL.fieldByAttr('snippets'))
      .hasText(expectedSnippet, 'it renders empty terraform snippet');
  });

  test('it does not call onPolicyChange when callback is not provided', async function (assert) {
    this.onPolicyChange = undefined;
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'test/path');
    await click(GENERAL.checkboxByAttr('read'));
    await click(GENERAL.button('Add rule'));
    assert.true(true, 'no errors are thrown when callback is undefined');
  });

  test('it adds a rule', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Add rule'));
    assert
      .dom(GENERAL.cardContainer())
      .exists({ count: 2 }, 'two templates render after clicking "Add rule"');
    await this.assertEmptyTemplate(assert, { index: '1' });
  });

  test('it deletes a rule', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
    // Fill in template
    await fillIn(GENERAL.inputByAttr('path'), 'some/api/path');
    await click(GENERAL.checkboxByAttr('patch'));
    // Delete the only rendered template
    await click(GENERAL.button('Delete'));
    // One template renders but content should reset
    assert
      .dom(GENERAL.cardContainer())
      .exists({ count: 1 }, 'it still renders one rule after deleting the only rule');
    await this.assertEmptyTemplate(assert);
  });

  test('it maintains state across multiple rules', async function (assert) {
    await this.renderComponent();
    // Set up first rule
    await fillIn(GENERAL.inputByAttr('path'), 'first/path');
    await click(GENERAL.checkboxByAttr('read'));
    // Add second rule
    await click(GENERAL.button('Add rule'));
    await fillIn(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`, 'second/path');
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('update')}`);

    assert.dom(`${GENERAL.cardContainer('0')} ${GENERAL.inputByAttr('path')}`).hasValue('first/path');
    assert.dom(`${GENERAL.cardContainer('0')} ${GENERAL.checkboxByAttr('read')}`).isChecked();
    assert.dom(`${GENERAL.cardContainer('0')} ${GENERAL.checkboxByAttr('update')}`).isNotChecked();
    assert.dom(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`).hasValue('second/path');
    assert.dom(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('update')}`).isChecked();
    assert.dom(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('read')}`).isNotChecked();
  });

  test('it deletes the correct rule when multiple exist', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'first/path');
    await click(GENERAL.checkboxByAttr('read'));
    // Second rule
    await click(GENERAL.button('Add rule'));
    await fillIn(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`, 'second/path');
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('update')}`);
    // Third rule
    await click(GENERAL.button('Add rule'));
    await fillIn(`${GENERAL.cardContainer('2')} ${GENERAL.inputByAttr('path')}`, 'third/path');
    await click(`${GENERAL.cardContainer('2')} ${GENERAL.checkboxByAttr('list')}`);
    assert.dom(GENERAL.cardContainer()).exists({ count: 3 });
    // Delete middle rule
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.button('Delete')}`);
    assert.dom(GENERAL.cardContainer()).exists({ count: 2 });
    assert.dom(`${GENERAL.cardContainer('0')} ${GENERAL.inputByAttr('path')}`).hasValue('first/path');
    assert.dom(`${GENERAL.cardContainer('0')} ${GENERAL.checkboxByAttr('read')}`).isChecked();
    assert.dom(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`).hasValue('third/path');
    assert.dom(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('list')}`).isChecked();
  });

  test('it updates snippets', async function (assert) {
    this.policyName = 'my-secure-policy';
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'my/super/secret/*');
    await click(GENERAL.checkboxByAttr('patch'));
    await click(GENERAL.revealButton('Automation snippets'));
    // Check terraform snippet
    let expectedSnippet = `resource "vault_policy" "<local identifier>" {
  name = "my-secure-policy"

  policy = <<EOT
  path "my/super/secret/*" {
    capabilities = ["patch"]
}
EOT
}`;
    assert.dom(GENERAL.fieldByAttr('snippets')).hasText(expectedSnippet, 'it renders terraform snippet');

    // Check CLI snippet
    expectedSnippet = `vault policy write my-secure-policy - <<EOT
  path "my/super/secret/*" {
    capabilities = ["patch"]
}
EOT`;
    await click(GENERAL.inputByAttr('cli'));
    assert.dom(GENERAL.fieldByAttr('snippets')).hasText(expectedSnippet, 'it renders cli snippet');
  });

  test('it passes policy updates as changes are made', async function (assert) {
    await this.renderComponent();
    // Inputting path triggers callback
    await fillIn(GENERAL.inputByAttr('path'), 'my/super/secret/*');
    let expectedPolicy = `path "my/super/secret/*" {
    capabilities = []
}`;

    this.assertPolicyUpdate(assert, expectedPolicy, 'when path changes');

    // Clicking checkbox triggers callback
    await click(GENERAL.checkboxByAttr('update'));
    expectedPolicy = `path "my/super/secret/*" {
    capabilities = ["update"]
}`;

    this.assertPolicyUpdate(assert, expectedPolicy, 'when a capability is selected');

    // Adding a rule triggers callback
    await click(GENERAL.button('Add rule'));
    expectedPolicy = `path "my/super/secret/*" {
    capabilities = ["update"]
}
path "" {
    capabilities = []
}`;
    this.assertPolicyUpdate(assert, expectedPolicy, 'when a rule is added');

    // Updating added rule triggers callback
    await fillIn(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`, 'prod/');
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('list')}`);
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('read')}`);
    expectedPolicy = `path "my/super/secret/*" {
    capabilities = ["update"]
}
path "prod/" {
    capabilities = ["read", "list"]
}`;
    this.assertPolicyUpdate(assert, expectedPolicy, 'when an additional rule updates');

    // Unchecking box triggers callback
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('read')}`);
    expectedPolicy = `path "my/super/secret/*" {
    capabilities = ["update"]
}
path "prod/" {
    capabilities = ["list"]
}`;
    this.assertPolicyUpdate(assert, expectedPolicy, 'when checkbox is unselected');

    // Deleting a rule triggers callback
    await click(GENERAL.button('Delete'));
    expectedPolicy = `path "prod/" {
    capabilities = ["list"]
}`;
    this.assertPolicyUpdate(assert, expectedPolicy, 'when a rule is deleted');
  });

  // These tests ensure paths are never used as input identifiers.
  // The policy generator may not render in a form and needs to be flexible so it intentionally supports
  // multiple templates with the same or no path.
  test('it supports multiple rules with the same path', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('path'), 'test/path');
    await click(GENERAL.checkboxByAttr('patch'));
    await click(GENERAL.button('Add rule'));
    await fillIn(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`, 'test/path');
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('update')}`);

    const expectedPolicy = `path "test/path" {
    capabilities = ["patch"]
}
path "test/path" {
    capabilities = ["update"]
}`;
    this.assertPolicyUpdate(assert, expectedPolicy, 'when rules have the same path');
  });

  test('it supports multiple rules with an empty path', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.checkboxByAttr('list'));
    await click(GENERAL.button('Add rule'));
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('delete')}`);

    const expectedPolicy = `path "" {
    capabilities = ["list"]
}
path "" {
    capabilities = ["delete"]
}`;
    this.assertPolicyUpdate(assert, expectedPolicy, 'when rules do have an empty path');
  });
});
