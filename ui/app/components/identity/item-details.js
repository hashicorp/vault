import Ember from 'ember';

const { inject } = Ember;

export default Ember.Component.extend({
  flashMessages: inject.service(),

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
            `There was a problem enabling the entity: ${model.id} - ${e.error.join(' ') || e.message}`
          );
        });
    },
  },
});
