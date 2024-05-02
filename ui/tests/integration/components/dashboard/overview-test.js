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
import { LICENSE_START } from 'vault/mirage/handlers/clients';

module('Integration | Component | dashboard/overview', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

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
    this.server.get('sys/internal/counters/config', function () {
      return {
        request_id: 'some-config-id',
        data: {
          billing_start_timestamp: LICENSE_START.toISOString(),
        },
      };
    });
  });

  test('it should show dashboard empty states', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1';
    this.isRootNamespace = true;
    await render(
      hbs`
        <Dashboard::Overview
          @version={{this.version}}
          @isRootNamespace={{this.isRootNamespace}}
          @refreshModel={{this.refreshModel}} />
      `
    );
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

  test('it should hide client count and replication card on community', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1';
    this.isRootNamespace = true;

    await render(
      hbs`
        <Dashboard::Overview
          @secretsEngines={{this.secretsEngines}}
          @vaultConfiguration={{this.vaultConfiguration}}
          @replication={{this.replication}}
          @version={{this.version}}
          @isRootNamespace={{this.isRootNamespace}}
          @refreshModel={{this.refreshModel}} />
      `
    );

    assert.dom(DASHBOARD.cardHeader('Vault version')).exists();
    assert.dom(DASHBOARD.cardName('secrets-engines')).exists();
    assert.dom(DASHBOARD.cardName('learn-more')).exists();
    assert.dom(DASHBOARD.cardName('quick-actions')).exists();
    assert.dom(DASHBOARD.cardName('configuration-details')).exists();
    assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
  });

  test('it should show client count on enterprise w/ license', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.version.type = 'enterprise';
    this.license = {
      autoloaded: {
        license_id: '7adbf1f4-56ef-35cd-3a6c-50ef2627865d',
      },
    };

    await render(
      hbs`
      <Dashboard::Overview
      @secretsEngines={{this.secretsEngines}}
      @vaultConfiguration={{this.vaultConfiguration}}
      @replication={{this.replication}}
      @version={{this.version}}
      @isRootNamespace={{true}}
      @refreshModel={{this.refreshModel}} />`
    );
    assert.dom(DASHBOARD.cardHeader('Vault version')).exists();
    assert.dom(DASHBOARD.cardName('secrets-engines')).exists();
    assert.dom(DASHBOARD.cardName('learn-more')).exists();
    assert.dom(DASHBOARD.cardName('quick-actions')).exists();
    assert.dom(DASHBOARD.cardName('configuration-details')).exists();
    assert.dom(DASHBOARD.cardName('client-count')).exists();
  });

  test('it should hide replication on enterprise not on root namespace', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.version.type = 'enterprise';
    this.isRootNamespace = false;

    await render(
      hbs`
      <Dashboard::Overview
        @version={{this.version}}
        @isRootNamespace={{this.isRootNamespace}}
        @secretsEngines={{this.secretsEngines}}
        @vaultConfiguration={{this.vaultConfiguration}}
        @replication={{this.replication}}
        @refreshModel={{this.refreshModel}} />`
    );

    assert.dom(DASHBOARD.cardHeader('Vault version')).exists();
    assert.dom('[data-test-badge-namespace]').exists();
    assert.dom(DASHBOARD.cardName('secrets-engines')).exists();
    assert.dom(DASHBOARD.cardName('learn-more')).exists();
    assert.dom(DASHBOARD.cardName('quick-actions')).exists();
    assert.dom(DASHBOARD.cardName('configuration-details')).exists();
    assert.dom(DASHBOARD.cardName('replication')).doesNotExist();
    assert.dom(DASHBOARD.cardName('client-count')).doesNotExist();
  });

  module('learn more card', function () {
    test('shows the learn more card on community', async function (assert) {
      await render(
        hbs`<Dashboard::Overview @secretsEngines={{this.secretsEngines}} @vaultConfiguration={{this.vaultConfiguration}} @replication={{this.replication}} @refreshModel={{this.refreshModel}} />`
      );

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
      this.version = this.owner.lookup('service:version');
      this.version.version = '1.13.1+ent';
      this.version.type = 'enterprise';
      this.version.features = [
        'Performance Replication',
        'DR Replication',
        'Namespaces',
        'Transform Secrets Engine',
      ];
      this.isRootNamespace = true;
      this.license = {
        autoloaded: {
          license_id: '7adbf1f4-56ef-35cd-3a6c-50ef2627865d',
        },
      };
      await render(
        hbs`
          <Dashboard::Overview
            @version={{this.version}}
            @isRootNamespace={{this.isRootNamespace}}
            @secretsEngines={{this.secretsEngines}}
            @vaultConfiguration={{this.vaultConfiguration}}
            @replication={{this.replication}}
            @refreshModel={{this.refreshModel}} />
        `
      );
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
