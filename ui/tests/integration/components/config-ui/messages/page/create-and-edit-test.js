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
import { PAGE } from 'vault/tests/helpers/config-ui/message-selectors';

module('Integration | Component | messages/page/create-and-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');
    this.message = this.store.createRecord('config-ui/message');
  });

  test('it should display all the create form fields and default radio button values', async function (assert) {
    assert.expect(17);

    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
      owner: this.engine,
    });

    assert.dom(PAGE.title).hasText('Create message');
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
    assert.dom('[data-test-kv-key="0"]').exists();
    assert.dom('[data-test-kv-value="0"]').exists();
    assert.dom(PAGE.input('startTime')).exists();
    assert
      .dom(PAGE.input('startTime'))
      .hasValue(format(addDays(startOfDay(new Date()), 1), datetimeLocalStringFormat));
    assert.dom(PAGE.input('endTime')).exists();
    assert.dom(PAGE.input('endTime')).hasValue('');
  });

  test('it should display validation errors for invalid form fields', async function (assert) {
    assert.expect(8);
    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
      owner: this.engine,
    });

    await fillIn(PAGE.input('startTime'), '2024-01-20T00:00');
    await fillIn(PAGE.input('endTime'), '2024-01-01T00:00');
    await click(PAGE.button('create-message'));
    assert.dom(PAGE.input('title')).hasClass('has-error-border');
    assert.dom(`${PAGE.fieldValidation('title')} ${PAGE.inlineErrorMessage}`).hasText('Title is required.');
    assert.dom(PAGE.input('message')).hasClass('has-error-border');
    assert
      .dom(`${PAGE.fieldValidation('message')} ${PAGE.inlineErrorMessage}`)
      .hasText('Message is required.');
    assert.dom(PAGE.input('startTime')).hasClass('has-error-border');
    assert
      .dom(`${PAGE.fieldValidation('startTime')} ${PAGE.inlineErrorMessage}`)
      .hasText('Start time is after end time.');
    assert.dom(PAGE.input('endTime')).hasClass('has-error-border');
    assert
      .dom(`${PAGE.fieldValidation('endTime')} ${PAGE.inlineErrorMessage}`)
      .hasText('End time is before start time.');
  });

  test('it should create new message', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/config/ui/custom-messages', () => {
      assert.ok(true, 'POST request made to create message');
    });

    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
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
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(PAGE.button('create-message'));
  });

  test('it should have form vaildations', async function (assert) {
    assert.expect(4);
    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
      owner: this.engine,
    });
    await click(PAGE.button('create-message'));
    assert.dom(PAGE.input('title')).hasClass('has-error-border', 'show error border for title field');
    assert.dom(`${PAGE.fieldValidation('title')} ${PAGE.inlineErrorMessage}`).hasText('Title is required.');
    assert.dom(PAGE.input('message')).hasClass('has-error-border', 'show error border for message field');
    assert
      .dom(`${PAGE.fieldValidation('message')} ${PAGE.inlineErrorMessage}`)
      .hasText('Message is required.');
  });

  test('it should prepopulate form if form is in edit mode', async function (assert) {
    assert.expect(13);
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: 'hhhhh-iiii-lllll-dddd',
      type: 'modal',
      authenticated: false,
      title: 'Hello world',
      message: 'Blah blah blah. Some super long message.',
      start_time: '2023-12-12T08:00:00.000Z',
      end_time: '2023-12-21T08:00:00.000Z',
      link: { 'Learn more': 'www.learnmore.com' },
    });
    this.message = this.store.peekRecord('config-ui/message', 'hhhhh-iiii-lllll-dddd');
    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
      owner: this.engine,
    });

    assert.dom(PAGE.title).hasText('Edit message');
    assert.dom(PAGE.radio('authenticated')).exists();
    assert.dom(PAGE.radio('unauthenticated')).isChecked();
    assert.dom(PAGE.radio('modal')).exists();
    assert.dom(PAGE.radio('modal')).isChecked();
    assert.dom(PAGE.input('title')).hasValue('Hello world');
    assert.dom(PAGE.input('message')).hasValue('Blah blah blah. Some super long message.');
    assert.dom('[data-test-kv-key="0"]').exists();
    assert.dom('[data-test-kv-key="0"]').hasValue('Learn more');
    assert.dom('[data-test-kv-value="0"]').exists();
    assert.dom('[data-test-kv-value="0"]').hasValue('www.learnmore.com');
    await click('#specificDate');
    assert
      .dom(PAGE.input('startTime'))
      .hasValue(format(new Date(this.message.startTime), datetimeLocalStringFormat));
    assert
      .dom(PAGE.input('endTime'))
      .hasValue(format(new Date(this.message.endTime), datetimeLocalStringFormat));
  });

  test('it should show a preview image modal when preview is clicked', async function (assert) {
    assert.expect(6);
    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
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
    assert.expect(4);
    await render(hbs`<Messages::Page::CreateAndEdit @message={{this.message}} />`, {
      owner: this.engine,
    });
    await click(PAGE.radio('modal'));
    await fillIn(PAGE.input('title'), 'Preview modal title');
    await fillIn(PAGE.input('message'), 'Some preview modal message thats super long.');
    await click(PAGE.button('preview'));
    assert.dom(PAGE.modal('preview modal')).exists();
    assert.dom(PAGE.modal('preview image')).doesNotExist();
    assert.dom(PAGE.modalTitle('Preview modal title')).hasText('Preview modal title');
    assert.dom(PAGE.modalBody('Preview modal title')).hasText('Some preview modal message thats super long.');
  });

  test('it should show multiple modal message', async function (assert) {
    assert.expect(2);

    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '01234567-89ab-cdef-0123-456789abcdef',
      active: true,
      type: 'modal',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { here: 'www.example.com' },
      startTime: '2021-08-01T00:00:00Z',
      endTime: '',
    });
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '01234567-89ab-vvvv-0123-456789abcdef',
      active: true,
      type: 'modal',
      authenticated: false,
      title: 'Message title 2',
      message: 'Some long long long message',
      link: { here: 'www.example.com' },
      startTime: '2021-08-01T00:00:00Z',
      endTime: '2090-08-01T00:00:00Z',
    });

    this.messages = this.store.peekAll('config-ui/message');

    await render(
      hbs`<Messages::Page::CreateAndEdit @message={{this.message}} @messages={{this.messages}} @hasSomeActiveModals={{true}} />`,
      {
        owner: this.engine,
      }
    );
    await fillIn(PAGE.input('title'), 'Awesome custom message title');
    await fillIn(
      PAGE.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
    );
    await click(PAGE.radio('modal'));
    await click(PAGE.button('create-message'));
    assert.dom(PAGE.modalTitle('Warning: more than one modal')).exists();
    assert
      .dom(PAGE.modalBody('Warning: more than one modal'))
      .hasText(
        'You have an active modal configured after the user logs in and are trying to create another one. It is recommended to avoid having more than one modal at once as it can be intrusive for users. Would you like to continue creating your message? Click “Confirm” to continue.'
      );
    await click(PAGE.modalButton('confirm'));
  });
});
