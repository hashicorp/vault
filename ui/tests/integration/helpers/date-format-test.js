import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | date-format', function(hooks) {
  setupRenderingTest(hooks);

  test('it is able to format a date object', async function(assert) {
    let today = new Date();
    this.set('today', today);

    await render(hbs`<p data-test-date-format>Date: {{date-format today "YYYY"}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText(today.getFullYear(), 'it renders the date in the year format');
  });

  test('it formats the date as specified', async function(assert) {
    let today = new Date();
    this.set('today', today);

    await render(hbs`<p class="date-format">{{date-format today 'hh:mm:ss'}}</p>`);
    let formattedDate = document.querySelector('.date-format').innerText;
    assert.ok(formattedDate.match(/^\d{2}:\d{2}:\d{2}$/));
  });

  test('it supports date strings', async function(assert) {
    let todayString = new Date().getFullYear().toString();
    this.set('todayString', todayString);

    await render(hbs`<p data-test-date-format>Date: {{date-format todayString}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText(todayString, 'it renders the a date if passed in as a string');
  });
});
