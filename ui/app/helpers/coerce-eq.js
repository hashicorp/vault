/*jshint eqeqeq: false */
import Ember from 'ember';

export function coerceEq(params) {
  return params[0] == params[1];
}

export default Ember.Helper.helper(coerceEq);
