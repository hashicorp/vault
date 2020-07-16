import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
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

    await render(hbs`<p data-test-date-format>Date: {{date-format todayString 'dd MM yyyy'}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText('16 July 2020', 'it renders the date if passed in as a string');
  });

  test('it supports ISO strings', async function(assert) {
    let iso = '2016-01-01';
    this.set('iso', iso);

    await render(hbs`<p data-test-date-format>Date: {{date-format iso}}</p>`);
    assert
      .dom('[data-test-date-format]')
      .includesText(iso, 'it renders the a date if passed in as an ISO string');
  });

  test('it fails gracefully', async function(assert) {
    let antiDate = 'lol';
    this.set('antiDate', antiDate);

    await render(hbs`<p data-test-date-format>Date: {{date-format antiDate}}</p>`);
    assert.dom('[data-test-date-format]').includesText(antiDate, 'it renders what it is passed');
  });
});
