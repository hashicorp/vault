import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  checkedValue: false,
  handleToggleOff: computed('checkedValue', function() {
    console.log(this.checkedValue, 'setup to handle change later');
    let status = this.checkedValue ? 'on' : 'off';
    return status;
  }),
});
