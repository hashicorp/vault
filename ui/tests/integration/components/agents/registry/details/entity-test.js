/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | agents/registry/details/entity', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = {
      agent: {
        id: '1',
        display_name: 'agent-1',
        entity_id: '6266cb5f-b9ff-4ac6-8a12-6fcb84c336ec',
        creation_time: '2026-06-26T18:56:38.507Z',
        owner: 'owner-1@example.com',
        last_updated_time: '2026-06-26T18:56:38.507Z',
        no_default_ceiling_policy: false,
      },
      entity: {
        id: '6266cb5f-b9ff-4ac6-8a12-6fcb84c336ec',
        metadata: {
          owner_team: 'platform',
          environment: 'production',
        },
        aliases: [
          {
            id: '946acaa2-0577-453b-a801-3884023dd128',
            name: 'alias-0',
            creation_time: '2023-02-21T14:52:57.822Z',
            last_update_time: '2025-12-17T00:49:04.138Z',
            mount_accessor: 'auth_token_d2eccb3b',
            mount_path: 'auth/token/',
            mount_type: 'token',
          },
          {
            id: '92c031d1-a6df-4f97-ade1-8fa6f7cc4c4b',
            name: 'alias-1',
            creation_time: '2025-06-29T02:08:22.947Z',
            last_update_time: '2026-05-17T18:18:36.835Z',
            mount_accessor: 'auth_userpass_a212c204',
            mount_path: 'auth/userpass/',
            mount_type: 'userpass',
          },
          {
            id: '3a5e317b-c850-4f70-affa-ac11760508a0',
            name: 'alias-2',
            creation_time: '2024-07-17T15:21:42.403Z',
            last_update_time: '2024-07-29T19:55:52.690Z',
            mount_accessor: 'auth_token_d2eccb3b',
            mount_path: 'auth/token/',
            mount_type: 'token',
          },
        ],
        name: 'entity-0',
        disabled: false,
        creation_time: '2025-05-25T22:35:56.701Z',
        last_update_time: '2025-06-29T15:59:34.503Z',
        group_ids: ['01'],
        description: 'Production deployment - credential rotation, TLS cert issuance. 15-min TTL.',
        policies: ['g-deploy-prod', 'g-deploy-test'],
      },
      groups: [{ id: '01', name: 'Engineering-core', policies: ['eng-1', 'eng-2'] }],
    };

    this.renderComponent = () => render(hbs`<Agents::Registry::Details::Entity @data={{this.data}} />`);
  });

  test('it renders merged ids fallback and entity metadata section', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('merged-ids'))
      .includesText('Merged IDs', 'renders merged ids card heading');
    assert
      .dom(GENERAL.cardContainer('merged-ids'))
      .includesText('None', 'renders none fallback when no ids exist');

    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Owner team', 'renders sentence-cased metadata label');
    assert
      .dom(`${GENERAL.widget('entity-metadata')} input[name="owner_team"]`)
      .hasValue('platform', 'renders owner_team metadata value');
    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Environment', 'renders second metadata label');
    assert
      .dom(`${GENERAL.widget('entity-metadata')} input[name="environment"]`)
      .hasValue('production', 'renders environment metadata value');
  });

  test('it renders merged entity ids when provided', async function (assert) {
    this.data.entity.merged_entity_ids = ['merged-id-1', 'merged-id-2'];

    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('merged-ids'))
      .includesText('merged-id-1', 'renders first merged id in the merged ids card');
    assert
      .dom(GENERAL.cardContainer('merged-ids'))
      .includesText('merged-id-2', 'renders second merged id in the merged ids card');
    assert
      .dom(GENERAL.cardContainer('merged-ids'))
      .doesNotIncludeText('None', 'does not render fallback text when merged ids exist');
  });

  test('it renders groups table', async function (assert) {
    await this.renderComponent();

    assert.dom().includesText('In group', 'renders the group column header');
    assert.dom().includesText('Policies', 'renders the policies column header');

    assert.dom().includesText('Engineering-core', 'renders the group name');
    assert.dom().includesText('eng-1', 'renders first policy badge');
    assert.dom().includesText('eng-2', 'renders second policy badge');

    assert.dom(GENERAL.pagination).doesNotExist('pagination is not rendered for the groups table');
  });

  test('it renders empty policy state when entity has no policies', async function (assert) {
    this.data.groups[0].policies = [];

    await this.renderComponent();

    assert
      .dom()
      .includesText('No policy in this group yet', 'renders empty policy state when policies is empty');
  });
  test('it renders entity metadata with keys and values', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.widget('entity-metadata')).exists('renders the entity metadata section');
    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Owner team', 'renders humanized owner_team label');
    assert
      .dom(`${GENERAL.widget('entity-metadata')} input[name="owner_team"]`)
      .hasValue('platform', 'renders owner_team value');
    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Environment', 'renders humanized environment label');
    assert
      .dom(`${GENERAL.widget('entity-metadata')} input[name="environment"]`)
      .hasValue('production', 'renders environment value');
  });

  test('it hides the entity metadata section when metadata is absent', async function (assert) {
    delete this.data.entity.metadata;

    await this.renderComponent();

    assert.dom().doesNotContainText('Entity metadata');
    assert.dom(GENERAL.widget('entity-metadata')).doesNotExist();
  });
});
