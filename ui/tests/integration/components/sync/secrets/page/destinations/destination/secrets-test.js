/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupModels } from 'vault/tests/helpers/sync/setup-models';
import hbs from 'htmlbars-inline-precompile';
import { click, render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module(
  'Integration | Component | sync | Secrets::Page::Destinations::Destination::Secrets',
  function (hooks) {
    setupRenderingTest(hooks);
    setupEngine(hooks, 'sync');
    setupMirage(hooks);
    setupModels(hooks);

    hooks.beforeEach(async function () {
      this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

      await render(
        hbs`
        <Secrets::Page::Destinations::Destination::Secrets
          @destination={{this.destination}}
          @associations={{this.associations}}
        />
      `,
        { owner: this.engine }
      );
    });

    test('it should render DestinationHeader component', async function (assert) {
      assert.dom(PAGE.title).includesText('us-west-1', 'DestinationHeader component renders');
    });

    test('it should render empty list state', async function (assert) {
      this.set('associations.meta.filteredTotal', 0);
      assert.dom(PAGE.emptyStateTitle).hasText('No synced secrets yet', 'Empty state title renders');
      assert
        .dom(PAGE.emptyStateMessage)
        .hasText(
          'Select secrets from existing K/V engines and sync them to the destination.',
          'Empty state message renders'
        );
      assert.dom(PAGE.emptyStateActions).hasText('Sync secret', 'Empty state action renders');
    });

    test('it should render list item details', async function (assert) {
      const { list } = PAGE.associations;
      assert.dom(list.name).hasText('kv/my-secret', 'Association mount/secret renders as name');
      assert.dom(list.status).hasText('SYNCED', 'Association status renders');
      assert
        .dom(list.updated)
        .hasText('last synced on September 20th 2023, 8:51:53 AM', 'Last synced datetime renders');
    });

    test('it should render list item menu actions', async function (assert) {
      assert.expect(5);

      this.server.post('/sys/sync/destinations/aws-sm/us-west-1/associations/:action', (schema, req) => {
        const { action } = req.params;
        const operation = { set: 'sync', remove: 'unsync' }[action] || null;
        assert.ok(operation, `Request made to ${operation} secret`);
      });

      await click(PAGE.menuTrigger);

      const { menu } = PAGE.associations.list;
      assert.dom(menu.sync).hasText('Sync now', 'Sync menu action renders');
      assert.dom(menu.edit).hasText('Edit secret', 'Edit secret menu action renders');
      assert.dom(menu.unsync).hasText('Unsync', 'Unsync menu action renders');

      await click(menu.sync);

      await click(menu.unsync);
      await click(PAGE.confirmButton);
    });
  }
);
