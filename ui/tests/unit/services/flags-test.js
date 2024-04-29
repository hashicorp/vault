/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

const ACTIVATED_FLAGS_RESPONSE = {
  data: {
    activated: ['secrets-sync'],
    unactivated: [],
  },
};

const FEATURE_FLAGS_RESPONSE = {
  data: {
    feature_flags: ['VAULT_CLOUD_ADMIN_NAMESPACE'],
  },
};

module('Unit | Service | flags', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:flags');
  });

  test('it loads with defaults', function (assert) {
    assert.deepEqual(this.service.featureFlags, [], 'Flags are empty until fetched');
    assert.deepEqual(this.service.activatedFlags, [], 'Activated flags are empty until fetched');
  });

  module('#fetchActivatedFlags', function (hooks) {
    hooks.beforeEach(function () {
      this.owner.lookup('service:version').type = 'enterprise';
    });

    test('it returns activated flags', async function (assert) {
      assert.expect(2);

      this.server.get('sys/activation-flags', () => {
        assert.true(true, 'GET request made to activation-flags endpoint');
        return ACTIVATED_FLAGS_RESPONSE;
      });

      await this.service.fetchActivatedFlags();
      assert.deepEqual(
        this.service.activatedFlags,
        ACTIVATED_FLAGS_RESPONSE.data.activated,
        'Activated flags are fetched and set'
      );
    });

    test('it returns an empty array if no flags are activated', async function (assert) {
      this.server.get('sys/activation-flags', () => {
        return {
          data: {
            activated: [],
            unactivated: [],
          },
        };
      });

      await this.service.fetchActivatedFlags();
      assert.deepEqual(this.service.activatedFlags, [], 'Activated flags are empty');
    });

    test('it returns an empty array if the cluster is OSS', async function (assert) {
      this.owner.lookup('service:version').type = 'community';

      await this.service.fetchActivatedFlags();
      assert.deepEqual(this.service.activatedFlags, [], 'Activated flags are empty');
    });
  });

  module('#fetchFeatureFlags', function (hooks) {
    hooks.beforeEach(function () {
      this.owner.lookup('service:version').type = 'enterprise';
    });

    test('it returns feature flags', async function (assert) {
      assert.expect(2);

      this.server.get('sys/internal/ui/feature-flags', () => {
        assert.true(true, 'GET request made to feature-flags endpoint');
        return FEATURE_FLAGS_RESPONSE;
      });

      await this.service.fetchFeatureFlags();
      assert.deepEqual(
        this.service.featureFlags,
        FEATURE_FLAGS_RESPONSE.data.feature_flags,
        'Feature flags are fetched and set'
      );
    });
  });

  module('#managedNamespaceRoot', function () {
    test('it returns null when flag is not present', function (assert) {
      assert.strictEqual(this.service.managedNamespaceRoot, null);
    });

    test('it returns the namespace root when flag is present', function (assert) {
      this.service.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      assert.strictEqual(
        this.service.managedNamespaceRoot,
        'admin',
        'Managed namespace is admin when flag present'
      );

      this.service.featureFlags = ['SOMETHING_ELSE'];
      assert.strictEqual(
        this.service.managedNamespaceRoot,
        null,
        'Flags were overwritten and root namespace is null again'
      );
    });
  });

  module('#secretsSyncActivated', function (hooks) {
    hooks.beforeEach(function () {
      this.owner.lookup('service:version').type = 'enterprise';
      this.service.activatedFlags = ACTIVATED_FLAGS_RESPONSE.data.activated;
    });

    test('it returns true when secrets sync is activated', function (assert) {
      assert.true(this.service.secretsSyncIsActivated);
    });

    test('it returns false when secrets sync is not activated', function (assert) {
      this.service.activatedFlags = [];
      assert.false(this.service.secretsSyncIsActivated);
    });
  });
});
