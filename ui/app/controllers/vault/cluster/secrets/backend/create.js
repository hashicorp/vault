import Ember from 'ember';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default Ember.Controller.extend(BackendCrumbMixin, {
  queryParams: ['initialKey'],

  initialKey: '',

  actions: {
    refresh: function() {
      this.send('refreshModel');
    },
    hasChanges(hasChanges) {
      this.send('hasDataChanges', hasChanges);
    },
  },
});
