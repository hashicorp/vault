import { helper as buildHelper } from '@ember/component/helper';

const MANAGED_AUTH_BACKENDS = ['cert', 'userpass', 'ldap', 'okta', 'radius'];

export function supportedManagedAuthBackends() {
  return MANAGED_AUTH_BACKENDS;
}

export default buildHelper(supportedManagedAuthBackends);
