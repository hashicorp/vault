/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentURL, settled, visit, waitFor } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { STATUS_DISABLED_RESPONSE, mockReplicationBlock } from 'vault/tests/helpers/replication';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const s = {
  title: (t) => `[data-test-replication-title="${t}"]`,
  detailLink: (mode) => `[data-test-replication-details-link="${mode}"]`,
  summaryCard: '[data-test-replication-summary-card]',
  dashboard: '[data-test-replication-dashboard]',
  enableForm: '[data-test-replication-enable-form]',
  knownSecondary: (name) => `[data-test-secondaries-node="${name}"]`,
};

// wait for specific title selector as an attempt to stabilize flaky tests
async function assertTitle(assert, title) {
  await waitFor(GENERAL.hdsPageHeaderTitle);
  assert.dom(GENERAL.hdsPageHeaderTitle).hasText(title);
}

module('Acceptance | Enterprise | replication modes', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.setupMocks = (payload) => {
      this.server.get('sys/replication/status', () => ({
        data: payload,
      }));
      return login();
    };
  });

  test('replication page when unsupported', async function (assert) {
    this.server.get('sys/replication/status', () => ({
      data: {
        mode: 'unsupported',
      },
    }));

    await login();
    await visit('/vault/replication');

    await assertTitle(assert, 'Replication unsupported');

    // Nav links
    assert.dom(GENERAL.navLink('Performance')).doesNotExist('hides performance link');
    assert.dom(GENERAL.navLink('Disaster Recovery')).doesNotExist('hides dr link');
  });

  test('replication page when disabled', async function (assert) {
    await this.setupMocks(STATUS_DISABLED_RESPONSE);
    await visit('/vault/replication');

    await assertTitle(assert, 'Enable Replication');

    // Nav links
    assert.dom(GENERAL.navLink('Performance')).exists('shows performance link');
    assert.dom(GENERAL.navLink('Disaster Recovery')).exists('shows dr link');

    await click(GENERAL.navLink('Performance'));
    assert.strictEqual(currentURL(), '/vault/replication/performance', 'it navigates to the correct page');
    await settled();
    assert.dom(s.enableForm).exists();

    await click(GENERAL.navLink('Disaster Recovery'));

    await assertTitle(assert, 'Enable Disaster Recovery Replication', 'dr');
  });

  ['primary', 'secondary'].forEach((mode) => {
    test(`replication page when perf ${mode} only`, async function (assert) {
      await this.setupMocks({
        dr: mockReplicationBlock(),
        performance: mockReplicationBlock(mode),
      });
      await visit('/vault/replication');

      await assertTitle(assert, 'Replication');
      assert.dom(s.detailLink('performance')).hasText('Details', 'CTA to see performance details');
      assert.dom(s.detailLink('dr')).hasText('Enable', 'CTA to enable dr');

      // Nav links
      assert.dom(GENERAL.navLink('Performance')).exists('shows performance link');
      assert.dom(GENERAL.navLink('Disaster Recovery')).exists('shows dr link');

      await click(GENERAL.navLink('Performance'));
      assert.strictEqual(currentURL(), `/vault/replication/performance`, `goes to correct URL`);
      await waitFor(s.dashboard);
      assert.dom(s.dashboard).exists(`it shows the replication dashboard`);

      await click(GENERAL.navLink('Disaster Recovery'));
      await assertTitle(assert, 'Enable Disaster Recovery Replication', 'dr');
      assert.dom(s.enableForm).exists('it shows the enable view for dr');
    });
  });

  test('replication page when dr primary only', async function (assert) {
    await this.setupMocks({
      dr: mockReplicationBlock('primary'),
      performance: mockReplicationBlock(),
    });
    await visit('/vault/replication');
    await assertTitle(assert, 'Replication');
    assert.dom(s.detailLink('performance')).hasText('Enable', 'CTA to enable performance');
    assert.dom(s.detailLink('dr')).hasText('Details', 'CTA to see dr details');

    // Nav links
    assert.dom(GENERAL.navLink('Performance')).exists('shows performance link');
    assert.dom(GENERAL.navLink('Disaster Recovery')).exists('shows dr link');

    await click(GENERAL.navLink('Performance'));
    assert.strictEqual(currentURL(), `/vault/replication/performance`, `goes to correct URL`);
    await waitFor(s.enableForm);
    assert.dom(s.enableForm).exists('it shows the enable view for performance');

    await click(GENERAL.navLink('Disaster Recovery'));
    await assertTitle(assert, 'Disaster Recovery', 'Disaster Recovery');
    assert.dom(GENERAL.badge('primary')).exists('shows primary badge for dr');
    assert.dom(s.dashboard).exists(`it shows the replication dashboard`);
  });

  test('replication page both primary', async function (assert) {
    await this.setupMocks({
      dr: mockReplicationBlock('primary'),
      performance: mockReplicationBlock('primary'),
    });
    await visit('/vault/replication');
    await assertTitle(assert, 'Disaster Recovery & Performance', 'Disaster Recovery & Performance');
    assert.dom(GENERAL.badge('primary')).exists('shows primary badge for dr');
    assert.dom(s.summaryCard).exists({ count: 2 }, 'shows 2 summary cards');

    await click(GENERAL.navLink('Performance'));
    await assertTitle(assert, 'Performance', 'Performance');
    assert.dom(GENERAL.badge('primary')).exists('shows primary badge for dr');
    assert.dom(s.enableForm).doesNotExist();

    await click(GENERAL.navLink('Disaster Recovery'));
    await assertTitle(assert, 'Disaster Recovery', 'Disaster Recovery');
    assert.dom(GENERAL.badge('primary')).exists('shows primary badge for dr');
    assert.dom(s.enableForm).doesNotExist();
  });
});
