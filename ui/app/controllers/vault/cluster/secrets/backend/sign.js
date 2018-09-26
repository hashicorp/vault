import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { set, get } from '@ember/object';

export default Controller.extend({
  store: service(),
  loading: false,
  emptyData: '{\n}',
  actions: {
    sign() {
      this.set('loading', true);
      this.model.save().finally(() => {
        this.set('loading', false);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.get('model'), attr, JSON.parse(val));
      }
    },

    newModel() {
      const model = this.get('model');
      const roleModel = model.get('role');
      model.unloadRecord();
      const newModel = this.get('store').createRecord('ssh-sign', {
        role: roleModel,
        id: `${get(roleModel, 'backend')}-${get(roleModel, 'name')}`,
      });
      this.set('model', newModel);
    },
  },
});
