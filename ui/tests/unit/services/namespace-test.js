/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { ADMINISTRATIVE_NAMESPACE, ROOT_NAMESPACE } from 'vault/services/namespace';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Unit | Service | namespace', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api');
    this.auth = this.owner.lookup('service:auth');
    this.authStub = sinon.stub(this.auth, 'authData');
    this.flags = this.owner.lookup('service:flags');
    this.namespaceService = this.owner.lookup('service:namespace');
  });

  hooks.afterEach(function () {
    this.authStub.restore();
  });

  module('#userRootNamespace', function () {
    test('self-managed clusters: returns userRootNamespace from auth data', function (assert) {
      this.authStub.value({ userRootNamespace: 'bahamas' });

      assert.strictEqual(this.namespaceService.userRootNamespace, 'bahamas', 'returns user root namespace');
    });

    test('self-managed clusters: it returns fallback when auth data is not set', function (assert) {
      this.authStub.value(undefined);

      assert.strictEqual(this.namespaceService.userRootNamespace, ROOT_NAMESPACE);
    });

    test('self-managed clusters: it returns fallback when userRootNamespace is undefined', function (assert) {
      this.authStub.value({ userRootNamespace: undefined });

      assert.strictEqual(this.namespaceService.userRootNamespace, ROOT_NAMESPACE);
    });

    test('self-managed clusters: it returns fallback when userRootNamespace is null', function (assert) {
      this.authStub.value({ userRootNamespace: null });

      assert.strictEqual(this.namespaceService.userRootNamespace, ROOT_NAMESPACE);
    });

    test('self-managed clusters: it returns empty string when userRootNamespace is an empty string', function (assert) {
      this.authStub.value({ userRootNamespace: '' });

      assert.strictEqual(this.namespaceService.userRootNamespace, ROOT_NAMESPACE);
    });

    module('HVD managed clusters', function (hooks) {
      hooks.beforeEach(function () {
        this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      });

      test('returns userRootNamespace from auth data', function (assert) {
        this.authStub.value({ userRootNamespace: 'bahamas' });

        assert.strictEqual(this.namespaceService.userRootNamespace, 'bahamas', 'returns user root namespace');
      });

      test('it returns fallback when auth data is not set', function (assert) {
        this.authStub.value(undefined);

        assert.strictEqual(this.namespaceService.userRootNamespace, ADMINISTRATIVE_NAMESPACE);
      });

      test('it returns fallback when userRootNamespace is undefined', function (assert) {
        this.authStub.value({ userRootNamespace: undefined });

        assert.strictEqual(this.namespaceService.userRootNamespace, ADMINISTRATIVE_NAMESPACE);
      });

      test('it returns fallback when userRootNamespace is null', function (assert) {
        this.authStub.value({ userRootNamespace: null });

        assert.strictEqual(this.namespaceService.userRootNamespace, ADMINISTRATIVE_NAMESPACE);
      });

      test('it still returns an empty string when userRootNamespace is an empty string', function (assert) {
        this.authStub.value({ userRootNamespace: '' });

        assert.strictEqual(
          this.namespaceService.userRootNamespace,
          '',
          'an empty string takes precedence if authData?.userRootNamespace is set as such'
        );
      });
    });
  });

  module('#inRootNamespace', function () {
    test('returns true when path is empty string', function (assert) {
      this.namespaceService.path = '';

      assert.true(this.namespaceService.inRootNamespace, 'returns true for root namespace');
    });

    test('returns false when path is not empty', function (assert) {
      this.namespaceService.path = 'admin';

      assert.false(this.namespaceService.inRootNamespace, 'returns false for non-root namespace');
    });
  });

  module('#inHvdAdminNamespace', function () {
    test('returns true when in HVD managed cluster is in admin namespace', function (assert) {
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespaceService.path = 'admin';

      assert.true(
        this.namespaceService.inHvdAdminNamespace,
        'returns true for HVD managed cluster at admin namespace'
      );
    });

    test('returns false when not in HVD managed cluster', function (assert) {
      this.flags.featureFlags = [];
      this.namespaceService.path = 'admin';

      assert.false(this.namespaceService.inHvdAdminNamespace, 'returns false when not HVD managed');
    });

    test('returns false when in HVD managed cluster but not in admin namespace', function (assert) {
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespaceService.path = 'other';

      assert.false(this.namespaceService.inHvdAdminNamespace, 'returns false when not at admin namespace');
    });
  });

  module('#currentNamespace', function () {
    test('returns "root" when in root namespace', function (assert) {
      this.namespaceService.path = '';

      assert.strictEqual(this.namespaceService.currentNamespace, 'root', 'returns "root" for empty path');
    });

    test('returns last part of path for nested namespace', function (assert) {
      this.namespaceService.path = 'parent/child/grandchild';

      assert.strictEqual(
        this.namespaceService.currentNamespace,
        'grandchild',
        'returns last segment of namespace path'
      );
    });

    test('returns namespace for single-level path', function (assert) {
      this.namespaceService.path = 'carrot';

      assert.strictEqual(this.namespaceService.currentNamespace, 'carrot', 'returns single namespace');
    });
  });

  module('#relativeNamespace', function () {
    test('returns relative path when userRootNamespace is set', function (assert) {
      this.authStub.value({ userRootNamespace: 'app' });
      this.namespaceService.path = 'app/staging/group1';

      assert.strictEqual(
        this.namespaceService.relativeNamespace,
        'staging/group1',
        'returns path relative to user root'
      );
    });

    test('returns full path when userRootNamespace is empty', function (assert) {
      this.authStub.value({ userRootNamespace: '' });
      this.namespaceService.path = 'app/staging';

      assert.strictEqual(
        this.namespaceService.relativeNamespace,
        'app/staging',
        'returns full path when user root is empty'
      );
    });

    test('returns empty string when path equals userRootNamespace', function (assert) {
      this.authStub.value({ userRootNamespace: 'app' });
      this.namespaceService.path = 'app';

      assert.strictEqual(
        this.namespaceService.relativeNamespace,
        '',
        'returns empty string when at user root'
      );
    });
  });

  module('#setNamespace', function () {
    test('sets path to empty string when path is undefined', function (assert) {
      this.namespaceService.setNamespace(undefined);

      assert.strictEqual(this.namespaceService.path, '', 'sets path to empty string for undefined');
    });

    test('sets path to empty string when path is "root"', function (assert) {
      this.namespaceService.setNamespace('root');

      assert.strictEqual(this.namespaceService.path, '', 'converts "root" to empty string');
    });

    test('it sets path to empty string when "root/" has a trailing slash', function (assert) {
      this.namespaceService.setNamespace('root/');

      assert.strictEqual(this.namespaceService.path, '', '');
    });

    test('sets path when valid namespace path is provided', function (assert) {
      this.namespaceService.setNamespace('admin/team');

      assert.strictEqual(this.namespaceService.path, 'admin/team', 'sets namespace path');
    });

    // Right now the auth cluster handles sanitizing the namespace input so setNamespace does NOT
    // manage that logic. We might want to consider moving it to setNamespace for consistency.
    test('it sets all non-root namespaces as is', function (assert) {
      this.namespaceService.setNamespace('admin/');

      assert.strictEqual(this.namespaceService.path, 'admin/', 'admin maintains trailing slash');
    });
  });

  module('#findNamespacesForUser', function (hooks) {
    hooks.beforeEach(function () {
      this.internalUiListNamespacesStub = sinon.stub(this.api.sys, 'internalUiListNamespaces');
      this.stubNamespaces = async (namespaces) => {
        this.internalUiListNamespacesStub.resolves({ keys: namespaces });
      };
      this.authStub.value({ userRootNamespace: '' });
    });

    hooks.afterEach(function () {
      this.internalUiListNamespacesStub.restore();
    });

    test('fetches namespaces and updates accessibleNamespaces', async function (assert) {
      this.stubNamespaces(['ns1', 'ns2', 'ns3']);
      await this.namespaceService.findNamespacesForUser.perform();

      assert.deepEqual(
        this.namespaceService.accessibleNamespaces,
        ['ns1', 'ns2', 'ns3'],
        'sets accessible namespaces'
      );
    });

    test('prepends userRootNamespace to each namespace key', async function (assert) {
      this.authStub.value({ userRootNamespace: 'parent' });
      this.stubNamespaces(['team1', 'team2']);

      await this.namespaceService.findNamespacesForUser.perform();

      assert.deepEqual(
        this.namespaceService.accessibleNamespaces,
        ['parent/team1', 'parent/team2'],
        'prepends user root to namespace keys'
      );
    });

    test('removes trailing slashes from namespace paths', async function (assert) {
      this.stubNamespaces(['ns1/', 'ns2/']);

      await this.namespaceService.findNamespacesForUser.perform();

      assert.deepEqual(
        this.namespaceService.accessibleNamespaces,
        ['ns1', 'ns2'],
        'removes trailing slashes'
      );
    });

    test('handles undefined keys', async function (assert) {
      this.stubNamespaces();
      await this.namespaceService.findNamespacesForUser.perform();

      assert.deepEqual(this.namespaceService.accessibleNamespaces, [], 'sets empty array');
    });

    test('it handles API error gracefully', async function (assert) {
      this.internalUiListNamespacesStub.rejects(getErrorResponse({}, 500));

      await this.namespaceService.findNamespacesForUser.perform();

      assert.strictEqual(
        this.namespaceService.accessibleNamespaces,
        null,
        'it does not set accessibleNamespaces when response errors'
      );
    });

    test('it does not modify accessibleNamespaces on error', async function (assert) {
      this.namespaceService.accessibleNamespaces = ['existing'];
      this.internalUiListNamespacesStub.rejects(getErrorResponse({}, 500));

      await this.namespaceService.findNamespacesForUser.perform();

      assert.deepEqual(
        this.namespaceService.accessibleNamespaces,
        ['existing'],
        'accessibleNamespaces is unchanged'
      );
    });
  });

  module('#reset', function () {
    test('resets accessibleNamespaces to null', function (assert) {
      this.namespaceService.accessibleNamespaces = ['ns1', 'ns2'];

      this.namespaceService.reset();

      assert.strictEqual(this.namespaceService.accessibleNamespaces, null, 'resets to null');
    });
  });

  module('#getOptions', function () {
    test('it adds empty string as userRootNamespace to options', function (assert) {
      this.namespaceService.accessibleNamespaces = ['ns1', 'ns2'];
      this.authStub.value({ userRootNamespace: '' });

      const options = this.namespaceService.getOptions();

      assert.deepEqual(
        options,
        [
          { path: '', label: 'root' },
          { path: 'ns1', label: 'ns1' },
          { path: 'ns2', label: 'ns2' },
        ],
        'returns namespace options with root'
      );
    });

    test('it adds admin userRootNamespace when not in accessible namespaces', function (assert) {
      this.namespaceService.accessibleNamespaces = ['admin/team1', 'admin/team2'];
      this.authStub.value({ userRootNamespace: 'admin' });

      const options = this.namespaceService.getOptions();

      assert.deepEqual(
        options,
        [
          { path: 'admin', label: 'admin' },
          { path: 'admin/team1', label: 'admin/team1' },
          { path: 'admin/team2', label: 'admin/team2' },
        ],
        'adds user root namespace at the beginning'
      );
    });

    test('labels empty string (root) namespace as "root"', function (assert) {
      this.namespaceService.accessibleNamespaces = ['ns1'];
      this.authStub.value({ userRootNamespace: '' });

      const options = this.namespaceService.getOptions();

      assert.strictEqual(options[0].path, '', 'root namespace path is empty string');
      assert.strictEqual(options[0].label, 'root', 'root namespace labeled as "root"');
    });

    test('does not duplicate userRootNamespace if already in list', function (assert) {
      this.namespaceService.accessibleNamespaces = ['admin', 'admin/team1'];
      this.authStub.value({ userRootNamespace: 'admin' });

      const options = this.namespaceService.getOptions();

      const adminOptions = options.filter((o) => o.path === 'admin');
      assert.strictEqual(adminOptions.length, 1, 'does not duplicate user root namespace');
    });

    test('returns fallback root namespace when accessibleNamespaces is empty', function (assert) {
      this.namespaceService.accessibleNamespaces = [];
      this.authStub.value({ userRootNamespace: '' });

      const options = this.namespaceService.getOptions();

      assert.deepEqual(options, [{ path: '', label: 'root' }], 'adds root namespace as fallback');
    });

    test('handles null accessibleNamespaces', function (assert) {
      this.namespaceService.accessibleNamespaces = null;
      this.authStub.value({ userRootNamespace: '' });

      const options = this.namespaceService.getOptions();

      assert.deepEqual(options, [{ path: '', label: 'root' }], 'adds root namespace as fallback');
    });

    test('handles undefined userRootNamespace', function (assert) {
      this.authStub.value({ userRootNamespace: undefined });

      const options = this.namespaceService.getOptions();

      assert.deepEqual(options, [{ path: '', label: 'root' }], 'adds root namespace as fallback');
    });

    test('handles null userRootNamespace', function (assert) {
      this.authStub.value({ userRootNamespace: null });

      const options = this.namespaceService.getOptions();

      assert.deepEqual(options, [{ path: '', label: 'root' }], 'adds root namespace as fallback');
    });
  });
});
