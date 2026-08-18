/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | agents/registry/details/alias', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = {
      agent: {
        id: 'agent-1',
        display_name: 'test-agent-1',
        entity_id: 'entity-1',
      },
      entity: {
        id: 'entity-1',
        name: 'test-entity-1',
        metadata: {
          owner_team: 'platform',
          environment: 'production',
        },
        aliases: [
          {
            id: 'alias-1',
            name: 'test-alias-1',
            mount_accessor: 'auth_token_d2eccb3b',
            mount_path: 'auth/token/',
            mount_type: 'token',
            custom_metadata: {
              client_id: 'client-1',
              region: 'us-east-1',
            },
          },
          {
            id: 'alias-2',
            name: 'test-alias-2',
            mount_accessor: 'auth_userpass_a212c204',
            mount_path: 'auth/userpass/',
            mount_type: 'userpass',
            custom_metadata: {
              client_id: 'client-2',
              region: 'eu-west-1',
            },
          },
        ],
      },
    };

    this.renderComponent = () => render(hbs`<Agents::Registry::Details::Alias @data={{this.data}} />`);
  });

  test('it applies initial alias selection from aliasId', async function (assert) {
    this.data.aliasId = 'alias-2';

    await this.renderComponent();

    assert
      .dom(GENERAL.linkTo('selected alias'))
      .hasText('alias-2', 'renders the selected alias id on first render');
    assert
      .dom(GENERAL.copySnippet('selected-alias-mount-path'))
      .includesText('auth/userpass/', 'renders fields for the selected alias');
    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Owner team', 'renders entity metadata label in sentence case');
    assert
      .dom(GENERAL.widget('alias-custom-metadata'))
      .includesText('Client id', 'renders alias metadata label in sentence case');
  });

  test('it re-syncs selected alias when aliasId changes', async function (assert) {
    this.data.aliasId = 'alias-1';

    await this.renderComponent();

    assert.dom(GENERAL.linkTo('selected alias')).hasText('alias-1', 'starts with the first selected alias');
    assert
      .dom(GENERAL.copySnippet('selected-alias-mount-path'))
      .includesText('auth/token/', 'starts with first alias mount path');

    this.set('data.aliasId', 'alias-2');

    assert.dom(GENERAL.linkTo('selected alias')).hasText('alias-2', 'updates to the newly selected alias');
    assert
      .dom(GENERAL.copySnippet('selected-alias-mount-path'))
      .includesText('auth/userpass/', 'updates mount path after selectedAliasId changes');
  });

  test('it sets previous/next disabled state at the first alias', async function (assert) {
    this.data.aliasId = 'alias-1';

    await this.renderComponent();

    assert.dom(GENERAL.button('alias-previous')).isDisabled('previous is disabled at first alias');
    assert.dom(GENERAL.button('alias-next')).isNotDisabled('next is enabled at first alias');
  });

  test('it updates previous/next disabled state after pagination', async function (assert) {
    this.data.aliasId = 'alias-1';

    await this.renderComponent();

    await click(GENERAL.button('alias-next'));

    assert.dom(GENERAL.button('alias-previous')).isNotDisabled('previous is enabled at the last alias');
    assert.dom(GENERAL.button('alias-next')).isDisabled('next is disabled at the last alias');

    await click(GENERAL.button('alias-previous'));

    assert.dom(GENERAL.button('alias-previous')).isDisabled('previous is disabled again at first alias');
    assert.dom(GENERAL.button('alias-next')).isNotDisabled('next is re-enabled after moving back');
  });

  test('it hides pagination controls when there is only one alias', async function (assert) {
    this.data.aliasId = 'alias-1';
    this.data.entity.aliases = [this.data.entity.aliases[0]];

    await this.renderComponent();

    assert.dom(GENERAL.button('alias-previous')).doesNotExist('previous is hidden for a single alias');
    assert.dom(GENERAL.button('alias-next')).doesNotExist('next is hidden for a single alias');
  });

  test('it hides pagination controls when there are no aliases', async function (assert) {
    this.data.entity.aliases = [];

    await this.renderComponent();

    assert.dom(GENERAL.button('alias-previous')).doesNotExist('previous is hidden when there are no aliases');
    assert.dom(GENERAL.button('alias-next')).doesNotExist('next is hidden when there are no aliases');
  });

  test('it renders entity metadata labels and values', async function (assert) {
    this.data.aliasId = 'alias-1';

    await this.renderComponent();

    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Owner team', 'renders sentence-cased entity metadata label');
    assert
      .dom(`${GENERAL.widget('entity-metadata')} input[name="owner_team"]`)
      .hasValue('platform', 'renders entity metadata value for owner_team');
    assert
      .dom(GENERAL.widget('entity-metadata'))
      .includesText('Environment', 'renders additional sentence-cased entity metadata label');
    assert
      .dom(`${GENERAL.widget('entity-metadata')} input[name="environment"]`)
      .hasValue('production', 'renders entity metadata value for environment');
  });

  test('it renders custom_metadata for selected alias and updates on alias change', async function (assert) {
    this.data.aliasId = 'alias-1';

    await this.renderComponent();

    assert
      .dom(GENERAL.widget('alias-custom-metadata'))
      .includesText('Client id', 'renders sentence-cased custom metadata label');
    assert
      .dom(`${GENERAL.widget('alias-custom-metadata')} input[name="client_id"]`)
      .hasValue('client-1', 'renders custom metadata value for first alias');
    assert
      .dom(GENERAL.widget('alias-custom-metadata'))
      .includesText('Region', 'renders additional custom metadata label');
    assert
      .dom(`${GENERAL.widget('alias-custom-metadata')} input[name="region"]`)
      .hasValue('us-east-1', 'renders additional custom metadata value for first alias');

    this.set('data.aliasId', 'alias-2');

    assert
      .dom(`${GENERAL.widget('alias-custom-metadata')} input[name="client_id"]`)
      .hasValue('client-2', 'updates custom metadata value when selected alias changes');
    assert
      .dom(`${GENERAL.widget('alias-custom-metadata')} input[name="region"]`)
      .hasValue('eu-west-1', 'updates additional custom metadata value when selected alias changes');
  });

  test('it hides entity metadata section when custom_metadata is not present', async function (assert) {
    this.data.entity.metadata = null;

    await this.renderComponent();

    assert.dom().includesText('Alias metadata');
    assert.dom(GENERAL.widget('entity-metadata')).doesNotExist();
  });

  test('it hides custom metadata section when missing', async function (assert) {
    this.data.entity.aliases[0].custom_metadata = null;

    await this.renderComponent();

    assert.dom().includesText('Alias metadata');
    assert.dom(GENERAL.widget('alias-custom-metadata')).doesNotExist();
  });

  test('metadata section is hidden when entity and alias metadata is missing', async function (assert) {
    this.data.entity.metadata = null;
    this.data.entity.aliases[0].custom_metadata = null;

    await this.renderComponent();

    assert.dom().doesNotIncludeText('Alias metadata');
    assert.dom(GENERAL.widget('entity-metadata')).doesNotExist();
    assert.dom(GENERAL.widget('alias-custom-metadata')).doesNotExist();
  });

  test('it should render placeholder when mount_path is not returned', async function (assert) {
    this.data.entity.aliases[0].mount_path = null;

    await this.renderComponent();
    assert
      .dom(GENERAL.textDisplay('selected-alias-mount-path'))
      .includesText('/', 'renders placeholder for mount_path');
  });
});
