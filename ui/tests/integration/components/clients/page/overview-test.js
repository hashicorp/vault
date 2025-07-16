/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { filterActivityResponse } from 'vault/mirage/handlers/clients';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | clients/page/overview', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.get('sys/internal/counters/activity', (_, req) => {
      const namespace = req.requestHeaders['X-Vault-Namespace'];
      if (namespace === 'no-data') {
        return {
          request_id: 'some-activity-id',
          data: {
            by_namespace: [],
            end_time: '2024-08-31T23:59:59Z',
            months: [],
            start_time: '2024-01-01T00:00:00Z',
            total: {
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 0,
              secret_syncs: 0,
            },
          },
        };
      }
      return {
        request_id: 'some-activity-id',
        data: filterActivityResponse(ACTIVITY_RESPONSE_STUB, namespace),
      };
    });
    this.store = this.owner.lookup('service:store');
    this.mountPath = '';
    this.namespace = '';
    this.versionHistory = '';
    this.activity = await this.store.queryRecord('clients/activity', {});

    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'aria-prohibited-attr': { enabled: false },
      },
    });
  });

  test('it shows empty state message upon initial load', async function (assert) {
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}}/>`);

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');

    assert.dom(CLIENT_COUNT.attribution.card).exists('shows card for table state');
    assert
      .dom(CLIENT_COUNT.attribution.card)
      .hasText(
        'Select a month to view client attribution View the namespace mount breakdown of clients by selecting a month. Client count documentation',
        'Show initial table state message'
      );
  });

  test('it shows correct state message when month selection has no data', async function (assert) {
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}} />`);

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '6/23');

    assert
      .dom(CLIENT_COUNT.attribution.card)
      .hasText(
        'No data is available for the selected month View the namespace mount breakdown of clients by selecting another month. Client count documentation',
        'Shows correct message for a month selection with no data'
      );
  });

  test('it shows table when month selection has data', async function (assert) {
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}} />`);

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '9/23');

    assert.dom(CLIENT_COUNT.attribution.card).doesNotExist('does not show card when table has data');
    assert.dom(CLIENT_COUNT.attribution.table).exists('shows table');
    assert.dom(CLIENT_COUNT.attribution.paginationInfo).hasText('1–3 of 6', 'shows correct pagination info');
  });

  test('it filters the table when a namespace filter is applied', async function (assert) {
    this.namespace = 'ns1';
    this.activity = await this.store.queryRecord('clients/activity', {
      namespace: this.namespace,
    });
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}} @namespace={{this.namespace}} />`);

    await fillIn(GENERAL.selectByAttr('attribution-month'), '9/23');

    assert.dom(CLIENT_COUNT.attribution.card).doesNotExist('does not show card when table has data');
    assert.dom(CLIENT_COUNT.attribution.table).exists();
    assert.dom(CLIENT_COUNT.attribution.paginationInfo).hasText('1–3 of 3', 'shows correct pagination info');
  });

  test('it hides the table when a mount filter is applied', async function (assert) {
    this.namespace = 'ns1';
    this.mountPath = 'auth/userpass-0';
    this.activity = await this.store.queryRecord('clients/activity', {
      namespace: this.namespace,
      mountPath: this.mountPath,
    });
    await render(
      hbs`<Clients::Page::Overview @activity={{this.activity}} @namespace={{this.namespace}} @mountPath={{this.mountPath}}/>`
    );
    assert.dom(CLIENT_COUNT.attribution.card).doesNotExist('does not show card when table has data');
    assert
      .dom(CLIENT_COUNT.attribution.table)
      .doesNotExist('does not show table when a mount filter is applied');
  });

  test('it paginates table data', async function (assert) {
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}}  />`);

    await fillIn(GENERAL.selectByAttr('attribution-month'), '9/23');

    assert
      .dom(CLIENT_COUNT.attribution.row)
      .exists({ count: 3 }, 'Correct number of table rows render based on page size');
    assert.dom(CLIENT_COUNT.attribution.counts(0)).hasText('96', 'First page shows data');
    assert.dom(CLIENT_COUNT.attribution.pagination).exists('shows pagination');
    assert.dom(CLIENT_COUNT.attribution.paginationInfo).hasText('1–3 of 6', 'shows correct pagination info');

    await click(GENERAL.pagination.next);

    assert.dom(CLIENT_COUNT.attribution.counts(0)).hasText('53', 'Second page shows new data');
    assert.dom(CLIENT_COUNT.attribution.paginationInfo).hasText('4–6 of 6', 'shows correct pagination info');
  });

  test('it shows correct month options for billing period', async function (assert) {
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}} />`);

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '');
    await triggerEvent(GENERAL.selectByAttr('attribution-month'), 'change');

    // assert that months options in select are those of selected billing period
    const expectedMonths = this.activity.byMonth.reverse().map((m) => m.month);

    // '' represents default state of 'Select month'
    const expectedOptions = ['', ...expectedMonths];
    const actualOptions = findAll(`${GENERAL.selectByAttr('attribution-month')} option`).map(
      (option) => option.value
    );
    assert.deepEqual(actualOptions, expectedOptions, 'All <option> values match expected list');
  });
});
