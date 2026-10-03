/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | replication-page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.model = {
      replicationMode: 'dr',
      dr: { mode: 'primary' },
      performance: { mode: 'primary' },
      replicationAttrs: {
        mode: 'secondary',
        clusterId: '12ab',
        replicationDisabled: false,
      },
    };
  });

  test('it renders', async function (assert) {
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-replication-page]').exists();
    assert.dom('[data-test-layout-loading]').doesNotExist();
  });

  test('it renders loader when either clusterId is unknown or mode is bootstrapping', async function (assert) {
    this.model.replicationAttrs.clusterId = '';
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-layout-loading]').exists();

    this.model.replicationAttrs.clusterId = '123456';
    this.model.replicationAttrs.mode = 'bootstrapping';
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-layout-loading]').exists();
  });

  test('it re-fetches data when replication mode changes', async function (assert) {
    assert.expect(4);
    this.server.get('sys/replication/:mode/status', (schema, req) => {
      assert.strictEqual(
        req.params.mode,
        this.model.replicationMode,
        `fetchStatus called with correct mode: ${this.model.replicationMode}`
      );
      return {
        data: {
          mode: 'primary',
        },
      };
    });
    await render(
      hbs`<ReplicationPage @model={{this.model}} as |Page|><Page.header @showTabs={{false}} /></ReplicationPage>`
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Disaster Recovery');

    this.model.replicationMode = 'performance';
    await render(
      hbs`<ReplicationPage @model={{this.model}} as |Page|><Page.header @showTabs={{false}} /></ReplicationPage>`
    );

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Performance');
  });
});
