import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  // ARG TODO: if other components need access to this response, will need to move it higher up to replication parent.
  checkedValue: false,
  handleToggleOff: computed('checkedValue', function() {
    console.log(this.checkedValue, 'setup to handle change later');
    let status = this.checkedValue ? 'on' : 'off';
    return status;
  }),
});
