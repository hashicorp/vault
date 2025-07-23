/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, findAll, currentURL, visit, settled, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import { addSecondary, disableReplication, enableReplication } from 'vault/tests/helpers/replication';
import { addDays } from 'date-fns';
import formatRFC3339 from 'date-fns/formatRFC3339';
import timestamp from 'core/utils/timestamp';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | Enterprise | replication', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    await settled();
    await disableReplication('dr');
    await settled();
    await disableReplication('performance');
    await settled();
  });

  hooks.afterEach(async function () {
    await disableReplication('dr');
    await settled();
    await disableReplication('performance');
    await settled();
  });

  test('DR primary: enables primary and adds secondary', async function (assert) {
    const secondaryName = 'drSecondary';

    await enableReplication('dr', 'primary');
    await pollCluster(this.owner);
    await settled();
    assert
      .dom('[data-test-replication-title="Disaster Recovery"]')
      .includesText('Disaster Recovery', 'it displays the replication type correctly');
    assert
      .dom('[data-test-replication-mode-display]')
      .includesText('primary', 'it displays the cluster mode correctly');

    await addSecondary(secondaryName);
    await pollCluster(this.owner);
    await settled();

    await click('[data-test-replication-link="secondaries"]');
    assert
      .dom('[data-test-secondary-name]')
      .includesText(secondaryName, 'it displays the secondary in the list of known secondaries');

    // verify overflow-wrap class is applied to secondary name
    assert
      .dom('[data-test-secondary-name]')
      .hasClass('overflow-wrap', 'it applies the overflow-wrap class to the secondary name');
  });

  test('DR primary: shows demotion warning when Performance replication is active', async function (assert) {
    await visit('vault/replication/performance');

    // enable perf replication
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click(GENERAL.submitButton);

    await pollCluster(this.owner);

    // enable dr replication
    await visit('/vault/replication/dr');

    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');

    await click(GENERAL.submitButton);

    await pollCluster(this.owner);
    await visit('/vault/replication/dr/manage');

    await click(GENERAL.button('demote'));
    assert.ok(findAll('[data-test-demote-warning]').length, 'displays the demotion warning');
  });

  test('DR primary: shows empty state when secondary mode is not enabled and we navigated to the secondary details page', async function (assert) {
    // enable dr replication

    await visit('/vault/replication/dr');

    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click(GENERAL.submitButton);
    await settled(); // eslint-disable-line
    await pollCluster(this.owner);
    await visit('/vault/replication-dr-promote/details');

    assert
      .dom('[data-test-component="empty-state"]')
      .exists('Empty state is shown when no secondary is configured');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'This Disaster Recovery secondary has not been enabled. You can do so from the Disaster Recovery Primary.',
        'Renders the correct message for when a primary is enabled but no secondary is configured and we have navigated to the secondary details page.'
      );
  });

  test('DR secondary: shows empty state when replication is not enabled', async function (assert) {
    await visit('/vault/replication-dr-promote/details');

    assert.dom('[data-test-component="empty-state"]').exists();
    assert
      .dom(GENERAL.emptyStateTitle)
      .includesText('Disaster Recovery secondary not set up', 'shows the correct title of the empty state');

    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'This cluster has not been enabled as a Disaster Recovery Secondary. You can do so by enabling replication and adding a secondary from the Disaster Recovery Primary.',
        'renders default message specific to when no replication is enabled'
      );
  });

  test('Performance primary: add secondary and delete config', async function (assert) {
    const secondaryName = `performanceSecondary`;
    const mode = 'deny';

    await enableReplication('performance', 'primary');
    await pollCluster(this.owner);
    await settled();
    // confirm that the details dashboard shows
    await waitFor('[data-test-replication-dashboard]', 2000);
    await addSecondary(secondaryName, mode);
    await pollCluster(this.owner);
    await settled();
    await click('[data-test-replication-link="secondaries"]');
    await click(GENERAL.menuTrigger);
    await click('[data-test-replication-path-filter-link]');
    assert.strictEqual(
      currentURL(),
      `/vault/replication/performance/secondaries/config/show/${secondaryName}`
    );
    assert.dom('[data-test-mount-config-mode]').includesText(mode, 'show page renders the correct mode');

    // delete config by choosing "no filter" on the edit screen
    await click('[data-test-replication-link="edit-mount-config"]');
    await click('#no-filtering');
    await click('[data-test-config-save]');
    await settled(); // eslint-disable-line
    assert.strictEqual(
      currentURL(),
      `/vault/replication/performance/secondaries`,
      'redirects to the secondaries page'
    );
    // nav back to details page and confirm secondary is in the known secondaries table
    await click('[data-test-replication-link="details"]');

    assert
      .dom(`[data-test-secondaries-node=${secondaryName}]`)
      .exists('shows a table row the recently added secondary');
  });

  test('Performance primary: manages secondary token generation and TTL configuration', async function (assert) {
    const secondaryNameFirst = 'firstSecondary';
    const secondaryNameSecond = 'secondSecondary';

    // enable perf replication
    await enableReplication('performance', 'primary');
    await pollCluster(this.owner);
    await settled();
    // confirm that the details dashboard shows
    await addSecondary(secondaryNameFirst);
    await pollCluster(this.owner);
    await settled();

    // checks on secondary token modal
    assert.dom('#replication-copy-token-modal').exists();
    assert.dom(GENERAL.inlineError).hasText('Copy token to dismiss modal');
    assert.dom(GENERAL.infoRowValue('TTL')).hasText('1800s', 'shows the correct TTL of 1800s');
    // click off the modal to make sure you don't just have to click on the copy-close button to copy the token
    assert.dom(GENERAL.cancelButton).isDisabled('cancel is disabled');
    await click(GENERAL.button('Copy token'));
    assert.dom(GENERAL.cancelButton).isEnabled('cancel is enabled after token is copied');
    await click(GENERAL.cancelButton);

    // add another secondary not using the default ttl
    await click('[data-test-secondary-add]');
    await fillIn('[data-test-input="Secondary ID"]', secondaryNameSecond);
    await click(GENERAL.toggleInput('Time to Live (TTL) for generated secondary token'));
    await fillIn('[data-test-ttl-value]', 3);
    await click('[data-test-secondary-add]');

    await pollCluster(this.owner);
    await settled();
    assert.dom(GENERAL.infoRowValue('TTL')).hasText('180s', 'shows the correct TTL of 180s');
    await click(GENERAL.button('Copy token'));
    await click(GENERAL.cancelButton);

    // confirm you were redirected to the secondaries page
    assert.strictEqual(
      currentURL(),
      `/vault/replication/performance/secondaries`,
      'redirects to the secondaries page'
    );
    assert
      .dom('[data-test-secondary-name]')
      .includesText(secondaryNameFirst, 'it displays the secondary in the list of secondaries');
  });

  test('Performance primary: demotes primary to secondary and displays correct status', async function (assert) {
    // enable perf replication
    await enableReplication('performance', 'primary');
    await pollCluster(this.owner);
    await settled();

    // demote perf primary to a secondary
    await click('[data-test-replication-link="manage"]');

    // open demote modal
    await click(GENERAL.button('demote'));

    // enter confirmation text
    await fillIn('[data-test-confirmation-modal-input="Demote to secondary?"]', 'Performance');
    // Click confirm button
    await click(GENERAL.confirmButton);

    await pollCluster(this.owner);
    await settled();

    await click('[data-test-replication-link="details"]');
    await waitFor('[data-test-replication-dashboard]');
    assert.dom('[data-test-replication-dashboard]').exists();
    assert.dom('[data-test-selectable-card-container="secondary"]').exists();
    assert
      .dom('[data-test-replication-mode-display]')
      .hasText('secondary', 'it displays the cluster mode correctly in header');
  });

  test('Replication Dashboard: displays summary cards for both Performance and DR primaries', async function (assert) {
    await enableReplication('performance', 'primary');
    await pollCluster(this.owner);
    await settled();

    await visit('/vault/replication');

    assert
      .dom(`[data-test-replication-summary-card]`)
      .doesNotExist(`does not render replication summary card when both modes are not enabled as primary`);

    // enable DR primary replication
    await click('[data-test-sidebar-nav-link="Disaster Recovery"]');
    // let the controller set replicationMode in afterModel
    await waitFor('[data-test-replication-enable-form]');
    await click(GENERAL.submitButton);

    await pollCluster(this.owner);
    await settled();

    // Breadcrumbs only load once we're in the summary mode after enabling
    await waitFor('[data-test-replication-breadcrumb]');
    // navigate using breadcrumbs back to replication.index
    assert.dom('[data-test-replication-breadcrumb]').exists('shows the replication breadcrumb (flaky)');
    await click('[data-test-replication-breadcrumb] a');

    assert
      .dom('[data-test-replication-summary-card]')
      .exists({ count: 2 }, 'renders two replication-summary-card components');

    // navigate to details page using the "Details" link
    await click('[data-test-manage-link="Disaster Recovery"]');

    assert
      .dom('[data-test-selectable-card-container="primary"]')
      .exists('shows the correct card on the details dashboard');
    assert.strictEqual(currentURL(), '/vault/replication/dr');
  });

  module('DR Secondary Mirage-only test meep', function () {
    setupMirage(hooks);
    // nested module to avoid running the beforeEach and afterEach hooks. QUnit will not apply the out hook unless specifically reused.
    test('does not run analytics service in DR secondary state', async function (assert) {
      // Override analytics tracking
      const analyticsService = this.owner.lookup('service:analytics');
      analyticsService.trackEvent = () => {
        throw new Error(
          'Analytics should not be called in DR secondary mode. Documenting this behavior for clarity as the analytics service requests a promise to resolve and that breaks the DR secondary flow.'
        );
      };

      // Override sys/health to be DR secondary
      this.server.get('/sys/health', () => {
        return new Response(
          200,
          {},
          {
            enterprise: true,
            initialized: true,
            sealed: false,
            standby: false,
            license: {
              expiry_time: formatRFC3339(addDays(timestamp.now(), 33)),
              state: 'stored',
            },
            performance_standby: false,
            replication_performance_mode: 'disabled',
            replication_dr_mode: 'secondary',
            server_time_utc: 1753289940,
            version: '1.21.0+ent',
            cluster_name: 'vault-cluster-e779cd7c',
            cluster_id: 'f877cae2-7a56-159a-7a53-73d273738256',
            last_wal: 121,
          }
        );
      });

      await visit('/vault/replication-dr-promote');

      assert
        .dom('[data-test-dr-secondary-banner]')
        .exists('Displays DR secondary banner or read-only UI element');
    });
  });
});
