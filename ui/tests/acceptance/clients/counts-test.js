/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, { STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { visit, click, currentURL, fillIn } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import timestamp from 'core/utils/timestamp';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | clients | counts', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  test('it should prompt user to query start time for community version', async function (assert) {
    assert.expect(2);
    this.owner.lookup('service:version').type = 'community';
    await visit('/vault/clients/counts/overview');

    assert.dom(GENERAL.emptyStateTitle).hasText('No data received');
    assert.dom(GENERAL.emptyStateMessage).hasText('Select a start date above to query client count data.');
  });

  test('it should redirect to counts overview route for transitions to parent', async function (assert) {
    await visit('/vault/clients');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Redirects to counts overview route');
  });

  test('it should persist filter query params between child routes', async function (assert) {
    await visit('/vault/clients/counts/overview');
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2020-03');
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2022-02');
    await click(GENERAL.saveButton);
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/overview?end_time=1706659200&start_time=1643673600',
      'Start and end times added as query params'
    );

    await click(GENERAL.tab('token'));
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/token?end_time=1706659200&start_time=1643673600',
      'Start and end times persist through child route change'
    );

    await click(GENERAL.navLink('Dashboard'));
    await click(GENERAL.navLink('Client Count'));
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/overview',
      'Query params are reset when exiting route'
    );
  });

  test('it should render empty state if no permission to query activity data', async function (assert) {
    assert.expect(2);
    server.get('/sys/internal/counters/activity', () => {
      return overrideResponse(403);
    });
    await visit('/vault/clients/counts/overview');
    assert.dom(GENERAL.emptyStateTitle).hasText('You are not authorized');
    assert
      .dom(GENERAL.emptyStateActions)
      .hasText(
        'You must be granted permissions to view this page. Ask your administrator if you think you should have access to the /v1/sys/internal/counters/activity endpoint.'
      );
  });
});
