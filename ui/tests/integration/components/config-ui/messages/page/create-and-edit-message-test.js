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
import { localDateTimeString } from 'vault/models/config-ui/message';
import { format, addDays, startOfDay } from 'date-fns';

const PAGE = {
  // General selectors that are common between pages
  radioGroup: (groupName) => `[data-test-radio-group="${groupName}"]`,
  radioField: (fieldName) => `[data-test-radio-field="${fieldName}"]`,
  field: (fieldName) => `[data-test-field="${fieldName}"]`,
  input: (input) => `[data-test-input="${input}"]`,
  button: (buttonName) => `[data-test-button="${buttonName}"]`,
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
    await render(
      hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} @authenticated={{true}} />`,
      {
        owner: this.engine,
      }
    );

    assert.dom('[data-test-page-title]').hasText('Create message');
    assert
      .dom('[data-test-form-subtext]')
      .hasText('Create a custom message for all users when they access a Vault system via the UI.');
    assert.dom(PAGE.radioGroup('authenticated')).exists();
    assert.dom(`${PAGE.radioGroup('authenticated')} ${PAGE.radioField('authenticated')}`).isChecked();
    assert.dom(PAGE.radioGroup('type')).exists();
    assert.dom(`${PAGE.radioGroup('type')} ${PAGE.radioField('banner')}`).isChecked();
    assert.dom(PAGE.field('title')).exists();
    assert.dom(PAGE.field('message')).exists();
    assert.dom(PAGE.field('linkTitle')).exists();
    assert.dom(PAGE.field('linkHref')).exists();
    assert.dom(PAGE.input('startTime')).exists();
    assert
      .dom(PAGE.input('startTime'))
      .hasValue(format(addDays(startOfDay(new Date()), 1), localDateTimeString));
    assert.dom(PAGE.input('endTime')).exists();
    assert.dom(PAGE.input('endTime')).hasValue('');
  });

  test('it should create new message', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/config/ui/custom-messages', () => {
      assert.ok(true, 'POST request made to create message');
    });

    await render(
      hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} @authenticated={{true}} />`,
      {
        owner: this.engine,
      }
    );
    await fillIn(PAGE.input('title'), 'Awesome custom message title');
    await fillIn(
      PAGE.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
    );
    await fillIn(
      PAGE.input('startTime'),
      format(addDays(startOfDay(new Date('2023-12-12')), 1), localDateTimeString)
    );
    await click('#specific-date');
    await fillIn(
      PAGE.input('endTime'),
      format(addDays(startOfDay(new Date('2023-12-12')), 10), localDateTimeString)
    );

    await click(PAGE.button('create-message'));
  });

  test('it should show a disabled create form button when there are no title or message', async function (assert) {
    await render(
      hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} @authenticated={{true}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(PAGE.button('create-message')).isDisabled();
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
    });
    this.message = this.store.peekRecord('config-ui/message', 'hhhhh-iiii-lllll-dddd');
    await render(
      hbs`<Messages::Page::CreateAndEditMessageForm @message={{this.message}} @authenticated={{true}} />`,
      {
        owner: this.engine,
      }
    );

    assert.dom('[data-test-page-title]').hasText('Edit message');
    assert
      .dom('[data-test-form-subtext]')
      .hasText('Edit a custom message for all users when they access a Vault system via the UI.');
    assert.dom(PAGE.radioGroup('authenticated')).exists();
    // assert.dom(`${PAGE.radioGroup('authenticated')} ${PAGE.radioField('unauthenticated')}`).isChecked();
    assert.dom(PAGE.radioGroup('type')).exists();
    assert.dom(`${PAGE.radioGroup('type')} ${PAGE.radioField('modal')}`).isChecked();
    assert.dom(PAGE.input('title')).hasValue('Hello world');
    assert.dom(PAGE.input('message')).hasValue('Blah blah blah. Some super long message.');
    assert.dom(PAGE.field('linkTitle')).exists();
    assert.dom(PAGE.field('linkHref')).exists();
    // assert.dom(PAGE.input('startTime')).exists();
    // await this.pauseTest();
    // assert
    //   .dom(PAGE.input('startTime'))
    //   .hasValue(format(addDays(startOfDay(new Date(this.message.startTime)), 1), localDateTimeString));
    // assert.dom(PAGE.input('endTime')).exists();
    // assert.dom(PAGE.input('endTime')).hasValue(format(new Date(this.message.endTime), localDateTimeString));
    // await this.pauseTest();
  });
});
