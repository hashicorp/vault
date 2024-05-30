/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import timestamp from 'core/utils/timestamp';

const SELECTORS = {
  monthDropdown: '[data-test-toggle-month]',
  specificMonth: (m) => `[data-test-dropdown-month="${m}"]`,
  yearDropdown: '[data-test-toggle-year]',
  specificYear: (y) => `[data-test-dropdown-year="${y}"]`,
  submitButton: '[data-test-date-dropdown-submit]',
  monthOptions: '[data-test-dropdown-month]',
};

module('Integration | Component | date-dropdown', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
  });

  test('it renders dropdown', async function (assert) {
    await render(hbs`
      <div class="has-padding-l">
        <DateDropdown/>
      </div>
    `);
    assert.dom(SELECTORS.submitButton).hasText('Submit', 'button renders default text');
  });

  test('it renders dropdown and selects month and year', async function (assert) {
    assert.expect(26);
    const parentAction = (args) => {
      assert.propEqual(
        args,
        {
          dateType: 'start',
          monthIdx: 1,
          monthName: 'February',
          year: 2016,
        },
        'sends correct args to parent'
      );
    };
    this.set('parentAction', parentAction);

    await render(hbs`
    <div class="has-padding-l">
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

    await click(SELECTORS.specificMonth('February'));
    assert.dom(SELECTORS.monthDropdown).hasText('February', 'dropdown shows selected month');
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting month');

    await click(SELECTORS.yearDropdown);

    assert.dom('[data-test-dropdown-year]').exists({ count: 5 }, 'dropdown has 5 years');
    for (const year of [2018, 2017, 2016, 2015, 2014]) {
      assert.dom(SELECTORS.specificYear(year)).exists();
    }

    await click('[data-test-dropdown-year="2016"]');
    assert.dom(SELECTORS.yearDropdown).hasText(`2016`, `dropdown shows selected year`);
    assert.dom('.ember-basic-dropdown-content').doesNotExist('dropdown closes after selecting year');
    assert.dom(SELECTORS.submitButton).isNotDisabled('button enabled when month and year selected');

    await click(SELECTORS.submitButton);
  });

  test('selecting month first: current year enabled when current month selected', async function (assert) {
    assert.expect(5);
    await render(hbs`
    <div class="has-padding-l">
      <DateDropdown/>
    </div>
    `);
    // select current month
    await click(SELECTORS.monthDropdown);
    await click(SELECTORS.specificMonth('January'));
    await click(SELECTORS.yearDropdown);
    // all years should be selectable
    for (const year of [2018, 2017, 2016, 2015, 2014]) {
      assert.dom(SELECTORS.specificYear(year)).isNotDisabled(`year ${year} is selectable`);
    }
  });

  test('selecting month first: it disables current year when future months selected', async function (assert) {
    assert.expect(5);
    await render(hbs`
    <div class="has-padding-l">
      <DateDropdown/>
    </div>
    `);

    // select future month
    await click(SELECTORS.monthDropdown);
    await click(SELECTORS.specificMonth('June'));
    await click(SELECTORS.yearDropdown);

    assert.dom(SELECTORS.specificYear(2018)).isDisabled(`current year is disabled`);
    // previous years should be selectable
    for (const year of [2017, 2016, 2015, 2014]) {
      assert.dom(SELECTORS.specificYear(year)).isNotDisabled(`year ${year} is selectable`);
    }
  });

  test('selecting year first: it disables future months when current year selected', async function (assert) {
    assert.expect(12);
    await render(hbs`
    <div class="has-padding-l">
      <DateDropdown/>
    </div>
    `);
    await click(SELECTORS.yearDropdown);
    await click(SELECTORS.specificYear(2018));
    await click(SELECTORS.monthDropdown);

    const expectedSelectable = ['January', 'February', 'March', 'April'];
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
    <div class="has-padding-l">
      <DateDropdown/>
    </div>
    `);

    await click(SELECTORS.yearDropdown);
    await click(SELECTORS.specificYear(2017));
    await click(SELECTORS.monthDropdown);

    ARRAY_OF_MONTHS.forEach((month) => {
      assert.dom(SELECTORS.specificMonth(month)).isNotDisabled(`${month} is selectable for previous year`);
    });
  });
});
