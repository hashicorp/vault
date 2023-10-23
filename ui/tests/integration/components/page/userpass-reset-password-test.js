import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const S = {
  infoBanner: '[data-test-current-user-banner]',
  save: '[data-test-reset-password-save]',
  error: '[data-test-reset-password-error]',
  input: '[data-test-textarea]',
  success: '[data-test-reset-password-success]',
  reset: '[data-test-reset-reset-password]',
};
module('Integration | Component | page/userpass-reset-password', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'userpass3';
    this.username = 'alice';
  });

  test('form works -- happy path', async function (assert) {
    assert.expect(8);
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

    assert.dom(S.success).includesText('Password set successfully', 'Shows success message');
    assert.dom(S.infoBanner).doesNotExist('Info banner no longer shown', 'info banner hidden');
    assert.dom(S.input).doesNotExist('Input no longer shown', 'input hidden');

    await click(S.reset);
    assert.dom(S.success).doesNotExist('Reset removes success banner');
    assert.dom(S.input).hasValue('', 'Reset shows input again with empty value');
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

    assert.dom(S.error).hasText('Error some error occurred', 'Shows error from API');
  });
});
