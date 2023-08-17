/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { CLUSTER_STATES } from 'core/helpers/cluster-states';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | replication-primary-card', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'replication');

  test('it renders', async function (assert) {
    const title = 'Last WAL';
    const description = 'WALL-E';
    const metric = '3000';

    this.set('title', title);
    this.set('description', description);
    this.set('metric', metric);

    await render(
      hbs`
      <ReplicationPrimaryCard
        @title={{this.title}}
        @description={{this.description}}
        @metric='3000'
      />`,
      { owner: this.engine }
    );

    assert.dom('[data-test-hasError]').doesNotExist('shows no error for non-State cards');

    assert.dom('.last-wal').includesText(title);
    assert.dom('[data-test-description]').includesText(description);
    assert.dom('[data-test-metric]').includesText(metric);
  });

  Object.keys(CLUSTER_STATES).forEach((state) => {
    test(`it renders a card when cluster has the ${state} state`, async function (assert) {
      this.set('glyph', CLUSTER_STATES[state].glyph);
      this.set('state', state);

      await render(
        hbs`
        <ReplicationPrimaryCard
          @title='State'
          @description='Updated every ten seconds.'
          @glyph={{this.glyph}}
          @metric={{this.state}}
        />`,
        { owner: this.engine }
      );

      if (CLUSTER_STATES[state].isOk) {
        assert.dom('[data-test-hasError]').doesNotExist();
        assert.dom('[data-test-cluster-state-icon]').exists('shows an icon if state is ok');
      } else {
        assert.dom('[data-test-hasError]').exists('shows an error if the cluster state is not ok');
        assert.dom('[data-test-cluster-state-icon]').doesNotExist('does not show an icon if state is not ok');
      }
    });
  });
});
