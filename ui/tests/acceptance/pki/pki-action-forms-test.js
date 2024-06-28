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
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { PKI_CONFIGURE_CREATE, PKI_GENERATE_ROOT } from 'vault/tests/helpers/pki/pki-selectors';

const { issuerPemBundle } = CERTIFICATES;

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
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
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
      await click(`${GENERAL.emptyStateActions} a`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      assert.dom(GENERAL.title).hasText('Configure PKI');
      assert.dom(GENERAL.emptyStateTitle).exists({ count: 1 }, 'Shows empty state by default');
      await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
      assert.dom(GENERAL.emptyStateTitle).doesNotExist();
      // Submit before filling out form shows an error
      await click(PKI_CONFIGURE_CREATE.importSubmit);
      assert.dom(GENERAL.messageError).hasText('Error please upload your PEM bundle');
      // Fill in form data
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click(PKI_CONFIGURE_CREATE.importSubmit);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(GENERAL.title).hasText('View imported items');
      assert.dom(PKI_CONFIGURE_CREATE.importForm).doesNotExist('import form is hidden after save');
      assert.dom(PKI_CONFIGURE_CREATE.importMapping).exists('import mapping is shown after save');
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
      await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
      assert.dom(PKI_CONFIGURE_CREATE.importForm).exists('import form is shown save');
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click(PKI_CONFIGURE_CREATE.importSubmit);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(GENERAL.title).hasText('View imported items');
      assert.dom(PKI_CONFIGURE_CREATE.importForm).doesNotExist('import form is hidden after save');
      assert.dom(PKI_CONFIGURE_CREATE.importMapping).exists('import mapping is shown after save');
      assert
        .dom(PKI_CONFIGURE_CREATE.importedIssuer)
        .hasText('my-imported-issuer', 'Issuer value is displayed');
      assert.dom(PKI_CONFIGURE_CREATE.importedKey).hasText('my-imported-key', 'Key value is displayed');
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
      await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
      assert.dom(PKI_CONFIGURE_CREATE.importForm).exists('import form is shown save');
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click(PKI_CONFIGURE_CREATE.importSubmit);

      assert.dom(PKI_CONFIGURE_CREATE.importForm).doesNotExist('import form is hidden after save');
      assert.dom(PKI_CONFIGURE_CREATE.importMapping).exists('import mapping is shown after save');
      assert
        .dom(PKI_CONFIGURE_CREATE.importedIssuer)
        .hasText('my-imported-issuer', 'Issuer value is displayed');
      assert.dom(PKI_CONFIGURE_CREATE.importedKey).hasText('None', 'Shows placeholder value for key');
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
      await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
      assert.dom(PKI_CONFIGURE_CREATE.importForm).exists('import form is shown save');
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', this.pemBundle);
      await click(PKI_CONFIGURE_CREATE.importSubmit);

      assert.dom(PKI_CONFIGURE_CREATE.importForm).doesNotExist('import form is hidden after save');
      assert.dom(PKI_CONFIGURE_CREATE.importMapping).exists('import mapping is shown after save');
      assert.dom(PKI_CONFIGURE_CREATE.importedIssuer).hasText('None', 'Shows placeholder value for issuer');
      assert.dom(PKI_CONFIGURE_CREATE.importedKey).hasText('None', 'Shows placeholder value for key');
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
      await click(`${GENERAL.emptyStateActions} a`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      assert.dom(GENERAL.title).hasText('Configure PKI');
      assert.dom(GENERAL.emptyStateTitle).exists({ count: 1 }, 'Shows empty state by default');
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
      assert.dom(GENERAL.emptyStateTitle).doesNotExist();
      // The URLs section is populated based on params returned from OpenAPI. This test will break when
      // the backend adds fields. We should update the count accordingly.
      assert.dom(PKI_GENERATE_ROOT.urlField).exists({ count: 4 });
      // Fill in form
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await typeIn(GENERAL.inputByAttr('commonName'), commonName);
      await typeIn(GENERAL.inputByAttr('issuerName'), issuerName);
      await click(PKI_GENERATE_ROOT.keyParamsGroupToggle);
      await typeIn(GENERAL.inputByAttr('keyName'), keyName);
      await click(GENERAL.saveButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(GENERAL.title).hasText('View Root Certificate');
      assert.dom(PKI_CONFIGURE_CREATE.nextStepsBanner).doesNotExist('no private key warning');
      assert.dom(GENERAL.title).hasText('View Root Certificate', 'Updates title on page');
      assert.dom(PKI_GENERATE_ROOT.saved.certificate).exists('Copyable certificate exists');
      assert.dom(PKI_GENERATE_ROOT.saved.issuerName).hasText(issuerName);
      assert.dom(PKI_GENERATE_ROOT.saved.issuerLink).exists('Issuer link exists');
      assert.dom(PKI_GENERATE_ROOT.saved.keyLink).exists('Key link exists');
      assert.dom(PKI_GENERATE_ROOT.saved.keyName).hasText(keyName);
      assert.dom('[data-test-done]').exists('Done button exists');
      // Check that linked issuer has correct common name
      await click(PKI_GENERATE_ROOT.saved.issuerLink);
      assert.dom(GENERAL.infoRowValue('Common name')).hasText(commonName);
    });
    test('type=exported', async function (assert) {
      const commonName = 'my-exported-name';
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
      // Fill in form
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await typeIn(GENERAL.inputByAttr('commonName'), commonName);
      await click(GENERAL.saveButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'stays on page on success'
      );
      assert.dom(GENERAL.title).hasText('View Root Certificate');
      assert
        .dom(PKI_CONFIGURE_CREATE.nextStepsBanner)
        .hasText('Next steps The private_key is only available once. Make sure you copy and save it now.');
      assert.dom(GENERAL.title).hasText('View Root Certificate', 'Updates title on page');
      assert.dom(PKI_GENERATE_ROOT.saved.certificate).exists('Copyable certificate exists');
      assert
        .dom(PKI_GENERATE_ROOT.saved.issuerName)
        .doesNotExist('Issuer name not shown because it was not named');
      assert.dom(PKI_GENERATE_ROOT.saved.issuerLink).exists('Issuer link exists');
      assert.dom(PKI_GENERATE_ROOT.saved.keyLink).exists('Key link exists');
      assert.dom(PKI_GENERATE_ROOT.saved.privateKey).exists('Copyable private key exists');
      assert.dom(PKI_GENERATE_ROOT.saved.keyName).doesNotExist('Key name not shown because it was not named');
      assert.dom('[data-test-done]').exists('Done button exists');
      // Check that linked issuer has correct common name
      await click(PKI_GENERATE_ROOT.saved.issuerLink);
      assert.dom(GENERAL.infoRowValue('Common name')).hasText(commonName);
    });
  });

  module('generate CSR', function () {
    test('happy path', async function (assert) {
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(`${GENERAL.emptyStateActions} a`);
      assert.dom(GENERAL.title).hasText('Configure PKI');
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-csr'));
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await fillIn(GENERAL.inputByAttr('commonName'), 'my-common-name');
      await click('[data-test-save]');
      assert.dom(GENERAL.title).hasText('View Generated CSR');
      await assert.dom(PKI_CONFIGURE_CREATE.csrDetails).exists('renders CSR details after save');
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
      await click(`${GENERAL.emptyStateActions} a`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-csr'));
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await fillIn(GENERAL.inputByAttr('commonName'), 'my-common-name');
      await click('[data-test-save]');
      await assert.dom(PKI_CONFIGURE_CREATE.csrDetails).exists('renders CSR details after save');
      assert.dom(GENERAL.title).hasText('View Generated CSR');
      assert
        .dom('[data-test-next-steps-csr]')
        .hasText(
          'Next steps Copy the CSR below for a parent issuer to sign and then import the signed certificate back into this mount. The private_key is only available once. Make sure you copy and save it now.'
        );
      assert.dom(PKI_GENERATE_ROOT.saved.privateKey).exists('Copyable private key exists');
      await click('[data-test-done]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'Transitions to overview after viewing csr details'
      );
    });
  });
});
