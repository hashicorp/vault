/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { toLabel } from 'core/helpers/to-label';

module('Integration | Component | SyncStatusBadge', function (hooks) {
  setupRenderingTest(hooks);

  const SYNC_STATUSES = {
    SYNCING: { icon: 'sync', color: 'neutral' },
    SYNCED: { icon: 'check-circle', color: 'success' },
    UNSYNCING: { icon: 'sync-reverse', color: 'neutral' },
    UNSYNCED: { icon: 'sync-alert', color: 'warning' },
    INTERNAL_VAULT_ERROR: { icon: 'x-circle', color: 'critical' },
    CLIENT_SIDE_ERROR: { icon: 'x-circle', color: 'critical' },
    EXTERNAL_SERVICE_ERROR: { icon: 'x-circle', color: 'critical' },
    UNKNOWN: { icon: 'help', color: 'neutral' },
  };
  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.16.0+ent';
    this.status = 'Some unaccounted for status';
    this.renderComponent = () => {
      return render(hbs`<SyncStatusBadge @status={{this.status}} data-test-badge />`);
    };
  });

  test('it should render when status does not exist', async function (assert) {
    assert.expect(2);
    await this.renderComponent();
    assert.dom(PAGE.badgeText.icon('help')).exists('renders help icon');
    assert.dom(PAGE.badgeText.text).hasText(this.status);
  });

  test('it renders badge and icon for each status type', async function (assert) {
    assert.expect(24);
    for (const status in SYNC_STATUSES) {
      this.status = status;
      const label = toLabel([status]);
      const { icon, color } = SYNC_STATUSES[status];
      await this.renderComponent();
      assert.dom(PAGE.badgeText.icon(icon)).exists(`status: ${status} renders icon: ${icon}`);
      assert.dom(PAGE.badgeText.text).hasText(label, `status: ${status} renders label: ${label}`);
      assert
        .dom('[data-test-badge]')
        .hasClass(`hds-badge--color-${color}`, `status: ${status} renders color: ${color}`);
    }
  });
});
