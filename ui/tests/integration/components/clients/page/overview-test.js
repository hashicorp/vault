/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { fillIn, findAll, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { filterActivityResponse } from 'vault/mirage/handlers/clients';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { flattenMounts } from 'core/utils/client-count-utils';

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

    const store = this.owner.lookup('service:store');
    this.activity = await store.queryRecord('clients/activity', {});
    this.mostRecentMonth = this.activity.byMonth[this.activity.byMonth.length - 1];
    this.onFilterChange = sinon.spy();
    this.filterQueryParams = { namespace_path: '', mount_path: '', mount_type: '' };
    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Overview 
        @activity={{this.activity}} 
        @onFilterChange={{this.onFilterChange}} 
        @filterQueryParams={{this.filterQueryParams}} 
      />`);
  });

  test('it shows correct state message when month selection has no data', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '2023-06-01T00:00:00Z');

    assert
      .dom(CLIENT_COUNT.card('table empty state'))
      .hasText('No data found Clear or update filters to view client count data. Client count documentation');
  });

  test('it shows table when month selection has data', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '9/23');

    assert.dom(CLIENT_COUNT.card('table empty state')).doesNotExist('does not show card when table has data');
    assert.dom(GENERAL.table('attribution')).exists('shows table');
    assert.dom(GENERAL.paginationInfo).hasText('1–6 of 6', 'shows correct pagination info');
    assert.dom(GENERAL.paginationSizeSelector).hasValue('10', 'page size selector defaults to "10"');
  });

  test('it shows correct month options for billing period', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '');
    await triggerEvent(GENERAL.selectByAttr('attribution-month'), 'change');

    // assert that months options in select are those of selected billing period
    // '' represents default state of 'Select month'
    const expectedOptions = ['', ...this.activity.byMonth.reverse().map((m) => m.timestamp)];
    const actualOptions = findAll(`${GENERAL.selectByAttr('attribution-month')} option`).map(
      (option) => option.value
    );
    assert.deepEqual(actualOptions, expectedOptions, 'All <option> values match expected list');
  });

  test('it initially renders attribution with by_namespace data', async function (assert) {
    await this.renderComponent();
    const topNamespace = this.activity.byNamespace[0];
    const topMount = topNamespace.mounts[0];
    // Assert table renders namespace with the highest counts at the top
    assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(topNamespace.label);
    assert.dom(GENERAL.tableData(0, 'clients')).hasText(`${topMount.clients}`);
  });

  test('it filters attribution table by month', async function (assert) {
    await this.renderComponent();
    const mostRecentMonth = this.mostRecentMonth;
    await fillIn(GENERAL.selectByAttr('attribution-month'), mostRecentMonth.timestamp);
    // Drill down to new_clients then grab the first namespace
    const sortedMounts = flattenMounts(mostRecentMonth.new_clients.namespaces).sort(
      (a, b) => b.clients - a.clients
    );
    const topMount = sortedMounts[0];
    assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(topMount.namespace_path);
    assert.dom(GENERAL.tableData(0, 'clients')).hasText(`${topMount.clients}`);
    assert.dom(GENERAL.tableData(0, 'mount_path')).hasText(topMount.mount_path);
  });
});
