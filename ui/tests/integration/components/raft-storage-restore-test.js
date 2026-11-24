/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, triggerEvent, click, waitFor, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { Response } from 'miragejs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | raft-storage-restore', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it should restore snapshot', async function (assert) {
    assert.expect(2);

    this.server.post('/sys/storage/raft/snapshot', () => {
      assert.ok(true, 'Request made to restore snapshot');
      return;
    });
    this.server.post('/sys/storage/raft/snapshot-force', () => {
      assert.ok(true, 'Request made to force restore snapshot');
      return;
    });

    await render(hbs`<RaftStorageRestore />`);
    await triggerEvent('[data-test-file-input]', 'change', {
      files: [new Blob(['Raft Snapshot'])],
    });
    await click(GENERAL.submitButton);
    await waitFor('#force-restore');
    await click('#force-restore');
    await click(GENERAL.submitButton);
  });

  test('it should show errors returned when restoring a snapshot', async function (assert) {
    this.server.post('/sys/storage/raft/snapshot', () => new Response(400, {}, { errors: ['bad snapshot'] }));

    await render(hbs`<RaftStorageRestore />`);
    await triggerEvent('[data-test-file-input]', 'change', {
      files: [new Blob(['Raft Snapshot'])],
    });
    await click(GENERAL.submitButton);
    await waitFor('[data-test-message-error]');
    assert.dom('[data-test-message-error]').includesText('bad snapshot', 'Error from the API is displayed');
    assert.notOk(find('.flash-message'), 'Success is not flashed');
  });
});
