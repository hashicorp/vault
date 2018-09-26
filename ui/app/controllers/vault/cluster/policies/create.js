import Controller from '@ember/controller';
import trimRight from 'vault/utils/trim-right';
import PolicyEditController from 'vault/mixins/policy-edit-controller';

export default Controller.extend(PolicyEditController, {
  showFileUpload: false,
  file: null,
  actions: {
    setPolicyFromFile(index, fileInfo) {
      let { value, fileName } = fileInfo;
      let model = this.get('model');
      model.set('policy', value);
      if (!model.get('name')) {
        let trimmedFileName = trimRight(fileName, ['.json', '.txt', '.hcl', '.policy']);
        model.set('name', trimmedFileName);
      }
      this.set('showFileUpload', false);
    },
  },
});
