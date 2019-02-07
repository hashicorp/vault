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
          throw err;
        });
    },
    tuneModel: function(model) {
      let data = model.config.serialize();
      data.description = model.description;
      return model
        .tune(data)
        .then(() => {
          let transition = this.transitionToRoute('vault.cluster.access.methods');
          return transition.followRedirects();
        })
        .catch(err => {
          throw err;
        });
    },
  },
});
