/**
 * Copyright IBM Corp. 2016, 2025
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
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import timestamp from 'core/utils/timestamp';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import CustomMessage from 'vault/forms/custom-message';

module('Integration | Component | messages/page/create-and-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const now = new Date('2023-07-02T00:00:00Z'); // stub "now" for testing
    sinon.replace(timestamp, 'now', sinon.fake.returns(now));

    this.message = new CustomMessage(
      {
        authenticated: true,
        type: 'banner',
        start_time: addDays(startOfDay(timestamp.now()), 1).toISOString(),
      },
      { isNew: true }
    );

    this.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: true } },
      { label: 'Create Message' },
    ];

    this.renderComponent = () =>
      render(
        hbs`<Messages::Page::CreateAndEdit @message={{this.message}} @messages={{this.messages}} @breadcrumbs={{this.breadcrumbs}} />`,
        {
          owner: this.engine,
        }
      );
  });

  test('it should display all the create form fields and default radio button values', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create message');
    assert.dom(GENERAL.fieldLabel('authenticated')).hasText('Where should we display this message?');
    assert.dom(CUSTOM_MESSAGES.radio('authenticated')).exists();
    assert.dom(CUSTOM_MESSAGES.radio('unauthenticated')).exists();
    assert.dom(CUSTOM_MESSAGES.radio('authenticated')).isChecked();
    assert.dom(CUSTOM_MESSAGES.radio('unauthenticated')).isNotChecked();
    assert.dom(GENERAL.fieldLabel('type')).hasText('Type');
    assert.dom(CUSTOM_MESSAGES.radio('banner')).exists();
    assert.dom(CUSTOM_MESSAGES.radio('modal')).exists();
    assert.dom(CUSTOM_MESSAGES.radio('banner')).isChecked();
    assert.dom(CUSTOM_MESSAGES.radio('modal')).isNotChecked();
    assert.dom(CUSTOM_MESSAGES.field('title')).exists();
    assert.dom(CUSTOM_MESSAGES.field('message')).exists();
    assert.dom('[data-test-kv-key="0"]').exists();
    assert.dom('[data-test-kv-value="0"]').exists();
    assert
      .dom(CUSTOM_MESSAGES.input('start_time'))
      .hasValue(
        format(addDays(startOfDay(timestamp.now()), 1), datetimeLocalStringFormat),
        `message start_time defaults to midnight of following day. test context start_time: ${
          this.message.start_time
        }, now: ${timestamp.now().toISOString()}`
      );
    assert.dom(CUSTOM_MESSAGES.input('end_time')).hasValue('');
  });

  test('it should display validation errors for invalid form fields', async function (assert) {
    assert.expect(8);

    await this.renderComponent();

    await fillIn(CUSTOM_MESSAGES.input('start_time'), '2024-01-20T00:00');
    await fillIn(CUSTOM_MESSAGES.input('end_time'), '2024-01-01T00:00');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('title'))
      .exists('Validation error for field `title` renders')
      .hasText('Title is required.');
    assert
      .dom(GENERAL.validationErrorByAttr('message'))
      .exists('Validation error for field `message` renders')
      .hasText('Message is required.');
    assert
      .dom(GENERAL.validationErrorByAttr('start_time'))
      .exists('Validation error for field `start_time` renders')
      .hasText('Start time is after end time.');
    assert
      .dom(GENERAL.validationErrorByAttr('end_time'))
      .exists('Validation error for field `end_time` renders')
      .hasText('End time is before start time.');
  });

  test('it should create new message', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/config/ui/custom-messages', () => {
      assert.true(true, 'POST request made to create message');
    });

    await this.renderComponent();

    await fillIn(CUSTOM_MESSAGES.input('title'), 'create new message title from component');
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await fillIn(
      CUSTOM_MESSAGES.input('start_time'),
      format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
    );
    await click('#specificDate');
    await fillIn(
      CUSTOM_MESSAGES.input('end_time'),
      format(addDays(startOfDay(new Date('2023-12-12')), 10), datetimeLocalStringFormat)
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(GENERAL.submitButton);
  });

  test('it should have form vaildations', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert.dom(`${GENERAL.validationErrorByAttr('title')}`).hasText('Title is required.');
    assert.dom(`${GENERAL.validationErrorByAttr('message')}`).hasText('Message is required.');
  });

  test('it should prepopulate form if form is in edit mode', async function (assert) {
    assert.expect(13);

    this.message = new CustomMessage({
      id: 'hhhhh-iiii-lllll-dddd',
      type: 'modal',
      authenticated: false,
      title: 'Hello world',
      message: 'Blah blah blah. Some super long message.',
      start_time: new Date('2023-12-12T08:00:00.000Z'),
      end_time: new Date('2023-12-21T08:00:00.000Z'),
      link: { 'Learn more': 'www.learnmore.com' },
    });

    await this.renderComponent();

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Edit message');
    assert.dom(CUSTOM_MESSAGES.radio('authenticated')).exists();
    assert.dom(CUSTOM_MESSAGES.radio('unauthenticated')).isChecked();
    assert.dom(CUSTOM_MESSAGES.radio('modal')).exists();
    assert.dom(CUSTOM_MESSAGES.radio('modal')).isChecked();
    assert.dom(CUSTOM_MESSAGES.input('title')).hasValue('Hello world');
    assert.dom(CUSTOM_MESSAGES.input('message')).hasValue('Blah blah blah. Some super long message.');
    assert.dom('[data-test-kv-key="0"]').exists();
    assert.dom('[data-test-kv-key="0"]').hasValue('Learn more');
    assert.dom('[data-test-kv-value="0"]').exists();
    assert.dom('[data-test-kv-value="0"]').hasValue('www.learnmore.com');
    await click('#specificDate');
    assert
      .dom(CUSTOM_MESSAGES.input('start_time'))
      .hasValue(format(new Date(this.message.start_time), datetimeLocalStringFormat));
    assert
      .dom(CUSTOM_MESSAGES.input('end_time'))
      .hasValue(format(new Date(this.message.end_time), datetimeLocalStringFormat));
  });

  test('it should show a preview image modal when preview is clicked', async function (assert) {
    await this.renderComponent();

    await fillIn(CUSTOM_MESSAGES.input('title'), 'preview modal component test');
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview modal')).doesNotExist();
    assert.dom(CUSTOM_MESSAGES.modal('preview image')).exists();
    assert
      .dom(CUSTOM_MESSAGES.alertTitle('preview modal component test'))
      .hasText('preview modal component test');
    assert
      .dom(CUSTOM_MESSAGES.alertDescription('preview modal component test'))
      .hasText('Lorem ipsum dolor sit amet, consectetur adipiscing elit.');
    assert.dom('img').hasAttribute('src', '/ui/images/custom-messages-dashboard.png');

    await click(GENERAL.button('Close preview'));
    await click('#unauthenticated');
    await click(GENERAL.button('preview'));
    assert.dom('img').hasAttribute('src', '/ui/images/custom-messages-login.png');
  });

  test('it should show a preview modal when preview is clicked', async function (assert) {
    await this.renderComponent();

    await click(CUSTOM_MESSAGES.radio('modal'));
    await fillIn(CUSTOM_MESSAGES.input('title'), 'Preview modal title');
    await fillIn(CUSTOM_MESSAGES.input('message'), 'Some preview modal message thats super long.');
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview modal')).exists();
    assert.dom(CUSTOM_MESSAGES.modal('preview image')).doesNotExist();
    assert.dom(CUSTOM_MESSAGES.modalTitle('Preview modal title')).hasText('Preview modal title');
    assert
      .dom(CUSTOM_MESSAGES.modalBody('Preview modal title'))
      .hasText('Some preview modal message thats super long.');
  });

  test('it should show multiple modal message', async function (assert) {
    this.messages = [
      {
        id: '01234567-89ab-cdef-0123-456789abcdef',
        active: true,
        type: 'modal',
        authenticated: true,
        title: 'Message title 1',
        message: 'Some long long long message',
        link: { here: 'www.example.com' },
        start_time: new Date('2021-08-01T00:00:00Z'),
        end_time: '',
      },
      {
        id: '01234567-89ab-vvvv-0123-456789abcdef',
        active: true,
        type: 'modal',
        authenticated: false,
        title: 'Message title 2',
        message: 'Some long long long message',
        link: { here: 'www.example.com' },
        start_time: new Date('2021-08-01T00:00:00Z'),
        end_time: new Date('2090-08-01T00:00:00Z'),
      },
    ];

    await this.renderComponent();

    await fillIn(CUSTOM_MESSAGES.input('title'), 'multiple modal message component test');
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await click(CUSTOM_MESSAGES.radio('modal'));
    await click(GENERAL.submitButton);
    assert.dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal')).exists();
    assert
      .dom(CUSTOM_MESSAGES.modalBody('Warning: more than one modal'))
      .hasText(
        'You have an active modal configured after the user logs in and are trying to create another one. It is recommended to avoid having more than one modal at once as it can be intrusive for users. Would you like to continue creating your message? Click “Confirm” to continue.'
      );

    await click(GENERAL.button('confirm-multiple'));
  });
});
