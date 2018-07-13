import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    onSave({ saveType, model }) {
      if (saveType === 'delete') {
        this.send('reload');
      }
    },
  },
});
