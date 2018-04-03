import Ember from 'ember';

export function includes([haystack, needle]) {
  return haystack.includes(needle);
}

export default Ember.Helper.helper(includes);
