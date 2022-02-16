import { clickable, create, isPresent } from 'ember-cli-page-object';

export default create({
  menuToggle: clickable('[data-test-popup-menu-trigger="true"]'),
  customEndMonthBtn: clickable('[data-test-show-calendar]'),
  showsCalendar: isPresent('[data-test-calendar-widget-container]'),

  async openCalendar() {
    await this.menuToggle();
    await this.customEndMonthBtn();
    return;
  },
});
