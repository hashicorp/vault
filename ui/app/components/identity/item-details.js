import { inject as service } from '@ember/service';
import Component from '@ember/component';

export default Component.extend({
  flashMessages: service(),

  actions: {
    enable(model) {
      model.set('disabled', false);

      model
        .save()
        .then(() => {
          this.get('flashMessages').success(`Successfully enabled entity: ${model.id}`);
        })
        .catch(e => {
          this.get('flashMessages').success(
            `There was a problem enabling the entity: ${model.id} - ${e.errors.join(' ') || e.message}`
          );
        });
    },
  },
});
