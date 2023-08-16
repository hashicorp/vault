import Controller from '@ember/controller';

export default Controller.extend({
  showLicenseError: false,

  actions: {
    transitionToCluster() {
      return this.model.reload().then(() => {
        return this.transitionToRoute('vault.cluster', this.model.name);
      });
    },

    isUnsealed(data) {
      return data.sealed === false;
    },

    handleLicenseError() {
      this.set('showLicenseError', true);
    },
  },
});
