/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import Sinon from 'sinon';

module('Integration | Component | code-generator/policy/flyout', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.assertSaveRequest = (assert, expectedPolicy, msg = 'policy content is correct') => {
      this.server.post('/sys/policies/acl/:name', (_, req) => {
        const { policy } = JSON.parse(req.requestBody);
        assert.true(true, 'it makes POST request to sys/policies/acl');
        assert.strictEqual(req.params.name, 'test-policy', 'policy name is correct');
        assert.strictEqual(policy, expectedPolicy, msg);
        return overrideResponse(204);
      });
    };
    this.renderComponent = async ({ open = true } = {}) => {
      await render(hbs`<CodeGenerator::Policy::Flyout />`);
      if (open) {
        await click(GENERAL.button('Generate policy'));
      }
    };
  });

  test('it renders button trigger and opens and closes the flyout', async function (assert) {
    await this.renderComponent({ open: false });
    assert.dom(GENERAL.button('Generate policy')).exists().hasText('Generate policy');
    assert.dom(GENERAL.flyout).doesNotExist();

    await click(GENERAL.button('Generate policy'));
    assert.dom(GENERAL.flyout).exists('flyout opens after clicking button');
    assert.dom(GENERAL.inputByAttr('name')).exists();
    assert.dom(GENERAL.fieldByAttr('visual editor')).exists();
    assert.dom(GENERAL.accordionButton('Automation snippets')).exists();
    assert.dom(GENERAL.submitButton).exists().hasText('Save');
    assert.dom(GENERAL.cancelButton).exists().hasText('Cancel');

    await click(GENERAL.cancelButton);
    assert.dom(GENERAL.flyout).doesNotExist('flyout closes after clicking cancel');
  });

  test('it preserves state when re-opened', async function (assert) {
    assert.expect(3);
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'secret/data/*');
    await click(GENERAL.checkboxByAttr('read'));
    await click(GENERAL.cancelButton);
    // Re-open flyout to confirm input values are preserved
    await click(GENERAL.button('Generate policy'));
    assert.dom(GENERAL.inputByAttr('name')).hasValue('test-policy');
    assert.dom(GENERAL.inputByAttr('path')).hasValue('secret/data/*');
    assert.dom(GENERAL.checkboxByAttr('read')).isChecked();
  });

  test('it updates automation snippets as policy changes', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'my-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'prod/app/*');
    await click(GENERAL.checkboxByAttr('update'));
    await click(GENERAL.accordionButton('Automation snippets'));
    const expectedTfvp = `resource "vault_policy" "<local identifier>" {
  name = "my-policy"

  policy = <<EOT
  path "prod/app/*" {
    capabilities = ["update"]
}
EOT
}`;
    assert.dom(GENERAL.fieldByAttr('terraform')).hasText(expectedTfvp);

    const expectedCli = `vault policy write my-policy - <<EOT
  path "prod/app/*" {
    capabilities = ["update"]
}
EOT`;
    assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli);
  });

  test('it converts policy name to lowercase', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'MyPolicy');
    assert.dom(GENERAL.inputByAttr('name')).hasValue('mypolicy', 'name is converted to lowercase');
  });

  test('it does not submit default stanza templates as policy payload', async function (assert) {
    assert.expect(3);
    const expectedPolicy = '';
    this.assertSaveRequest(assert, expectedPolicy, 'policy payload is empty when visual editor is untouched');
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await click(GENERAL.submitButton);
  });

  test('it saves a policy', async function (assert) {
    assert.expect(7);
    const flashSuccessSpy = Sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    const expectedPolicy = `path "secret/data/*" {\n    capabilities = ["read"]\n}`;
    this.assertSaveRequest(assert, expectedPolicy);
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'secret/data/*');
    await click(GENERAL.checkboxByAttr('read'));
    await click(GENERAL.submitButton);

    assert.true(flashSuccessSpy.calledOnce, 'flash success is called once');
    const [message, options] = flashSuccessSpy.lastCall.args;
    assert.strictEqual(message, 'ACL policy "test-policy" saved successfully.', 'flash message is correct');
    assert.propEqual(
      options,
      {
        link: {
          text: 'View policy',
          route: 'vault.cluster.policy.show',
          models: ['acl', 'test-policy'],
        },
      },
      'flash options include title and link to view policy'
    );
    assert.dom(GENERAL.flyout).doesNotExist('flyout closes after successful save');
  });

  test('it resets after saving a policy', async function (assert) {
    assert.expect(11);
    const expectedPolicy = `path "secret/data/*" {\n    capabilities = ["read"]\n}`;
    this.assertSaveRequest(assert, expectedPolicy);
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'secret/data/*');
    await click(GENERAL.checkboxByAttr('read'));
    await click(GENERAL.submitButton);
    // Re-open flyout to confirm it resets after saving
    await click(GENERAL.button('Generate policy'));
    assert.dom(GENERAL.inputByAttr('name')).hasValue('', 'name is cleared');
    assert.dom(GENERAL.inputByAttr('path')).hasValue('', 'path is cleared');
    assert.dom(GENERAL.checkboxByAttr('read')).isNotChecked('capabilities are unchecked');
    await click(GENERAL.accordionButton('Automation snippets'));
    const expectedTfvp = `resource "vault_policy" "<local identifier>" {
  name = "<policy name>"

  policy = <<EOT
  path "" {
    capabilities = []
}
EOT
}`;
    assert.dom(GENERAL.fieldByAttr('terraform')).hasText(expectedTfvp);

    const expectedCli = `vault policy write <policy name> - <<EOT
  path "" {
    capabilities = []
}
EOT`;
    assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli);
    // Fill in name and save again to make sure policyContent is reset
    this.assertSaveRequest(assert, '', 'policy content is empty after a successful save');
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await click(GENERAL.submitButton);
  });

  test('it displays error message when save fails', async function (assert) {
    this.server.post('/sys/policies/acl/:name', () => {
      return overrideResponse(400, { errors: ["'policy' parameter not supplied or empty"] });
    });
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'empty-policy');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).exists().hasText("Error 'policy' parameter not supplied or empty");
    assert.dom(GENERAL.flyout).exists('flyout remains open after error');
  });

  test('it handles multiple rules in the policy', async function (assert) {
    assert.expect(3);
    const expectedPolicy = `path "first/path" {\n    capabilities = ["read"]\n}\npath "second/path" {\n    capabilities = ["update"]\n}`;
    this.assertSaveRequest(assert, expectedPolicy);
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'first/path');
    await click(GENERAL.checkboxByAttr('read'));

    await click(GENERAL.button('Add rule'));
    await fillIn(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`, 'second/path');
    await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('update')}`);

    await click(GENERAL.submitButton);
  });

  test('it disables buttons while saving', async function (assert) {
    assert.expect(2);
    this.server.post('/sys/policies/acl/:name', () => {
      // Assert button states while the request is in-flight
      assert.dom(GENERAL.submitButton).isDisabled();
      assert.dom(GENERAL.cancelButton).isDisabled();
      return overrideResponse(204);
    });
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'secret/*');
    await click(GENERAL.checkboxByAttr('read'));
    await click(GENERAL.submitButton);
  });

  test('it renders validation errors', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).exists().hasText('Error There is an error with this form.');
    assert.dom(GENERAL.inputByAttr('name')).hasClass('hds-form-text-input--is-invalid');
    assert.dom(GENERAL.validationErrorByAttr('name')).hasText('Name is required.');
  });

  test('it resets errors after saving', async function (assert) {
    const expectedPolicy = `path "secret/*" {\n    capabilities = ["read"]\n}`;
    this.assertSaveRequest(assert, expectedPolicy);
    await this.renderComponent();

    // First attempt without name
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).exists().hasText('Error There is an error with this form.');
    assert.dom(GENERAL.validationErrorByAttr('name')).exists('validation error shows');

    // Second attempt with name
    await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
    await fillIn(GENERAL.inputByAttr('path'), 'secret/*');
    await click(GENERAL.checkboxByAttr('read'));
    await click(GENERAL.submitButton);

    // Reopen flyout to check error state has reset
    await click(GENERAL.button('Generate policy'));
    assert.dom(GENERAL.messageError).doesNotExist('error banner is cleared');
    assert.dom(GENERAL.validationErrorByAttr('name')).doesNotExist('validation error is cleared');
  });

  test('it resets errors if flyout is closed and policy is NOT saved', async function (assert) {
    await this.renderComponent();
    // Attempt to save
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).exists().hasText('Error There is an error with this form.');
    assert.dom(GENERAL.validationErrorByAttr('name')).exists('validation error shows');
    // Cancel and close flyout
    await click(GENERAL.cancelButton);
    // Reopen flyout to check error state has reset
    await click(GENERAL.button('Generate policy'));
    assert.dom(GENERAL.messageError).doesNotExist('error banner is cleared');
    assert.dom(GENERAL.validationErrorByAttr('name')).doesNotExist('validation error is cleared');
  });
});
