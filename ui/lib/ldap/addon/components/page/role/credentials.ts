/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type { LdapStaticRoleCredentials, LdapDynamicRoleCredentials } from 'vault/secrets/ldap';
import type { Breadcrumb } from 'vault/vault/app-types';
import { ResponseError } from '@hashicorp/vault-client-typescript';

interface Args {
  credentials: LdapStaticRoleCredentials | LdapDynamicRoleCredentials;
  error?: ResponseError;
  breadcrumbs: Array<Breadcrumb>;
}

// Vault returns Go's zero time (0001-01-01T00:00:00Z) for a static role that has never been rotated.
const isZeroTime = (timestamp?: string | null) => !!timestamp && new Date(timestamp).getUTCFullYear() <= 1;

export default class LdapRoleCredentialsPageComponent extends Component<Args> {
  staticFields = [
    {
      label: 'Last Vault rotation',
      key: 'last_vault_rotation',
      formatDate: 'MMM d yyyy, h:mm:ss aaa',
      alwaysRender: true,
      defaultShown: 'Never',
    },
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

  get displayCredentials() {
    const { credentials } = this.args;
    if (credentials && 'last_vault_rotation' in credentials && isZeroTime(credentials.last_vault_rotation)) {
      return { ...credentials, last_vault_rotation: null };
    }
    return credentials;
  }
}
