/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, focus, triggerKeyEvent } from '@ember/test-helpers';
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
    {
      id: kvMetadataPath('my-engine', 'beep/boop/blah'),
      path: 'beep/boop/blah',
      fullSecretPath: 'beep/boop/blah',
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
  /* TODO:
  - type my- and it filters to my-secrets and stays focused
  -- backspace in this state to "my" and one shows and still focused
  -- tab to complete
  -- escape and it deletes and sets to list

  - filter to beep/ and it shows two?
  - tab to complete
  -backspace in directory
  - type and hit enter goes to secret

  -- type non -existing thing and it goes create


  **/

  test('it renders', async function (assert) {
    assert.expect(4);
    // mirage hooks
    this.owner.lookup('service:router').reopen({
      transitionTo(route, { queryParams: { pageFilter } }) {
        assert.strictEqual(route, `${MOUNT_POINT}.list`, 'List route sent when TAB on empty input.');
        assert.deepEqual(pageFilter, 'my-secret', 'Filters to the first secret in the list.');
      },
    });

    await render(hbs`<KvListFilter @model={{this.model}} @mountPoint={{this.mountPoint}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-component="kv-list-filter"]')
      .hasAttribute('placeholder', 'Filter secrets', 'Placeholder applied to input.');

    await focus('[data-test-component="kv-list-filter"]');

    assert.dom('[data-test-help-tab]').exists('on focus with no filterValue is shows tab help text');

    await triggerKeyEvent('[data-test-component="kv-list-filter"]', 'keydown', 9);

    // await fillIn('[data-test-component="kv-list-filter"]', 'my');
    // await this.pauseTest();
  });
});
