/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import syncHandler from 'vault/mirage/handlers/sync';
import { setupModels } from 'vault/tests/helpers/sync/setup-models';
import hbs from 'htmlbars-inline-precompile';
import { click, render } from '@ember/test-helpers';
import sinon from 'sinon';

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
      syncHandler(this.server);
      this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
      sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

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
          'Select secrets from existing KV version 2 engines and sync them to the destination.',
          'Empty state message renders'
        );
      assert.dom(PAGE.emptyStateActions).hasText('Sync secrets', 'Empty state action renders');
    });

    test('it should render list item details', async function (assert) {
      const { list } = PAGE.associations;
      assert.dom(list.name).hasText('my-secret', 'Association mount/secret renders as name');
      assert.dom(list.status).hasText('Synced', 'Association status renders');
      assert
        .dom(list.updated)
        .hasText('last updated on September 20th 2023, 10:51:53 AM', 'Last synced datetime renders');
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
      assert.dom(menu.view).hasText('View secret', 'View secret menu action renders');
      assert.dom(menu.unsync).hasText('Unsync', 'Unsync menu action renders');

      await click(menu.sync);

      await click(menu.unsync);
      await click(PAGE.confirmButton);
    });

    module('flash messages', function (hooks) {
      hooks.beforeEach(function () {
        const flashMessages = this.owner.lookup('service:flash-messages');

        this.flashSuccessSpy = sinon.spy(flashMessages, 'success');
        this.flashDangerSpy = sinon.spy(flashMessages, 'danger');
      });

      test('unsync should render flash messages', async function (assert) {
        await click(PAGE.menuTrigger);

        const { menu } = PAGE.associations.list;
        await click(menu.unsync);
        await click(PAGE.confirmButton);

        assert.true(
          this.flashSuccessSpy.calledWith('Unsync operation initiated.'),
          'Success message is displayed'
        );
        assert.true(this.flashDangerSpy.notCalled);
      });

      test('sync now should render flash messages', async function (assert) {
        await click(PAGE.menuTrigger);

        const { menu } = PAGE.associations.list;
        await click(menu.sync);

        assert.true(
          this.flashSuccessSpy.calledWith('Sync operation initiated.'),
          'Success message is displayed'
        );
        assert.true(this.flashDangerSpy.notCalled);
      });

      test('it should show an error message when sync fails', async function (assert) {
        this.server.post('/sys/sync/destinations/:type/:name/associations/set', () => {
          return new Response(500);
        });

        await click(PAGE.menuTrigger);

        const { menu } = PAGE.associations.list;
        await click(menu.sync);

        assert.true(this.flashSuccessSpy.notCalled);
        assert.true(this.flashDangerSpy.calledOnce);
      });
    });
  }
);
