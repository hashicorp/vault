/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { PKI_TIDY } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | Page::PkiTidyStatus', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';

    this.store.createRecord('pki/issuer', { issuerId: 'abcd-efgh' });
    this.store.createRecord('pki/tidy', { backend: this.secretMountPath.currentPath, tidyType: 'auto' });

    this.autoTidyConfig = this.store.peekAll('pki/tidy');
    this.tidyStatus = {
      acme_account_deleted_count: 0,
      acme_account_revoked_count: 0,
      acme_account_safety_buffer: 2592000,
      acme_orders_deleted_count: 0,
      cert_store_deleted_count: 0,
      cross_revoked_cert_deleted_count: 0,
      current_cert_store_count: null,
      current_revoked_cert_count: null,
      error: null,
      internal_backend_uuid: '9d3bd186-0fdd-9ca4-f298-2e180536b743',
      issuer_safety_buffer: 31536000,
      last_auto_tidy_finished: '2023-05-18T13:27:36.390785-07:00',
      message: 'Tidying certificate store: checking entry 0 of 1',
      missing_issuer_cert_count: 0,
      pause_duration: '15s',
      revocation_queue_deleted_count: 0,
      revocation_queue_safety_buffer: 36000,
      revoked_cert_deleted_count: 0,
      safety_buffer: 2073600,
      state: 'Running',
      tidy_acme: false,
      tidy_cert_store: true,
      tidy_cross_cluster_revoked_certs: false,
      tidy_expired_issuers: false,
      tidy_move_legacy_ca_bundle: false,
      time_started: '2023-05-18T13:27:36.390959-07:00',
    };
    this.engineId = 'pki';
  });

  test('shows the correct titles for the alert banner based on states', async function (assert) {
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    // running state
    assert.dom(PKI_TIDY.hdsAlertTitle).hasText('Tidy in progress');
    assert.dom(PKI_TIDY.cancelTidyAction).exists();
    assert.dom(PKI_TIDY.hdsAlertButtonText).hasText('Cancel tidy');
    // inactive state
    this.tidyStatus.state = 'Inactive';
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    assert.dom(PKI_TIDY.hdsAlertTitle).hasText('Tidy is inactive');
    // finished state
    this.tidyStatus.state = 'Finished';
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    assert.dom(PKI_TIDY.hdsAlertTitle).hasText('Tidy operation finished');
    // error state
    this.tidyStatus.state = 'Error';
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    assert.dom(PKI_TIDY.hdsAlertTitle).hasText('Tidy operation failed');
    // cancelling state
    this.tidyStatus.state = 'Cancelling';
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    assert.dom(PKI_TIDY.hdsAlertTitle).hasText('Tidy operation cancelling');
    // cancelled state
    this.tidyStatus.state = 'Cancelled';
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    assert.dom(PKI_TIDY.hdsAlertTitle).hasText('Tidy operation cancelled');
  });
  test('shows the fields even if the data returns null values', async function (assert) {
    this.tidyStatus.time_started = null;
    this.tidyStatus.time_finished = null;
    await render(
      hbs`<Page::PkiTidyStatus @autoTidyConfig={{this.autoTidyConfig}} @tidyStatus={{this.tidyStatus}} />`,
      { owner: this.engine }
    );
    assert.dom(PKI_TIDY.timeStartedRow).exists();
    assert.dom(PKI_TIDY.timeFinishedRow).exists();
  });
});
