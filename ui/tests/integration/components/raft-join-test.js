/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | raft-join', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    sinon.stub(this.router, 'transitionTo');
  });

  test('it renders', async function (assert) {
    await render(hbs`<RaftJoin />`);
    assert.dom('[data-test-join-choice]').exists();
  });

  test('it shows the join form when clicking next', async function (assert) {
    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    assert.dom('[data-test-join-header]').exists();
  });
  test('it returns to the first screen when clicking back', async function (assert) {
    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    assert.dom('[data-test-join-header]').exists();
    await click(GENERAL.cancelButton);
    assert.dom('[data-test-join-choice]').exists();
  });

  test('it calls onDismiss when a user chooses to init', async function (assert) {
    const spy = sinon.spy();
    this.set('onDismiss', spy);
    await render(hbs`<RaftJoin @onDismiss={{this.onDismiss}} />`);

    await click('[data-test-join-init]');
    await click('[data-test-next]');
    assert.ok(spy.calledOnce, 'it calls the passed onDismiss');
  });

  test('it shows a validation error when leader api address is blank', async function (assert) {
    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('leader_api_addr'))
      .hasText("Leader API Address can't be blank.");
  });

  test('it clears the validation error once the field is corrected and resubmitted', async function (assert) {
    this.server.post('/sys/storage/raft/join', () => ({}));

    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('leader_api_addr'))
      .exists('validation error shows on blank submit');

    await fillIn('[data-test-input="leader_api_addr"]', 'https://127.0.0.1:8200');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('leader_api_addr'))
      .doesNotExist('validation error clears once the field is valid and resubmitted');
  });

  test('it joins an existing raft cluster', async function (assert) {
    assert.expect(2);

    this.server.post('/sys/storage/raft/join', (schema, req) => {
      const body = JSON.parse(req.requestBody);
      assert.strictEqual(body.leader_api_addr, 'https://127.0.0.1:8200', 'join request sent with form data');
      return {};
    });

    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    await fillIn('[data-test-input="leader_api_addr"]', 'https://127.0.0.1:8200');
    await click(GENERAL.submitButton);
    assert.ok(
      this.router.transitionTo.calledWith('vault.cluster.unseal'),
      'transitions to unseal on success'
    );
  });

  test('it shows the API error via MessageError at the top of the form when the join request fails', async function (assert) {
    this.server.post('/sys/storage/raft/join', () =>
      overrideResponse(400, { errors: ['connection refused'] })
    );

    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    await fillIn('[data-test-input="leader_api_addr"]', 'https://127.0.0.1:8200');
    await click(GENERAL.submitButton);

    assert.dom('[data-test-message-error]').exists('MessageError renders the API failure');
    assert.dom('[data-test-message-error-description]').hasTextContaining('connection refused');
    assert.notOk(this.router.transitionTo.called, 'does not transition to unseal on failure');
  });
});
