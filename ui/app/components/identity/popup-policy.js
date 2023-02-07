import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import Base from './_popup-base';

export default Base.extend({
  model: alias('params.firstObject'),
  policyName: computed('params', function () {
    return this.params.objectAt(1);
  }),

  messageArgs(model, policyName) {
    return [model, policyName];
  },

  successMessage(model, policyName) {
    return `Successfully removed '${policyName}' policy from ${model.id} `;
  },

  errorMessage(e, model, policyName) {
    const error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem removing '${policyName}' policy - ${error}`;
  },

  transaction(model, policyName) {
    const policies = model.get('policies');
    model.set('policies', policies.without(policyName));
    return model.save();
  },
});
