import Ember from 'ember';

export function reduceToArray(params) {
  return params.reduce(function(result, param) {
    if (Ember.isNone(param)) {
      return result;
    }
    if (Ember.typeOf(param) === 'array') {
      return result.concat(param);
    } else {
      return result.concat([param]);
    }
  }, []);
}

export default Ember.Helper.helper(reduceToArray);
