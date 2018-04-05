import Ember from 'ember';

const TOOLS_ACTIONS = ['wrap', 'lookup', 'unwrap', 'rewrap', 'random', 'hash'];

export function toolsActions() {
  return TOOLS_ACTIONS;
}

export default Ember.Helper.helper(toolsActions);
