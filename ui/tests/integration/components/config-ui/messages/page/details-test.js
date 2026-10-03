/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { dateFormat } from 'core/helpers/date-format';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const allFields = [
  { label: 'Active', key: 'active' },
  { label: 'Type', key: 'type' },
  { label: 'Authenticated', key: 'authenticated' },
  { label: 'Title', key: 'title' },
  { label: 'Message', key: 'message' },
  { label: 'Start time', key: 'start_time' },
  { label: 'End time', key: 'end_time' },
  { label: 'Link', key: 'link' },
];

module('Integration | Component | messages/page/details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.context = { owner: this.engine };

    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));

    this.message = {
      id: '01234567-89ab-cdef-0123-456789abcdef',
      active: true,
      type: 'banner',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { here: 'www.example.com' },
      start_time: new Date('2021-08-01T00:00:00Z'),
      end_time: undefined,
    };
    this.capabilities = { canDelete: true, canUpdate: true };
    this.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: this.message.authenticated } },
      { label: this.message.title },
    ];
  });

  test('it should show the message details', async function (assert) {
    await render(
      hbs`<Messages::Page::Details @message={{this.message}} @capabilities={{this.capabilities}} @breadcrumbs={{this.breadcrumbs}} />`,
      this.context
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Message title 1');
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: allFields.length }, 'Correct number of filtered fields render');

    allFields.forEach((field) => {
      assert.dom(GENERAL.infoRowLabel(field.label)).hasText(field.label, `${field.label} label renders`);
      if (field.key === 'start_time' || field.key === 'end_time') {
        const formattedDate = dateFormat([this.message[field.key], 'MMM d, yyyy hh:mm aaa'], {
          withTimeZone: true,
        });
        assert
          .dom(GENERAL.infoRowValue(field.label))
          .hasText(formattedDate || 'Never', `${field.label} value renders`);
      } else if (field.key === 'authenticated' || field.key === 'active') {
        assert
          .dom(GENERAL.infoRowValue(field.label))
          .hasText(this.message[field.key] ? 'Yes' : 'No', `${field.label} value renders`);
      } else if (field.key === 'link') {
        assert.dom(GENERAL.infoRowValue('Link')).hasText('here');
      } else {
        assert
          .dom(GENERAL.infoRowValue(field.label))
          .hasText(this.message[field.key], `${field.label} value renders`);
      }
    });

    assert.dom(GENERAL.confirmTrigger).exists('delete button exists');
    assert.dom(GENERAL.linkTo('edit')).exists();
  });
});
