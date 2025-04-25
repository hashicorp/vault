/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn } from '@ember/test-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { AUTH_METHOD_MAP } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

/*
NOTE: In the app these components are actually rendered dynamically by Auth::FormTemplate
and so the components rendered in these tests does not represent "real world" situations.
This is intentional to test component logic specific to auth/form/base or auth/form/<type> 
separately from auth/form-template.
*/

export default (test, { standardSubmit = true } = {}) => {
  test('it renders fields', async function (assert) {
    await this.renderComponent();
    assert.dom(AUTH_FORM.authForm(this.authType)).exists(`${this.authType}: it renders form component`);
    this.expectedFields.forEach((field) => {
      assert.dom(GENERAL.inputByAttr(field)).exists(`${this.authType}: it renders ${field}`);
    });
  });

  test('it fires onError callback', async function (assert) {
    this.authenticateStub.throws('permission denied');
    await this.renderComponent();
    await click(AUTH_FORM.login);

    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(
      actual,
      'Authentication failed: permission denied: Sinon-provided permission denied',
      'it calls onError'
    );
  });

  test('it fires onSuccess callback', async function (assert) {
    this.authenticateStub.returns('success!');
    await this.renderComponent();
    await click(AUTH_FORM.login);

    const [actual] = this.onSuccess.lastCall.args;
    assert.strictEqual(actual, 'success!', 'it calls onSuccess');
  });

  // some methods are tested separately because they have more complex submit logic
  if (standardSubmit) {
    test('it submits form data with defaults', async function (assert) {
      await this.renderComponent();
      const { options } = AUTH_METHOD_MAP.find((m) => m.authType === this.authType);
      const { loginData } = options;

      for (const [field, value] of Object.entries(loginData)) {
        await fillIn(GENERAL.inputByAttr(field), value);
      }
      await click(AUTH_FORM.login);
      const [actual] = this.authenticateStub.lastCall.args;
      assert.propEqual(
        actual.data,
        this.expectedSubmit.default,
        'auth service "authenticate" method is called with form data'
      );
    });

    // not for testing real-world submit, that happens in acceptance tests.
    // component here just yields <:advancedSettings> to test form submits data from yielded inputs
    test('it submits form data from yielded inputs', async function (assert) {
      await this.renderComponent({ yieldBlock: true });
      const { options } = AUTH_METHOD_MAP.find((m) => m.authType === this.authType);
      const { loginData } = options;

      for (const [field, value] of Object.entries(loginData)) {
        await fillIn(GENERAL.inputByAttr(field), value);
      }

      if (this.authType === 'token') {
        // token doesn't support custom paths, so just test yielding functionality
        await fillIn(GENERAL.inputByAttr('yield'), `yield-${this.authType}`);
      } else {
        await fillIn(GENERAL.inputByAttr('path'), `custom-${this.authType}`);
      }

      await click(AUTH_FORM.login);
      const [actual] = this.authenticateStub.lastCall.args;
      assert.propEqual(
        actual.data,
        this.expectedSubmit.custom,
        'auth service "authenticate" method is called with yielded form data'
      );
    });
  }
};
