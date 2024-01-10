/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const SELECTORS = {
  modalToggle: '[data-test-confirm-action-trigger]',
  title: '[data-test-confirm-action-title]',
  message: '[data-test-confirm-action-message]',
  confirm: '[data-test-confirm-button]',
  cancel: '[data-test-confirm-cancel-button]',
};
module('Integration | Component | confirm-action', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onConfirm = sinon.spy();
  });

  test('it renders defaults and calls onConfirmAction', async function (assert) {
    await render(hbs`
      <ConfirmAction
        @buttonText="DELETE"
        @onConfirmAction={{this.onConfirm}}
      />
      `);

    assert.dom(SELECTORS.modalToggle).hasText('DELETE', 'renders button text');
    await click(SELECTORS.modalToggle);
    // hasClass assertion wasn't working so this is the workaround
    assert.strictEqual(
      find('#confirm-action-modal').className,
      'hds-modal hds-modal--size-small hds-modal--color-critical has-text-left',
      'renders critical modal color by default'
    );
    assert.strictEqual(
      find(SELECTORS.confirm).className,
      'hds-button hds-button--size-medium hds-button--color-critical',
      'renders critical confirm button'
    );
    assert.dom(SELECTORS.title).hasText('Are you sure?', 'renders default title');
    assert
      .dom(SELECTORS.message)
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click(SELECTORS.cancel);
    assert.false(this.onConfirm.called, 'does not call the action when Cancel is clicked');
    await click(SELECTORS.modalToggle);
    await click(SELECTORS.confirm);
    assert.true(this.onConfirm.called, 'calls the action when Confirm is clicked');
    assert.dom(SELECTORS.title).doesNotExist('modal closes after confirm is clicked');
  });

  test('it renders isInDropdown defaults and calls onConfirmAction', async function (assert) {
    setRunOptions({
      rules: {
        // this component breaks this rule because it expects to be rendered within <ul>
        listitem: { enabled: false },
      },
    });
    await render(hbs`
      <ConfirmAction
        @buttonText="DELETE"
        @onConfirmAction={{this.onConfirm}}
        @isInDropdown={{true}}
      />
      `);

    assert.dom(`li ${SELECTORS.modalToggle}`).exists('element renders inside <li>');
    assert.dom(SELECTORS.modalToggle).hasClass('hds-confirm-action-critical', 'button has dropdown styling');
    await click(SELECTORS.modalToggle);
    assert.dom(SELECTORS.title).hasText('Are you sure?', 'renders default title');
    assert
      .dom(SELECTORS.message)
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click('[data-test-confirm-cancel-button]');
    assert.false(this.onConfirm.called, 'does not call the action when Cancel is clicked');
    await click(SELECTORS.modalToggle);
    await click(SELECTORS.confirm);
    assert.true(this.onConfirm.called, 'calls the action when Confirm is clicked');
    assert.dom(SELECTORS.title).doesNotExist('modal closes after confirm is clicked');
  });

  test('it renders loading state', async function (assert) {
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @isRunning={{true}}
      />
      `);

    await click(SELECTORS.modalToggle);

    assert.dom(SELECTORS.confirm).isDisabled('disables confirm button when loading');
    assert.dom('[data-test-confirm-button] [data-test-icon="loading"]').exists('it renders loading icon');
  });

  test('it renders disabledMessage modal', async function (assert) {
    this.condition = true;
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @confirmTitle="Do this?"
        @confirmMessage="Are you really, really sure?"
        @disabledMessage={{if this.condition "This is the reason you cannot do the thing"}}
      />
      `);

    await click(SELECTORS.modalToggle);
    assert.strictEqual(
      find('#confirm-action-modal').className,
      'hds-modal hds-modal--size-small hds-modal--color-neutral has-text-left',
      'renders critical modal color by default'
    );
    assert.dom(SELECTORS.title).hasText('Not allowed', 'renders disabled title');
    assert
      .dom(SELECTORS.message)
      .hasText('This is the reason you cannot do the thing', 'renders disabled message as body text');
    assert.dom(SELECTORS.confirm).doesNotExist('does not render confirm action button');
    assert.dom(SELECTORS.cancel).hasText('Close');
  });

  test('it renders passed args', async function (assert) {
    this.condition = false;
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @modalColor="warning"
        @buttonColor="secondary"
        @confirmTitle="Do this?"
        @confirmMessage="Are you really, really sure?"
        @disabledMessage={{if this.condition "This is the reason you cannot do the thing"}}
      />
      `);

    // hasClass assertion wasn't working so this is the workaround
    assert.strictEqual(
      find(SELECTORS.modalToggle).className,
      'hds-button hds-button--size-medium hds-button--color-secondary',
      'renders @buttonColor classes'
    );
    await click(SELECTORS.modalToggle);
    assert.strictEqual(
      find('#confirm-action-modal').className,
      'hds-modal hds-modal--size-small hds-modal--color-warning has-text-left',
      'renders warning modal'
    );
    assert.strictEqual(
      find(SELECTORS.confirm).className,
      'hds-button hds-button--size-medium hds-button--color-primary',
      'renders primary confirm button'
    );
    assert.dom(SELECTORS.title).hasText('Do this?', 'renders passed title');
    assert.dom(SELECTORS.message).hasText('Are you really, really sure?', 'renders passed body text');
    assert.dom(SELECTORS.confirm).hasText('Confirm');
  });
});
