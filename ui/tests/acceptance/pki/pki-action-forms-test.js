/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, skip, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, fillIn, typeIn, visit } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS as S } from 'vault/tests/helpers/pki/workflow';
import { issuerPemBundle } from 'vault/tests/helpers/pki/values';

module('Acceptance | pki action forms test', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await logout.visit();
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
    await logout.visit();
  });

  module('import', function (hooks) {
    setupMirage(hooks);

    hooks.beforeEach(function () {
      this.pemBundle = issuerPemBundle;
    });

    test('happy path', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(S.emptyStateLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      assert.dom(S.configuration.title).hasText('Configure PKI');
      assert.dom(S.configuration.emptyState).exists({ count: 1 }, 'Shows empty state by default');
      await click(S.configuration.optionByKey('import'));
      assert.dom(S.configuration.emptyState).doesNotExist();
      // Submit before filling out form shows an error
      await click('[data-test-pki-import-pem-bundle]');
      assert.dom('[data-test-alert-banner="alert"]').hasText('Error please upload your PEM bundle');
      // Fill in form data
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click('[data-test-pki-import-pem-bundle]');

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(S.configuration.title).hasText('View imported items');
      assert.dom(S.configuration.importForm).doesNotExist('import form is hidden after save');
      assert.dom(S.configuration.importMapping).exists('import mapping is shown after save');
      await click('[data-test-done]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'redirects to overview when done'
      );
    });
    skip('with many imports', async function (assert) {
      // TODO VAULT-14791
      this.server.post(`${this.mountPath}/config/ca`, () => {
        return {
          request_id: 'some-config-id',
          data: {
            imported_issuers: ['my-imported-issuer', 'imported2'],
            imported_keys: ['my-imported-key', 'imported3'],
            mapping: {
              'my-imported-issuer': 'my-imported-key',
              imported2: '',
            },
          },
        };
      });
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(S.configuration.optionByKey('import'));
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click('[data-test-pki-import-pem-bundle]');

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(S.configuration.title).hasText('View imported items');
      assert.dom(S.configuration.importForm).doesNotExist('import form is hidden after save');
      assert.dom(S.configuration.importMapping).exists('import mapping is shown after save');
      assert.dom(S.configuration.importedIssuer).hasText('my-imported-issuer', 'Issuer value is displayed');
      assert.dom(S.configuration.importedKey).hasText('my-imported-key', 'Key value is displayed');
      await click('[data-test-done]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'redirects to overview when done'
      );
    });
    skip('shows imported items when keys is empty', async function (assert) {
      // TODO VAULT-14791
      this.server.post(`${this.mountPath}/config/ca`, () => {
        return {
          request_id: 'some-config-id',
          data: {
            imported_issuers: ['my-imported-issuer', 'my-imported-issuer2'],
            imported_keys: null,
            mapping: {
              'my-imported-issuer': '',
              'my-imported-issuer2': '',
            },
          },
        };
      });
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(S.configuration.optionByKey('import'));
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click('[data-test-pki-import-pem-bundle]');

      assert.dom(S.configuration.importForm).doesNotExist('import form is hidden after save');
      assert.dom(S.configuration.importMapping).exists('import mapping is shown after save');
      assert.dom(S.configuration.importedIssuer).hasText('my-imported-issuer', 'Issuer value is displayed');
      assert.dom(S.configuration.importedKey).hasText('my-imported-key', 'Key value is displayed');
    });
  });

  module('generate root', function () {
    test('happy path', async function (assert) {
      const commonName = 'my-common-name';
      const issuerName = 'my-first-issuer';
      const keyName = 'my-first-key';
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(S.emptyStateLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      assert.dom(S.configuration.title).hasText('Configure PKI');
      assert.dom(S.configuration.emptyState).exists({ count: 1 }, 'Shows empty state by default');
      await click(S.configuration.optionByKey('generate-root'));
      assert.dom(S.configuration.emptyState).doesNotExist();
      // The URLs section is populated based on params returned from OpenAPI. This test will break when
      // the backend adds fields. We should update the count accordingly.
      assert.dom(S.configuration.urlField).exists({ count: 4 });
      // Fill in form
      await fillIn(S.configuration.typeField, 'internal');
      await typeIn(S.configuration.inputByName('commonName'), commonName);
      await typeIn(S.configuration.inputByName('issuerName'), issuerName);
      await click(S.configuration.keyParamsGroupToggle);
      await typeIn(S.configuration.inputByName('keyName'), keyName);
      await click(S.configuration.generateRootSave);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(S.configuration.title).hasText('View root certificate');
      assert.dom('[data-test-alert-banner="alert"]').doesNotExist('no private key warning');
      assert.dom(S.configuration.title).hasText('View root certificate', 'Updates title on page');
      assert.dom(S.configuration.saved.certificate).hasClass('allow-copy', 'copyable certificate is masked');
      assert.dom(S.configuration.saved.issuerName).hasText(issuerName);
      assert.dom(S.configuration.saved.issuerLink).exists('Issuer link exists');
      assert.dom(S.configuration.saved.keyLink).exists('Key link exists');
      assert.dom(S.configuration.saved.keyName).hasText(keyName);
      assert.dom('[data-test-done]').exists('Done button exists');
      // Check that linked issuer has correct common name
      await click(S.configuration.saved.issuerLink);
      assert.dom(S.issuerDetails.valueByName('Common name')).hasText(commonName);
    });
    test('type=exported', async function (assert) {
      const commonName = 'my-exported-name';
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(S.configuration.optionByKey('generate-root'));
      // Fill in form
      await fillIn(S.configuration.typeField, 'exported');
      await typeIn(S.configuration.inputByName('commonName'), commonName);
      await click(S.configuration.generateRootSave);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(S.configuration.title).hasText('View root certificate');
      assert
        .dom('[data-test-alert-banner="alert"]')
        .hasText(
          'Next steps This private key material will only be available once. Copy or download it now.'
        );
      assert.dom(S.configuration.title).hasText('View root certificate', 'Updates title on page');
      assert
        .dom(S.configuration.saved.certificate)
        .hasClass('allow-copy', 'copyable masked certificate exists');
      assert
        .dom(S.configuration.saved.issuerName)
        .doesNotExist('Issuer name not shown because it was not named');
      assert.dom(S.configuration.saved.issuerLink).exists('Issuer link exists');
      assert.dom(S.configuration.saved.keyLink).exists('Key link exists');
      assert
        .dom(S.configuration.saved.privateKey)
        .hasClass('allow-copy', 'copyable masked private key exists');
      assert.dom(S.configuration.saved.keyName).doesNotExist('Key name not shown because it was not named');
      assert.dom('[data-test-done]').exists('Done button exists');
      // Check that linked issuer has correct common name
      await click(S.configuration.saved.issuerLink);
      assert.dom(S.issuerDetails.valueByName('Common name')).hasText(commonName);
    });
  });

  module('generate CSR', function () {
    test('happy path', async function (assert) {
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(S.emptyStateLink);
      assert.dom(S.configuration.title).hasText('Configure PKI');
      await click(S.configuration.optionByKey('generate-csr'));
      await fillIn(S.configuration.typeField, 'exported');
      await fillIn(S.configuration.inputByName('commonName'), 'my-common-name');
      await click('[data-test-save]');
      assert.dom(S.configuration.title).hasText('View generated CSR');
      await assert.dom(S.configuration.csrDetails).exists('renders CSR details after save');
      await click('[data-test-done]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'Transitions to overview after viewing csr details'
      );
    });
    test('type = exported', async function (assert) {
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(S.emptyStateLink);
      await click(S.configuration.optionByKey('generate-csr'));
      await fillIn(S.configuration.typeField, 'exported');
      await fillIn(S.configuration.inputByName('commonName'), 'my-common-name');
      await click('[data-test-save]');
      await assert.dom(S.configuration.csrDetails).exists('renders CSR details after save');
      assert.dom(S.configuration.title).hasText('View generated CSR');
      assert
        .dom('[data-test-alert-banner="alert"]')
        .hasText(
          'Next steps This private key material will only be available once. Copy or download it now.'
        );
      assert
        .dom(S.configuration.saved.privateKey)
        .hasClass('allow-copy', 'copyable masked private key exists');
      await click('[data-test-done]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'Transitions to overview after viewing csr details'
      );
    });
  });
});
