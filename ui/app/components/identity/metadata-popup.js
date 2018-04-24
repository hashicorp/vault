import Ember from 'ember';

export default Ember.Component.extend({
  tagName: '',
  flashMessages: Ember.inject.service(),
  model: null,
  key: null,

  actions: {
    removeKey(model, key) {
      let metadata = model.get('metadata');
      delete metadata[key];
      model.set('metadata', null);
      model.set('metadata', metadata);

      return model
        .save()
        .then(() => {
          this.get('flashMessages').success(`Successfully removed '${key}' from metadata`);
        })
        .catch(e => {
          model.rollbackAttributes();
          this.get('flashMessages').success(
            `There was a problem removing '${key}' from the metadata - ${e.error.join(' ') || e.message}`
          );
        });
    },
  },
});
