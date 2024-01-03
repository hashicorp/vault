/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Serializer | config-ui/message', function (hooks) {
  setupTest(hooks);

  test('it should always encode message when creating/updating a message', function (assert) {
    const store = this.owner.lookup('service:store');
    const record = store.createRecord('config-ui/message', {
      id: '01234567-89ab-cdef-0123-456789abcdef',
      active: true,
      type: 'banner',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { title: '', href: '' },
      startTime: '2024-01-03T20:54:29.802Z',
      endTime: '',
    });
    const expectedResult = {
      authenticated: true,
      end_time: null,
      link: {
        href: '',
        title: '',
      },
      message: 'U29tZSBsb25nIGxvbmcgbG9uZyBtZXNzYWdl',
      start_time: '2024-01-03T20:54:29.802Z',
      title: 'Message title 1',
      type: 'banner',
    };

    const serializedRecord = record.serialize();
    assert.deepEqual(serializedRecord, expectedResult, 'encode the message string');
  });

  test('it should always use ISO date format when creating/updating a message', function (assert) {
    const store = this.owner.lookup('service:store');
    const date = new Date();
    const record = store.createRecord('config-ui/message', {
      id: '01234567-89ab-cdef-0123-456789abcdef',
      active: true,
      type: 'banner',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { title: '', href: '' },
      startTime: date,
      endTime: '',
    });
    const expectedResult = {
      authenticated: true,
      end_time: null,
      link: {
        href: '',
        title: '',
      },
      message: 'U29tZSBsb25nIGxvbmcgbG9uZyBtZXNzYWdl',
      start_time: date.toISOString(),
      title: 'Message title 1',
      type: 'banner',
    };

    const serializedRecord = record.serialize();
    assert.deepEqual(serializedRecord, expectedResult, 'uses ISO date string');
  });
});
