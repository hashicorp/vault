/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { filterActivityResponse } from 'vault/mirage/handlers/clients';
import { CHARTS, CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

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
  });

  test('it hides attribution data when mount filter applied', async function (assert) {
    this.mountPath = '';
    this.activity = await this.store.queryRecord('clients/activity', {
      namespace: 'ns1',
    });
    await render(
      hbs`<Clients::Page::Overview @activity={{this.activity}} @namespace="ns1" @mountPath={{this.mountPath}} />`
    );

    assert.dom(CHARTS.container('Vault client counts')).exists('shows running totals');
    assert.dom(CLIENT_COUNT.attributionBlock('namespace')).exists();
    assert.dom(CLIENT_COUNT.attributionBlock('mount')).exists();

    this.set('mountPath', 'auth/authid/0');
    assert.dom(CHARTS.container('Vault client counts')).exists('shows running totals');
    assert.dom(CLIENT_COUNT.attributionBlock('namespace')).doesNotExist();
    assert.dom(CLIENT_COUNT.attributionBlock('mount')).doesNotExist();
  });

  test('it hides attribution data when no data returned', async function (assert) {
    this.mountPath = '';
    this.activity = await this.store.queryRecord('clients/activity', {
      namespace: 'no-data',
    });
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}} />`);
    assert.dom(CLIENT_COUNT.usageStats('Total usage')).exists();
    assert.dom(CHARTS.container('Vault client counts')).doesNotExist('usage stats instead of running totals');
    assert.dom(CLIENT_COUNT.attributionBlock('namespace')).doesNotExist();
    assert.dom(CLIENT_COUNT.attributionBlock('mount')).doesNotExist();
  });

  test('it shows the correct mount attributions', async function (assert) {
    this.nsService = this.owner.lookup('service:namespace');
    const rootActivity = await this.store.queryRecord('clients/activity', {});
    this.activity = rootActivity;
    await render(hbs`<Clients::Page::Overview @activity={{this.activity}} />`);
    // start at "root" namespace
    let expectedMounts = rootActivity.byNamespace.find((ns) => ns.label === 'root').mounts;
    assert
      .dom(`${CLIENT_COUNT.attributionBlock('mount')} [data-test-group="y-labels"] text`)
      .exists({ count: expectedMounts.length });
    assert
      .dom(`${CLIENT_COUNT.attributionBlock('mount')} [data-test-group="y-labels"]`)
      .includesText(expectedMounts[0].label);

    // now pretend we're querying within a child namespace
    this.nsService.path = 'ns1';
    this.activity = await this.store.queryRecord('clients/activity', {
      namespace: 'ns1',
    });
    expectedMounts = rootActivity.byNamespace.find((ns) => ns.label === 'ns1').mounts;
    assert
      .dom(`${CLIENT_COUNT.attributionBlock('mount')} [data-test-group="y-labels"] text`)
      .exists({ count: expectedMounts.length });
    assert
      .dom(`${CLIENT_COUNT.attributionBlock('mount')} [data-test-group="y-labels"]`)
      .includesText(expectedMounts[0].label);
  });
});
