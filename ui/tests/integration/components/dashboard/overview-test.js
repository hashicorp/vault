/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';

module('Integration | Component | dashboard/overview', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flags = this.owner.lookup('service:flags');
    this.namespace = this.owner.lookup('service:namespace');
    this.permissions = this.owner.lookup('service:permissions');
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.version.type = 'enterprise';
    this.isRootNamespace = true;
    this.replication = {
      dr: {
        clusterId: '123',
        state: 'running',
      },
      performance: {
        clusterId: 'abc-1',
        state: 'running',
        isPrimary: true,
      },
    };
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_f3400dee',
        path: 'kv-test/',
        type: 'kv',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_f3300dee',
        path: 'kv-1/',
        type: 'kv',
      },
    });
    this.secretsEngines = this.store.peekAll('secret-engine', {});
    this.vaultConfiguration = {
      api_addr: 'http://127.0.0.1:8200',
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      listeners: [
        {
          config: {
            address: '127.0.0.1:8200',
            tls_disable: 1,
          },
          type: 'tcp',
        },
      ],
    };
    this.refreshModel = () => {};
    this.renderComponent = async () => {
      return await render(
        hbs`
        <Dashboard::Overview
          @secretsEngines={{this.secretsEngines}}
          @vaultConfiguration={{this.vaultConfiguration}}
          @replication={{this.replication}}
          @version={{this.version}}
          @isRootNamespace={{this.isRootNamespace}}
          @refreshModel={{this.refreshModel}} 
          @replicationUpdatedAt={{this.replicationUpdatedAt}}
          />
      `
      );
    };
  });

  test('it should show dashboard empty states in root namespace', async function (assert) {
    this.version.version = '1.13.1';
    this.secretsEngines = null;
    this.replication = null;
    this.vaultConfiguration = null;
    await this.renderComponent();
    assert.dom(DASHBOARD.cardHeader('Vault version')).exists();
    assert.dom(DASHBOARD.cardName('secrets-engines')).exists();
    assert.dom(DASHBOARD.emptyState('secrets-engines')).exists();
    assert.dom(DASHBOARD.cardName('learn-more')).exists();
    assert.dom(DASHBOARD.cardName('quick-actions')).exists();
    assert.dom(DASHBOARD.emptyState('quick-actions')).exists();
    assert.dom(DASHBOARD.cardName('configuration-details')).doesNotExist();
    assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
  });

  module('client count and replication card', function () {
    test('it should hide cards on community in root namespace', async function (assert) {
      this.version.version = '1.13.1';
      this.version.type = 'community';
      this.server.get(
        'sys/internal/counters/activity',
        () => new Error('uh oh! a request was made to sys/internal/counters/activity')
      );
      await this.renderComponent();

      assert.dom(DASHBOARD.cardHeader('Vault version')).exists();
      assert.dom(DASHBOARD.cardName('secrets-engines')).exists();
      assert.dom(DASHBOARD.cardName('learn-more')).exists();
      assert.dom(DASHBOARD.cardName('quick-actions')).exists();
      assert.dom(DASHBOARD.cardName('configuration-details')).exists();
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
      assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
    });

    test('it should hide cards on enterprise if permission but not in root namespace', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };
      this.isRootNamespace = false;
      await this.renderComponent();
      assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    });

    test('it should show cards on enterprise if has permission and in root namespace', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };
      await this.renderComponent();
      assert.dom(DASHBOARD.cardHeader('Vault version')).exists();
      assert.dom(DASHBOARD.cardName('secrets-engines')).exists();
      assert.dom(DASHBOARD.cardName('learn-more')).exists();
      assert.dom(DASHBOARD.cardName('quick-actions')).exists();
      assert.dom(DASHBOARD.cardName('configuration-details')).exists();
      assert.dom(DASHBOARD.cardName('client-count')).exists();
      assert.dom(DASHBOARD.cardName('replication')).exists();
    });

    test('it should show client count on enterprise in admin namespace when running a managed mode', async function (assert) {
      this.permissions.exactPaths = {
        'admin/sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'admin/sys/replication/status': {
          capabilities: ['read'],
        },
      };

      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = 'admin';
      this.isRootNamespace = false;

      await this.renderComponent();

      assert.dom(DASHBOARD.cardName('client-count')).exists();
    });

    test('it should hide client count on enterprise in any other namespace when running a managed mode', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };

      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = 'groceries';
      this.isRootNamespace = false;

      await this.renderComponent();

      assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
    });

    test('it should hide cards on enterprise in root namespace but no permission', async function (assert) {
      await this.renderComponent();
      assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    });

    test('it should hide cards on enterprise if no permission and not in root namespace', async function (assert) {
      this.isRootNamespace = false;
      await this.renderComponent();
      assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    });

    test('it should hide client count on enterprise in root namespace if no activity permission', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['deny'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };

      await this.renderComponent();
      assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
      assert.dom(DASHBOARD.cardName('replication')).exists();
    });

    test('it should hide replication on enterprise in root namespace if no replication status permission', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['deny'],
        },
      };

      await this.renderComponent();
      assert.dom(DASHBOARD.cardName('client-count')).exists();
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    });

    test('it should hide replication on enterprise if has permission and in root namespace but is empty', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };
      this.replication = {};
      await this.renderComponent();
      assert.dom(DASHBOARD.cardName('client-count')).exists();
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    });
  });

  module('learn more card', function () {
    test('shows the learn more card on community', async function (assert) {
      this.version.version = '1.13.1';
      this.version.type = 'community';
      await this.renderComponent();

      assert.dom('[data-test-learn-more-title]').hasText('Learn more');
      assert
        .dom('[data-test-learn-more-subtext]')
        .hasText(
          'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
        );
      assert.dom('[data-test-learn-more-links] a').exists({ count: 3 });
      assert
        .dom('[data-test-feedback-form]')
        .hasText("Don't see what you're looking for on this page? Let us know via our feedback form .");
    });
    test('shows the learn more card on enterprise', async function (assert) {
      this.version.features = [
        'Performance Replication',
        'DR Replication',
        'Namespaces',
        'Transform Secrets Engine',
      ];
      await this.renderComponent();
      assert.dom('[data-test-learn-more-title]').hasText('Learn more');
      assert
        .dom('[data-test-learn-more-subtext]')
        .hasText(
          'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
        );
      assert.dom('[data-test-learn-more-links] a').exists({ count: 4 });
      assert
        .dom('[data-test-feedback-form]')
        .hasText("Don't see what you're looking for on this page? Let us know via our feedback form .");
    });
  });
});
