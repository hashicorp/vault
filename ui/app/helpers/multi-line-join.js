import Ember from 'ember';

export function multiLineJoin([arr]) {
  return arr.join('\n');
}

export default Ember.Helper.helper(multiLineJoin);
