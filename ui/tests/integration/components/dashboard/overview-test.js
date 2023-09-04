/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';

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
    assert.dom('[data-test-dashboard-version-header]').exists();
    assert.dom(SELECTORS.cardName('secrets-engines')).exists();
    assert.dom(SELECTORS.emptyState('secrets-engines')).exists();
    assert.dom(SELECTORS.cardName('learn-more')).exists();
    assert.dom(SELECTORS.cardName('quick-actions')).exists();
    assert.dom(SELECTORS.emptyState('quick-actions')).exists();
    assert.dom(SELECTORS.cardName('configuration-details')).doesNotExist();
    assert.dom(SELECTORS.cardName('replication')).doesNotExist();
    assert.dom(SELECTORS.cardName('client-count')).doesNotExist();
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

    assert.dom('[data-test-dashboard-version-header]').exists();
    assert.dom(SELECTORS.cardName('secrets-engines')).exists();
    assert.dom(SELECTORS.cardName('learn-more')).exists();
    assert.dom(SELECTORS.cardName('quick-actions')).exists();
    assert.dom(SELECTORS.cardName('configuration-details')).exists();
    assert.dom(SELECTORS.cardName('replication')).doesNotExist();
    assert.dom(SELECTORS.cardName('client-count')).doesNotExist();
  });

  test('it should show client count and replication card on enterprise w/ license + namespace enabled', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
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
      @license={{this.license}}
      @refreshModel={{this.refreshModel}} />`
    );

    assert.dom('[data-test-dashboard-version-header]').exists();
    assert.dom(SELECTORS.cardName('secrets-engines')).exists();
    assert.dom(SELECTORS.cardName('learn-more')).exists();
    assert.dom(SELECTORS.cardName('quick-actions')).exists();
    assert.dom(SELECTORS.cardName('configuration-details')).exists();
    assert.dom(SELECTORS.cardName('replication')).exists();
    assert.dom(SELECTORS.cardName('client-count')).exists();
  });

  test('it should hide client count on enterprise w/o license ', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.isRootNamespace = true;

    await render(
      hbs`
      <Dashboard::Overview
        @secretsEngines={{this.secretsEngines}}
        @vaultConfiguration={{this.vaultConfiguration}}
        @replication={{this.replication}}
        @version={{this.version}}
        @isRootNamespace={{this.isRootNamespace}}
        @refreshModel={{this.refreshModel}}
      />`
    );

    assert.dom('[data-test-dashboard-version-header]').exists();
    assert.dom('[data-test-badge-namespace]').exists();
    assert.dom(SELECTORS.cardName('secrets-engines')).exists();
    assert.dom(SELECTORS.cardName('learn-more')).exists();
    assert.dom(SELECTORS.cardName('quick-actions')).exists();
    assert.dom(SELECTORS.cardName('configuration-details')).exists();
    assert.dom(SELECTORS.cardName('replication')).exists();
    assert.dom(SELECTORS.cardName('client-count')).doesNotExist();
  });

  test('it should hide replication on enterprise not on root namespace', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.isRootNamespace = false;
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
        @license={{this.license}}
        @refreshModel={{this.refreshModel}} />`
    );

    assert.dom('[data-test-dashboard-version-header]').exists();
    assert.dom('[data-test-badge-namespace]').exists();
    assert.dom(SELECTORS.cardName('secrets-engines')).exists();
    assert.dom(SELECTORS.cardName('learn-more')).exists();
    assert.dom(SELECTORS.cardName('quick-actions')).exists();
    assert.dom(SELECTORS.cardName('configuration-details')).exists();
    assert.dom(SELECTORS.cardName('replication')).doesNotExist();
    assert.dom(SELECTORS.cardName('client-count')).exists();
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
            @license={{this.license}}
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
