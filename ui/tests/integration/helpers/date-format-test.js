import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { format } from 'date-fns';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | date-format', function(hooks) {
  setupRenderingTest(hooks);

  test('it is able to format a date object', async function(assert) {
    let today = new Date();
    this.set('today', today);

    await render(hbs`<p data-test-date-format>Date: {{date-format today "yyyy"}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText(today.getFullYear(), 'it renders the date in the year format');
  });

  test('it supports date timestamps', async function(assert) {
    let today = new Date().getTime();
    this.set('today', today);

    await render(hbs`<p class="date-format">{{date-format today 'hh:mm:ss'}}</p>`);
    let formattedDate = document.querySelector('.date-format').innerText;
    assert.ok(formattedDate.match(/^\d{2}:\d{2}:\d{2}$/));
  });

  test('it supports date strings', async function(assert) {
    let todayString = new Date().getFullYear().toString();
    this.set('todayString', todayString);

    await render(hbs`<p data-test-date-format>Date: {{date-format todayString "yyyy"}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText(todayString, 'it renders the a date if passed in as a string');
  });

  test('it supports ten digit dates', async function(assert) {
    let tenDigitDate = 1621785298;
    this.set('tenDigitDate', tenDigitDate);

    await render(hbs`<p data-test-date-format>Date: {{date-format tenDigitDate "MM/dd/yyyy"}}</p>`);
    assert.dom('[data-test-date-format]').includesText('05/23/2021');
  });

  test('it supports already formatted dates', async function(assert) {
    let formattedDate = new Date();
    this.set('formattedDate', formattedDate);

    await render(
      hbs`<p data-test-date-format>Date: {{date-format formattedDate 'MMMM dd, yyyy hh:mm:ss a' isFormatted=true}}</p>`
    );
    assert.dom('[data-test-date-format]').includesText(format(formattedDate, 'MMMM dd, yyyy hh:mm:ss a'));
  });

  test('displays correct date when timestamp is in ISO 8601 format', async function(assert) {
    let timestampDate = '2021-09-01T00:00:00Z';
    this.set('timestampDate', timestampDate);

    await render(
      hbs`<p data-test-date-format>Date: {{date-format timestampDate 'MMM dd, yyyy' dateOnly=true}}</p>`
    );
    assert.dom('[data-test-date-format]').includesText('Date: Sep 01, 2021');
  });
});
