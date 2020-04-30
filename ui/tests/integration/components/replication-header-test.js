import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

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

const TITLE = 'Disaster Recovery';

module('Integration | Enterprise | Component | replication-header', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', DATA);
    this.set('title', TITLE);
    this.set('isSecondary', true);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}}/>`);

    assert.dom('[data-test-replication-header]').exists();
  });

  test('it renders with clusterId and mode when set', async function(assert) {
    await render(hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}}/>`);

    assert
      .dom('[data-test-clusterId]')
      .includesText(DATA.dr.clusterIdDisplay, `shows the correct clusterId value`);

    assert.dom('[data-test-mode]').includesText(DATA.dr.mode, `shows the correct mode value`);
  });

  test('it does not show tabs when showTabs is not set', async function(assert) {
    await render(hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}}/>`);

    assert.dom('[data-test-tabs]').doesNotExist();
  });
});
