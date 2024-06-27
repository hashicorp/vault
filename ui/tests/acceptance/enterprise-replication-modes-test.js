/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentURL, settled, visit, waitFor } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { STATUS_DISABLED_RESPONSE, mockReplicationBlock } from 'vault/tests/helpers/replication';

const s = {
  navLink: (title) => `[data-test-sidebar-nav-link="${title}"]`,
  title: '[data-test-replication-title]',
  detailLink: (mode) => `[data-test-replication-details-link="${mode}"]`,
  summaryCard: '[data-test-replication-summary-card]',
  dashboard: '[data-test-replication-dashboard]',
  enableForm: '[data-test-replication-enable-form]',
  knownSecondary: (name) => `[data-test-secondaries-node="${name}"]`,
};

module('Acceptance | Enterprise | replication modes', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.setupMocks = (payload) => {
      this.server.get('sys/replication/status', () => ({
        data: payload,
      }));
      return authPage.login();
    };
  });

  test('replication page when unsupported', async function (assert) {
    this.server.get('sys/replication/status', () => ({
      data: {
        mode: 'unsupported',
      },
    }));

    await authPage.login();
    await visit('/vault/replication');
    assert.dom(s.title).hasText('Replication unsupported', 'it shows the unsupported view');

    // Nav links
    assert.dom(s.navLink('Performance')).doesNotExist('hides performance link');
    assert.dom(s.navLink('Disaster Recovery')).doesNotExist('hides dr link');
  });

  test('replication page when disabled', async function (assert) {
    await this.setupMocks(STATUS_DISABLED_RESPONSE);
    await visit('/vault/replication');
    assert.dom(s.title).hasText('Enable Replication', 'it shows the enable view');

    // Nav links
    assert.dom(s.navLink('Performance')).exists('shows performance link');
    assert.dom(s.navLink('Disaster Recovery')).exists('shows dr link');

    await click(s.navLink('Performance'));
    assert.strictEqual(currentURL(), '/vault/replication/performance', 'it navigates to the correct page');
    await settled();
    assert.dom(s.enableForm).exists();

    await click(s.navLink('Disaster Recovery'));
    assert.dom(s.title).hasText('Enable Disaster Recovery Replication', 'it shows the enable view for dr');
  });

  ['primary', 'secondary'].forEach((mode) => {
    test(`replication page when perf ${mode} only`, async function (assert) {
      await this.setupMocks({
        dr: mockReplicationBlock(),
        performance: mockReplicationBlock(mode),
      });
      await visit('/vault/replication');

      assert.dom(s.title).hasText('Replication', 'it shows default view');
      assert.dom(s.detailLink('performance')).hasText('Details', 'CTA to see performance details');
      assert.dom(s.detailLink('dr')).hasText('Enable', 'CTA to enable dr');

      // Nav links
      assert.dom(s.navLink('Performance')).exists('shows performance link');
      assert.dom(s.navLink('Disaster Recovery')).exists('shows dr link');

      await click(s.navLink('Performance'));
      assert.strictEqual(currentURL(), `/vault/replication/performance`, `goes to correct URL`);
      await waitFor(s.dashboard);
      assert.dom(s.dashboard).exists(`it shows the replication dashboard`);

      await click(s.navLink('Disaster Recovery'));
      assert.dom(s.title).hasText('Enable Disaster Recovery Replication', 'it shows the dr title');
      assert.dom(s.enableForm).exists('it shows the enable view for dr');
    });
  });
  // DR secondary mode is a whole other thing, test primary only here
  test(`replication page when dr primary only`, async function (assert) {
    await this.setupMocks({
      dr: mockReplicationBlock('primary'),
      performance: mockReplicationBlock(),
    });
    await visit('/vault/replication');
    assert.dom(s.title).hasText('Replication', 'it shows default view');
    assert.dom(s.detailLink('performance')).hasText('Enable', 'CTA to enable performance');
    assert.dom(s.detailLink('dr')).hasText('Details', 'CTA to see dr details');

    // Nav links
    assert.dom(s.navLink('Performance')).exists('shows performance link');
    assert.dom(s.navLink('Disaster Recovery')).exists('shows dr link');

    await click(s.navLink('Performance'));
    assert.strictEqual(currentURL(), `/vault/replication/performance`, `goes to correct URL`);
    await waitFor(s.enableForm);
    assert.dom(s.enableForm).exists('it shows the enable view for performance');

    await click(s.navLink('Disaster Recovery'));
    assert.dom(s.title).hasText(`Disaster Recovery primary`, 'it shows the dr title');
    assert.dom(s.dashboard).exists(`it shows the replication dashboard`);
  });

  test(`replication page both primary`, async function (assert) {
    await this.setupMocks({
      dr: mockReplicationBlock('primary'),
      performance: mockReplicationBlock('primary'),
    });
    await visit('/vault/replication');
    assert.dom(s.title).hasText('Disaster Recovery & Performance primary', 'it shows primary view');
    assert.dom(s.summaryCard).exists({ count: 2 }, 'shows 2 summary cards');

    await click(s.navLink('Performance'));
    assert.dom(s.title).hasText(`Performance primary`, `it shows the performance mode details`);
    assert.dom(s.enableForm).doesNotExist();

    await click(s.navLink('Disaster Recovery'));
    assert.dom(s.title).hasText(`Disaster Recovery primary`, 'it shows the dr mode details');
    assert.dom(s.enableForm).doesNotExist();
  });
});
