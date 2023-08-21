/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, focus, triggerKeyEvent, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath } from 'vault/utils/kv-path';

const MODELS = {
  secrets: [
    {
      id: kvMetadataPath('my-engine', 'my-secret'),
      path: 'my-secret',
      fullSecretPath: 'my-secret',
    },
    {
      id: kvMetadataPath('my-engine', 'my'),
      path: 'my',
      fullSecretPath: 'my',
    },
    {
      id: kvMetadataPath('my-engine', 'beep/boop/bop'),
      path: 'beep/boop/bop',
      fullSecretPath: 'beep/boop/bop',
    },
    {
      id: kvMetadataPath('my-engine', 'beep/boop-1'),
      path: 'beep/boop-1',
      fullSecretPath: 'beep/boop-1',
    },
  ],
};
const MOUNT_POINT = 'vault.cluster.secrets.backend.kv';

module('Integration | Component | kv | kv-list-filter', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.model = MODELS;
    this.mountPoint = MOUNT_POINT;
  });

  test('it renders and TAB defaults to first secret in list', async function (assert) {
    assert.expect(4);
    // mirage hook for TAB
    this.owner.lookup('service:router').reopen({
      transitionTo(route, { queryParams: { pageFilter } }) {
        assert.strictEqual(route, `${MOUNT_POINT}.list`, 'List route sent when TAB on empty input.');
        assert.deepEqual(pageFilter, 'my-secret', 'Filters to the first secret in the list.');
      },
    });

    await render(hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-component="kv-list-filter"]')
      .hasAttribute('placeholder', 'Filter secrets', 'Placeholder applied to input.');

    await focus('[data-test-component="kv-list-filter"]');
    assert.dom('[data-test-help-tab]').exists('on focus, with no filterValue, displays help text');
    // trigger tab
    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 9);
  });

  test('it filters partial matches', async function (assert) {
    assert.expect(2);
    // mirage hook for TAB
    this.owner.lookup('service:router').reopen({
      transitionTo(route, { queryParams: { pageFilter } }) {
        assert.strictEqual(route, `${MOUNT_POINT}.list`, 'List route sent when no pathToSecret.');
        assert.deepEqual(pageFilter, 'my-secret', 'Sets page filter to my-secret.');
      },
    });

    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @pageFilter="my-" />`,
      {
        owner: this.engine,
      }
    );
    // focus on input and trigger TAB
    await focus('[data-test-component="kv-list-filter"]');
    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 9);
  });

  test('it clears last item on backspace and clears to directory on esc', async function (assert) {
    assert.expect(8);
    // mirage hook for filling in the input
    this.owner.lookup('service:router').reopen({
      transitionTo(route, pathToSecret, { queryParams: { pageFilter } }) {
        assert.deepEqual(pageFilter, 'boop-', 'Sends the correct pageFilter on fillIn.');
      },
    });

    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @filterValue="beep/" @pageFilter=""/>`,
      {
        owner: this.engine,
      }
    );
    // focus on input and trigger backspace
    await focus('[data-test-component="kv-list-filter"]');
    await fillIn('[data-test-component="kv-list-filter"]', 'beep/boop-');

    this.owner.lookup('service:router').reopen({
      transitionTo(route, pathToSecret, { queryParams: { pageFilter } }) {
        assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'Correct route sent.');
        assert.strictEqual(pathToSecret, 'beep/', 'PathToSecret is the parent directory.');
        assert.deepEqual(pageFilter, 'boop', 'Clears last item in pageFilter on backspace.');
      },
    });
    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 8);
    assert.strictEqual(
      document.activeElement.id,
      'secret-filter',
      'the input still remains focused after delete.'
    );

    this.owner.lookup('service:router').reopen({
      transitionTo(route, pathToSecret, { queryParams: { pageFilter } }) {
        assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'Still on a directory route.');
        assert.strictEqual(pathToSecret, 'beep/', 'Parent directory still shown.');
        assert.deepEqual(pageFilter, null, 'Clears pageFilter on escape.');
      },
    });
    // trigger escape
    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 27);
  });

  test('it transitions create page on enter when secret path is new', async function (assert) {
    assert.expect(5);
    // mirage hook for fillIn
    this.owner.lookup('service:router').reopen({
      transitionTo(route, pathToSecret, { queryParams: { pageFilter } }) {
        assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'Still on a directory route.');
        assert.strictEqual(pathToSecret, 'beep/boop/', 'Parent directory still shown.');
        assert.deepEqual(pageFilter, 'new-secret', 'Sends correct pageFilter.');
      },
    });

    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @filterValue="beep/boop/"/>`,
      {
        owner: this.engine,
      }
    );
    // focus on input, fillIn and then trigger enter
    await focus('[data-test-component="kv-list-filter"]');
    await fillIn('[data-test-component="kv-list-filter"]', 'beep/boop/new-secret');

    // mirage hook for entering to create
    this.owner.lookup('service:router').reopen({
      transitionTo(route, { queryParams: { initialKey } }) {
        assert.strictEqual(
          route,
          `${MOUNT_POINT}.create`,
          'Sends to create route when secret does not exists.'
        );
        assert.deepEqual(initialKey, 'beep/boop/new-secret', 'It sends full secret path.');
      },
    });
    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 13);
  });

  test('it transitions details page on enter when secret path exists', async function (assert) {
    assert.expect(1);
    // mirage hook for entering to details
    this.owner.lookup('service:router').reopen({
      transitionTo(route) {
        assert.strictEqual(
          route,
          `${MOUNT_POINT}.secret.details`,
          'Sends to details route when secret does exists.'
        );
      },
    });

    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @filterValue="beep/boop/bop"/>`,
      {
        owner: this.engine,
      }
    );
    // focus on input, fillIn and then trigger enter
    await focus('[data-test-component="kv-list-filter"]');
    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 13);
  });
});
