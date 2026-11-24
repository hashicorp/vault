/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | Enterprise | /access/namespaces/create', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    await visit('/vault/access/namespaces/create');
  });

  test('it renders a permission denied error and stays on the create page', async function (assert) {
    this.server.post('/sys/namespaces/denied-namespace', () => overrideResponse(403));

    await fillIn(GENERAL.inputByAttr('path'), 'denied-namespace');
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.messageDescription).hasText('permission denied', 'server error is rendered');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.create',
      'stays on the create page so the entered path is not lost'
    );
  });

  test('it renders the message returned by the server when create fails', async function (assert) {
    this.server.post('/sys/namespaces/bad-namespace', () =>
      overrideResponse(400, { errors: ['namespace path is invalid'] })
    );

    await fillIn(GENERAL.inputByAttr('path'), 'bad-namespace');
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.messageDescription)
      .hasText('namespace path is invalid', 'the specific server message is surfaced, not a generic one');
  });

  test('it clears a previous error when the form is resubmitted', async function (assert) {
    let attempt = 0;
    this.server.post('/sys/namespaces/retry-namespace', () => {
      attempt += 1;
      return attempt === 1 ? overrideResponse(400, { errors: ['transient failure'] }) : overrideResponse(204);
    });

    await fillIn(GENERAL.inputByAttr('path'), 'retry-namespace');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageDescription).hasText('transient failure', 'first attempt surfaces the error');

    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).doesNotExist('the stale error is cleared once the retry succeeds');
  });

  test('it does not submit when the path fails client-side validation', async function (assert) {
    this.server.post('/sys/namespaces/:path', () => {
      assert.true(false, 'an invalid path must never reach the server');
      return overrideResponse(204);
    });

    await fillIn(GENERAL.inputByAttr('path'), 'invalid namespace');
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('path'))
      .hasText("Path can't contain whitespace.", 'validation error is rendered');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.create',
      'stays on the create page'
    );
  });
});
