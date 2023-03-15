/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import timestamp from 'core/utils/timestamp';

module('Integration | Helper | date-format', function (hooks) {
  setupRenderingTest(hooks);

  test('it is able to format a date object', async function (assert) {
    const today = timestamp.now();
    this.set('today', today);

    await render(hbs`{{date-format this.today "yyyy"}}`);
    assert.dom(this.element).includesText('2018', 'it renders the date in the year format');
  });

  test('it supports date timestamps', async function (assert) {
    const today = this.today.getTime();
    this.set('today', today);

    await render(hbs`{{date-format this.today 'hh:mm:ss'}}`);
    const formattedDate = this.element.innerText;
    assert.strictEqual(formattedDate, '2:15:30');
  });

  test('it supports date strings', async function (assert) {
    const todayString = this.today.getFullYear().toString();
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
    const formattedDate = this.today;
    this.set('formattedDate', formattedDate);

    await render(hbs`{{date-format this.formattedDate 'MMMM dd, yyyy hh:mm:ss a' isFormatted=true}}`);
    assert.dom(this.element).hasText('April 03, 2018 02:15:20 PM');
  });

  test('displays time zone if withTimeZone=true', async function (assert) {
    const timestampDate = '2022-12-06T11:29:15-08:00';
    const zone = timestamp.now().toLocaleTimeString(undefined, { timeZoneName: 'short' }).split(' ')[2];
    this.set('timestampDate', timestampDate);

    await render(hbs`{{date-format this.timestampDate 'MMM d yyyy, h:mm:ss aaa' withTimeZone=true}}`);
    assert.dom(this.element).hasTextContaining(`${zone}`);
  });

  test('it returns the date passed in if it cannot be parsed', async function (assert) {
    this.set('timestampDate', 'foobar');

    await render(hbs`{{date-format this.timestampDate 'MMM d yyyy, h:mm:ss aaa' withTimeZone=true}}`);
    assert.dom(this.element).hasText('foobar');
  });
});
