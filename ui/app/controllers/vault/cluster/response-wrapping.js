import Ember from 'ember';

export default Ember.Controller.extend({
  queryParams: {
    selectedAction: 'action',
  },

  selectedAction: 'wrap',
});
