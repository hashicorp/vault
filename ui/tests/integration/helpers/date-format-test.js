/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { find, render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const TEST_DATE = new Date('2018-04-03T14:15:30');

module('Integration | Helper | date-format', function (hooks) {
  setupRenderingTest(hooks);

  hooks.before(function () {
    this.today = TEST_DATE;
  });

  [
    { type: 'date object', value: TEST_DATE },
    { type: 'date strings', value: TEST_DATE.toDateString() },
    { type: 'ISO strings', value: TEST_DATE.toISOString() },
    { type: 'UTC strings', value: TEST_DATE.toUTCString() },
    { type: 'millis from epoch as string', value: TEST_DATE.getTime().toString() },
    { type: 'seconds from epoch as string', value: '1522782930' },
    { type: 'millis from epoch', value: TEST_DATE.getTime() },
    { type: 'seconds from epoch', value: 1522782930 },
  ].forEach((testCase) => {
    test(`it supports formatting ${testCase.type}`, async function (assert) {
      this.set('value', testCase.value);

      await render(hbs`{{date-format this.value "MM/dd/yyyy"}}`);
      assert
        .dom(this.element)
        .hasText(
          '04/03/2018',
          `it renders the date if passed in as a ${testCase.type} (eg. ${testCase.value})`
        );
    });
  });

  test('displays time zone if withTimeZone=true', async function (assert) {
    this.set('withTimezone', true);
    this.set('timestampDate', TEST_DATE);

    await render(
      hbs`<span data-test-formatted>{{date-format this.timestampDate 'yyyy' withTimeZone=this.withTimezone}}</span>`
    );
    const result = find('[data-test-formatted]');
    // Compare to with timezone, which should add 4 characters
    // Testing the difference because depending on the time of year the value may change
    const resultLengthWithTimezone = result.innerText.length;
    assert.strictEqual(resultLengthWithTimezone - 4, 4, 'Adds 4 characters for timezone');
  });

  test('fails gracefully if given a non-date value', async function (assert) {
    this.set('value', 'not a date');

    await render(hbs`{{date-format this.value 'yyyy'}}`);
    assert.dom(this.element).hasText('not a date', 'renders the value passed when non-date string');

    this.set('value', { date: 7 });
    await settled();
    assert.dom(this.element).hasText('[object Object]', 'renders object when non-date object');

    this.set('value', undefined);
    await settled();
    assert.dom(this.element).hasText('', 'renders empty string when falsey');
  });
});
