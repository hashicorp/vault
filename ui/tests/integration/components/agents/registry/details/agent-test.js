/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | agents/registry/details/agent', function (hooks) {
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
        description: 'This is a test agent',
      },
      entity: {
        id: '6266cb5f-b9ff-4ac6-8a12-6fcb84c336ec',
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
      },
    };

    this.renderComponent = () => render(hbs`<Agents::Registry::Details::Agent @data={{this.data}} />`);
  });

  test('it renders the primary agent data', async function (assert) {
    await this.renderComponent();

    assert
      .dom('[data-test-agent-description]')
      .includesText(this.data.agent.description, 'renders the agent description');

    assert
      .dom(GENERAL.cardContainer('agent-status'))
      .includesText('Enabled', 'renders enabled status for an active entity');
    assert
      .dom(GENERAL.cardContainer('agent-status'))
      .includesText('owner-1@example.com', 'renders the agent owner');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('Agent registry ID', 'renders the agent id label');
    assert.dom(GENERAL.cardContainer('agent-details')).includesText('1', 'renders the agent id value');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('entity-0', 'renders the linked entity name');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('6266cb5f-b9ff-4ac6-8a12-6fcb84c336ec', 'renders the entity id');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('alias-0', 'renders the first alias name');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('auth/token/', 'renders the first alias mount path');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('auth_userpass_a212c204', 'renders alias accessors');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('92c031d1-a6df-4f97-ade1-8fa6f7cc4c4b', 'renders alias ids');
  });

  test('it renders disabled status and optional alias fields when present', async function (assert) {
    this.data.entity = {
      ...this.data.entity,
      disabled: true,
      aliases: [
        {
          ...this.data.entity.aliases[0],
          issuer_id: 'issuer-1',
          external_id: 'external-1',
          profile_name: 'profile-name-1',
          profile_id: 'profile-id-1',
          namespace: 'admin/',
        },
      ],
    };

    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('agent-status'))
      .includesText('Disabled', 'renders disabled status when the entity is disabled');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('issuer-1', 'renders the optional issuer id');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('external-1', 'renders the optional external id');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('profile-name-1', 'renders the optional profile name');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('profile-id-1', 'renders the optional profile id');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('admin/', 'renders the optional namespace');
  });

  test('it renders oauth aliases fields: mount_type, namespace, profile_name and profile_id', async function (assert) {
    this.data.entity = {
      ...this.data.entity,
      aliases: [
        {
          id: 'oauth-alias-1',
          name: 'repo:my-org/my-repo:ref:refs/heads/main',
          creation_time: '2025-01-10T09:00:00.000Z',
          last_update_time: '2025-06-01T12:00:00.000Z',
          mount_accessor: 'oauth-resource-server_root_b2f5c891-3a7d-4e12-9f8a-1c6d4e7b2a03',
          mount_path: 'auth/oauth-resource-server/',
          // mount_type is set to 'oauth' by the enrichment logic in the page component
          mount_type: 'oauth',
          namespace: 'root',
          profile_name: 'github-actions',
          profile_id: 'b2f5c891-3a7d-4e12-9f8a-1c6d4e7b2a03',
        },
      ],
    };

    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('repo:my-org/my-repo:ref:refs/heads/main', 'renders the OAuth alias name');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('oauth', 'renders the oauth mount_type badge');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('root', 'renders the namespace from the OAuth profile');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('github-actions', 'renders the profile_name from the OAuth profile');
    assert
      .dom(GENERAL.cardContainer('agent-details'))
      .includesText('b2f5c891-3a7d-4e12-9f8a-1c6d4e7b2a03', 'renders the profile_id from the OAuth profile');
  });

  test('it renders policy allowed actions when aggregatedPolicy is provided', async function (assert) {
    this.data.aggregatePolicy = {
      policy: {
        'identity/*': ['read'],
      },
      policyString: `path "identity/*" {
  capabilities = ["read"]
}`,
    };

    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('policy-allowed-actions'))
      .includesText('Allowed actions', 'renders policy allowed actions card title');
    assert.dom('[data-test-path-entry="identity/*"]').exists('renders allowed actions path entry');
  });
});
