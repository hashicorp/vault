/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Service from '@ember/service';
import { run } from '@ember/runloop';
import { reject, resolve } from 'rsvp';
import { SHAMIR_FORM } from 'vault/tests/helpers/components/shamir-selectors';

const licenseError = { httpStatus: 500, errors: ['failed because licensing is in an invalid state'] };
const response = {
  progress: 1,
  required: 3,
  complete: false,
};

const adapter = {
  foo() {
    return resolve(response);
  },
  responseWithErrors() {
    return reject({ httpStatus: 400, errors: ['something is wrong', 'seriously wrong'] });
  },
  responseWithLicense() {
    return reject(licenseError);
  },
};

const storeStub = Service.extend({
  adapterFor() {
    return adapter;
  },
});

// Checks that the correct data were passed around happens in the integration test
// this one is checking that things happen at the right time
module('Integration | Component | shamir/flow', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.keyPart = 'some-key-partition';
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeStub);
      this.storeService = this.owner.lookup('service:store');
    });
  });

  test('it sends data to the passed action and calls updateProgress', async function (assert) {
    const updateSpy = sinon.spy();
    const completeSpy = sinon.spy();
    this.set('updateProgress', updateSpy);
    this.set('checkComplete', () => false);
    this.set('onSuccess', completeSpy);
    this.set('progress', 0);

    await render(hbs`
      <Shamir::Flow
        @action="foo"
        @threshold={{3}}
        @progress={{this.progress}}
        @updateProgress={{this.updateProgress}}
        @checkComplete={{this.checkComplete}}
        @onShamirSuccess={{this.onSuccess}}
      />`);

    await fillIn(SHAMIR_FORM.input, this.keyPart);
    await click(SHAMIR_FORM.submitButton);

    assert.ok(completeSpy.notCalled, 'onShamirSuccess was not called');
    assert.ok(updateSpy.calledOnce, 'updateProgress was called');
    // Default shamir flow expects the updated values to be passed
    // in from parent model, so this approximates the update happening
    // from a side effect of the updateProgress call
    this.set('progress', 2);
    // Pretend the next call will mean completion
    this.set('checkComplete', () => true);
    await settled();

    await fillIn(SHAMIR_FORM.input, this.keyPart);
    await click(SHAMIR_FORM.submitButton);

    assert.ok(completeSpy.calledOnce, 'onShamirSuccess was called');
    assert.ok(updateSpy.calledTwice, 'updateProgress was called again');
  });

  test('it shows the error when adapter fails with 400 httpStatus', async function (assert) {
    assert.expect(3);
    const updateSpy = sinon.spy();
    const completeSpy = sinon.spy();
    this.set('updateProgress', updateSpy);
    this.set('checkComplete', completeSpy);
    await render(hbs`
      <Shamir::Flow
        @action="response-with-errors"
        @threshold={{3}}
        @progress={{2}}
        @updateProgress={{this.updateProgress}}
        @checkComplete={{this.checkComplete}}
      />`);

    await fillIn(SHAMIR_FORM.input, this.keyPart);
    await click(SHAMIR_FORM.submitButton);
    assert.dom(SHAMIR_FORM.error).exists({ count: 2 }, 'renders errors');
    assert.ok(completeSpy.notCalled, 'checkComplete was not called');
    assert.ok(updateSpy.notCalled, 'updateProgress was not called');
  });

  test.skip('it throws the error when adapter fails with license error', async function (assert) {
    assert.expect(2);
    try {
      const licenseSpy = sinon.spy();
      this.set('onLicenseError', licenseSpy);
      await render(hbs`
        <Shamir::Flow
          @action="response-with-license"
          @threshold={{3}}
          @progress={{2}}
          @onLicenseError={{this.onLicenseError}}
        />`);
      await fillIn(SHAMIR_FORM.input, this.keyPart);
      await click(SHAMIR_FORM.submitButton);
      assert.ok(licenseSpy.calledOnce, 'license error triggered');
    } catch (e) {
      assert.deepEqual(e, licenseError, 'throws the error');
    }
  });
});
