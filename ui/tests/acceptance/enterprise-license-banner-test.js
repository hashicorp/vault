/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import formatRFC3339 from 'date-fns/formatRFC3339';
import { addDays, subDays } from 'date-fns';
import timestamp from 'core/utils/timestamp';

const generateHealthResponse = (now, state) => {
  let expiry;
  switch (state) {
    case 'expired':
      expiry = subDays(now, 2);
      break;
    case 'expiring':
      expiry = addDays(now, 10);
      break;
    default:
      expiry = addDays(now, 33);
      break;
  }
  return {
    initialized: true,
    sealed: false,
    standby: false,
    license: {
      expiry_time: formatRFC3339(expiry),
      state: 'stored',
    },
    performance_standby: false,
    replication_performance_mode: 'disabled',
    replication_dr_mode: 'disabled',
    server_time_utc: 1622562585,
    version: '1.9.0+ent',
    cluster_name: 'vault-cluster-e779cd7c',
    cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
    last_wal: 121,
  };
};

module('Acceptance | Enterprise | License banner warnings', function (hooks) {
  setupApplicationTest(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => new Date('2018-04-03T14:15:30'));
  });
  hooks.beforeEach(function () {
    this.now = timestamp.now();
  });
  hooks.afterEach(function () {
    this.server.shutdown();
  });
  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it shows no license banner if license expires in > 30 days', async function (assert) {
    const healthResp = generateHealthResponse(this.now);
    this.server = new Pretender(function () {
      this.get('/v1/sys/health', (response) => {
        return [response, { 'Content-Type': 'application/json' }, JSON.stringify(healthResp)];
      });
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
      this.get('/v1/sys/internal/ui/mounts', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.get('/v1/sys/license/features', this.passthrough);
    });
    await visit('/vault/auth');
    assert.dom('[data-test-license-banner]').doesNotExist('license banner does not show');
    this.server.shutdown();
  });
  test('it shows license banner warning if license expires within 30 days', async function (assert) {
    const healthResp = generateHealthResponse(this.now, 'expiring');
    this.server = new Pretender(function () {
      this.get('/v1/sys/health', (response) => {
        return [response, { 'Content-Type': 'application/json' }, JSON.stringify(healthResp)];
      });
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
      this.get('/v1/sys/internal/ui/mounts', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.get('/v1/sys/license/features', this.passthrough);
    });
    await visit('/vault/auth');
    assert.dom('[data-test-license-banner-warning]').exists('license warning shows');
    this.server.shutdown();
  });

  test('it shows license banner alert if license has already expired', async function (assert) {
    const healthResp = generateHealthResponse(this.now, 'expired');
    this.server = new Pretender(function () {
      this.get('/v1/sys/health', (response) => {
        return [response, { 'Content-Type': 'application/json' }, JSON.stringify(healthResp)];
      });
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
      this.get('/v1/sys/internal/ui/mounts', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.get('/v1/sys/license/features', this.passthrough);
    });
    await visit('/vault/auth');
    assert.dom('[data-test-license-banner-expired]').exists('expired license message shows');
    this.server.shutdown();
  });
});
