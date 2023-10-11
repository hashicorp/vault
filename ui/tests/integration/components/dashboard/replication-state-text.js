/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import SELECTORS from 'vault/tests/helpers/components/dashboard/replication-card';

module('Integration | Component | dashboard/replication-state-text', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.name = 'DR Primary';
    this.clusterState = {
      glyph: 'circle-check',
      isOk: true,
    };
  });

  test('it displays replication states', async function (assert) {
    await render(
      hbs`
        <Dashboard::ReplicationStateText 
          @name={{this.name}} 
          @version={{this.version}} 
          @subText={{this.subText}} 
          @clusterStates={{this.clusterStates}} />
          `
    );
    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('running');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle')).exists();

    this.name = 'DR Primary';
    this.clusterState = {
      glyph: 'x-circle',
      isOk: false,
    };
    await render(
      hbs`
        <Dashboard::ReplicationStateText 
          @name={{this.name}} 
          @version={{this.version}} 
          @subText={{this.subText}} 
          @clusterStates={{this.clusterStates}} />
          `
    );
    assert
      .dom(SELECTORS.getReplicationTitle('dr-perf', 'Performance primary'))
      .hasText('Performance primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'Performance primary')).hasText('running');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Performance primary', 'x-circle')).exists();
    assert
      .dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Performance primary', 'x-circle'))
      .hasClass('has-text-danger');
  });
});
