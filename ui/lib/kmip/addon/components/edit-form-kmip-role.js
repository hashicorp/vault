import EditForm from 'core/components/edit-form';
import layout from '../templates/components/edit-form-kmip-role';
import { Promise } from 'rsvp';
import { computed } from '@ember/object';

export default EditForm.extend({
  layout,
  model: null,

  init() {
    this._super(...arguments);

    if (this.model.isNew) {
      this.model.set('operationAll', true);
    }
  },

  actions: {
    switchUpdated(checked) {
      this.model.set('operationNone', !checked);
      this.model.set('operationAll', checked);
    },
  },
});
