import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, triggerEvent, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

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
    await click('[data-test-edit-form-submit]');
    await click('#force-restore');
    await click('[data-test-edit-form-submit]');
  });
});
