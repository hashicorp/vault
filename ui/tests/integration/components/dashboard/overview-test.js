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
    this.refreshModel = () => {};
    this.model = {
      secretsEngines: this.secretsEngines,
      vaultConfiguration: {
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
      },
    };
  });

  test('it should hide client count and replication card on community', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1';
    this.namespace = this.owner.lookup('service:namespace');
    this.namespace.isRootNamespace = true;
    this.model.version = this.version;
    this.model.isRootNamespace = this.namespace;

    await render(hbs`<Dashboard::Overview @model={{this.model}} @refreshModel={{this.refreshModel}} />`);

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
    this.model.version = this.version;
    this.namespace = this.owner.lookup('service:namespace');
    this.namespace.isRootNamespace = true;
    this.model.isRootNamespace = this.namespace;

    this.license = {
      autoloaded: {
        license_id: '7adbf1f4-56ef-35cd-3a6c-50ef2627865d',
      },
    };
    this.model.license = this.license;

    await render(hbs`<Dashboard::Overview @model={{this.model}} @refreshModel={{this.refreshModel}} />`);

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
    this.model.version = this.version;
    this.namespace = this.owner.lookup('service:namespace');
    this.namespace.isRootNamespace = true;
    this.model.isRootNamespace = this.namespace;

    await render(hbs`<Dashboard::Overview @model={{this.model}} @refreshModel={{this.refreshModel}} />`);

    assert.dom('[data-test-dashboard-version-header]').exists();
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
    this.model.version = this.version;
    this.namespace = this.owner.lookup('service:namespace');
    this.namespace.path = 'hello';
    this.model.isRootNamespace = false;
    this.license = {
      autoloaded: {
        license_id: '7adbf1f4-56ef-35cd-3a6c-50ef2627865d',
      },
    };
    this.model.license = this.license;
    await render(hbs`<Dashboard::Overview @model={{this.model}} @refreshModel={{this.refreshModel}} />`);

    assert.dom('[data-test-dashboard-version-header]').exists();
    assert.dom(SELECTORS.cardName('secrets-engines')).exists();
    assert.dom(SELECTORS.cardName('learn-more')).exists();
    assert.dom(SELECTORS.cardName('quick-actions')).exists();
    assert.dom(SELECTORS.cardName('configuration-details')).exists();
    assert.dom(SELECTORS.cardName('replication')).doesNotExist();
    assert.dom(SELECTORS.cardName('client-count')).exists();
  });

  module('learn more card', function () {
    test('shows the learn more card on community', async function (assert) {
      await render(hbs`<Dashboard::Overview @model={{this.model}} @refreshModel={{this.refreshModel}} />`);

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
      this.model.version = this.version;
      this.namespace = this.owner.lookup('service:namespace');
      this.namespace.isRootNamespace = true;
      this.model.isRootNamespace = this.namespace;
      this.license = {
        autoloaded: {
          license_id: '7adbf1f4-56ef-35cd-3a6c-50ef2627865d',
        },
      };
      this.model.license = this.license;
      await render(hbs`<Dashboard::Overview @model={{this.model}} @refreshModel={{this.refreshModel}} />`);
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
