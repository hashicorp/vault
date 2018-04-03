import Ember from 'ember';

export function stringify([target], { skipFormat }) {
  if (skipFormat) {
    return JSON.stringify(target);
  }
  return JSON.stringify(target, null, 2);
}

export default Ember.Helper.helper(stringify);
