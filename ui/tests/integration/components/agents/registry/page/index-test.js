/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import { dateFormat } from 'core/helpers/date-format';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | agents/page/registry', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.downloadStub = Sinon.stub(this.owner.lookup('service:download'), 'csv');
    this.agent = {
      display_name: 'agent-1',
      name: 'cool agent name',
      entity_id: 'entity-1',
      no_default_ceiling_policy: false,
      entity: {
        name: 'test-entity',
        disabled: true,
        creation_time: '2025-06-01T13:02:03Z',
        last_update_time: '2025-06-01T13:02:03Z',
        aliases: [
          {
            id: 'alias-id',
            name: 'test-alias',
            creation_time: '2025-06-01T13:22:03Z',
            last_update_time: '2025-06-01T13:22:03Z',
            mount_accessor: 'auth_token_d2eccb3b',
            mount_path: 'auth/token/',
            mount_type: 'token',
          },
        ],
      },
    };
    this.model = { agents: [this.agent] };
    this.breadcrumbs = [];
  });

  test('it renders empty state when no agents are present', async function (assert) {
    this.model.agents = [];

    await render(hbs`
      <Agents::Registry::Page @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} />
    `);

    assert.dom(GENERAL.emptyStateTitle).hasText('No agents yet');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Your agents will be listed here. Register an agent in CLI or API to get started.');
    assert.dom(GENERAL.button('refresh-agents-list')).exists();
  });

  test('it exports table data as CSV on click', async function (assert) {
    await render(hbs`
      <Agents::Registry::Page @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} />
    `);

    await click('[data-test-agents-registry-export]');

    assert.true(this.downloadStub.calledOnce, 'Clicking Export triggers a CSV download');
    const [filename, content] = this.downloadStub.lastCall.args;
    const { creation_time, last_update_time } = this.agent.entity;
    const { creation_time: aliasCreationTime, last_update_time: aliasLastUpdateTime } =
      this.agent.entity.aliases[0];
    const entryCreatedAt = dateFormat([creation_time, 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });
    const entryUpdatedAt = dateFormat([last_update_time, 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });
    const aliasCreatedAt = dateFormat([aliasCreationTime, 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });
    const aliasUpdatedAt = dateFormat([aliasLastUpdateTime, 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });

    assert.strictEqual(filename, 'agent_registry', 'It uses the expected CSV filename');
    assert.strictEqual(
      content,
      `Agent name,Agentic entity in Vault,Entity / Alias ID,Entity status,Entity created at,Entity updated at\nagent-1,test-entity,entity-1,Disabled,"${entryCreatedAt}","${entryUpdatedAt}"\n,Alias: test-alias,alias-id,/,"${aliasCreatedAt}","${aliasUpdatedAt}"`,
      'The CSV includes table rows and escaped values in display order'
    );
  });
});
