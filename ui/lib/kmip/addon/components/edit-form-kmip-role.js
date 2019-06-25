import EditForm from 'core/components/edit-form';
import layout from '../templates/components/edit-form-kmip-role';
import { Promise } from 'rsvp';

export default EditForm.extend({
  layout,
  display: null,
  init() {
    this._super(...arguments);
    let display = 'operationAll';
    if (this.model.operationNone) {
      display = 'operationNone';
    }
    if (!this.model.isNew && !this.model.operationNone && !this.model.operationAll) {
      display = 'choose';
    }
    this.set('display', display);
  },

  actions: {
    updateModel(val) {
      if (val === 'operationAll') {
        this.model.set('operationNone', false);
        this.model.set('operationAll', true);
      }
      if (val === 'operationNone') {
        this.model.set('operationNone', true);
        this.model.set('operationAll', false);
      }
    },

    preSave(model) {
      let { display } = this;

      return new Promise(function(resolve) {
        if (display === 'choose') {
          model.set('operationNone', null);
          model.set('operationAll', null);
          return resolve(model);
        }
        model.newFields.without('role').forEach(field => {
          // this will set operationAll or operationNone to true
          if (field === display) {
            model.set(field, true);
          } else {
            model.set(field, null);
          }
        });
        return resolve(model);
      });
    },
  },
});
