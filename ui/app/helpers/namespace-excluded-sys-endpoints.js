import Ember from 'ember';

// these endpoints will always be at the root
// namespace of the sys mount
const ENDPOINTS = ['/sys/health', '/sys/seal-status', '/sys/license/features'];

export function namespaceExcludedSysEndpoint() {
  return ENDPOINTS;
}

export default Ember.Helper.helper(namespaceExcludedSysEndpoint);
