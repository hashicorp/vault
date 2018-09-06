import Ember from 'ember';

export function pathOrArray([maybeArray, target]) {
  if (Array.isArray(maybeArray)) {
    return maybeArray;
  }
  return target.get(maybeArray);
}

export default Ember.Helper.helper(pathOrArray);
