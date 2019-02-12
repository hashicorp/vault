import Base from './_popup-base';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';

export default Base.extend({
  model: alias('params.firstObject'),
  key: computed('params', function() {
    return this.params.objectAt(1);
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
    let metadata = model.metadata;
    delete metadata[key];
    model.set('metadata', { ...metadata });
    return model.save();
  },
});
