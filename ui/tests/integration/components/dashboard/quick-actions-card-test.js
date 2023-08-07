/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { fillIn } from '@ember/test-helpers';

import { selectChoose } from 'ember-power-select/test-support/helpers';

// import SELECTORS from 'vault/tests/helpers/components/dashboard/quick-actions-card';

module('Integration | Component | dashboard/quick-actions-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kubernetes_f3400dee',
        path: 'kubernetes-test/',
        type: 'kubernetes',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'database_f3400dee',
        path: 'database-test/',
        type: 'database',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'pki_i1234dd',
        path: 'apki-test/',
        type: 'pki',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'secrets_j2350ii',
        path: 'secrets-test/',
        type: 'kv',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'nomad_123hh',
        path: 'nomad/',
        type: 'nomad',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'pki_f3400dee',
        path: 'pki-0-test/',
        type: 'pki',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'pki_i1234dd',
        path: 'pki-1-test/',
        description: 'pki-1-path-description',
        type: 'pki',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'secrets_j2350ii',
        path: 'secrets-1-test/',
        type: 'kv',
      },
    });

    this.secretsEngines = this.store.peekAll('secret-engine', {});

    this.renderComponent = () => {
      return render(hbs`<Dashboard::QuickActionsCard @secretsEngines={{this.secretsEngines}} />`);
    };
  });

  test('it should show quick action empty state if no engine is selected', async function (assert) {
    await this.renderComponent();
    assert.dom('.title').hasText('Quick actions');
    assert.dom('[data-test-secrets-engines-select]').exists({ count: 1 });
    assert.dom('[data-test-component="empty-state"]').exists({ count: 1 });
  });

  test('it should show correct actions for pki', async function (assert) {
    await this.renderComponent();
    await selectChoose('.search-select', 'pki-0-test');
    await fillIn('[data-test-select="action-select"]', 'Issue certificate');
    assert.dom('[data-test-component="empty-state"]').doesNotExist();
    await fillIn('[data-test-select="action-select"]', 'Issue certificate');
    assert.dom('[data-test-button="Issue leaf certificate"]').exists({ count: 1 });
    assert.dom('[data-test-search-select-params-title]').hasText('Role to use');
    await fillIn('[data-test-select="action-select"]', 'View certificate');
    assert.dom('[data-test-search-select-params-title]').hasText('Certificate serial number');
    assert.dom('[data-test-button="View certificate"]').exists({ count: 1 });
    await fillIn('[data-test-select="action-select"]', 'View issuer');
    assert.dom('[data-test-search-select-params-title]').hasText('Issuer');
    assert.dom('[data-test-button="View issuer"]').exists({ count: 1 });
  });
  test('it should show correct actions for database', async function (assert) {
    await this.renderComponent();
    await selectChoose('.search-select', 'database-test');
    assert.dom('[data-test-component="empty-state"]').doesNotExist();
    await fillIn('[data-test-select="action-select"]', 'Generate credentials for database');
    assert.dom('[data-test-search-select-params-title]').hasText('Role to use');
    assert.dom('[data-test-button="Generate credentials"]').exists({ count: 1 });
  });
  test('it should show correct actions for kv', async function (assert) {
    await this.renderComponent();
    await selectChoose('.search-select', 'secrets-1-test');
    assert.dom('[data-test-component="empty-state"]').doesNotExist();
    await fillIn('[data-test-select="action-select"]', 'Find KV secrets');
    assert.dom('[data-test-search-select-params-title]').hasText('Secret Path');
    assert.dom('[data-test-button="Read secrets"]').exists({ count: 1 });
  });
});
