import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { pauseTest, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DATATEST = EmberObject.create({
  dr: {
    mode: 'secondary',
    rm: {
      mode: 'dr',
    },
    clusterIdDisplay: 12345,
  },
  unsealed: 'good',
});
const DATA = {
  dr: {
    mode: 'secondary',
    rm: {
      mode: 'dr',
    },
    clusterIdDisplay: 12345,
  },
  unsealed: 'good',
};

module('Integration | Enterprise | Component | replication-header', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', DATA);
  });
  console.log(DATATEST);
  test('it renders', async function(assert) {
    await render(
      hbs`<ReplicationHeader @data={{data}} @isSecondary={{true}} @title={{'Disaster Recover'}}/>`
    );
    //  ARG I should be able to remove isSecondary, but I get mode undefined, when it should be defined.
    await pauseTest();
    assert.dom('[data-test-replication-header]').exists();
  });

  test('it renders with clusterId and mode when set', async function(assert) {
    await render(hbs`<ReplicationHeader
      @data={{data}}
      @isSecondary={{true}}
      @showTabs={{true}}
      @title={{'Disaster Recover'}}/>`);

    assert
      .dom('[data-test-clusterId]')
      .includesText(DATA.dr.clusterIdDisplay, `shows the correct clusterId value`);
    assert.dom('[data-test-mode]').includesText(DATA.dr.mode, `shows the correct mode value`);
  });

  test('it does not show tabs when showTabs is not set', async function(assert) {
    await render(hbs`<ReplicationHeader
      @data={{data}}
      @isSecondary={{true}}
      @showTabs={{false}}
      @title={{'Disaster Recover'}}/>`);
    await pauseTest();

    assert.dom('[data-test-tabs]').doesNotExist();
  });
});
