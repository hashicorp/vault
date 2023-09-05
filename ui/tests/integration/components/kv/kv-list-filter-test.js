/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, triggerKeyEvent, typeIn, click } from '@ember/test-helpers';
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

  test('it transitions correctly for query without slash', async function (assert) {
    assert.expect(3);
    const routerSvc = this.owner.lookup('service:router');
    sinon.stub(routerSvc, 'transitionTo').callsFake((route, params) => {
      assert.strictEqual(route, `${MOUNT_POINT}.list`, 'List route sent when no pathToSecret.');
      assert.deepEqual(
        params,
        { queryParams: { pageFilter: 'my-secret' } },
        'Sends correct transition params.'
      );
    });

    await render(hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-kv-list-filter]')
      .hasAttribute('placeholder', 'Search secret path', 'Placeholder applied to input.');

    await typeIn('[data-test-kv-list-filter]', 'my-secret');
    await click('[data-test-kv-list-filter-submit]');
  });

  test('it transitions correctly for query ending in /', async function (assert) {
    assert.expect(3);
    const routerSvc = this.owner.lookup('service:router');
    sinon.stub(routerSvc, 'transitionTo').callsFake((route, params) => {
      assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'List route sent when params');
      assert.deepEqual(params, 'beep/', 'Sends directory as param');
    });

    await render(hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-kv-list-filter]')
      .hasAttribute('placeholder', 'Search secret path', 'Placeholder applied to input.');

    await typeIn('[data-test-kv-list-filter]', 'beep/');
    await click('[data-test-kv-list-filter-submit]');
  });

  test('it transitions correctly for nested query', async function (assert) {
    assert.expect(4);
    const routerSvc = this.owner.lookup('service:router');
    sinon.stub(routerSvc, 'transitionTo').callsFake((route, params, { queryParams }) => {
      assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'List route sent when params');
      assert.deepEqual(params, 'beep/', 'Sends directory as url param');
      assert.deepEqual(queryParams, { pageFilter: 'boo' }, 'Sends directory as query param');
    });

    await render(hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-kv-list-filter]')
      .hasAttribute('placeholder', 'Search secret path', 'Placeholder applied to input.');

    await typeIn('[data-test-kv-list-filter]', 'beep/boo');
    await click('[data-test-kv-list-filter-submit]');
  });

  test('it prefills filterbar from pageFilter', async function (assert) {
    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @filterValue="beep/boop/bop" />`,
      {
        owner: this.engine,
      }
    );
    assert.dom('[data-test-kv-list-filter]').hasValue('beep/boop/bop');
  });

  test('it clears to directory on esc', async function (assert) {
    assert.expect(3);
    const routerSvc = this.owner.lookup('service:router');
    sinon.stub(routerSvc, 'transitionTo').callsFake((route, params, { queryParams }) => {
      assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'List route sent when params');
      assert.deepEqual(params, 'beep/boop/', 'Sends base directory as url param');
      assert.deepEqual(queryParams, { pageFilter: null }, 'clears pageFilter param');
    });

    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @filterValue="beep/boop/bop" />`,
      {
        owner: this.engine,
      }
    );
    // trigger esc
    await triggerKeyEvent('[data-test-kv-list-filter]', 'keydown', 27);
  });

  test('it clears to previous directory on esc', async function (assert) {
    assert.expect(3);
    const routerSvc = this.owner.lookup('service:router');
    sinon.stub(routerSvc, 'transitionTo').callsFake((route, params, { queryParams }) => {
      assert.strictEqual(route, `${MOUNT_POINT}.list-directory`, 'List route sent when params');
      assert.deepEqual(params, 'beep/', 'Sends base directory as url param');
      assert.deepEqual(queryParams, { pageFilter: null }, 'clears pageFilter param');
    });

    await render(
      hbs`<KvListFilter @secrets={{this.model.secrets}} @mountPoint={{this.mountPoint}} @filterValue="beep/boop/" />`,
      {
        owner: this.engine,
      }
    );
    // trigger esc
    await triggerKeyEvent('[data-test-kv-list-filter]', 'keydown', 27);
  });
});
