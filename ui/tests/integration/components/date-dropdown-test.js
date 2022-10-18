/* eslint-disable qunit/no-conditional-assertions */
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
    let dropdownListMonths = findAll('[data-test-month-list] button');

    assert.strictEqual(dropdownListMonths.length, 12, 'dropdown has 12 months');
    for (let [index, month] of ARRAY_OF_MONTHS.entries()) {
      assert.dom(dropdownListMonths[index]).hasText(`${month}`, `dropdown includes ${month}`);
    }

    await click(dropdownListMonths[0]);
    assert.dom(monthDropdown).hasText('January', 'dropdown selects January');
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting month');

    await click(yearDropdown);
    let dropdownListYears = findAll('[data-test-year-list] button');
    assert.strictEqual(dropdownListYears.length, 5, 'dropdown has 5 years');

    for (let [index, year] of dropdownListYears.entries()) {
      let comparisonYear = CURRENT_YEAR - index;
      assert.dom(year).hasText(`${comparisonYear}`, `dropdown includes ${comparisonYear}`);
    }

    await click(dropdownListYears[0]);
    assert.dom(yearDropdown).hasText(`${CURRENT_YEAR}`, `dropdown selects ${CURRENT_YEAR}`);
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting year');
    assert.false(submitButton.disabled, 'button enabled when month and year selected');

    await click(submitButton);
  });

  test('it disables correct years when selecting month first', async function (assert) {
    assert.expect(60);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');

    // select each month and assert year is enabled/disabled correctly
    for (let monthIdx = 0; monthIdx < 12; monthIdx++) {
      await click(monthDropdown);
      let dropdownListMonths = findAll('[data-test-month-list] button');
      await click(dropdownListMonths[monthIdx]);
      await click(yearDropdown);
      let dropdownListYears = findAll('[data-test-year-list] button');

      if (monthIdx <= CURRENT_MONTH) {
        for (let year of dropdownListYears) {
          assert.false(year.disabled, `${ARRAY_OF_MONTHS[monthIdx]} ${year.innerText} enabled`);
        }
      } else {
        for (let [yearIndex, year] of dropdownListYears.entries()) {
          if (yearIndex === 0) {
            assert.true(year.disabled, `${ARRAY_OF_MONTHS[monthIdx]} ${year.innerText} disabled`);
          } else {
            assert.false(year.disabled, `${ARRAY_OF_MONTHS[monthIdx]} ${year.innerText} enabled`);
          }
        }
      }
      await click(yearDropdown);
    }
  });

  test('it disables correct months when selecting year first', async function (assert) {
    assert.expect(60);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    const monthDropdown = find('[data-test-popup-menu-trigger="month"]');
    const yearDropdown = find('[data-test-popup-menu-trigger="year"]');

    // select each year and assert each month is enabled/disabled correctly
    for (let yearIdx = 0; yearIdx < 5; yearIdx++) {
      await click(yearDropdown);
      let dropdownListYears = findAll('[data-test-year-list] button');
      await click(dropdownListYears[yearIdx]);

      await click(monthDropdown);
      let dropdownListMonths = findAll('[data-test-month-list] button');

      if (yearIdx === 0) {
        // current year is selected
        for (let [monthIndex, month] of dropdownListMonths.entries()) {
          if (monthIndex <= CURRENT_MONTH) {
            assert.false(
              month.disabled,
              `${ARRAY_OF_MONTHS[monthIndex]} ${dropdownListYears[yearIdx].innerText.trim()} enabled`
            );
          } else {
            assert.true(
              month.disabled,
              `${ARRAY_OF_MONTHS[monthIndex]} ${dropdownListYears[yearIdx].innerText.trim()} disabled`
            );
          }
        }
      } else {
        // past year is selected
        for (let [monthIndex, month] of dropdownListMonths.entries()) {
          assert.false(
            month.disabled,
            `${ARRAY_OF_MONTHS[monthIndex]} ${dropdownListYears[yearIdx].innerText.trim()} enabled`
          );
        }
      }
      await click(monthDropdown);
    }
  });
});
