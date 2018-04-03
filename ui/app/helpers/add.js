import Ember from 'ember';

export function add(params) {
  return params.reduce((sum, param) => parseInt(param, 0) + sum, 0);
}

export default Ember.Helper.helper(add);
