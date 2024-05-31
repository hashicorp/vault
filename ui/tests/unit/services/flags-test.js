/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';

const ACTIVATED_FLAGS_RESPONSE = {
  data: {
    activated: ['secrets-sync'],
    unactivated: [],
  },
};

const FEATURE_FLAGS_RESPONSE = {
  feature_flags: ['VAULT_CLOUD_ADMIN_NAMESPACE'],
};

module('Unit | Service | flags', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:flags');
    this.version = this.owner.lookup('service:version');
    this.permissions = this.owner.lookup('service:permissions');
  });

  test('it loads with defaults', function (assert) {
    assert.deepEqual(this.service.featureFlags, [], 'Flags are empty until fetched');
    assert.deepEqual(this.service.activatedFlags, [], 'Activated flags are empty until fetched');
  });

  module('#fetchActivatedFlags', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
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

    test('it does not call activation-flags endpoint if the cluster is OSS', async function (assert) {
      this.version.type = 'community';

      this.server.get(
        'sys/activation-flags',
        () =>
          new Error(
            'uh oh! a request was made to sys/activation-flags, this should not happen for community versions'
          )
      );

      await this.service.fetchActivatedFlags();
      assert.deepEqual(this.service.activatedFlags, [], 'Activated flags are empty');
    });
  });

  module('#fetchFeatureFlags', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
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
        FEATURE_FLAGS_RESPONSE.feature_flags,
        'Feature flags are fetched and set'
      );
    });
  });

  module('#hvdManagedNamespaceRoot', function () {
    test('it returns null when flag is not present', function (assert) {
      assert.strictEqual(this.service.hvdManagedNamespaceRoot, null);
    });

    test('it returns the namespace root when flag is present', function (assert) {
      this.service.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      assert.strictEqual(
        this.service.hvdManagedNamespaceRoot,
        'admin',
        'Managed namespace is admin when flag present'
      );

      this.service.featureFlags = ['SOMETHING_ELSE'];
      assert.strictEqual(
        this.service.hvdManagedNamespaceRoot,
        null,
        'Flags were overwritten and root namespace is null again'
      );
    });
  });

  module('#secretsSyncIsActivated', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
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

  module('#showSecretsSync', function () {
    test('it returns false when version is community', function (assert) {
      this.version.type = 'community';
      assert.false(this.service.showSecretsSync);
    });

    module('isHvdManaged', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
        this.service.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      });

      test('it returns true when not activated', function (assert) {
        this.service.activatedFlags = [];
        assert.true(this.service.showSecretsSync);
      });

      test('it returns true when activated', function (assert) {
        this.service.activatedFlags = ACTIVATED_FLAGS_RESPONSE.data.activated;
        assert.true(this.service.showSecretsSync);
      });
    });

    module('is Enterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });

      test('it returns false when not on license ', function (assert) {
        this.version.features = ['replication'];
        assert.false(this.service.showSecretsSync);
      });

      module('no permissions to sys/sync', function (hooks) {
        hooks.beforeEach(function () {
          this.version.features = ['Secrets Sync'];
          const hasNavPermission = sinon.stub(this.permissions, 'hasNavPermission');
          hasNavPermission.returns(false);
        });

        test('it returns false when activated ', function (assert) {
          this.service.activatedFlags = ACTIVATED_FLAGS_RESPONSE.data.activated;
          assert.false(this.service.showSecretsSync);
        });

        test('it returns true when not activated ', function (assert) {
          // the activate endpoint is located at a different path than sys/sync.
          // the expected UX experience: if the feature is not activated, regardless of permissions
          // the user should see the landing page and a banner that tells them to either have an admin activate the feature or activate it themselves
          this.service.activatedFlags = [];
          assert.true(this.service.showSecretsSync);
        });
      });

      module('user has permissions to sys/sync', function (hooks) {
        hooks.beforeEach(function () {
          this.version.features = ['Secrets Sync'];
          const hasNavPermission = sinon.stub(this.permissions, 'hasNavPermission');
          hasNavPermission.returns(true);
        });

        test('it returns true when activated ', function (assert) {
          this.service.activatedFlags = ACTIVATED_FLAGS_RESPONSE.data.activated;
          assert.true(this.service.showSecretsSync);
        });

        test('it returns true when not activated ', function (assert) {
          this.service.activatedFlags = [];
          assert.true(this.service.showSecretsSync);
        });
      });
    });
  });
});
