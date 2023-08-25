import Component from '@glimmer/component';

import type {
  LdapStaticRoleCredentials,
  LdapDynamicRoleCredentials,
} from 'ldap/routes/roles/role/credentials';
import { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  credentials: LdapStaticRoleCredentials | LdapDynamicRoleCredentials;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapRoleCredentialsPageComponent extends Component<Args> {
  staticFields = [
    { label: 'Last Vault rotation', key: 'last_vault_rotation', formatDate: 'MMM d yyyy, h:mm:ss aaa' },
    { label: 'Password', key: 'password', hasBlock: 'masked' },
    { label: 'Username', key: 'username' },
    { label: 'Rotation period', key: 'rotation_period', formatTtl: true },
    { label: 'Time remaining', key: 'ttl', formatTtl: true },
  ];
  dynamicFields = [
    { label: 'Distinguished Name', key: 'distinguished_names' },
    { label: 'Username', key: 'username', hasBlock: 'masked' },
    { label: 'Password', key: 'password', hasBlock: 'masked' },
    { label: 'Lease ID', key: 'lease_id' },
    { label: 'Lease duration', key: 'lease_duration', formatTtl: true },
    { label: 'Lease renewable', key: 'renewable', hasBlock: 'check' },
  ];
  get fields() {
    return this.args.credentials.type === 'dynamic' ? this.dynamicFields : this.staticFields;
  }
}
