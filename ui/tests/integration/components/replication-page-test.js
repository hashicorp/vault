/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupMirage } from 'ember-cli-mirage/test-support';

const MODEL = {
  replicationMode: 'dr',
  dr: { mode: 'primary' },
  performance: { mode: 'primary' },
  replicationAttrs: {
    mode: 'secondary',
    clusterId: '12ab',
    replicationDisabled: false,
  },
};

module('Integration | Component | replication-page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.set('model', MODEL);
  });

  test('it renders', async function (assert) {
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-replication-page]').exists();
    assert.dom('[data-test-layout-loading]').doesNotExist();
  });

  test('it renders loader when either clusterId is unknown or mode is bootstrapping', async function (assert) {
    this.set('model.replicationAttrs.clusterId', '');
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-layout-loading]').exists();

    this.set('model.replicationAttrs.clusterId', '123456');
    this.set('model.replicationAttrs.mode', 'bootstrapping');
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-layout-loading]').exists();
  });

  test.skip('it re-fetches data when replication mode changes', async function (assert) {
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
      hbs`<ReplicationPage @model={{this.model}} as |Page|><Page.header @showTabs={{true}} /></ReplicationPage>`
    );
    await waitFor('[data-test-replication-title]');
    // Title has spaces and newlines, so we can't use hasText because it won't match exactly
    assert.dom('[data-test-replication-title]').includesText('Disaster Recovery');
    this.set('model', { ...MODEL, replicationMode: 'performance' });
    await settled();
    assert.dom('[data-test-replication-title]').includesText('Performance');
  });
});
