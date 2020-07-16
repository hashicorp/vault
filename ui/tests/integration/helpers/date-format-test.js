import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find } from '@ember/test-helpers';
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
    let todayString = 'Thu Jul 16 2020 09:13:57';
    this.set('todayString', todayString);

    await render(hbs`<p data-test-date-format>Date: {{date-format todayString 'dd MMMM yyyy'}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText('16 July 2020', 'it renders the date if passed in as a string');
  });

  test('it supports ISO strings', async function(assert) {
    let iso = '2014-02-11T11:30:30';
    this.set('iso', iso);

    await render(hbs`<p data-test-date-format>Date: {{date-format iso 'dd MMMM yyyy'}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText('11 February 2014', 'it renders the formatted date if passed in as an ISO string');
  });

  test('it fails gracefully with strings', async function(assert) {
    let antiDate = 'lol';
    this.set('antiDate', antiDate);

    await render(hbs`<p data-test-date-format>Date: {{date-format antiDate}}</p>`);
    assert.dom('[data-test-date-format]').includesText(antiDate, 'it renders what it is passed');
  });

  test('it fails gracefully with non-strings', async function(assert) {
    let nonDate = { text: 'lol' };
    this.set('nonDate', nonDate);

    await render(hbs`<p data-test-date-format>Date: -{{date-format nonDate}}-</p>`);
    assert.equal(find('[data-test-date-format]').innerText, 'Date: --', 'it renders an empty string');
  });

  test('it renders a formatted date if no format is passed', async function(assert) {
    let date = new Date(2020, 0, 20);
    this.set('date', date);

    await render(hbs`<p data-test-date-format>Date: {{date-format date}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText('20 Jan 2020', 'it renders the date in a default format');
  });
});
