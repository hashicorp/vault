import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_DETAILS = {
  state: 'stream-wals',
  connection_state: 'ready',
  lastWAL: 0,
  lastRemoteWAL: 10,
  delta: 10,
};

const STATE_ERROR = {
  state: 'idle',
  connection_state: 'ready',
  lastWAL: 0,
  lastRemoteWAL: 10,
  delta: 10,
};

const CONNECTION_ERROR = {
  state: 'idle',
  connection_state: 'transient-failure',
  lastWAL: 0,
  lastRemoteWAL: 10,
  delta: 10,
};

module('Integration | Enterprise | Component | replication-secondary-card', function(hooks) {
  setupRenderingTest(hooks);
  const title = 'States';
  const hasErrorClass = false;

  hooks.beforeEach(function() {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('stateError', STATE_ERROR);
    this.set('connectionError', CONNECTION_ERROR);
    this.set('hasErrorClass', hasErrorClass);
    this.set('title', title);
  });

  test('it renders', async function(assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{replicationDetails}} @title={{title}} @hasErrorClass={{hasErrorClass}}/>`
    );

    assert.dom('[data-test-replication-secondary-card]').exists();
  });
  test('it renders with state, and connection set', async function(assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{replicationDetails}} @title={{title}} @hasErrorClass={{hasErrorClass}}/>`
    );
    assert.dom('[data-test-state]').includesText(REPLICATION_DETAILS.state, `shows the correct state value`);

    assert
      .dom('[data-test-connection]')
      .includesText(REPLICATION_DETAILS.connection_state, `shows the correct connection value`);
  });

  test('it renders with lastWAL, lastRemoteWAL and delta set when title is not States', async function(assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{replicationDetails}} @hasErrorClass={{hasErrorClass}}/>`
    );
    assert
      .dom('[data-test-lastWAL]')
      .includesText(REPLICATION_DETAILS.lastWAL, `shows the correct lastWAL value`);
    assert
      .dom('[data-test-lastRemoteWAL]')
      .includesText(REPLICATION_DETAILS.lastRemoteWAL, `shows the correct lastRemoteWAL value`);
    assert.dom('[data-test-delta]').includesText(REPLICATION_DETAILS.delta, `shows the correct delta value`);
  });

  test('it renders tooltip with check-circle-outline when state is stream-wals', async function(assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{replicationDetails}} @title={{title}} />`
    );
    assert.dom('[data-test-glyph]').hasClass('has-text-success', `shows success icon`);
  });

  test('it renders hasErrorMessage when state is idle', async function(assert) {
    await render(hbs`<ReplicationSecondaryCard @replicationDetails={{stateError}} @title={{title}} />`);
    assert.dom('[data-test-error]').hasClass('has-text-danger', `shows error class`);
  });

  test('it renders hasErrorMessage when connection is shutdown', async function(assert) {
    await render(hbs`<ReplicationSecondaryCard @replicationDetails={{connectionError}} @title={{title}} />`);
    assert.dom('[data-test-error]').hasClass('has-text-danger', `shows error class`);
  });
});
