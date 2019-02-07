import Controller from '@ember/controller';
import { inject as service } from '@ember/service';

export default Controller.extend({
  wizard: service(),
  actions: {
    saveModel: function(model) {
      return model
        .save()
        .then(() => {
          let transition = this.transitionToRoute('vault.cluster.access.methods');
          return transition.followRedirects();
        })
        .catch(err => {
          debugger; // eslint-disable-line
          throw err;
        });
    },
  },
});
