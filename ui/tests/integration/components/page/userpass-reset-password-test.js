/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { click, fillIn, render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const S = {
  infoBanner: '[data-test-current-user-banner]',
  save: '[data-test-reset-password-save]',
  error: '[data-test-reset-password-error]',
  input: '[data-test-textarea]',
};
module('Integration | Component | page/userpass-reset-password', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'userpass3';
    this.username = 'alice';
  });

  test('form works -- happy path', async function (assert) {
    assert.expect(4);
    this.server.post(`/auth/${this.backend}/users/${this.username}/password`, (schema, req) => {
      const body = JSON.parse(req.requestBody);
      assert.ok(true, 'correct endpoint called for update (once)');
      assert.deepEqual(body, { password: 'new' }, 'request body is correct');
      return {};
    });
    await render(hbs`<Page::UserpassResetPassword @backend={{this.backend}} @username={{this.username}} />`);

    assert
      .dom(S.infoBanner)
      .hasText(
        `You are updating the password for ${this.username} on the ${this.backend} auth mount.`,
        'info text correct'
      );

    await fillIn(S.input, 'new');
    await click(S.save);
    // eslint-disable-next-line ember/no-settled-after-test-helper
    await settled();

    assert.dom(S.input).hasValue('', 'After successful save shows input again with empty value');
  });

  test('form works -- handles error', async function (assert) {
    this.server.post(`/auth/${this.backend}/users/${this.username}/password`, () => {
      return new Response(403, {}, { errors: ['some error occurred'] });
    });
    await render(hbs`<Page::UserpassResetPassword @backend={{this.backend}} @username={{this.username}} />`);

    assert
      .dom(S.infoBanner)
      .hasText(`You are updating the password for ${this.username} on the ${this.backend} auth mount.`);

    await click(S.save);
    assert.dom(S.error).hasText('Error Please provide a new password.');

    await fillIn(S.input, 'invalid-pw');
    await click(S.save);

    // eslint-disable-next-line ember/no-settled-after-test-helper
    await settled();
    assert.dom(S.error).hasText('Error some error occurred', 'Shows error from API');
  });
});
