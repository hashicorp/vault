import Actions from 'core/components/replication-actions-single';
import layout from '../templates/components/replication-action-recover';
import keys from 'vault/lib/keycodes';

export default Actions.extend({
  layout,
  onSubmit() {},
  actions: {
    onSubmit() {
      if (event.keyCode === keys.ESC) {
        // if escape close modal and return
        this.toggleProperty('isModalActive');
        return;
      }
      return this.onSubmit(...arguments);
    },
  },
});
