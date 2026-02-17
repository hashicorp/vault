/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentRouteName, currentURL } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from '../helpers/general-selectors';
import sinon from 'sinon';
import * as displayNavItem from 'core/helpers/display-nav-item';

module('Acceptance | Enterprise | resilience-recovery', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.permissionsStub = sinon.stub(displayNavItem, 'computeNavBar');

    await login();
  });

  hooks.afterEach(function () {
    this.permissionsStub.restore();
  });

  test('it should redirect to recovery snapshots route when replication is disabled', async function (assert) {
    this.permissionsStub.callsFake((context, routeName) => {
      if (routeName === displayNavItem.RouteName.SECRETS_RECOVERY) {
        return true;
      }

      return false;
    });

    await click(GENERAL.navLink('Resilience and recovery'));

    assert.strictEqual(currentURL(), '/vault/recovery/snapshots', 'snapshots url renders');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.recovery.snapshots.index',
      'snapshots route renders'
    );
  });

  test('it should redirect to seal route when user has permission and secrets recovery is not supported', async function (assert) {
    this.permissionsStub.callsFake((context, routeName) => {
      if (routeName === displayNavItem.RouteName.SEAL) {
        return true;
      }

      return false;
    });

    await click(GENERAL.navLink('Resilience and recovery'));

    assert.strictEqual(currentURL(), '/vault/settings/seal', 'seal route renders');
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.seal', 'seal route renders');
  });

  test('it should redirect to replication route when secrets recovery or seal is not supported', async function (assert) {
    this.permissionsStub.callsFake((context, routeName) => {
      if (routeName === displayNavItem.RouteName.REPLICATION) {
        return true;
      }

      return false;
    });

    await click(GENERAL.navLink('Resilience and recovery'));

    assert.strictEqual(currentURL(), '/vault/replication', 'replication route renders');
    assert.strictEqual(currentRouteName(), 'vault.cluster.replication.index', 'replication route renders');
  });
});
