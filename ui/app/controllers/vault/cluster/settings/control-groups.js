import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    onSave({ saveType }) {
      if (saveType === 'destroyRecord') {
        this.send('reload');
      }
    },
  },
});
