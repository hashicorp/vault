import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, find, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';

const CURRENT_DATE = new Date();
const CURRENT_YEAR = CURRENT_DATE.getFullYear(); // integer of year
const CURRENT_MONTH = CURRENT_DATE.getMonth(); // index of month

module('Integration | Component | date-dropdown', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders dropdown', async function (assert) {
    await render(hbs`
      <div class="is-flex-align-baseline">
        <DateDropdown/>
      </div>
    `);
    assert.dom('[data-test-date-dropdown-submit]').hasText('Submit', 'button renders default text');
    assert
      .dom('[data-test-date-dropdown-cancel]')
      .doesNotExist('it does not render cancel button by default');
  });

  test('it fires off cancel callback', async function (assert) {
    assert.expect(2);
    const onCancel = () => {
      assert.ok('fires onCancel callback');
    };
    this.set('onCancel', onCancel);
    await render(hbs`
      <div class="is-flex-align-baseline">
        <DateDropdown @handleCancel={{this.onCancel}} @submitText="Save"/>
      </div>
    `);
    assert.dom('[data-test-date-dropdown-submit]').hasText('Save', 'button renders passed in text');
    await click(find('[data-test-date-dropdown-cancel]'));
  });

  test('it renders dropdown and selects month and year', async function (assert) {
    assert.expect(26);
    const parentAction = (args) => {
      assert.propEqual(
        args,
        {
          dateType: 'start',
          monthIdx: 0,
          monthName: 'January',
          year: 2022,
        },
        'sends correct args to parent'
      );
    };
    this.set('parentAction', parentAction);

    await render(hbs`
    <div class="is-flex-align-baseline">
    <DateDropdown 
      @handleSubmit={{this.parentAction}} 
      @dateType="start"
    />
    </div>
    `);

    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');
    const submitButton = find('[data-test-date-dropdown-submit]');

    assert.true(submitButton.disabled, 'button is disabled when no month or year selected');

    await click(monthDropdown);
    const dropdownListMonths = findAll('[data-test-month-list] button');

    assert.strictEqual(dropdownListMonths.length, 12, 'dropdown has 12 months');
    for (const [index, month] of ARRAY_OF_MONTHS.entries()) {
      assert.dom(dropdownListMonths[index]).hasText(`${month}`, `dropdown includes ${month}`);
    }

    await click(dropdownListMonths[0]);
    assert.dom(monthDropdown).hasText('January', 'dropdown selects January');
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting month');

    await click(yearDropdown);
    const dropdownListYears = findAll('[data-test-year-list] button');
    assert.strictEqual(dropdownListYears.length, 5, 'dropdown has 5 years');

    for (const [index, year] of dropdownListYears.entries()) {
      const comparisonYear = CURRENT_YEAR - index;
      assert.dom(year).hasText(`${comparisonYear}`, `dropdown includes ${comparisonYear}`);
    }

    await click(dropdownListYears[0]);
    assert.dom(yearDropdown).hasText(`${CURRENT_YEAR}`, `dropdown selects ${CURRENT_YEAR}`);
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting year');
    assert.false(submitButton.disabled, 'button enabled when month and year selected');

    await click(submitButton);
  });

  test('selecting month first: it enables current year when selecting valid months', async function (assert) {
    // the date dropdown displays 5 years, multiply by month to calculate how many assertions to expect
    const datesEnabled = (CURRENT_MONTH + 1) * 5;
    assert.expect(datesEnabled);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');

    // select months before or equal to current month and assert year is enabled
    for (let monthIdx = 0; monthIdx < CURRENT_MONTH + 1; monthIdx++) {
      await click(monthDropdown);
      const dropdownListMonths = findAll('[data-test-month-list] button');
      await click(dropdownListMonths[monthIdx]);
      await click(yearDropdown);
      const dropdownListYears = findAll('[data-test-year-list] button');
      for (const year of dropdownListYears) {
        assert.false(year.disabled, `${ARRAY_OF_MONTHS[monthIdx]} ${year.innerText} enabled`);
      }
      await click(yearDropdown);
    }
  });

  test('selecting month first: it disables current year and enables past years when selecting future month', async function (assert) {
    assert.expect(CURRENT_MONTH - 11);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');

    // select future months and assert current year (index=0) is disabled and past years enabled
    for (let monthIdx = CURRENT_MONTH + 1; monthIdx < 12; monthIdx++) {
      await click(monthDropdown);
      const dropdownListMonths = findAll('[data-test-month-list] button');
      await click(dropdownListMonths[monthIdx]);
      await click(yearDropdown);
      const dropdownListYears = findAll('[data-test-year-list] button');
      const currentYear = dropdownListYears.entries()[0];
      assert.true(currentYear.disabled, `${ARRAY_OF_MONTHS[monthIdx]} ${currentYear.innerText} disabled`);

      // remove current year, assert rest of years enabled
      dropdownListYears.entries()[0].splice(0, 1);
      for (const year of dropdownListYears.entries()) {
        assert.false(year.disabled, `${ARRAY_OF_MONTHS[monthIdx]} ${year.innerText} enabled`);
      }
      await click(yearDropdown);
    }
  });

  test('selecting year first: it disables future months when current year selected', async function (assert) {
    assert.expect(12);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);
    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');
    await click(yearDropdown);
    await click(`[data-test-dropdown-year="${CURRENT_YEAR}"]`);
    await click(monthDropdown);
    const dropdownListMonths = findAll('[data-test-month-list] button');
    const enabledMonths = dropdownListMonths.slice(0, CURRENT_MONTH + 1);
    const disabledMonths = dropdownListMonths.slice(CURRENT_MONTH + 1);
    for (const [monthIndex, month] of enabledMonths.entries()) {
      assert.false(month.disabled, `${ARRAY_OF_MONTHS[monthIndex]} ${CURRENT_YEAR} enabled`);
    }
    for (const [monthIndex, month] of disabledMonths.entries()) {
      assert.true(month.disabled, `${ARRAY_OF_MONTHS[monthIndex]} ${CURRENT_YEAR} disabled`);
    }
  });

  test('selecting year first: it enables all months when past year is selected', async function (assert) {
    assert.expect(48);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');

    // start at 1 because current year (index=0) is accounted for in previous test
    for (let yearIdx = 1; yearIdx < 5; yearIdx++) {
      await click(yearDropdown);
      const dropdownListYears = findAll('[data-test-year-list] button');
      await click(dropdownListYears[yearIdx]);
      await click(monthDropdown);
      const dropdownListMonths = findAll('[data-test-month-list] button');
      for (const [monthIndex, month] of dropdownListMonths.entries()) {
        assert.false(
          month.disabled,
          `${ARRAY_OF_MONTHS[monthIndex]} ${dropdownListYears[yearIdx].innerText.trim()} enabled`
        );
      }
      await click(monthDropdown);
    }
  });
});
