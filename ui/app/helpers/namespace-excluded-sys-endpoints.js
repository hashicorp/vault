import Ember from 'ember';

// these endpoints will always be at the root
// namespace of the sys mount
const ENDPOINTS = [
  '/sys/audit',
  '/sys/audit-hash',
  '/sys/config/auditing',
  '/sys/config/cors',
  '/sys/config/ui',
  '/sys/generate-root',
  '/sys/health',
  '/sys/init',
  '/sys/key-status',
  '/sys/leader',
  '/sys/license',
  '/sys/raw',
  '/sys/rekey',
  '/sys/rekey-recovery-key',
  '/sys/replication',
  '/sys/rotate',
  '/sys/seal',
  '/sys/seal-status',
  '/sys/step-down',
];

export function namespaceExcludedSysEndpoint() {
  return ENDPOINTS;
}

export default Ember.Helper.helper(namespaceExcludedSysEndpoint);
