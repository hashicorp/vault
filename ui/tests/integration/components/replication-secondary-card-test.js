import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const TITLE = 'States';

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
  state: 'stream-wals',
  connection_state: 'transient_failure',
  lastWAL: 0,
  lastRemoteWAL: 10,
  delta: 10,
};

module('Integration | Enterprise | Component | replication-secondary-card', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('title', TITLE);
  });

  test('it renders', async function(assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{replicationDetails}} @title={{title}} />`
    );
    assert.dom('[data-test-replication-secondary-card]').exists();
    assert.dom('[data-test-state]').includesText(REPLICATION_DETAILS.state, `shows the correct state value`);
    assert
      .dom('[data-test-connection]')
      .includesText(REPLICATION_DETAILS.connection_state, `shows the correct connection value`);
  });

  test('it renders with lastWAL, lastRemoteWAL and delta set when title is not States', async function(assert) {
    await render(hbs`<ReplicationSecondaryCard @replicationDetails={{replicationDetails}} />`);
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
    this.set('stateError', STATE_ERROR);
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{stateError}} @title={{title}} @hasErrorClass={{true}} />`
    );
    assert.dom('[data-test-error]').includesText('state', 'show correct error title');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('Please check your server logs.', 'show correct error message');
  });

  test('it renders hasErrorMessage when connection is transient_failure', async function(assert) {
    this.set('connectionError', CONNECTION_ERROR);
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{connectionError}} @title={{title}} @hasErrorClass={{true}} />`
    );
    assert.dom('[data-test-error]').includesText('connection_state', 'show correct error title');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText(
        'There has been some transient failure.  Your cluster will eventually switch back to connection and try to establish a connection again.',
        'show correct error message'
      );
  });
});
