/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
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
      assert.dom(S.configuration.importError).hasText('Error please upload your PEM bundle');
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
    test('with many imports', async function (assert) {
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
      assert.dom(S.configuration.importForm).exists('import form is shown save');
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
    test('shows imported items when keys is empty', async function (assert) {
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
      assert.dom(S.configuration.importForm).exists('import form is shown save');
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click('[data-test-pki-import-pem-bundle]');

      assert.dom(S.configuration.importForm).doesNotExist('import form is hidden after save');
      assert.dom(S.configuration.importMapping).exists('import mapping is shown after save');
      assert.dom(S.configuration.importedIssuer).hasText('my-imported-issuer', 'Issuer value is displayed');
      assert.dom(S.configuration.importedKey).hasText('None', 'Shows placeholder value for key');
    });
    test('shows None for imported items if nothing new imported', async function (assert) {
      this.server.post(`${this.mountPath}/config/ca`, () => {
        return {
          request_id: 'some-config-id',
          data: {
            imported_issuers: null,
            imported_keys: null,
            mapping: {},
          },
        };
      });
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(S.configuration.optionByKey('import'));
      assert.dom(S.configuration.importForm).exists('import form is shown save');
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click('[data-test-pki-import-pem-bundle]');

      assert.dom(S.configuration.importForm).doesNotExist('import form is hidden after save');
      assert.dom(S.configuration.importMapping).exists('import mapping is shown after save');
      assert.dom(S.configuration.importedIssuer).hasText('None', 'Shows placeholder value for issuer');
      assert.dom(S.configuration.importedKey).hasText('None', 'Shows placeholder value for key');
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
      assert.dom(S.configuration.title).hasText('View Root Certificate');
      assert.dom(S.configuration.nextStepsBanner).doesNotExist('no private key warning');
      assert.dom(S.configuration.title).hasText('View Root Certificate', 'Updates title on page');
      assert.dom(S.configuration.saved.certificate).exists('Copyable certificate exists');
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
      assert.dom(S.configuration.title).hasText('View Root Certificate');
      assert
        .dom(S.configuration.nextStepsBanner)
        .hasText('Next steps The private_key is only available once. Make sure you copy and save it now.');
      assert.dom(S.configuration.title).hasText('View Root Certificate', 'Updates title on page');
      assert.dom(S.configuration.saved.certificate).exists('Copyable certificate exists');
      assert
        .dom(S.configuration.saved.issuerName)
        .doesNotExist('Issuer name not shown because it was not named');
      assert.dom(S.configuration.saved.issuerLink).exists('Issuer link exists');
      assert.dom(S.configuration.saved.keyLink).exists('Key link exists');
      assert.dom(S.configuration.saved.privateKey).exists('Copyable private key exists');
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
      await fillIn(S.configuration.typeField, 'internal');
      await fillIn(S.configuration.inputByName('commonName'), 'my-common-name');
      await click('[data-test-save]');
      assert.dom(S.configuration.title).hasText('View Generated CSR');
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
      assert.dom(S.configuration.title).hasText('View Generated CSR');
      assert
        .dom('[data-test-next-steps-csr]')
        .hasText(
          'Next steps Copy the CSR below for a parent issuer to sign and then import the signed certificate back into this mount. The private_key is only available once. Make sure you copy and save it now.'
        );
      assert.dom(S.configuration.saved.privateKey).exists('Copyable private key exists');
      await click('[data-test-done]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'Transitions to overview after viewing csr details'
      );
    });
  });
});
