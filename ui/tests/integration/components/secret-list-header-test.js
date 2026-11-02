/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { getContext, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const mirageToModels = (data) => {
  const context = getContext();
  const store = context.owner.lookup('service:store');
  const modelName = Array.isArray(data) ? data[0].modelName : data.modelName;
  const json = context.server.serializerOrRegistry.serialize(data);
  store.push(json);
  return Array.isArray(data)
    ? data.map(({ id }) => store.peekRecord(modelName, id))
    : store.peekRecord(modelName, data.id);
};

module('Integration | Component | secret-list-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  for (const rolePath of ['roles', 'static-roles']) {
    for (const capability of ['list', 'create', 'update']) {
      test(`database Roles tab is shown with only ${capability} on ${rolePath}`, async function (assert) {
        this.server.post('/sys/capabilities-self', (_, request) => {
          const { paths } = JSON.parse(request.requestBody);
          return {
            data: Object.fromEntries(
              paths.map((path) => [path, path === `custom-database/${rolePath}` ? [capability] : ['deny']])
            ),
          };
        });
        this.model = { id: 'custom-database', engineType: 'database' };

        await render(hbs`<SecretListHeader @model={{this.model}} />`);

        assert.dom(GENERAL.secretTab('Roles')).hasText('Roles', 'Either role type can make the tab visible');
        assert
          .dom(GENERAL.secretTab('Connections'))
          .doesNotExist('Role permissions do not grant access to the Connections tab');
      });
    }
  }

  for (const roleCapabilities of [['deny'], ['read', 'delete']]) {
    test(`database Roles tab is hidden with only ${roleCapabilities.join(
      '/'
    )} on both role types`, async function (assert) {
      this.server.post('/sys/capabilities-self', (_, request) => {
        const { paths } = JSON.parse(request.requestBody);
        return {
          data: Object.fromEntries(
            paths.map((path) => [path, path === 'custom-database/config' ? ['list'] : roleCapabilities])
          ),
        };
      });
      this.model = { id: 'custom-database', engineType: 'database' };

      await render(hbs`<SecretListHeader @model={{this.model}} />`);

      assert.dom(GENERAL.secretTab('Roles')).doesNotExist('Neither role type grants tab access');
      assert
        .dom(GENERAL.secretTab('Connections'))
        .hasText('Connections', 'Connection permissions still independently control their tab');
    });
  }

  test('it should render version badge for kv and generic engine types', async function (assert) {
    const backends = supportedSecretBackends();
    assert.expect(backends.length);

    this.server.post('/sys/capabilities-self', () => {});

    for (const type of backends) {
      const data = this.server.create('secret-engine', 2, { type });
      this.model = mirageToModels(data);
      await render(hbs`
        <SecretListHeader
          @model={{this.model}}
        />
      `);

      if (['kv', 'generic'].includes(type)) {
        assert
          .dom(GENERAL.badge('kv version'))
          .hasText(
            `version ${this.model.version || 1}`,
            `Badge renders with correct version for ${type} engine type`
          );
      } else {
        assert
          .dom(GENERAL.badge('kv version'))
          .doesNotExist(`Version badge does not render for ${type} engine type`);
      }
    }
  });
});
