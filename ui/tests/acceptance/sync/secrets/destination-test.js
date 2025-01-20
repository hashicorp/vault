/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import syncScenario from 'vault/mirage/scenarios/sync';
import syncHandlers from 'vault/mirage/handlers/sync';
import authPage from 'vault/tests/pages/auth';
import { click, visit, currentURL, fillIn } from '@ember/test-helpers';
import { PAGE as ts } from 'vault/tests/helpers/sync/sync-selectors';

module('Acceptance | enterprise | sync | destination', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    syncScenario(this.server);
    syncHandlers(this.server);
    return authPage.login();
  });

  test('it should transition to overview route via breadcrumb', async function (assert) {
    await visit('vault/sync/secrets/destinations/aws-sm/destination-aws/secrets');
    await click(ts.breadcrumbAtIdx(0));
    assert.strictEqual(
      currentURL(),
      '/vault/sync/secrets/overview',
      'Transitions to overview on breadcrumb click'
    );
  });

  test('it should transition to correct routes when performing actions', async function (assert) {
    await click(ts.navLink('Secrets Sync'));
    await click(ts.tab('Destinations'));
    await click(ts.listItem);
    assert.dom(ts.tab('Secrets')).hasClass('active', 'Secrets tab is active');

    await click(ts.tab('Details'));
    assert.dom(ts.infoRowLabel('Name')).exists('Destination details display');

    await click(ts.toolbar('Sync secrets'));
    await click(ts.destinations.sync.cancel);

    await click(ts.toolbar('Edit destination'));
    assert.dom(ts.inputByAttr('name')).isDisabled('Edit view renders with disabled name field');
    await click(ts.cancelButton);
    assert.dom(ts.tab('Details')).hasClass('active', 'Details view is active');
  });

  test('it should delete destination', async function (assert) {
    await visit('vault/sync/secrets/destinations/aws-sm/destination-aws/details');
    await click(ts.toolbar('Delete destination'));
    await fillIn(ts.confirmModalInput, 'DELETE');
    await click(ts.confirmButton);
    assert.dom(ts.destinations.deleteBanner).exists('Delete banner renders');
  });

  test('it should not save placeholder values for credentials and only save when there are changes', async function (assert) {
    assert.expect(2);

    const handler = this.server.patch(
      '/sys/sync/destinations/vercel-project/destination-vercel',
      (schema, req) => {
        assert.deepEqual(
          JSON.parse(req.requestBody),
          { access_token: 'foobar' },
          'Updated access token sent in patch request'
        );
        const { deployment_environments, project_id, team_id, name, type, secret_name_template } =
          this.server.create('sync-destination', 'vercel-project');
        return {
          data: {
            connection_details: { access_token: '*****', deployment_environments, project_id, team_id },
            name,
            options: { custom_tags: {}, secret_name_template },
            type,
          },
        };
      }
    );

    await visit('vault/sync/secrets/destinations/vercel-project/destination-vercel/edit');
    await click(ts.enableField('accessToken'));
    await fillIn(ts.inputByAttr('accessToken'), 'foobar');
    await click(ts.saveButton);
    await click(ts.toolbar('Edit destination'));
    await click(ts.saveButton);
    assert.strictEqual(
      handler.numberOfCalls,
      1,
      'Model is not dirty after server returns masked value for credentials and save request is not made when there are no changes'
    );
  });
});
