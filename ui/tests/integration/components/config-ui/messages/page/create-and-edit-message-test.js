/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { datetimeLocalStringFormat } from 'core/utils/date-formatters';
import { format, addDays, startOfDay } from 'date-fns';

const PAGE = {
  radio: (radioName) => `[data-test-radio="${radioName}"]`,
  field: (fieldName) => `[data-test-field="${fieldName}"]`,
  input: (input) => `[data-test-input="${input}"]`,
  button: (buttonName) => `[data-test-button="${buttonName}"]`,
  inlineErrorMessage: `[data-test-inline-error-message]`,
  fieldVaildation: (fieldName) => `[data-test-field-validation="${fieldName}"]`,
  modal: (name) => `[data-test-modal="${name}"]`,
  modalTitle: (title) => `[data-test-modal-title="${title}"]`,
  modalBody: '[data-test-modal-body]',
  modalButton: (name) => `[data-test-modal-button="${name}"]`,
  alertTitle: (name) => `[data-test-alert-title="${name}"]`,
  alertDescription: (name) => `[data-test-alert-description="${name}"]`,
};

module('Integration | Component | messages/page/create-and-edit-message', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');
    this.message = this.store.createRecord('config-ui/message');
  });

  test('it should display all the create form fields and default radio button values', async function (assert) {
    await render(hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} />`, {
      owner: this.engine,
    });

    assert.dom('[data-test-page-title]').hasText('Create message');
    assert
      .dom('[data-test-form-subtext]')
      .hasText('Create a custom message for all users when they access a Vault system via the UI.');
    assert.dom(PAGE.radio('authenticated')).exists();
    assert.dom(PAGE.radio('unauthenticated')).exists();
    assert.dom(PAGE.radio('authenticated')).isChecked();
    assert.dom(PAGE.radio('unauthenticated')).isNotChecked();
    assert.dom(PAGE.radio('banner')).exists();
    assert.dom(PAGE.radio('modal')).exists();
    assert.dom(PAGE.radio('banner')).isChecked();
    assert.dom(PAGE.radio('modal')).isNotChecked();
    assert.dom(PAGE.field('title')).exists();
    assert.dom(PAGE.field('message')).exists();
    assert.dom(PAGE.input('linkTitle')).exists();
    assert.dom(PAGE.input('linkHref')).exists();
    assert.dom(PAGE.input('startTime')).exists();
    assert
      .dom(PAGE.input('startTime'))
      .hasValue(format(addDays(startOfDay(new Date()), 1), datetimeLocalStringFormat));
    assert.dom(PAGE.input('endTime')).exists();
    assert.dom(PAGE.input('endTime')).hasValue('');
  });

  test('it should create new message', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/config/ui/custom-messages', () => {
      assert.ok(true, 'POST request made to create message');
    });

    await render(hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} />`, {
      owner: this.engine,
    });
    await fillIn(PAGE.input('title'), 'Awesome custom message title');
    await fillIn(
      PAGE.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
    );
    await fillIn(
      PAGE.input('startTime'),
      format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
    );
    await click('#specificDate');
    await fillIn(
      PAGE.input('endTime'),
      format(addDays(startOfDay(new Date('2023-12-12')), 10), datetimeLocalStringFormat)
    );

    await click(PAGE.button('create-message'));
  });

  test('it should have form vaildations', async function (assert) {
    await render(hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} />`, {
      owner: this.engine,
    });
    await click(PAGE.button('create-message'));
    assert.dom(PAGE.input('title')).hasClass('has-error-border', 'show error border for title field');
    assert.dom(`${PAGE.fieldVaildation('title')} ${PAGE.inlineErrorMessage}`).hasText('Title is required.');
    assert.dom(PAGE.input('message')).hasClass('has-error-border', 'show error border for message field');
    assert
      .dom(`${PAGE.fieldVaildation('message')} ${PAGE.inlineErrorMessage}`)
      .hasText('Message is required.');
  });

  test('it should prepopulate form if form is in edit mode', async function (assert) {
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: 'hhhhh-iiii-lllll-dddd',
      type: 'modal',
      authenticated: false,
      title: 'Hello world',
      message: 'Blah blah blah. Some super long message.',
      start_time: '2023-12-12T08:00:00.000Z',
      end_time: '2023-12-21T08:00:00.000Z',
      link_title: 'Learn more',
      link_href: 'www.learnmore.com',
    });
    this.message = this.store.peekRecord('config-ui/message', 'hhhhh-iiii-lllll-dddd');
    await render(hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} />`, {
      owner: this.engine,
    });

    assert.dom('[data-test-page-title]').hasText('Edit message');
    assert
      .dom('[data-test-form-subtext]')
      .hasText('Edit a custom message for all users when they access a Vault system via the UI.');
    assert.dom(PAGE.radio('authenticated')).exists();
    assert.dom(PAGE.radio('unauthenticated')).isChecked();
    assert.dom(PAGE.radio('modal')).exists();
    assert.dom(PAGE.radio('modal')).isChecked();
    assert.dom(PAGE.input('title')).hasValue('Hello world');
    assert.dom(PAGE.input('message')).hasValue('Blah blah blah. Some super long message.');
    assert.dom(PAGE.input('linkTitle')).exists();
    assert.dom(PAGE.input('linkTitle')).hasValue('Learn more');
    assert.dom(PAGE.input('linkHref')).exists();
    assert.dom(PAGE.input('linkHref')).hasValue('www.learnmore.com');
    await click('#specificDate');
    assert
      .dom(PAGE.input('startTime'))
      .hasValue(format(new Date(this.message.startTime), datetimeLocalStringFormat));
    assert
      .dom(PAGE.input('endTime'))
      .hasValue(format(new Date(this.message.endTime), datetimeLocalStringFormat));
  });

  test('it should show a preview image modal when preview is clicked', async function (assert) {
    await render(hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} />`, {
      owner: this.engine,
    });
    await fillIn(PAGE.input('title'), 'Awesome custom message title');
    await fillIn(
      PAGE.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
    );
    await click(PAGE.button('preview'));
    assert.dom(PAGE.modal('preview modal')).doesNotExist();
    assert.dom(PAGE.modal('preview image')).exists();
    assert.dom(PAGE.alertTitle('Awesome custom message title')).hasText('Awesome custom message title');
    assert
      .dom(PAGE.alertDescription('Awesome custom message title'))
      .hasText(
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
    assert.dom('img').hasAttribute('src', '/ui/images/custom-messages-dashboard.png');
    await click(PAGE.modalButton('Close'));
    await click('#unauthenticated');
    await click(PAGE.button('preview'));
    assert.dom('img').hasAttribute('src', '/ui/images/custom-messages-login.png');
  });

  test('it should show a preview modal when preview is clicked', async function (assert) {
    await render(hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} />`, {
      owner: this.engine,
    });
    await click(PAGE.radio('modal'));
    await fillIn(PAGE.input('title'), 'Preview modal title');
    await fillIn(PAGE.input('message'), 'Some preview modal message thats super long.');
    await click(PAGE.button('preview'));
    assert.dom(PAGE.modal('preview modal')).exists();
    assert.dom(PAGE.modal('preview image')).doesNotExist();
    assert.dom(PAGE.modalTitle('Preview modal title')).hasText('Preview modal title');
    assert.dom(PAGE.modalBody).hasText('Some preview modal message thats super long.');
  });
});
