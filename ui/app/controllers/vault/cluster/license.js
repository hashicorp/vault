import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    saveModel({ text }) {
      this.get('model')
        .save({ text })
        .then(() => {
          this.send('doRefresh');
        });
    },
  },
});
