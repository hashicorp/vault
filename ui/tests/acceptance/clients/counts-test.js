/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, { STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { visit, click, currentURL } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { SELECTORS as ts } from 'vault/tests/helpers/clients';
import timestamp from 'core/utils/timestamp';

module('Acceptance | clients | counts', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it should prompt user to query start time for community version', async function (assert) {
    assert.expect(2);
    this.owner.lookup('service:version').type = 'community';
    await visit('/vault/clients/counts/overview');

    assert.dom(ts.emptyStateTitle).hasText('No data received');
    assert.dom(ts.emptyStateMessage).hasText('Select a start date above to query client count data.');
  });

  test('it should redirect to counts overview route for transitions to parent', async function (assert) {
    await visit('/vault/clients');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Redirects to counts overview route');
  });

  test('it should persist filter query params between child routes', async function (assert) {
    await visit('/vault/clients/counts/overview');
    await click(ts.rangeDropdown);
    await click(ts.currentBillingPeriod);
    const timeQueryRegex = /end_time=\d+&start_time=\d+/g;
    assert.ok(currentURL().match(timeQueryRegex).length, 'Start and end times added as query params');

    await click(ts.tab('token'));
    assert.ok(
      currentURL().match(timeQueryRegex).length,
      'Start and end times persist through child route change'
    );

    await click(ts.navLink('Dashboard'));
    await click(ts.navLink('Client Count'));
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/overview',
      'Query params are reset when exiting route'
    );
  });
});
