/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { clickable, create, isPresent } from 'ember-cli-page-object';

export default create({
  clickPreviousYear: clickable('[data-test-previous-year]'),
  clickCurrentMonth: clickable('[data-test-current-month]'),
  clickCurrentBillingPeriod: clickable('[data-test-current-billing-period]'),
  customEndMonthBtn: clickable('[data-test-show-calendar]'),
  menuToggle: clickable('[data-test-calendar-widget-trigger]'),
  showsCalendar: isPresent('[data-test-calendar-widget-container]'),
  async openCalendar() {
    await this.menuToggle();
    await this.customEndMonthBtn();
    return;
  },
});
