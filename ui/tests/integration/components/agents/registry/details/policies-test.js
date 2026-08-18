/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | agents/registry/details/policies', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = {
      agent: {
        id: '1',
        display_name: 'agent-1',
        entity_id: 'entity-1',
        ceiling_policy: ['ceiling-policy'],
      },
      entity: {
        id: 'entity-1',
        name: 'entity-1',
        policies: ['policy-a', 'policy-b'],
      },
      groups: [
        {
          id: 'group-1',
          name: 'group-1',
          policies: ['group-policy'],
        },
      ],
      aggregatePolicy: null,
    };

    this.data.aggregatePolicy = {
      policy: {
        'identity/*': ['read'],
        'kv/data/appB/*': ['create', 'delete', 'list', 'read', 'update'],
      },
      policyString: `path "identity/*" {
  capabilities = ["read"]
}

path "kv/data/appB/*" {
  capabilities = ["update", "create", "list", "delete", "read"]
}`,
    };
  });

  test('it renders ceiling policy table and snippets tabs', async function (assert) {
    await render(hbs`
      <Agents::Registry::Details::Policies @data={{this.data}} />
    `);

    assert.dom(GENERAL.table('policy-list')).includesText('Policy list', 'renders table caption');
    assert.dom(GENERAL.table('policy-list')).includesText('Ceiling', 'renders ceiling level rows');
    assert
      .dom(GENERAL.table('policy-list'))
      .includesText('agent-1', 'renders display name as ceiling policy source');
    assert.dom(GENERAL.table('policy-list')).includesText('ceiling-policy', 'renders ceiling policy row');
    assert.dom(GENERAL.table('policy-list')).includesText('Entity', 'renders entity level rows');
    assert.dom(GENERAL.table('policy-list')).includesText('entity-1', 'renders entity source');
    assert.dom(GENERAL.table('policy-list')).includesText('policy-a', 'renders policy-a row');
    assert.dom(GENERAL.table('policy-list')).includesText('policy-b', 'renders policy-b row');
    assert.dom(GENERAL.table('policy-list')).includesText('Group', 'renders group level rows');
    assert.dom(GENERAL.table('policy-list')).includesText('group-1', 'renders group source');
    assert.dom(GENERAL.table('policy-list')).includesText('group-policy', 'renders group policy row');

    assert.dom(GENERAL.hdsTab()).exists({ count: 3 });
    assert
      .dom(GENERAL.hdsTab('terraform'))
      .hasAttribute('aria-selected', 'true', 'terraform tab is selected by default');

    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .includesText(
        'resource "vault_policy" "<local identifier>"',
        'renders terraform snippet on default tab'
      );

    await click(GENERAL.hdsTab('allowed-actions'));
    assert
      .dom(GENERAL.cardContainer('policy-allowed-actions'))
      .includesText('Allowed actions', 'renders allowed actions content when tab is clicked');
    assert.dom('[data-test-path-entry="identity/*"]').exists('renders first allowed path');

    await click(GENERAL.hdsTab('terraform'));
    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .includesText('resource "vault_policy" "<local identifier>"', 'renders terraform snippet');
    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .includesText('name = "agent-1"', 'renders policy name in terraform snippet');

    await click(GENERAL.hdsTab('cli'));
    assert
      .dom(GENERAL.fieldByAttr('cli'))
      .includesText('vault policy write agent-1', 'renders CLI snippet command');
  });

  test('it renders entity and group policies when ceiling policy is empty', async function (assert) {
    this.data = {
      ...this.data,
      agent: {
        ...this.data.agent,
        ceiling_policy: [],
      },
    };

    await render(hbs`
      <Agents::Registry::Details::Policies @data={{this.data}} />
    `);

    assert
      .dom(GENERAL.table('policy-list'))
      .doesNotIncludeText('ceiling-policy', 'does not render ceiling policy row');
    assert.dom(GENERAL.table('policy-list')).includesText('Entity', 'renders entity level rows');
    assert.dom(GENERAL.table('policy-list')).includesText('entity-1', 'renders entity name as policy source');
    assert.dom(GENERAL.table('policy-list')).includesText('policy-a', 'renders first entity policy');
    assert.dom(GENERAL.table('policy-list')).includesText('group-policy', 'renders group policy row');
  });
});
