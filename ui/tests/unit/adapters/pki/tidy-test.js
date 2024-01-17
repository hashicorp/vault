/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Unit | Adapter | pki/tidy', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
  });

  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/tidy');
    assert.ok(adapter);
  });

  test('it calls the correct endpoint when tidyType = manual', async function (assert) {
    assert.expect(1);

    this.server.post(`${this.backend}/tidy`, () => {
      assert.ok(true, 'request made to correct endpoint on create');
      return {};
    });
    this.payload = {
      tidy_cert_store: true,
      tidy_revocation_queue: false,
      safetyBuffer: '120h',
      backend: this.backend,
    };
    await this.store.createRecord('pki/tidy', this.payload).save({ adapterOptions: { tidyType: 'manual' } });
  });

  test('it should make a request to correct endpoint for findRecord', async function (assert) {
    assert.expect(1);
    this.server.get(`${this.backend}/config/auto-tidy`, () => {
      assert.ok(true, 'request made to correct endpoint on create');
      return {
        request_id: '2a4a1f36-20df-e71c-02d6-be15a09656f9',
        lease_id: '',
        renewable: false,
        lease_duration: 0,
        data: {
          acme_account_safety_buffer: 2592000,
          enabled: false,
          interval_duration: 43200,
          issuer_safety_buffer: 31536000,
          maintain_stored_certificate_counts: false,
          pause_duration: '0s',
          publish_stored_certificate_count_metrics: false,
          revocation_queue_safety_buffer: 172800,
          safety_buffer: 259200,
          tidy_acme: false,
          tidy_cert_store: false,
          tidy_cross_cluster_revoked_certs: false,
          tidy_expired_issuers: false,
          tidy_move_legacy_ca_bundle: false,
          tidy_revocation_queue: false,
          tidy_revoked_cert_issuer_associations: false,
          tidy_revoked_certs: false,
        },
        wrap_info: null,
        warnings: null,
        auth: null,
      };
    });

    this.store.findRecord('pki/tidy', this.backend);
  });
});
