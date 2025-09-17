/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | message-error', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.errorMessage = undefined;
    this.errors = undefined;
    this.model = undefined;
    this.onDismiss = sinon.spy();

    this.renderComponent = () => {
      return render(hbs`
        <MessageError 
        @errorMessage={{this.errorMessage}} 
        @errors={{this.errors}} 
        @model={{this.model}} 
        @onDismiss={{this.onDismiss}} 
        />`);
    };
  });

  test('it does not render if args have no value', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.messageError).doesNotExist();
  });

  test('it renders a single error message', async function (assert) {
    this.errorMessage = 'Something went wrong';
    await this.renderComponent();

    assert.dom(GENERAL.messageError).exists({ count: 1 });
    assert.dom(GENERAL.messageDescription).hasText('Something went wrong');
  });

  test('it renders an array of error objects', async function (assert) {
    this.errors = [new Error('oh dear'), new Error('uh oh')];
    await render(hbs`<MessageError @errors={{this.errors}} />`);

    assert.dom(GENERAL.messageError).exists({ count: 2 });
    findAll(GENERAL.messageDescription).forEach((e, idx) => {
      assert.dom(e).hasText(`Error: ${this.errors[idx].message}`);
    });
  });

  test('it renders an array of error strings', async function (assert) {
    this.errors = ['problem 1', 'problem 2'];
    await render(hbs`<MessageError @errors={{this.errors}} />`);

    assert.dom(GENERAL.messageError).exists({ count: 2 });
    findAll(GENERAL.messageDescription).forEach((e, idx) => {
      assert.dom(e).hasText(this.errors[idx]);
    });
  });

  test('it handles model errors', async function (assert) {
    this.model = {
      isError: true,
      adapterError: {
        errors: [new Error('there was a problem')],
      },
    };
    await this.renderComponent();

    assert.dom(GENERAL.messageError).exists({ count: 1 });
    assert.dom(GENERAL.messageDescription).hasText('there was a problem');
  });

  test('it formats CLI-style errors', async function (assert) {
    const apiError =
      'failed to satisfy enforcement just-duo. error: 2 errors occurred:\n\t* duo authentication failed: "Login request denied."\n\t* login MFA validation failed';
    const expectedDescriptions = [
      'duo authentication failed: "Login request denied."',
      'login MFA validation failed',
    ];
    this.errorMessage = apiError;
    await this.renderComponent();

    assert.dom(GENERAL.messageError).exists({ count: 1 });
    assert.dom(GENERAL.messageDescription).hasTextContaining('failed to satisfy enforcement just-duo');
    assert.dom(`${GENERAL.messageDescription} li`).exists({ count: 2 });
    findAll(`${GENERAL.messageDescription} li`).forEach((e, idx) => {
      assert.dom(e).hasText(expectedDescriptions[idx]);
    });
  });

  test('it returns null for invalid formatted errors', async function (assert) {
    const invalidFormat = '* some error without proper format';
    this.errorMessage = invalidFormat;
    await this.renderComponent();

    assert.dom(GENERAL.messageError).exists({ count: 1 });
    assert.dom(GENERAL.messageDescription).hasText(invalidFormat);
  });

  test('it fires on dismiss callback', async function (assert) {
    this.errorMessage = 'Some issue occurred';
    await this.renderComponent();
    await click(GENERAL.icon('x'));
    assert.true(this.onDismiss.calledOnce, '@onDismiss is called once');
  });
});
