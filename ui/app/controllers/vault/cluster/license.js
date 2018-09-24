import Ember from 'ember';

const { Controller } = Ember;
export default Controller.extend({
  actions: {
    saveModel({ text }) {
      this.get('model').save({ text }).then(() => {
        this.send('doRefresh');
      });
    },
  },
});
