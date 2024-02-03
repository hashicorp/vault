/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { dateFormat } from 'core/helpers/date-format';

const allFields = [
  { label: 'Active', key: 'active' },
  { label: 'Type', key: 'type' },
  { label: 'Authenticated', key: 'authenticated' },
  { label: 'Title', key: 'title' },
  { label: 'Message', key: 'message' },
  { label: 'Start time', key: 'startTime' },
  { label: 'End time', key: 'endTime' },
  { label: 'Link', key: 'link' },
];

module('Integration | Component | messages/page/details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');

    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));

    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '01234567-89ab-cdef-0123-456789abcdef',
      active: true,
      type: 'banner',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { here: 'www.example.com' },
      start_time: '2021-08-01T00:00:00Z',
      end_time: '',
      canDeleteCustomMessages: true,
      canEditCustomMessages: true,
    });
  });

  test('it should show the message details', async function (assert) {
    this.message = await this.store.peekRecord('config-ui/message', '01234567-89ab-cdef-0123-456789abcdef');

    await render(hbs`<Messages::Page::Details @message={{this.message}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-page-title]').hasText('Message title 1');
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: allFields.length }, 'Correct number of filtered fields render');

    allFields.forEach((field) => {
      assert
        .dom(`[data-test-row-label="${field.label}"]`)
        .hasText(field.label, `${field.label} label renders`);
      if (field.key === 'startTime' || field.key === 'endTime') {
        const formattedDate = dateFormat([this.message[field.key], 'MMM d, yyyy hh:mm aaa'], {
          withTimeZone: true,
        });
        assert
          .dom(`[data-test-row-value="${field.label}"]`)
          .hasText(formattedDate || 'Never', `${field.label} value renders`);
      } else if (field.key === 'authenticated' || field.key === 'active') {
        assert
          .dom(`[data-test-value-div="${field.label}"]`)
          .hasText(this.message[field.key] ? 'Yes' : 'No', `${field.label} value renders`);
      } else if (field.key === 'link') {
        assert.dom('[data-test-value-div="Link"]').exists();
        assert.dom('[data-test-value-div="Link"] [data-test-link="message link"]').hasText('here');
      } else {
        assert
          .dom(`[data-test-row-value="${field.label}"]`)
          .hasText(this.message[field.key], `${field.label} value renders`);
      }
    });
  });
});
