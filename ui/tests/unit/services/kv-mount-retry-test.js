/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | kv-mount-retry', function (hooks) {
  setupTest(hooks);

  const parseErrorFor = (status, message) => async () => ({ status, message });

  test('it returns the value on first success without retrying', async function (assert) {
    const service = this.owner.lookup('service:kv-mount-retry');
    let calls = 0;

    const result = await service.request.perform(async () => {
      calls++;
      return 'ok';
    }, parseErrorFor());

    assert.strictEqual(result, 'ok', 'returns the request value');
    assert.strictEqual(calls, 1, 'calls the request once');
  });

  test('it returns the successful value on the final allowed attempt', async function (assert) {
    const service = this.owner.lookup('service:kv-mount-retry');
    let calls = 0;

    const result = await service.request.perform(
      async () => {
        calls++;
        if (calls < 3) {
          throw new Error('transient');
        }
        return 'ok';
      },
      parseErrorFor(400, 'Upgrading from non-versioned to versioned data.'),
      { maxAttempts: 3, delayMs: 0 }
    );

    assert.strictEqual(result, 'ok', 'returns after the transient error clears');
    assert.strictEqual(calls, 3, 'retries until the request succeeds');
  });

  test('it rethrows immediately for a non-matching error without retrying', async function (assert) {
    const service = this.owner.lookup('service:kv-mount-retry');
    let calls = 0;

    await assert.rejects(
      service.request.perform(
        async () => {
          calls++;
          throw new Error('permission denied');
        },
        parseErrorFor(403, 'permission denied')
      ),
      'rejects an error that is not safe to retry'
    );
    assert.strictEqual(calls, 1, 'does not retry an unrelated error');
  });

  test('it stops retrying once maxAttempts is reached', async function (assert) {
    const service = this.owner.lookup('service:kv-mount-retry');
    let calls = 0;

    await assert.rejects(
      service.request.perform(
        async () => {
          calls++;
          throw new Error('still upgrading');
        },
        parseErrorFor(400, 'Upgrading from non-versioned to versioned data.'),
        { maxAttempts: 3, delayMs: 0 }
      ),
      'rejects after the bounded retry period'
    );
    assert.strictEqual(calls, 3, 'does not exceed the configured attempt limit');
  });

  test('it stops retrying when the task is canceled', async function (assert) {
    const service = this.owner.lookup('service:kv-mount-retry');
    let calls = 0;
    const taskInstance = service.request.perform(
      async () => {
        calls++;
        throw new Error('still upgrading');
      },
      parseErrorFor(400, 'Upgrading from non-versioned to versioned data.'),
      { delayMs: 1000 }
    );

    taskInstance.cancel();
    await taskInstance.catch(() => undefined);

    assert.true(taskInstance.isCanceled, 'the retry sequence is canceled');
    assert.strictEqual(calls, 1, 'the canceled request is not retried');
  });

  test('it cancels a stale KV list when a newer list starts', async function (assert) {
    const service = this.owner.lookup('service:kv-mount-retry');
    let calls = 0;
    service.api.secrets.kvV2List = async () => {
      calls++;
      return calls === 1 ? new Promise(() => undefined) : { keys: ['new-result'] };
    };

    const staleTask = service.listKvSecrets.perform('', 'kv-one');
    const currentTask = service.listKvSecrets.perform('', 'kv-two');
    const result = await currentTask;
    await staleTask.catch(() => undefined);

    assert.deepEqual(result, { keys: ['new-result'] }, 'the latest list request completes');
    assert.true(staleTask.isCanceled, 'the stale list request is canceled');
  });
});
