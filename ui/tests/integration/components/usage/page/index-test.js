/**
 * Copyright (c) HashiCorp, Inc.
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

module('Integration | Component | usage | Page::Usage', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.api = this.owner.lookup('service:api');
    this.generateUtilizationReportStub = sinon.stub(this.api.sys, 'generateUtilizationReport').resolves({});
  });

  hooks.afterEach(function () {
    this.generateUtilizationReportStub.restore();
  });

  test('it provides the correct fetch function to the dashboard component', async function (assert) {
    await render(hbs`<Usage::Page />`);
    assert.true(this.generateUtilizationReportStub.calledOnce, 'fetch function is called on render');
  });

  test('it remaps data to friendly names if available', async function (assert) {
    this.generateUtilizationReportStub.resolves({
      auth_methods: { alicloud: 2, cert: 2, userpass: 2, 'unknown-random-method': 1 },
      kvv1_secrets: 15,
      kvv2_secrets: 146,
      lease_count_quotas: {
        global_lease_count_quota: {
          capacity: 300000,
          count: 244121,
          name: 'default',
        },
        total_lease_count_quotas: 2,
      },
      namespaces: 10,
      secrets_sync: 79,
      pki: { total_issuers: 2, total_roles: 6 },
      replication_status: {
        dr_primary: false,
        dr_state: 'disabled',
        pr_primary: false,
        pr_state: 'enabled',
      },
      secret_engines: {
        keymgmt: 5,
        gcpkms: 10,
        pki: 11,
        'unknown-random-engine': 1,
      },
    });
    await render(hbs`<Usage::Page />`);

    const engineLabels = [...findAll('[data-test-dashboard-secret-engines] .axis g text')].map(
      (label) => label.textContent
    );
    const authMethodLabels = [...findAll('[data-test-dashboard-auth-methods] .axis g text')].map(
      (label) => label.textContent
    );
    assert.deepEqual(
      engineLabels,
      ['PKI Certificates', 'Google Cloud KMS', 'Key Management', 'unknown-random-engine'],
      'Engine labels are correct (sorted DESC)'
    );

    assert.deepEqual(
      authMethodLabels,
      ['AliCloud', 'TLS Certificates', 'Username & Password', 'unknown-random-method'],
      'Auth method labels are correct (sorted DESC)'
    );
  });
});
