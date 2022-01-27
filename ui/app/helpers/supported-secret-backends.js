import { helper as buildHelper } from '@ember/component/helper';

const SUPPORTED_SECRET_BACKENDS = [
  'aws',
  'database',
  'cubbyhole',
  'generic',
  'kv',
  'pki',
  'ssh',
  'transit',
  'kmip',
  'transform',
];

export function supportedSecretBackends() {
  return SUPPORTED_SECRET_BACKENDS;
}

export default buildHelper(supportedSecretBackends);
