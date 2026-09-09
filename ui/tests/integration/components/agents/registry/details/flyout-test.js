/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | agents/registry/details/flyout', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClose = sinon.spy();
    this.api = this.owner.lookup('service:api');

    this.data = {
      agent: {
        id: 'agent-1',
        display_name: 'test-agent-1',
        entity_id: 'entity-1',
        owner: 'test-owner',
        no_default_ceiling_policy: false,
        creation_time: '2026-01-01T10:00:00Z',
        last_updated_time: '2026-01-02T11:00:00Z',
      },
      entity: {
        id: 'entity-1',
        name: 'test-entity-1',
        disabled: false,
        creation_time: '2026-01-01T10:00:00Z',
        last_update_time: '2026-01-02T11:00:00Z',
        policies: ['policy-a'],
        aliases: [
          {
            id: 'alias-1',
            name: 'test-alias-1',
            mount_accessor: 'accessor-1',
            mount_path: 'token/',
            mount_type: 'token',
            creation_time: '2026-01-03T12:00:00Z',
            last_update_time: '2026-01-04T13:00:00Z',
          },
        ],
      },
      groups: [],
      aggregatePolicy: {
        policy: { '*': ['read'] },
        policyString: 'path "*" {\n    capabilities = ["read"]\n}',
      },
    };
    this.dataWithoutAliases = {
      ...this.data,
      agent: {
        ...this.data.agent,
        id: 'agent-2',
        display_name: 'test-agent-2',
        entity_id: 'entity-2',
      },
      entity: {
        ...this.data.entity,
        id: 'entity-2',
        name: 'test-entity-2',
        aliases: [],
      },
    };

    this.renderComponent = () =>
      render(hbs`<Agents::Registry::Details::Flyout @data={{this.data}} @onClose={{this.onClose}} />`);
  });

  test('it renders flyout header and tabs for agent with aliases', async function (assert) {
    await this.renderComponent();

    assert.dom('.hds-flyout').exists('renders the flyout');
    assert.dom('.hds-flyout__header').containsText('test-agent-1', 'renders the agent name in the header');
    assert.dom(GENERAL.hdsTab('agent')).exists('renders the agent details tab');
    assert.dom(GENERAL.hdsTab('policies')).exists('renders the policies tab');
    assert.dom(GENERAL.hdsTab('entity')).exists('renders the entity details tab');
    assert.dom(GENERAL.hdsTab('alias')).exists('renders the alias details tab when aliases exist');
  });

  test('it hides the policies tab when the agent has no policies', async function (assert) {
    this.data.aggregatePolicy = { policy: {}, policyString: '' };
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('agent')).exists('renders the base tabs');
    assert.dom(GENERAL.hdsTab('policies')).doesNotExist('does not render the policies tab');
  });

  test('it hides the alias details tab when the agent has no aliases', async function (assert) {
    this.data = this.dataWithoutAliases;

    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('agent')).exists('renders the base tabs');
    assert.dom(GENERAL.hdsTab('alias')).doesNotExist('does not render the alias details tab');
  });

  test('it defaults to alias tab when aliasId is defined', async function (assert) {
    this.data.aliasId = 'alias-1';
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('alias')).hasAttribute('aria-selected', 'true', 'defaults to alias tab');
    assert
      .dom(GENERAL.hdsTabPanel('agent'))
      .hasAttribute('hidden', '', 'hides the default tab panel when another initial tab is selected');
    assert.dom(GENERAL.hdsTabPanel('alias')).doesNotHaveAttribute('hidden', 'shows the requested tab panel');
  });

  test('it switches tabs when clicked', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('agent')).hasAttribute('aria-selected', 'true', 'defaults to the first tab');
    assert
      .dom(GENERAL.hdsTabPanel('agent'))
      .doesNotHaveAttribute('hidden', 'shows the first panel by default');

    await click(GENERAL.hdsTab('entity'));

    assert.dom(GENERAL.hdsTab('entity')).hasAttribute('aria-selected', 'true', 'selects the clicked tab');
    assert.dom(GENERAL.hdsTabPanel('entity')).doesNotHaveAttribute('hidden', 'shows the clicked tab panel');
    assert.dom(GENERAL.hdsTabPanel('agent')).hasAttribute('hidden', '', 'hides the previous panel');
  });

  test('policy detail inner tabs are all clickable', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.hdsTab('policies'));
    assert.dom(GENERAL.hdsTab('policies')).hasAttribute('aria-selected', 'true', 'policies tab is selected');

    // terraform tab is selected by default inside the policies panel
    assert
      .dom(GENERAL.hdsTab('terraform'))
      .hasAttribute('aria-selected', 'true', 'terraform tab is selected by default');
    assert
      .dom(GENERAL.hdsTabPanel('terraform'))
      .doesNotHaveAttribute('hidden', 'terraform panel is visible by default');

    await click(GENERAL.hdsTab('allowed-actions'));
    assert
      .dom(GENERAL.hdsTab('allowed-actions'))
      .hasAttribute('aria-selected', 'true', 'allowed actions tab becomes selected');
    assert
      .dom(GENERAL.hdsTabPanel('terraform'))
      .hasAttribute('hidden', '', 'terraform panel is hidden after switching away');

    await click(GENERAL.hdsTab('cli'));
    assert.dom(GENERAL.hdsTab('cli')).hasAttribute('aria-selected', 'true', 'CLI tab becomes selected');
    assert
      .dom(GENERAL.hdsTabPanel('cli'))
      .doesNotHaveAttribute('hidden', 'CLI panel is visible after clicking CLI tab');

    await click(GENERAL.hdsTab('terraform'));
    assert
      .dom(GENERAL.hdsTab('terraform'))
      .hasAttribute('aria-selected', 'true', 'terraform tab can be re-selected');
    assert
      .dom(GENERAL.hdsTabPanel('terraform'))
      .doesNotHaveAttribute('hidden', 'terraform panel is visible again');
  });
});
