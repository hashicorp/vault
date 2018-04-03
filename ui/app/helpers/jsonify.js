import Ember from 'ember';

export function jsonify([target]) {
  return JSON.parse(target);
}

export default Ember.Helper.helper(jsonify);
