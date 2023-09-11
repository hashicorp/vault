/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { find, render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import timestamp from 'core/utils/timestamp';

module('Integration | Helper | date-format', function (hooks) {
  setupRenderingTest(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => new Date('2018-04-03T14:15:30'));
  });
  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it is able to format a date object', async function (assert) {
    const today = timestamp.now();
    this.set('today', today);

    await render(hbs`{{date-format this.today "yyyy"}}`);
    assert.dom(this.element).includesText('2018', 'it renders the date in the year format');
  });

  test('it supports date timestamps', async function (assert) {
    const today = timestamp.now().getTime();
    this.set('today', today);

    await render(hbs`{{date-format this.today 'hh:mm:ss'}}`);
    const formattedDate = this.element.innerText;
    assert.strictEqual(formattedDate, '02:15:30');
  });

  test('it supports date strings', async function (assert) {
    const todayString = timestamp.now().getFullYear().toString();
    this.set('todayString', todayString);

    await render(hbs`{{date-format this.todayString "yyyy"}}`);
    assert.dom(this.element).includesText(todayString, 'it renders the a date if passed in as a string');
  });

  test('it supports ten digit dates', async function (assert) {
    const tenDigitDate = 1621785298;
    this.set('tenDigitDate', tenDigitDate);

    await render(hbs`{{date-format this.tenDigitDate "MM/dd/yyyy"}}`);
    assert.dom(this.element).includesText('05/23/2021');
  });

  test('it supports already formatted dates', async function (assert) {
    const formattedDate = timestamp.now();
    this.set('formattedDate', formattedDate);

    await render(hbs`{{date-format this.formattedDate 'MMMM dd, yyyy hh:mm:ss a' isFormatted=true}}`);
    assert.dom(this.element).hasText('April 03, 2018 02:15:30 PM');
  });

  test('displays time zone if withTimeZone=true', async function (assert) {
    const timestampDate = '2022-12-06T11:29:15-08:00';
    this.set('withTimezone', false);
    this.set('timestampDate', timestampDate);

    await render(
      hbs`<span data-test-formatted>{{date-format this.timestampDate 'MMM d yyyy, h:mm:ss aaa' withTimeZone=this.withTimezone}}</span>`
    );
    const result = find('[data-test-formatted]');
    const resultLength = result.innerText.length;
    // Compare to with timezone, which should add 4 characters
    // Testing the difference because depending on the day the length may change.
    this.set('withTimezone', true);
    await settled();
    const resultLengthWithTimezone = result.innerText.length;
    assert.strictEqual(resultLengthWithTimezone - resultLength, 4);
  });
});
