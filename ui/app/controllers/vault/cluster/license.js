import Controller from '@ember/controller';

export default Controller.extend({
  licenseSuccess() {
    this.send('doRefresh');
  },
  licenseError() {
    //eat the error (handled in MessageError component)
  },
  actions: {
    saveModel({ text }) {
      this.model.save({ text }).then(() => this.licenseSuccess(), () => this.licenseError());
    },
  },
});
