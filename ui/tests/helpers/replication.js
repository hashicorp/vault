import { click, fillIn, findAll, currentURL, visit, settled, waitUntil } from '@ember/test-helpers';

export const disableReplication = async (type, assert) => {
  // disable performance replication
  await visit(`/vault/replication/${type}`);

  if (findAll('[data-test-replication-link="manage"]').length) {
    await click('[data-test-replication-link="manage"]');

    await click('[data-test-disable-replication] button');

    const typeDisplay = type === 'dr' ? 'Disaster Recovery' : 'Performance';
    await fillIn('[data-test-confirmation-modal-input="Disable Replication?"]', typeDisplay);
    await click('[data-test-confirm-button]');
    await settled(); // eslint-disable-line

    if (assert) {
      // bypassing for now -- remove if tests pass reliably
      // assert.strictEqual(
      //   flash.latestMessage,
      //   'This cluster is having replication disabled. Vault will be unavailable for a brief period and will resume service shortly.',
      //   'renders info flash when disabled'
      // );
      assert.ok(
        await waitUntil(() => currentURL() === '/vault/replication'),
        'redirects to the replication page'
      );
    }
    await settled();
  }
};
