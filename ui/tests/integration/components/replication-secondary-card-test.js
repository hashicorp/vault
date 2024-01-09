/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const TITLE = 'Status';

const REPLICATION_DETAILS = {
  state: 'stream-wals',
  connection_state: 'ready',
  lastRemoteWAL: 10,
  knownPrimaryClusterAddrs: ['https://127.0.0.1:8201', 'https://127.0.0.1:8202'],
  primaries: [
    {
      api_address: 'http://127.0.0.1:8201',
    },
  ],
};

module('Integration | Component | replication-secondary-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('title', TITLE);
    // TODO: remove Tooltip/ember-basic-dropdown
    setRunOptions({
      rules: {
        'aria-command-name': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it renders', async function (assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.replicationDetails}} @title={{this.title}} />`
    );
    assert.dom('[data-test-replication-secondary-card]').exists();
    assert.dom('[data-test-state]').includesText(REPLICATION_DETAILS.state, `shows the correct state value`);
    assert
      .dom('[data-test-connection]')
      .includesText(REPLICATION_DETAILS.connection_state, `shows the correct connection value`);
  });

  test('it renders the Primary Cluster card when title is not Status', async function (assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.replicationDetails}} @title='Primary cluster'/>`
    );

    assert.dom('[data-test-info-table]').exists('it shows the known primary cluster details');

    const expectedUrl = `${REPLICATION_DETAILS.primaries[0].api_address}/ui/`;
    assert
      .dom('[data-test-primary-link]')
      .hasAttribute('href', expectedUrl, 'it renders a link to the primary cluster UI');
  });

  test('it does not render a link to the primary cluster UI when the primary api address or known primaries are unknown', async function (assert) {
    this.set('replicationDetails', {});
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.replicationDetails}} @title='Primary cluster'/>`
    );

    assert.dom('[data-test-primary-link]').doesNotExist();
  });

  test('it renders with emptyState if no knownPrimaryClusterAddrs are set', async function (assert) {
    this.set('replicationDetails', []);
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.replicationDetails}} @title='Primary cluster'/>`
    );
    assert.dom('[data-test-empty-state]').exists('shows empty state');
  });

  test('it renders tooltip with check-circle-outline when state is stream-wals', async function (assert) {
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.replicationDetails}} @title={{this.title}} />`
    );
    assert.dom('[data-test-glyph]').hasClass('has-text-success', `shows success icon`);
  });

  test('it renders hasErrorMessage when state is idle', async function (assert) {
    const stateError = {
      state: 'idle',
      connection_state: 'ready',
      lastRemoteWAL: 10,
    };

    this.set('stateError', stateError);
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.stateError}} @title={{this.title}} />`
    );
    assert.dom('[data-test-error]').includesText('state', 'show correct error title');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('Please check your server logs.', 'show correct error message');
  });

  test('it renders hasErrorMessage when connection is transient_failure', async function (assert) {
    const connectionError = {
      state: 'stream-wals',
      connection_state: 'transient_failure',
      lastRemoteWAL: 10,
    };

    this.set('connectionError', connectionError);
    await render(
      hbs`<ReplicationSecondaryCard @replicationDetails={{this.connectionError}} @title={{this.title}} />`
    );
    assert.dom('[data-test-error]').includesText('connection_state', 'show correct error title');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText(
        'There has been some transient failure. Your cluster will eventually switch back to connection and try to establish a connection again.',
        'show correct error message'
      );
  });
});
