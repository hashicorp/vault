/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { formatRFC3339 } from 'date-fns';
import subMonths from 'date-fns/subMonths';
import timestamp from 'core/utils/timestamp';
import { SERIALIZED_ACTIVITY_RESPONSE } from 'vault/tests/helpers/clients/client-count-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_TYPES } from 'core/utils/client-count-utils';

const CLIENTS_ATTRIBUTION = {
  title: '[data-test-attribution-title]',
  description: '[data-test-attribution-description]',
  subtext: '[data-test-attribution-subtext]',
  timestamp: '[data-test-attribution-timestamp]',
  chart: '[data-test-horizontal-bar-chart]',
  topItem: '[data-test-top-attribution]',
  topItemCount: '[data-test-attribution-clients]',
  yLabel: '[data-test-group="y-labels"]',
  yLabels: '[data-test-group="y-labels"] text',
};
module('Integration | Component | clients/attribution', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    this.timestampStub = sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
  });

  hooks.beforeEach(function () {
    const mockNow = this.timestampStub();
    this.mockNow = mockNow;
    this.startTimestamp = formatRFC3339(subMonths(mockNow, 6));
    this.timestamp = formatRFC3339(mockNow);
    this.selectedNamespace = null;
    this.namespaceAttribution = SERIALIZED_ACTIVITY_RESPONSE.by_namespace;
    this.authMountAttribution = SERIALIZED_ACTIVITY_RESPONSE.by_namespace.find(
      (ns) => ns.label === 'ns1'
    ).mounts;
  });

  test('it renders empty state with no data', async function (assert) {
    await render(hbs`
      <Clients::Attribution />
    `);

    assert.dom(GENERAL.emptyStateTitle).hasText('No data found');
    assert.dom(CLIENTS_ATTRIBUTION.title).hasText('Namespace attribution', 'uses default noun');
    assert.dom(CLIENTS_ATTRIBUTION.timestamp).hasNoText();
  });

  test('it updates language based on noun', async function (assert) {
    this.noun = '';
    await render(hbs`
      <Clients::Attribution
        @noun={{this.noun}}
        @attribution={{this.namespaceAttribution}}
        @responseTimestamp={{this.timestamp}}
        />
    `);
    assert.dom(CLIENTS_ATTRIBUTION.timestamp).includesText('Updated Apr 3');

    // when noun is blank, uses default
    assert.dom(CLIENTS_ATTRIBUTION.title).hasText('Namespace attribution');
    assert
      .dom(CLIENTS_ATTRIBUTION.description)
      .hasText(
        'This data shows the top ten namespaces by total clients and can be used to understand where clients are originating. Namespaces are identified by path.'
      );
    assert
      .dom(CLIENTS_ATTRIBUTION.subtext)
      .hasText('This data shows the top ten namespaces by total clients for the date range selected.');

    // when noun is mount
    this.set('noun', 'mount');
    assert.dom(CLIENTS_ATTRIBUTION.title).hasText('Mount attribution');
    assert
      .dom(CLIENTS_ATTRIBUTION.description)
      .hasText(
        'This data shows the top ten mounts by client count within this namespace, and can be used to understand where clients are originating. Mounts are organized by path.'
      );
    assert
      .dom(CLIENTS_ATTRIBUTION.subtext)
      .hasText(
        'The total clients used by the mounts for this date range. This number is useful for identifying overall usage volume.'
      );

    // when noun is namespace
    this.set('noun', 'namespace');
    assert.dom(CLIENTS_ATTRIBUTION.title).hasText('Namespace attribution');
    assert
      .dom(CLIENTS_ATTRIBUTION.description)
      .hasText(
        'This data shows the top ten namespaces by total clients and can be used to understand where clients are originating. Namespaces are identified by path.'
      );
    assert
      .dom(CLIENTS_ATTRIBUTION.subtext)
      .hasText('This data shows the top ten namespaces by total clients for the date range selected.');
  });

  test('it renders with data for namespaces', async function (assert) {
    await render(hbs`
      <Clients::Attribution
        @attribution={{this.namespaceAttribution}}
        @responseTimestamp={{this.timestamp}}
        />
    `);

    assert.dom(GENERAL.emptyStateTitle).doesNotExist();
    assert.dom(CLIENTS_ATTRIBUTION.chart).exists();
    assert.dom(CLIENTS_ATTRIBUTION.topItem).includesText('namespace').includesText('ns1');
    assert.dom(CLIENTS_ATTRIBUTION.topItemCount).includesText('namespace').includesText('18,903');
    assert
      .dom(CLIENTS_ATTRIBUTION.yLabels)
      .exists({ count: 2 }, 'bars reflect number of namespaces in single month');
    assert.dom(CLIENTS_ATTRIBUTION.yLabel).hasText('ns1root');
  });

  test('it renders with data for mounts', async function (assert) {
    await render(hbs`
      <Clients::Attribution
        @noun="mount"
        @attribution={{this.authMountAttribution}}
        />
    `);

    assert.dom(GENERAL.emptyStateTitle).doesNotExist();
    assert.dom(CLIENTS_ATTRIBUTION.chart).exists();
    assert.dom(CLIENTS_ATTRIBUTION.topItem).includesText('mount').includesText('auth/authid/0');
    assert.dom(CLIENTS_ATTRIBUTION.topItemCount).includesText('mount').includesText('8,394');
    assert
      .dom(CLIENTS_ATTRIBUTION.yLabels)
      .exists({ count: 3 }, 'bars reflect number of mounts in single month');
    assert.dom(CLIENTS_ATTRIBUTION.yLabel).hasText('auth/authid/0pki-engine-0kvv2-engine-0');
  });

  test('it shows secret syncs when flag is on', async function (assert) {
    this.isSecretsSyncActivated = true;
    await render(hbs`
      <Clients::Attribution
        @attribution={{this.namespaceAttribution}}
        @responseTimestamp={{this.timestamp}}
        @isSecretsSyncActivated={{true}}
        />
    `);

    assert.dom('[data-test-group="secret_syncs"] rect').exists({ count: 2 });
  });

  test('it hids secret syncs when flag is off or missing', async function (assert) {
    this.isSecretsSyncActivated = true;
    await render(hbs`
      <Clients::Attribution
        @attribution={{this.namespaceAttribution}}
        @responseTimestamp={{this.timestamp}}
        />
    `);

    assert.dom('[data-test-group="secret_syncs"]').doesNotExist();
  });

  test('it sorts and limits before rendering bars', async function (assert) {
    this.tooManyAttributions = Array(15)
      .fill(null)
      .map((_, idx) => {
        const attr = { label: `ns${idx}` };
        CLIENT_TYPES.forEach((type) => {
          attr[type] = 10 + idx;
        });
        return attr;
      });
    await render(hbs`
      <Clients::Attribution
        @attribution={{this.tooManyAttributions}}
        />
    `);
    assert.dom(CLIENTS_ATTRIBUTION.yLabels).exists({ count: 10 }, 'only 10 bars are shown');
    assert.dom(CLIENTS_ATTRIBUTION.topItem).includesText('ns14');
  });
});
