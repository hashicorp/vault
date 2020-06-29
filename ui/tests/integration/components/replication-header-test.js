import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DATA = {
  anyReplicationEnabled: true,
  dr: {
    mode: 'secondary',
    rm: {
      mode: 'dr',
    },
  },
  unsealed: 'good',
};

const TITLE = 'Disaster Recovery';
const SECONDARY_ID = '123abc';

module('Integration | Enterprise | Component | replication-header', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', DATA);
    this.set('title', TITLE);
    this.set('isSecondary', true);
    this.set('secondaryId', SECONDARY_ID);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}}/>`);

    assert.dom('[data-test-replication-header]').exists();
  });

  test('it renders with mode and secondaryId when set', async function(assert) {
    await render(
      hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}} @secondaryId={{secondaryId}}/>`
    );

    assert.dom('[data-test-secondaryId]').includesText(SECONDARY_ID, `shows the correct secondaryId value`);
    assert.dom('[data-test-mode]').includesText('secondary', `shows the correct mode value`);
  });

  test('it does not render mode or secondaryId when replication is not enabled', async function(assert) {
    const notEnabled = { anyReplicationEnabled: false };
    const noId = null;
    this.set('data', notEnabled);
    this.set('secondaryId', noId);

    await render(
      hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}} @secondaryId={{secondaryId}}/>`
    );

    assert.dom('[data-test-secondaryId]').doesNotExist();
    assert.dom('[data-test-mode]').doesNotExist();
  });

  test('it does not show tabs when showTabs is not set', async function(assert) {
    await render(hbs`<ReplicationHeader @data={{data}} @isSecondary={{isSecondary}} @title={{title}}/>`);

    assert.dom('[data-test-tabs]').doesNotExist();
  });
});
