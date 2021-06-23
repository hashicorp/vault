import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { subMonths, startOfToday, format, endOfMonth } from 'date-fns';

module('Integration | Component | pricing-metrics-dates', function(hooks) {
  setupRenderingTest(hooks);

  test('by default it sets the start and end inputs', async function(assert) {
    const expectedEnd = subMonths(startOfToday(), 1);
    const expectedStart = subMonths(expectedEnd, 12);
    await render(hbs`
      <PricingMetricsDates />
    `);
    assert.dom('[data-test-end-input]').hasValue(format(expectedEnd, 'MM/yyyy'), 'End input is last month');
    assert
      .dom('[data-test-start-input]')
      .hasValue(format(expectedStart, 'MM/yyyy'), 'Start input is 12 months before last month');
  });

  test('On init if end date passed, start is calculated', async function(assert) {
    const expectedStart = subMonths(new Date(2020, 8, 15), 12);
    this.set('queryEnd', '09-2020');
    await render(hbs`
      <PricingMetricsDates @queryEnd={{queryEnd}} />
    `);
    assert.dom('[data-test-end-input]').hasValue('09/2020', 'End input matches query');
    assert
      .dom('[data-test-start-input]')
      .hasValue(format(expectedStart, 'MM/yyyy'), 'Start input is 12 months before end input');
  });

  test('On init if query start date passed, end is default', async function(assert) {
    const expectedEnd = subMonths(startOfToday(), 1);
    this.set('queryStart', '01-2020');
    await render(hbs`
      <PricingMetricsDates @queryStart={{queryStart}} />
    `);
    assert.dom('[data-test-end-input]').hasValue(format(expectedEnd, 'MM/yyyy'), 'End input is last month');
    assert.dom('[data-test-start-input]').hasValue('01/2020', 'Start input matches query');
  });

  test('If result and query dates are within 1 day, warning is not shown', async function(assert) {
    this.set('resultStart', new Date(2020, 1, 1));
    this.set('resultEnd', new Date(2020, 9, 31));
    await render(hbs`
      <PricingMetricsDates
        @queryStart="2-2020"
        @queryEnd="10-2020"
        @resultStart={{resultStart}}
        @resultEnd={{resultEnd}}
      />
    `);
    assert.dom('[data-test-results-date-warning]').doesNotExist('Does not show result states warning');
  });

  test('If result and query start dates are > 1 day apart, warning is shown', async function(assert) {
    this.set('resultStart', new Date(2020, 1, 20));
    this.set('resultEnd', new Date(2020, 9, 31));
    await render(hbs`
      <PricingMetricsDates
        @queryStart="2-2020"
        @queryEnd="10-2020"
        @resultStart={{resultStart}}
        @resultEnd={{resultEnd}}
      />
    `);
    assert.dom('[data-test-results-date-warning]').exists('shows states warning');
  });

  test('If result and query end dates are > 1 day apart, warning is shown', async function(assert) {
    this.set('resultStart', new Date(2020, 1, 1));
    this.set('resultEnd', new Date(2020, 9, 15));
    await render(hbs`
      <PricingMetricsDates
        @queryStart="2-2020"
        @queryEnd="10-2020"
        @resultStart={{resultStart}}
        @resultEnd={{resultEnd}}
      />
    `);
    assert.dom('[data-test-results-date-warning]').exists('shows states warning');
  });

  test('it shows appropriate errors on input form', async function(assert) {
    const lastAvailable = endOfMonth(subMonths(startOfToday(), 1));
    const firstAvailable = subMonths(lastAvailable, 12);
    await render(hbs`
      <PricingMetricsDates @retentionMonths=12 @defaultSpan=6 />
    `);
    assert.dom('[data-test-form-error]').doesNotExist('No form error shows by default');

    await fillIn('[data-test-start-input]', format(subMonths(firstAvailable, 1), 'MM/yyyy'));
    assert
      .dom('[data-test-form-error]')
      .includesText(
        `No data retained before ${format(firstAvailable, 'MM/yyyy')}`,
        'shows the correct error message for starting before the configured retainment period'
      );

    await fillIn('[data-test-end-input]', format(subMonths(lastAvailable, -1), 'MM/yyyy'));
    assert
      .dom('[data-test-form-error]')
      .includesText(
        'Data is not available until the end of the month',
        'shows the correct error message for ending after the end of the last month'
      );

    await fillIn('[data-test-end-input]', 'not/date');
    assert
      .dom('[data-test-form-error]')
      .includesText(
        'End date is invalid. Please use format MM/yyyy',
        'shows the correct error message for non-date input'
      );

    await fillIn('[data-test-start-input]', `13/${format(lastAvailable, 'yyyy')}`);
    assert
      .dom('[data-test-form-error]')
      .includesText(
        'Start date is invalid. Please use format MM/yyyy',
        'shows the correct error message for an invalid month'
      );
  });
});
