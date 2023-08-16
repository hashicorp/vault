/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const TITLE = 'Disaster Recovery';

const REPLICATION_DETAILS = {
  dr: {
    state: 'running',
    lastWAL: 10,
    knownSecondaries: ['https://127.0.0.1:8201', 'https://127.0.0.1:8202'],
  },
  performance: {
    state: 'running',
    lastWAL: 20,
    knownSecondaries: ['https://127.0.0.1:8201'],
  },
};

module('Integration | Component | replication-summary-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('title', TITLE);
  });

  test('it renders', async function (assert) {
    await render(
      hbs`<ReplicationSummaryCard @replicationDetails={{this.replicationDetails}} @title={{this.title}} />`
    );
    assert.dom('[data-test-replication-summary-card]').exists();
    assert
      .dom('[data-test-lastWAL]')
      .includesText(REPLICATION_DETAILS.dr.lastWAL, `shows the correct lastWAL value`);

    const knownSecondaries = REPLICATION_DETAILS.dr.knownSecondaries.length;
    assert
      .dom('[data-test-known-secondaries]')
      .includesText(knownSecondaries, `shows the correct computed value of the known secondaries count`);
  });

  test('it shows the correct lastWAL and knownSecondaries when title is Performance', async function (assert) {
    await render(
      hbs`<ReplicationSummaryCard @replicationDetails={{this.replicationDetails}} @title="Performance" />`
    );
    assert
      .dom('[data-test-lastWAL]')
      .includesText(REPLICATION_DETAILS.performance.lastWAL, `shows the correct lastWAL value`);

    const knownSecondaries = REPLICATION_DETAILS.performance.knownSecondaries.length;
    assert
      .dom('[data-test-known-secondaries]')
      .includesText(knownSecondaries, `shows the correct computed value of the known secondaries count`);
  });
});
