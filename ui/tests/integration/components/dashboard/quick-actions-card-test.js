/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, findAll, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { fillIn } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';

import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';
import { setRunOptions } from 'ember-a11y-testing/test-support';

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
        path: 'kv-v2-test/',
        options: {
          version: 2,
        },
        type: 'kv',
      },
    });

    this.secretsEngines = this.store.peekAll('secret-engine', {});

    this.renderComponent = () => {
      return render(hbs`<Dashboard::QuickActionsCard @secretsEngines={{this.secretsEngines}} />`);
    };

    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  test('it should show quick action empty state if no engine is selected', async function (assert) {
    await this.renderComponent();
    assert.dom('.title').hasText('Quick actions');
    assert.dom(DASHBOARD.searchSelect('secrets-engines')).exists({ count: 1 });
    assert.dom(DASHBOARD.emptyState('no-mount-selected')).exists({ count: 1 });
  });

  test('it should show correct actions for pki', async function (assert) {
    await this.renderComponent();
    await selectChoose(DASHBOARD.searchSelect('secrets-engines'), 'pki-0-test');
    await fillIn(DASHBOARD.selectEl, 'Issue certificate');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    await fillIn(DASHBOARD.selectEl, 'Issue certificate');
    assert.dom(DASHBOARD.actionButton('Issue leaf certificate')).exists({ count: 1 });
    assert.dom(DASHBOARD.subtitle('param')).hasText('Role to use');
    await fillIn(DASHBOARD.selectEl, 'View certificate');
    assert.dom(DASHBOARD.subtitle('param')).hasText('Certificate serial number');
    assert.dom(DASHBOARD.actionButton('View certificate')).exists({ count: 1 });
    await fillIn(DASHBOARD.selectEl, 'View issuer');
    assert.dom(DASHBOARD.subtitle('param')).hasText('Issuer');
    assert.dom(DASHBOARD.actionButton('View issuer')).exists({ count: 1 });
  });
  test('it should show correct actions for database', async function (assert) {
    await this.renderComponent();
    await selectChoose(DASHBOARD.searchSelect('secrets-engines'), 'database-test');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    await fillIn(DASHBOARD.selectEl, 'Generate credentials for database');
    assert.dom(DASHBOARD.subtitle('param')).hasText('Role to use');
    assert.dom(DASHBOARD.actionButton('Generate credentials')).exists({ count: 1 });
  });
  test('it should show correct actions for kv', async function (assert) {
    await this.renderComponent();
    await click('[data-test-component="search-select"]#secrets-engines-select .ember-basic-dropdown-trigger');
    assert.strictEqual(
      findAll('li.ember-power-select-option').length,
      5,
      'renders only kv v2, pki and db engines'
    );
    await selectChoose(DASHBOARD.searchSelect('secrets-engines'), 'kv-v2-test');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    await fillIn(DASHBOARD.selectEl, 'Find KV secrets');
    assert.dom(DASHBOARD.kvSearchSelect).exists('Shows option to search fo KVv2 secret');
    assert.dom(DASHBOARD.actionButton('Read secrets')).exists({ count: 1 });
  });
});
