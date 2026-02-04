/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { computeNavBar, NavSection, RouteName } from 'core/helpers/display-nav-item';
import { setupRenderingTest } from 'ember-qunit';
import { module, test } from 'qunit';
import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import sinon from 'sinon';
import { ROOT_NAMESPACE } from 'vault/services/namespace';

class PermissionsService extends Service {
  @tracked globPaths = null;
  @tracked exactPaths = null;
  @tracked canViewAll = false;
  @tracked chrootNamespace = null;

  hasNavPermission(...args) {
    const [route, routeParams, requireAll] = args;
    void (typeof route === 'string');
    void (routeParams === undefined || Array.isArray(routeParams));
    void (requireAll === undefined || typeof requireAll === 'boolean');

    if (this.canViewAll) return true;
    if (this.globPaths || this.exactPaths) return true;
    return false;
  }
}

module('Unit | Helper | displayNavItem', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:permissions', PermissionsService);
    this.permissions = this.owner.lookup('service:permissions');
    this.namespace = this.owner.lookup('service:namespace');
    this.currentCluster = this.owner.lookup('service:current-cluster');
    this.flags = this.owner.lookup('service:flags');
    this.version = this.owner.lookup('service:version');

    this.permissionsStub = sinon.stub(this.permissions, 'hasNavPermission');
  });

  hooks.afterEach(function () {
    this.permissionsStub.restore();
  });

  module('secrets sync', function () {
    test('it returns true when it is hvd managed', function (assert) {
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

      const supportsClientCount = computeNavBar(this, RouteName.SECRETS_SYNC);

      assert.true(supportsClientCount);
    });

    test('it returns true when it is secrets sync activated but not hvd managed', function (assert) {
      this.flags.featureFlags = [];
      this.flags.activatedFlags = [RouteName.SECRETS_SYNC];
      this.permissionsStub.returns(true);

      const supportsClientCount = computeNavBar(this, RouteName.SECRETS_SYNC);

      assert.true(supportsClientCount);
    });

    test('it returns false when it is secrets sync activated but not hvd managed and permissions is false', function (assert) {
      this.flags.featureFlags = [];
      this.flags.activatedFlags = [RouteName.SECRETS_SYNC];
      this.permissionsStub.returns(false);

      const supportsClientCount = computeNavBar(this, RouteName.SECRETS_SYNC);

      assert.false(supportsClientCount);
    });

    test('it returns false when it is enterprise', function (assert) {
      this.flags.featureFlags = [];
      this.version.type = 'community';
      this.features = [];

      const supportsClientCount = computeNavBar(this, RouteName.SECRETS_SYNC);

      assert.false(supportsClientCount);
    });
  });

  module('client count', function () {
    test('it returns true when there are permissions and cluster is not secondary', function (assert) {
      this.permissionsStub.returns(true);
      this.currentCluster.setCluster({
        name: 'cluster-0',
        hasChrootNamespace: false,
        dr: { isSecondary: false },
      });

      this.version.features = [];

      const supportsClientCount = computeNavBar(this, NavSection.CLIENT_COUNT);

      assert.true(supportsClientCount);
    });

    test('it returns false when there are no permissions and cluster is secondary', function (assert) {
      this.permissionsStub.returns(false);

      this.currentCluster.setCluster({
        hasChrootNamespace: true,
        dr: { isSecondary: true },
      });

      this.version.features = ['PKI-only Secrets'];

      const supportsClientCount = computeNavBar(this, NavSection.CLIENT_COUNT);

      assert.false(supportsClientCount);
    });
  });

  module('vault usage', function () {
    test('it returns true when there are permissions and is enterprise', function (assert) {
      this.permissionsStub.returns(true);
      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

      const supportsClientCount = computeNavBar(this, RouteName.VAULT_USAGE);

      assert.true(supportsClientCount);
    });

    test('it returns false when there are no permissions', function (assert) {
      this.permissionsStub.returns(false);

      this.version.type = 'community';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

      const supportsClientCount = computeNavBar(this, RouteName.VAULT_USAGE);

      assert.false(supportsClientCount);
    });
  });

  module('license', function () {
    test('it returns true when there are permissions and is enterprise', function (assert) {
      this.permissionsStub.returns(true);
      this.version.features = ['Secrets sync'];
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = ROOT_NAMESPACE;
      this.currentCluster.setCluster({
        hasChrootNamespace: false,
        dr: { isSecondary: false },
      });

      const supportsClientCount = computeNavBar(this, RouteName.LICENSE);

      assert.true(supportsClientCount);
    });

    test('it returns false when there are no permissions', function (assert) {
      this.permissionsStub.returns(false);
      this.version.features = ['Secrets sync'];
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

      const supportsClientCount = computeNavBar(this, RouteName.LICENSE);

      assert.false(supportsClientCount);
    });
  });
  module('reporting', function () {
    test('it returns true when user can access vault usage', function (assert) {
      this.permissionsStub.returns(true);
      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

      const supportsClientCount = computeNavBar(this, RouteName.REPORTING);

      assert.true(supportsClientCount);
    });

    test('it returns true when can support license but not vault usage', function (assert) {
      this.permissionsStub.returns(true);
      this.version.features = ['Secrets sync'];
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = ROOT_NAMESPACE;
      this.currentCluster.setCluster({
        hasChrootNamespace: false,
        dr: { isSecondary: false },
      });

      const supportsClientCount = computeNavBar(this, RouteName.REPORTING);

      assert.true(supportsClientCount);
    });

    test('it returns false when user cannot access vault usage or cannot support license', function (assert) {
      this.permissionsStub.returns(false);
      this.version.features = [];
      this.version.type = 'community';

      const supportsClientCount = computeNavBar(this, RouteName.LICENSE);

      assert.false(supportsClientCount);
    });
  });

  module('secrets recovery', function () {
    test('it returns true when dr is not secondary', function (assert) {
      this.currentCluster.setCluster({ dr: { isSecondary: false } });
      const supportsSecretsRecovery = computeNavBar(this, RouteName.SECRETS_RECOVERY);

      assert.true(supportsSecretsRecovery);
    });

    test('it returns true when flags is not hvd managed and dr is secondary', function (assert) {
      this.currentCluster.setCluster({ dr: { isSecondary: true } });
      this.flags.featureFlags = [];
      const supportsSecretsRecovery = computeNavBar(this, RouteName.SECRETS_RECOVERY);

      assert.true(supportsSecretsRecovery);
    });

    test('it returns false when flags is hvd managed and dr is secondary', function (assert) {
      this.currentCluster.setCluster({ dr: { isSecondary: false } });
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      const supportsSecretsRecovery = computeNavBar(this, RouteName.SECRETS_RECOVERY);

      assert.true(supportsSecretsRecovery);
    });
  });

  module('seal', function () {
    test('it returns true when in root namespace and has permissions', function (assert) {
      this.namespace.path = ROOT_NAMESPACE;
      this.permissionsStub.returns(true);
      this.currentCluster.setCluster({ dr: { isSecondary: false } });

      const supportsSecretsRecovery = computeNavBar(this, RouteName.SEAL);

      assert.true(supportsSecretsRecovery);
    });

    test('it returns false when flags is hvd managed and dr is secondary', function (assert) {
      this.namespace.path = 'child-namespace';
      this.permissionsStub.returns(false);
      this.currentCluster.setCluster({ dr: { isSecondary: true } });

      const supportsSecretsRecovery = computeNavBar(this, RouteName.SEAL);

      assert.false(supportsSecretsRecovery);
    });
  });

  module('replication', function () {
    test('it returns true when in root namespace, on enterprise, and has permissions', function (assert) {
      this.version.type = 'enterprise';
      this.namespace.path = ROOT_NAMESPACE;
      this.currentCluster.setCluster({
        replicationRedacted: false,
        hasChrootNamespace: false,
      });
      this.permissionsStub.returns(true);

      const supportsSecretsRecovery = computeNavBar(this, RouteName.REPLICATION);

      assert.true(supportsSecretsRecovery);
    });

    test('it returns false when hasChrootNamespace is true and type is community', function (assert) {
      this.version.type = 'community';
      this.namespace.path = ROOT_NAMESPACE;
      this.currentCluster.setCluster({
        replicationRedacted: false,
        hasChrootNamespace: true,
      });
      this.permissionsStub.returns(false);

      const supportsSecretsRecovery = computeNavBar(this, RouteName.REPLICATION);

      assert.false(supportsSecretsRecovery);
    });
  });

  module('resilience and recovery', function () {
    test('it returns true when supports secrets recovery', function (assert) {
      this.currentCluster.setCluster({ dr: { isSecondary: false } });

      const supportsSecretsRecovery = computeNavBar(this, NavSection.RESILIENCE_AND_RECOVERY);

      assert.true(supportsSecretsRecovery);
    });

    test('it returns true when supports replication but does not support secrets recovery or cannot seal', function (assert) {
      this.version.type = 'enterprise';
      this.namespace.path = ROOT_NAMESPACE;
      this.currentCluster.setCluster({ dr: { isSecondary: true }, replicationRedacted: false });

      this.permissionsStub.returns(true);

      const supportsSecretsRecovery = computeNavBar(this, NavSection.RESILIENCE_AND_RECOVERY);

      assert.true(supportsSecretsRecovery);
    });

    test('it returns false when does not supports replication, support secrets recovery and cannot seal', function (assert) {
      this.version.type = 'community';
      this.namespace.path = ROOT_NAMESPACE;
      this.currentCluster.setCluster({ dr: { isSecondary: true }, replicationRedacted: true });
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

      this.permissionsStub.returns(false);

      const supportsSecretsRecovery = computeNavBar(this, NavSection.RESILIENCE_AND_RECOVERY);

      assert.false(supportsSecretsRecovery);
    });
  });
});
