import Ember from 'ember';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default Ember.Controller.extend(BackendCrumbMixin, {
  queryParams: {
    selectedAction: 'action',
  },

  actions: {
    refresh: function() {
      // closure actions don't bubble to routes,
      // so we have to manually bubble here
      this.send('refreshModel');
    },

    hasChanges(hasChanges) {
      this.send('hasDataChanges', hasChanges);
    },
  },
});
