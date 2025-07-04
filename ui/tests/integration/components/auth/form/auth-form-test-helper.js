/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, findAll } from '@ember/test-helpers';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

/*
NOTE: In the app these components are actually rendered dynamically by Auth::FormTemplate
and so the components rendered in these tests does not represent "real world" situations.
This is intentional to test component logic specific to auth/form/base or auth/form/<type> 
separately from auth/form-template.

See beforeEach hooks in auth/form/base-test to see setup for each method.
*/

export default (test) => {
  test('it renders fields', async function (assert) {
    const expectedFields = Object.keys(this.loginData);
    await this.renderComponent();
    assert.dom(AUTH_FORM.authForm(this.authType)).exists(`${this.authType}: it renders form component`);
    const fields = findAll('input');
    for (const field of fields) {
      assert.true(expectedFields.includes(field.name), `it renders field: ${field.name}`);
    }
  });

  test('it fires onError callback', async function (assert) {
    this.authenticateStub.rejects(getErrorResponse({ errors: ['uh oh!'] }, 500));
    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(actual, 'Authentication failed: uh oh!', 'it calls onError');
  });

  test('it fires onSuccess callback', async function (assert) {
    this.authenticateStub.resolves(this.authResponse);
    await this.renderComponent();
    await click(GENERAL.submitButton);

    // Only checking for authMethodType because this test just asserts the onSuccess callback fires.
    const [{ authMethodType }] = this.onSuccess.lastCall.args;
    assert.strictEqual(authMethodType, this.authType, 'it calls onSuccess');
  });

  test('it submits form data with defaults', async function (assert) {
    await this.renderComponent();
    await fillInLoginFields(this.loginData);
    await click(GENERAL.submitButton);

    // Since each login method accepts different args
    // the submit assertion is setup in each method's beforeEach hook
    this.assertSubmit(assert, this.authenticateStub.lastCall.args, this.loginData);
  });

  // not for testing real-world submit, that happens in acceptance tests.
  // component here just yields <:advancedSettings> to test form submits data from yielded inputs
  test('it submits form data from yielded inputs', async function (assert) {
    const customPath = `custom-${this.authType}`;
    const loginData = { ...this.loginData, path: customPath };
    await this.renderComponent({ yieldBlock: true });
    await fillInLoginFields(loginData);
    await click(GENERAL.submitButton);
    // Since each login method accepts different args
    // the submit assertion is setup in each method's beforeEach hook
    this.assertSubmit(assert, this.authenticateStub.lastCall.args, loginData);
  });
};
