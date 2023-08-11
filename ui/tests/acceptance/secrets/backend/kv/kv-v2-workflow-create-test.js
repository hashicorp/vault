import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { personas } from 'vault/tests/helpers/policy-generator/kv';
import { setupControlGroup, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { click, currentURL, fillIn, typeIn, visit } from '@ember/test-helpers';

/**
 * This test set is for testing the flow for creating new secrets and versions
 */
module('Acceptance | kv-v2 workflow | secret and version create', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `kv-create-${uuidv4()}`;
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeVersionedSecret(this.backend, 'app/first', 'foo', 'bar', 2);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(tokenWithPolicyCmd('admin', personas.admin(this.backend)));
      await authPage.login(token);
    });
    test('create & update root secret with default metadata', async function (assert) {
      const backend = this.backend;
      const secretPath = 'some secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(PAGE.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      // TODO: clear invalid form alert on form change?
      assert
        .dom(FORM.validationWarning)
        .hasText(
          "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests."
        );
      assert.dom(PAGE.create.metadataSection).doesNotExist('Hides metadata section by default');

      // Submit with API errors
      await click(FORM.saveBtn);
      assert.dom(FORM.messageError).hasText('Error no data provided', 'API error shows on form');

      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);

      // Details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details?version=1`,
        'Goes to details page after save'
      );
      assert.dom(PAGE.detail.versionCreated).includesText('Version 1 created');
      assert.dom(PAGE.secretRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(PAGE.infoRowValue('api_key')).hasText('***********');
      await click(PAGE.infoRowToggleMasked('api_key'));
      assert.dom(PAGE.infoRowValue('api_key')).hasText('partyparty', 'secret value shows after toggle');

      // Metadata page
      await click(PAGE.secretTab('Metadata'));
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('No custom metadata', 'No custom metadata empty state');
      assert
        .dom(`${PAGE.metadata.secretMetadataSection} ${PAGE.secretRow}`)
        .exists({ count: 3 }, '3 metadata rows show');
      // await this.pauseTest();
      assert.dom(PAGE.infoRowValue('Maximum versions')).hasText('0', 'max versions shows 0');
      assert.dom(PAGE.infoRowValue('Check-and-Set required')).hasText('No', 'cas not enforced');
      assert
        .dom(PAGE.infoRowValue('Delete version after'))
        .hasText('Never delete', 'Delete version after has default 0s');

      // Add new version
      await click(PAGE.secretTab('Secret'));
      await click(PAGE.detail.createNewVersion);
      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue(secretPath);
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('api_key');
      assert.dom(FORM.maskedValueInput()).hasValue('partyparty');
      await fillIn(FORM.keyInput(1), 'api_url');
      await fillIn(FORM.maskedValueInput(1), 'hashicorp.com');
      await click(FORM.saveBtn);

      // Back to details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details?version=2`
      );
      assert.dom(PAGE.detail.versionCreated).includesText('Version 2 created');
      assert.dom(PAGE.secretRow).exists({ count: 2 }, '2 rows of data shows');
      assert.dom(PAGE.infoRowValue('api_key')).hasText('***********');
      assert.dom(PAGE.infoRowValue('api_url')).hasText('***********');
      await click(PAGE.infoRowToggleMasked('api_key'));
      await click(PAGE.infoRowToggleMasked('api_url'));
      assert.dom(PAGE.infoRowValue('api_key')).hasText('partyparty', 'secret value shows after toggle');
      assert.dom(PAGE.infoRowValue('api_url')).hasText('hashicorp.com', 'secret value shows after toggle');
    });
    test('create nested secret with metadata', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(PAGE.list.createSecret);

      // Create secret
      await typeIn(FORM.inputByAttr('path'), 'my/');
      assert.dom(FORM.validation('path')).hasText("Path can't end in forward slash '/'.");
      await typeIn(FORM.inputByAttr('path'), 'secret');
      assert.dom(FORM.validation('path')).doesNotExist('form validation goes away');
      await fillIn(FORM.keyInput(), 'password');
      await fillIn(FORM.maskedValueInput(), 'kittens1234');

      await click(FORM.toggleMetadata);
      assert.dom(PAGE.create.metadataSection).exists('Shows metadata section after toggled');
      // Check initial values
      assert.dom(FORM.inputByAttr('maxVersions')).hasValue('0');
      assert.dom(FORM.inputByAttr('casRequired')).isNotChecked();
      assert.dom(FORM.toggleByLabel('Automate secret deletion')).isNotChecked();
      // MaxVersions validation
      await fillIn(FORM.inputByAttr('maxVersions'), 'seven');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('Maximum versions must be a number.');
      await fillIn(FORM.inputByAttr('maxVersions'), '99999999999999999');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('You cannot go over 16 characters.');
      await fillIn(FORM.inputByAttr('maxVersions'), '7');

      // Fill in other metadata
      await click(FORM.inputByAttr('casRequired'));
      await click(FORM.toggleByLabel('Automate secret deletion'));
      await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');

      // Fill in custom metadata
      await fillIn(`${PAGE.create.metadataSection} ${FORM.keyInput()}`, 'team');
      await fillIn(`${PAGE.create.metadataSection} ${FORM.valueInput()}`, 'UI');
      // Fill in metadata
      await click(FORM.saveBtn);

      // Details
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('my/secret')}/details?version=1`
      );
      assert.dom(PAGE.detail.versionCreated).includesText('Version 1 created');
      assert.dom(PAGE.secretRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(PAGE.infoRowValue('password')).hasText('***********');
      await click(PAGE.infoRowToggleMasked('password'));
      assert.dom(PAGE.infoRowValue('password')).hasText('kittens1234', 'secret value shows after toggle');

      // Metadata
      await click(PAGE.secretTab('Metadata'));
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.secretRow}`)
        .exists({ count: 1 }, 'One custom metadata row shows');
      assert.dom(`${PAGE.metadata.customMetadataSection} ${PAGE.infoRowValue('team')}`).hasText('UI');

      assert
        .dom(`${PAGE.metadata.secretMetadataSection} ${PAGE.secretRow}`)
        .exists({ count: 3 }, '3 metadata rows show');
      // await this.pauseTest();
      assert.dom(PAGE.infoRowValue('Maximum versions')).hasText('7', 'max versions shows 0');
      assert.dom(PAGE.infoRowValue('Check-and-Set required')).hasText('Yes', 'cas not enforced');
      assert
        .dom(PAGE.infoRowValue('Delete version after'))
        .hasText('16 minutes 40 seconds', 'Delete version after has custom value');
    });
    test('creates a secret at a sub-directory', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2F/directory`);
      assert.dom(PAGE.list.item('first')).exists('Lists first sub-secret');
      assert.dom(PAGE.list.item('new')).doesNotExist('Does not show new secret');
      await click(PAGE.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=app%2F`,
        'Goes to create page with initialKey'
      );
      await typeIn(FORM.inputByAttr('path'), 'new');
      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('app/new')}/details?version=1`,
        'Redirects to detail after save'
      );
      await click(PAGE.breadcrumbAtIdx(2));
      // TODO: Pagefilter should not be present if empty
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2F/directory?pageFilter=`,
        'sub-dir page'
      );
      assert.dom(PAGE.list.item('new')).exists('Lists new secret in sub-dir');
    });
    test('create new version of secret from older version', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details`);
      await click(PAGE.detail.versionDropdown);
      await click(`${PAGE.detail.version(1)} a`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details?version=1`,
        'goes to version 1'
      );
      assert.dom(PAGE.detail.versionCreated).includesText('Version 1 created');
      await click(PAGE.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details/edit?version=1`,
        'Goes to new version page'
      );
      assert
        .dom(FORM.versionAlert)
        .hasText(
          'Warning You are creating a new version based on data from Version 1. The current version for app/first is Version 2.',
          'Shows version warning'
        );
      assert.dom(FORM.keyInput()).hasValue('key-1', 'Key input has old value');
      assert.dom(FORM.maskedValueInput()).hasValue('val-1', 'Val input has old value');

      await fillIn(FORM.keyInput(), 'my-key');
      await fillIn(FORM.maskedValueInput(), 'my-value');
      await click(FORM.saveBtn);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details?version=3`,
        'goes to latest version 3'
      );
      await click(PAGE.infoRowToggleMasked('my-key'));
      assert.dom(PAGE.infoRowValue('my-key')).hasValue('my-value', 'has new value');
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([tokenWithPolicyCmd('data-reader', personas.dataReader(this.backend))]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        tokenWithPolicyCmd('data-list-reader', personas.dataListReader(this.backend)),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        tokenWithPolicyCmd('metadata-maintainer', personas.metadataMaintainer(this.backend)),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        tokenWithPolicyCmd('secret-creator', personas.secretCreator(this.backend)),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      const userPolicy = `
path "${this.backend}/data/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
  control_group = {
    max_ttl = "24h"
    factor "approver" {
      controlled_capabilities = ["write"]
      identity {
          group_names = ["managers"]
          approvals = 1
      }
    }
  }
}

path "${this.backend}/*" {
  capabilities = ["list"]
}

// Can we allow this so user can self-authorize?
path "sys/control-group/authorize" {
  capabilities = ["update"]
}

path "sys/control-group/request" {
  capabilities = ["update"]
}
`;
      const { userToken } = await setupControlGroup({ userPolicy });
      this.userToken = userToken;
      return authPage.login(userToken);
    });
    // Copy test outline from admin persona
  });
});
