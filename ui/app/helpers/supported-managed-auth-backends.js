import { helper as buildHelper } from '@ember/component/helper';

const MANAGED_AUTH_BACKENDS = ['okta', 'radius', 'ldap', 'cert', 'userpass'];

export function supportedManagedAuthBackends() {
  return MANAGED_AUTH_BACKENDS;
}

export default buildHelper(supportedManagedAuthBackends);
