import Ember from 'ember';
import PolicyEditController from 'vault/mixins/policy-edit-controller';

export default Ember.Controller.extend(PolicyEditController, {
  showFileUpload: false,
  file: null,

  actions: {
    setPolicyFromFile(index, fileInfo) {
      let { value, fileName } = fileInfo;
      let model = this.get('model');
      model.set('policy', value);
      if (!model.get('name')) {
        model.set('name', fileName);
      }
      this.set('showFileUpload', false);
    },
  },
});
