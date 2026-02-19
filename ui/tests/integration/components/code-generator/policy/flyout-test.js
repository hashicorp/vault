/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, typeIn, waitUntil, find } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import Sinon from 'sinon';

const SELECTORS = {
  pathByContainer: (idx) => `${GENERAL.cardContainer(idx)} ${GENERAL.inputByAttr('path')}`,
  checkboxByContainer: (idx, cap) => `${GENERAL.cardContainer(idx)} ${GENERAL.checkboxByAttr(cap)}`,
};

module('Integration | Component | code-generator/policy/flyout', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise'; // the flyout is only available for enterprise versions
    this.onClose = undefined;
    this.policyPaths = undefined;
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
      await render(
        hbs`<CodeGenerator::Policy::Flyout @onClose={{this.onClose}} @policyPaths={{this.policyPaths}} />`
      );
      if (open) {
        await click(GENERAL.button('Generate policy'));
      }
    };
  });

  test('it calls onClose callback', async function (assert) {
    this.onClose = Sinon.spy();
    await this.renderComponent();
    await click(GENERAL.cancelButton);
    assert.true(this.onClose.calledOnce, 'onClose callback is called');
  });

  test('it does not render for community versions', async function (assert) {
    this.version.type = 'community';
    await this.renderComponent({ open: false });
    assert.dom(GENERAL.button('Generate policy')).doesNotExist('Button does not render for CE version');
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

  test('it presets with paths from @policyPaths array', async function (assert) {
    this.policyPaths = ['some/preset/path'];
    await this.renderComponent();
    assert.dom(SELECTORS.pathByContainer(0)).hasValue('some/preset/path');
    assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
  });

  test('it handles empty @policyPaths array', async function (assert) {
    this.policyPaths = [];
    await this.renderComponent();

    assert.dom(SELECTORS.pathByContainer(0)).hasValue('', 'does not prepopulate with empty array');
  });

  test('it yields custom trigger component', async function (assert) {
    await render(hbs`<Hds::Dropdown as |D|>
      <D.ToggleButton @text="Toolbox" data-test-dropdown="Toolbox" />
      <CodeGenerator::Policy::Flyout>
        <:customTrigger as |openFlyout|>
          <D.Interactive @icon="shield-check" {{on "click" openFlyout}} data-test-button="Make me a policy!">
            Make me a policy!
          </D.Interactive>
        </:customTrigger>
      </CodeGenerator::Policy::Flyout>
      <D.Interactive @icon="wand" data-test-button="Magic stuff">Magic stuff</D.Interactive>
    </Hds::Dropdown>`);
    await click(GENERAL.dropdownToggle('Toolbox'));
    assert.dom(GENERAL.flyout).doesNotExist();
    assert
      .dom(GENERAL.button('Make me a policy!'))
      .exists()
      .hasText('Make me a policy!', 'custom trigger renders');
    await click(GENERAL.button('Make me a policy!'));
    assert.dom(GENERAL.flyout).exists('flyout opens after clicking custom trigger');
  });

  // This test is to demonstrate how to implement closing the dropdown when the flyout trigger is a dropdown element
  test('it closes dropdown if custom trigger is a dropdown item', async function (assert) {
    await render(hbs`<Hds::Dropdown as |D|>
      <D.ToggleButton @text="Toolbox" data-test-dropdown="Toolbox" />
      <CodeGenerator::Policy::Flyout @onClose={{D.close}} >
        <:customTrigger as |openFlyout|>
          <D.Interactive @icon="shield-check" {{on "click" openFlyout}} data-test-button="Make me a policy!">
            Make me a policy!
          </D.Interactive>
        </:customTrigger>
      </CodeGenerator::Policy::Flyout>
      <D.Interactive @icon="wand" data-test-button="Magic stuff">Magic stuff</D.Interactive>
    </Hds::Dropdown>`);
    await click(GENERAL.dropdownToggle('Toolbox'));
    assert.dom(GENERAL.dropdownToggle('Toolbox')).hasAttribute('aria-expanded', 'true');
    await click(GENERAL.button('Make me a policy!'));
    assert.dom(GENERAL.flyout).exists('flyout is open');
    await click(GENERAL.cancelButton);
    assert.dom(GENERAL.flyout).doesNotExist('flyout is closed');
    const dropdown = find(GENERAL.dropdownToggle('Toolbox'));
    await waitUntil(() => dropdown.ariaExpanded === 'false');
    assert
      .dom(GENERAL.dropdownToggle('Toolbox'))
      .hasAttribute('aria-expanded', 'false', 'dropdown closes when flyout is closed');
  });

  test('it does not render yielded custom trigger component on community', async function (assert) {
    this.version.type = 'community';
    await this.renderComponent({ open: false });
    await render(hbs`<Hds::Dropdown as |D|>
      <D.ToggleButton @text="Toolbox" data-test-dropdown="Toolbox" />
      <CodeGenerator::Policy::Flyout>
        <:customTrigger as |openFlyout|>
          <D.Interactive @icon="shield-check" {{on "click" openFlyout}} data-test-button="Make me a policy!">
            Make me a policy!
          </D.Interactive>
        </:customTrigger>
      </CodeGenerator::Policy::Flyout>
      <D.Interactive @icon="wand" data-test-button="Magic stuff">Magic stuff</D.Interactive>
    </Hds::Dropdown>`);
    await click(GENERAL.dropdownToggle('Toolbox'));
    assert.dom(GENERAL.button('Magic stuff')).exists('dropdown opens');
    assert.dom(GENERAL.button('Make me a policy!')).doesNotExist();
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
    await fillIn(SELECTORS.pathByContainer(1), 'second/path');
    await click(SELECTORS.checkboxByContainer(1, 'update'));

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

  module('capabilities service prepopulating', function (hooks) {
    hooks.beforeEach(function () {
      this.capabilities = this.owner.lookup('service:capabilities');
      const router = this.owner.lookup('service:router');
      this.currentRouteNameStub = Sinon.stub(router, 'currentRouteName');
      this.cacheCapabilityPaths = (route, paths) => {
        this.capabilities.cacheRoutePaths(route, paths);
      };
    });

    hooks.afterEach(function () {
      this.currentRouteNameStub.restore();
    });

    test('it handles null currentRouteName gracefully', async function (assert) {
      this.currentRouteNameStub.value(null);
      await this.renderComponent();
      assert.dom(SELECTORS.pathByContainer(0)).hasValue('', 'does not prepopulate when route name is null');
      assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
    });

    test('it does not prepopulate when no paths have been cached', async function (assert) {
      await this.renderComponent();
      assert.dom(SELECTORS.pathByContainer(0)).hasValue('');
    });

    test('it does not prepopulate paths when cached capabilities route is unrelated to the current route', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.secret');
      this.cacheCapabilityPaths('vault.cluster.settings', ['some/settings']);
      await this.renderComponent();
      assert.dom(SELECTORS.pathByContainer(0)).hasValue('');
    });

    test('it prepopulates paths when cached capabilities route equals current route', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.secret');
      this.cacheCapabilityPaths('vault.cluster.secrets.secret', ['super-secret/data']);
      await this.renderComponent();
      assert.dom(SELECTORS.pathByContainer(0)).hasValue('super-secret/data');
      assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
    });

    test('it prepopulates paths from longest matching parent route', async function (assert) {
      // Cache paths for parent route
      this.cacheCapabilityPaths('vault.cluster.secrets.backend.kv.secret', [
        'kv/data/my-secret',
        'kv/metadata/my-secret',
      ]);
      this.cacheCapabilityPaths('vault.cluster.secrets.backend.kv', ['should/not/cache']);
      // Current route is a child (e.g., secret.details)
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.kv.secret.details');
      await this.renderComponent();
      assert.dom(SELECTORS.pathByContainer(0)).hasValue('kv/data/my-secret', 'uses parent paths');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('kv/metadata/my-secret', 'includes all parent paths');
      assert.dom(GENERAL.cardContainer()).exists({ count: 2 });
    });

    // All of these tests run with the current route stubbed and cached paths
    module('when the flyout is prepopulated', function (hooks) {
      hooks.beforeEach(function () {
        this.cacheCapabilityPaths('vault.cluster.secrets.secret', ['super-secret/data']);
        this.currentRouteNameStub.value('vault.cluster.secrets.secret');
      });

      test('paths from arg take precedence over capabilities service', async function (assert) {
        this.policyPaths = ['super-explicit/path'];
        await this.renderComponent();
        assert.dom(SELECTORS.pathByContainer(0)).hasValue('super-explicit/path');
        assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
      });

      test('it prepopulates with a single capability path', async function (assert) {
        await this.renderComponent();
        assert.dom(SELECTORS.pathByContainer(0)).hasValue('super-secret/data');
        assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
      });

      test('it prepopulates with multiple capability paths', async function (assert) {
        this.cacheCapabilityPaths('vault.cluster.secrets.secret', ['path/one', 'path/two']);
        await this.renderComponent();
        assert.dom(SELECTORS.pathByContainer(0)).hasValue('path/one');
        assert.dom(SELECTORS.pathByContainer(1)).hasValue('path/two');
        assert.dom(GENERAL.cardContainer()).exists({ count: 2 });
      });

      test('it does not override user changes to a preset path on reopen', async function (assert) {
        await this.renderComponent();
        // User updates path
        await typeIn(SELECTORS.pathByContainer(0), '/*');
        // Close and reopen
        await click(GENERAL.cancelButton);
        await click(GENERAL.button('Generate policy'));
        assert
          .dom(SELECTORS.pathByContainer(0))
          .hasValue('super-secret/data/*', 'user path changes are preserved');
        assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
      });

      test('it does not override user capabilities selection for a preset path on reopen', async function (assert) {
        await this.renderComponent();

        // User updates path
        await click(SELECTORS.checkboxByContainer(0, 'read'));
        // Close and reopen
        await click(GENERAL.cancelButton);
        await click(GENERAL.button('Generate policy'));
        assert
          .dom(SELECTORS.checkboxByContainer(0, 'read'))
          .isChecked('user capabilities changes are preserved');
        assert.dom(GENERAL.cardContainer()).exists({ count: 1 });
      });

      test('it does not override user added stanza on reopen', async function (assert) {
        await this.renderComponent();
        await click(GENERAL.button('Add rule'));
        await fillIn(SELECTORS.pathByContainer(1), 'new/path/*');
        // Close and reopen
        await click(GENERAL.cancelButton);
        await click(GENERAL.button('Generate policy'));
        assert.dom(GENERAL.cardContainer()).exists({ count: 2 }, 'it renders two stanzas after reopening');
        assert.dom(SELECTORS.pathByContainer(0)).hasValue('super-secret/data', 'preset path still exists');
        assert.dom(SELECTORS.pathByContainer(1)).hasValue('new/path/*', 'user added path still exists');
      });

      test('it does not save prepopulated paths as policy content', async function (assert) {
        assert.expect(3);
        this.cacheCapabilityPaths('vault.cluster.secrets.secret', ['path/one', 'path/two']);
        await this.renderComponent();
        // Fill in name and save to make sure policyContent is empty
        this.assertSaveRequest(assert, '', 'policy content is empty despite pre-filled paths');
        await fillIn(GENERAL.inputByAttr('name'), 'test-policy');
        await click(GENERAL.submitButton);
      });
    });
  });
});
