import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { staticNow } from 'vault/tests/helpers/stubs';

const SELECTORS = {
  monthDropdown: '[data-test-popup-menu-trigger="month"]',
  specificMonth: (m) => `[data-test-dropdown-month="${m}"]`,
  yearDropdown: '[data-test-popup-menu-trigger="year"]',
  specificYear: (y) => `[data-test-dropdown-year="${y}"]`,

  submitButton: '[data-test-date-dropdown-submit]',
  cancelButton: '[data-test-date-dropdown-cancel]',
  monthOptions: '[data-test-month-list] button',
};

module('Integration | Component | date-dropdown', function (hooks) {
  setupRenderingTest(hooks);

  hooks.before(function () {
    sinon.stub(Date, 'now').returns(staticNow.getTime());
  });

  test('it renders dropdown', async function (assert) {
    await render(hbs`
      <div class="is-flex-align-baseline">
        <DateDropdown/>
      </div>
    `);
    assert.dom(SELECTORS.submitButton).hasText('Submit', 'button renders default text');
    assert.dom(SELECTORS.cancelButton).doesNotExist('it does not render cancel button by default');
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
    assert.dom(SELECTORS.submitButton).hasText('Save', 'button renders passed in text');
    await click(SELECTORS.cancelButton);
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
          year: 2021,
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
    assert.dom(SELECTORS.submitButton).isDisabled('button is disabled when no month or year selected');

    await click(SELECTORS.monthDropdown);

    assert.dom(SELECTORS.monthOptions).exists({ count: 12 }, 'dropdown has 12 months');
    ARRAY_OF_MONTHS.forEach((month) => {
      assert.dom(SELECTORS.specificMonth(month)).hasText(`${month}`, `dropdown includes ${month}`);
    });

    await click(SELECTORS.specificMonth('January'));
    assert.dom(SELECTORS.monthDropdown).hasText('January', 'dropdown shows selected January');
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting month');

    await click(SELECTORS.yearDropdown);

    assert.dom('[data-test-dropdown-year]').exists({ count: 5 }, 'dropdown has 5 years');
    for (const year of [2023, 2022, 2021, 2020, 2019]) {
      assert.dom(SELECTORS.specificYear(year)).exists();
    }

    await click('[data-test-dropdown-year="2021"]');
    assert.dom(SELECTORS.yearDropdown).hasText(`2021`, `dropdown shows selected 2021`);
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting year');
    assert.dom(SELECTORS.submitButton).isNotDisabled('button enabled when month and year selected');

    await click(SELECTORS.submitButton);
  });

  test('selecting month first: current year enabled when current month selected', async function (assert) {
    assert.expect(5);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);
    // select current month
    await click(SELECTORS.monthDropdown);
    await click(SELECTORS.specificMonth('January'));
    await click(SELECTORS.yearDropdown);
    // all years should be selectable
    for (const year of [2023, 2022, 2021, 2020, 2019]) {
      assert.dom(SELECTORS.specificYear(year)).isNotDisabled(`year ${year} is selectable`);
    }
  });

  test('selecting month first: it disables current year when future months selected', async function (assert) {
    assert.expect(5);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    // select future month
    await click(SELECTORS.monthDropdown);
    await click(SELECTORS.specificMonth('June'));
    await click(SELECTORS.yearDropdown);

    assert.dom(SELECTORS.specificYear(2023)).isDisabled(`year ${2023} is disabled`);
    // previous years should be selectable
    for (const year of [2022, 2021, 2020, 2019]) {
      assert.dom(SELECTORS.specificYear(year)).isNotDisabled(`year ${year} is selectable`);
    }
  });

  test('selecting year first: it disables future months when current year selected', async function (assert) {
    assert.expect(12);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);
    await click(SELECTORS.yearDropdown);
    await click(SELECTORS.specificYear(2023));
    await click(SELECTORS.monthDropdown);

    const expectedSelectable = ['January'];
    ARRAY_OF_MONTHS.forEach((month) => {
      if (expectedSelectable.includes(month)) {
        assert.dom(SELECTORS.specificMonth(month)).isNotDisabled(`${month} is selectable for current year`);
      } else {
        assert.dom(SELECTORS.specificMonth(month)).isDisabled(`${month} is disabled for current year`);
      }
    });
  });

  test('selecting year first: it enables all months when past year is selected', async function (assert) {
    assert.expect(12);
    await render(hbs`
    <div class="is-flex-align-baseline">
      <DateDropdown/>
    </div>
    `);

    await click(SELECTORS.yearDropdown);
    await click(SELECTORS.specificYear(2022));
    await click(SELECTORS.monthDropdown);

    ARRAY_OF_MONTHS.forEach((month) => {
      assert.dom(SELECTORS.specificMonth(month)).isNotDisabled(`${month} is selectable for previous year`);
    });
  });
});
