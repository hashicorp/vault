import Base from './_popup-base';
import Ember from 'ember';
const { computed } = Ember;

export default Base.extend({
  model: computed.alias('params.firstObject'),
  key: computed('params', function() {
    return this.get('params').objectAt(1);
  }),

  messageArgs(model, key) {
    return [model, key];
  },

  successMessage(model, key) {
    return `Successfully removed '${key}' from metadata`;
  },
  errorMessage(e, model, key) {
    let error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem removing '${key}' from the metadata - ${error}`;
  },

  transaction(model, key) {
    let metadata = model.get('metadata');
    delete metadata[key];
    model.set('metadata', { ...metadata });
    return model.save();
  },
});
