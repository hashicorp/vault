import Controller, { inject as controller } from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default Controller.extend(BackendCrumbMixin, {
  backendController: controller('vault.cluster.secrets.backend'),
  queryParams: ['initialKey', 'itemType'],

  initialKey: '',
  itemType: '',

  actions: {
    refresh: function () {
      this.send('refreshModel');
    },
    toggleAdvancedEdit(bool) {
      this.set('preferAdvancedEdit', bool);
      this.backendController.set('preferAdvancedEdit', bool);
    },
  },
});
