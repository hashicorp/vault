import EditForm from 'core/components/edit-form';
import layout from '../templates/components/edit-form-kmip-role';

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

    preSave(model) {
      // if we have operationAll or operationNone, we want to clear
      // out the others so that display shows the right data
      if (model.operationAll || model.operationNone) {
        model.operationFieldsWithoutSpecial.forEach(field => model.set(field, null));
      }
    },
  },
});
