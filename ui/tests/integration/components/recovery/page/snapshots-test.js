/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

import { setupMirage } from 'ember-cli-mirage/test-support';

import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { dateFormat } from 'core/helpers/date-format';

const SELECTORS = {
  badge: (name) => `[data-test-badge="${name}"]`,
};

module('Integration | Component | recovery/snapshots', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.model = {
      snapshots: [],
      canLoadSnapshot: false,
    };
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
  });

  test('it displays empty state in CE', async function (assert) {
    this.version.type = 'community';
    await render(hbs`<Recovery::Page::Snapshots @model={{this.model}}/>`);
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('Secrets Recovery is an enterprise feature', 'CE empty state title renders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Secrets Recovery allows you to restore accidentally deleted or lost secrets from a snapshot. The snapshots can be provided via upload or loaded from external storage.',
        'CE empty state message renders'
      );
    assert
      .dom(GENERAL.emptyStateActions)
      .hasText('Learn more about upgrading', 'CE empty state action renders');
  });

  test('it displays empty state in non root namespace', async function (assert) {
    this.nsService = this.owner.lookup('service:namespace');
    this.nsService.path = 'parent1/child1';
    this.server.get('/sys/internal/ui/namespaces', () => {
      return { data: { keys: ['parent1/', 'parent1/child1/'] } };
    });

    await render(hbs`<Recovery::Page::Snapshots @model={{this.model}}/>`);

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('Snapshot upload is restricted', 'non root namespace empty state title renders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Snapshot uploading is only available in root namespace. Please navigate to root and upload your snapshot. ',
        'non root empty state message renders'
      );
    assert
      .dom(GENERAL.emptyStateActions)
      .hasText('Take me to root namespace', 'non root empty state action renders');

    this.nsService.path = '';
  });

  test('it displays empty state when user cannot load snapshot', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots @model={{this.model}}/>`);
    assert.dom(GENERAL.emptyStateTitle).hasText('No snapshot available', 'empty state title renders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Ready to restore secrets? Please contact your admin to either upload a snapshot or grant you uploading permissions to get started.',
        'empty state message renders'
      );
    assert
      .dom(GENERAL.emptyStateActions)
      .hasText('Learn more about Secrets Recovery', 'empty state action renders');
  });

  test('it displays empty state when user can load snapshot', async function (assert) {
    this.model.canLoadSnapshot = true;
    await render(hbs`<Recovery::Page::Snapshots @model={{this.model}}/>`);
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('Upload a snapshot to get started', 'empty state title renders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Secrets Recovery allows you to restore accidentally deleted or lost secrets from a snapshot. The snapshots can be provided via upload or loaded from external storage.',
        'empty state message renders'
      );
    assert.dom(GENERAL.emptyStateActions).hasText('Upload snapshot', 'empty state action renders');
  });

  test('it displays loaded snapshot card', async function (assert) {
    const expiryDate = dateFormat(
      [new Date('2023-09-20T10:51:53.961861096-04:00'), 'MMMM do yyyy, h:mm:ss a'],
      {}
    );

    this.model = {
      snapshot_id: '6ecc06a9-3592-26f9-40cb-8d1b511890a6',
      status: 'ready',
      expires_at: expiryDate,
    };

    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);

    assert.dom(SELECTORS.badge('status')).hasText('ready', 'status badge renders');
  });
});
