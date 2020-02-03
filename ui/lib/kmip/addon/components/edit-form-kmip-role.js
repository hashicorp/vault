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
    toggleOperationSpecial(checked) {
      this.model.set('operationNone', !checked);
      this.model.set('operationAll', checked);
    },

    // when operationAll is true, we want all of the items
    // to appear checked, but we don't want to override what items
    // a user has selected - so this action creates an object that we
    // pass to the FormField component as the model instead of the real model
    placeholderOrModel(isOperationAll, attr) {
      return isOperationAll ? { [attr.name]: true } : this.model;
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
