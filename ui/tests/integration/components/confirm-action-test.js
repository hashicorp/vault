/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  title: '[data-test-confirm-action-title]',
  message: '[data-test-confirm-action-message]',
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

    assert.dom(GENERAL.confirmTrigger).hasText('DELETE', 'renders button text');
    await click(GENERAL.confirmTrigger);
    assert
      .dom('#confirm-action-modal')
      .hasClass('hds-modal--color-critical', 'renders critical modal color by default');
    assert
      .dom(GENERAL.confirmButton)
      .hasClass('hds-button--color-critical', 'renders critical confirm button');
    assert.dom(SELECTORS.title).hasText('Are you sure?', 'renders default title');
    assert
      .dom(SELECTORS.message)
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click(GENERAL.cancelButton);
    assert.false(this.onConfirm.called, 'does not call the action when Cancel is clicked');
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
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

    assert.dom(`li ${GENERAL.confirmTrigger}`).exists('element renders inside <li>');
    assert.dom(GENERAL.confirmTrigger).hasClass('hds-confirm-action-critical', 'button has dropdown styling');
    await click(GENERAL.confirmTrigger);
    assert.dom(SELECTORS.title).hasText('Are you sure?', 'renders default title');
    assert
      .dom(SELECTORS.message)
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click(GENERAL.cancelButton);
    assert.false(this.onConfirm.called, 'does not call the action when Cancel is clicked');
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
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

    await click(GENERAL.confirmTrigger);

    assert.dom(GENERAL.confirmButton).isDisabled('disables confirm button when loading');
    assert.dom('[data-test-confirm-button] [data-test-icon="loading"]').exists('it renders loading icon');
  });

  test('it renders disabledMessage modal', async function (assert) {
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @confirmTitle="Do this?"
        @confirmMessage="Are you really, really sure?"
        @disabledMessage="This is the reason you cannot do the thing"
      />
      `);

    await click(GENERAL.confirmTrigger);
    assert
      .dom('#confirm-action-modal')
      .hasClass('hds-modal--color-neutral', 'renders neutral modal because disabledMessage is present');
    assert.dom(SELECTORS.title).hasText('Not allowed', 'renders disabled title');
    assert
      .dom(SELECTORS.message)
      .hasText('This is the reason you cannot do the thing', 'renders disabled message as body text');
    assert.dom(GENERAL.confirmButton).doesNotExist('does not render confirm action button');
    assert.dom(GENERAL.cancelButton).hasText('Close');
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

    assert
      .dom(GENERAL.confirmTrigger)
      .hasClass('hds-button--color-secondary', 'renders @buttonColor classes');
    await click(GENERAL.confirmTrigger);
    assert.dom('#confirm-action-modal').hasClass('hds-modal--color-warning', 'renders warning modal');
    assert.dom(GENERAL.confirmButton).hasClass('hds-button--color-primary', 'renders primary confirm button');
    assert.dom(SELECTORS.title).hasText('Do this?', 'renders passed title');
    assert.dom(SELECTORS.message).hasText('Are you really, really sure?', 'renders passed body text');
    assert.dom(GENERAL.confirmButton).hasText('Confirm');
  });
});
