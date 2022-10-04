import { clickable, create, isPresent } from 'ember-cli-page-object';

export default create({
  clickPreviousYear: clickable('[data-test-previous-year]'),
  clickCurrentBillingPeriod: clickable('[data-test-current-billing-period]'),
  customEndMonthBtn: clickable('[data-test-show-calendar]'),
  menuToggle: clickable('[data-test-popup-menu-trigger="true"]'),
  showsCalendar: isPresent('[data-test-calendar-widget-container]'),

  async openCalendar() {
    await this.menuToggle();
    await this.customEndMonthBtn();
    return;
  },
});
