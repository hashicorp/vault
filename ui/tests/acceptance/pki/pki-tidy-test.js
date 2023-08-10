/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';

import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-tidy';

module('Acceptance | pki tidy', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await runCommands([
      `write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test" name="Hashicorp Test"`,
    ]);
    await logout.visit();
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
  });

  test('it configures a manual tidy operation and shows its details and tidy states', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/tidy`);
    await click(SELECTORS.tidyEmptyStateConfigure);
    assert.dom(SELECTORS.tidyConfigureModal.configureTidyModal).exists('Configure tidy modal exists');
    assert.dom(SELECTORS.tidyConfigureModal.tidyModalAutoButton).exists('Configure auto tidy button exists');
    assert
      .dom(SELECTORS.tidyConfigureModal.tidyModalManualButton)
      .exists('Configure manual tidy button exists');
    await click(SELECTORS.tidyConfigureModal.tidyModalManualButton);
    assert.dom(SELECTORS.tidyForm.tidyFormName('manual')).exists('Manual tidy form exists');
    await click(SELECTORS.tidyForm.inputByAttr('tidyCertStore'));
    await fillIn(SELECTORS.tidyForm.tidyPauseDuration, '10');
    await click(SELECTORS.tidyForm.tidySave);
    await click(SELECTORS.cancelTidyAction);
    assert.dom(SELECTORS.cancelTidyModalBackground).exists('Confirm cancel tidy modal exits');
    await click(SELECTORS.tidyConfigureModal.tidyModalCancelButton);
    // we can't properly test the background refresh fetching of tidy status in testing
    this.server.get(`${this.mountPath}/tidy-status`, () => {
      return {
        request_id: 'dba2d42d-1a6e-1551-80f8-4ddb364ede4b',
        lease_id: '',
        renewable: false,
        lease_duration: 0,
        data: {
          acme_account_deleted_count: 0,
          acme_account_revoked_count: 0,
          acme_account_safety_buffer: 2592000,
          acme_orders_deleted_count: 0,
          cert_store_deleted_count: 0,
          cross_revoked_cert_deleted_count: 0,
          current_cert_store_count: null,
          current_revoked_cert_count: null,
          error: null,
          internal_backend_uuid: '964a41f7-a159-53aa-d62e-fc1914e4a7e1',
          issuer_safety_buffer: 31536000,
          last_auto_tidy_finished: '2023-05-19T10:27:11.721825-07:00',
          message: 'Tidying certificate store: checking entry 0 of 1',
          missing_issuer_cert_count: 0,
          pause_duration: '1m40s',
          revocation_queue_deleted_count: 0,
          revocation_queue_safety_buffer: 36000,
          revoked_cert_deleted_count: 0,
          safety_buffer: 2073600,
          state: 'Cancelled',
          tidy_acme: false,
          tidy_cert_store: true,
          tidy_cross_cluster_revoked_certs: false,
          tidy_expired_issuers: false,
          tidy_move_legacy_ca_bundle: false,
          tidy_revocation_queue: false,
          tidy_revoked_cert_issuer_associations: false,
          tidy_revoked_certs: false,
          time_finished: '2023-05-19T10:28:51.733092-07:00',
          time_started: '2023-05-19T10:27:11.721846-07:00',
          total_acme_account_count: 0,
        },
        wrap_info: null,
        warnings: null,
        auth: null,
      };
    });
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await visit(`/vault/secrets/${this.mountPath}/pki/tidy`);
    assert.dom(SELECTORS.hdsAlertTitle).hasText('Tidy operation cancelled');
    assert
      .dom(SELECTORS.hdsAlertDescription)
      .hasText(
        'Your tidy operation has been cancelled. If this was a mistake configure and run another tidy operation.'
      );
    assert.dom(SELECTORS.alertUpdatedAt).exists();
  });

  test('it configures an auto tidy operation and shows its details', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/tidy`);
    await click(SELECTORS.tidyEmptyStateConfigure);
    assert.dom(SELECTORS.tidyConfigureModal.configureTidyModal).exists('Configure tidy modal exists');
    assert.dom(SELECTORS.tidyConfigureModal.tidyModalAutoButton).exists('Configure auto tidy button exists');
    assert
      .dom(SELECTORS.tidyConfigureModal.tidyModalManualButton)
      .exists('Configure manual tidy button exists');
    await click(SELECTORS.tidyConfigureModal.tidyModalAutoButton);
    assert.dom(SELECTORS.tidyForm.tidyFormName('auto')).exists('Auto tidy form exists');
    await click(SELECTORS.tidyForm.tidyCancel);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.tidy.index');
    await click(SELECTORS.tidyEmptyStateConfigure);
    await click(SELECTORS.tidyConfigureModal.tidyModalAutoButton);
    assert.dom(SELECTORS.tidyForm.tidyFormName('auto')).exists('Auto tidy form exists');
    await click(SELECTORS.tidyForm.toggleLabel('Automatic tidy disabled'));
    assert
      .dom(SELECTORS.tidyForm.tidySectionHeader('ACME operations'))
      .exists('Auto tidy form enabled shows ACME operations field');
    await click(SELECTORS.tidyForm.inputByAttr('tidyCertStore'));
    await click(SELECTORS.tidyForm.tidySave);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.tidy.auto.index');
    await click(SELECTORS.tidyForm.editAutoTidyButton);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.tidy.auto.configure');
    await click(SELECTORS.tidyForm.inputByAttr('tidyRevokedCerts'));
    await click(SELECTORS.tidyForm.tidySave);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.tidy.auto.index');
  });

  test('it opens a tidy modal when the user clicks on the tidy toolbar action', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/tidy`);
    await click(SELECTORS.tidyConfigureModal.tidyOptionsModal);
    assert.dom(SELECTORS.tidyConfigureModal.configureTidyModal).exists('Configure tidy modal exists');
    assert.dom(SELECTORS.tidyConfigureModal.tidyModalAutoButton).exists('Configure auto tidy button exists');
    assert
      .dom(SELECTORS.tidyConfigureModal.tidyModalManualButton)
      .exists('Configure manual tidy button exists');
    await click(SELECTORS.tidyConfigureModal.tidyModalCancelButton);
    assert.dom(SELECTORS.tidyEmptyState).exists();
  });

  test('it should show correct toolbar action depending on whether auto tidy is enabled', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/tidy`);
    assert
      .dom(SELECTORS.tidyConfigureModal.tidyOptionsModal)
      .exists('Configure tidy modal options button exists');
    await click(SELECTORS.tidyConfigureModal.tidyOptionsModal);
    assert.dom(SELECTORS.tidyConfigureModal.configureTidyModal).exists('Configure tidy modal exists');
    await click(SELECTORS.tidyConfigureModal.tidyOptionsModal);
    await click(SELECTORS.tidyConfigureModal.tidyModalAutoButton);
    await click(SELECTORS.tidyForm.toggleLabel('Automatic tidy disabled'));
    await click(SELECTORS.tidyForm.inputByAttr('tidyCertStore'));
    await click(SELECTORS.tidyForm.inputByAttr('tidyRevokedCerts'));
    await click(SELECTORS.tidyForm.tidySave);
    await visit(`/vault/secrets/${this.mountPath}/pki/tidy`);
    assert
      .dom(SELECTORS.manualTidyToolbar)
      .exists('Manual tidy toolbar action exists if auto tidy is configured');
    assert
      .dom(SELECTORS.autoTidyToolbar)
      .exists('Auto tidy toolbar action exists if auto tidy is configured');
  });
});
