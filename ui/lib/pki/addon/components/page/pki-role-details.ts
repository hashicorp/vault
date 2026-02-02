/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { toLabel } from 'core/helpers/to-label';

import type { PkiReadRoleResponse } from '@hashicorp/vault-client-typescript';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type VersionService from 'vault/services/version';
import type ApiService from 'vault/services/api';

interface Args {
  role: PkiReadRoleResponse & { name: string };
  capabilities: { canDelete: boolean; canEdit: boolean; canGenerateCert: boolean; canSign: boolean };
}

export default class DetailsPage extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly version: VersionService;
  @service declare readonly api: ApiService;

  label = (field: string) => {
    const label = toLabel([field]);
    return (
      {
        name: 'Role name',
        issuer_ref: 'Issuer',
        custom_ttl: 'Issued certificates expire after',
        not_before_duration: 'Issued certificate backdating',
        no_store: 'Store in storage backend',
        no_store_metadata: 'Store metadata in storage backend',
        max_ttl: 'Max TTL',
        generate_lease: 'Generate lease with certificate',
        basic_constraints_valid_for_non_ca: 'Add basic constraints',
        allowed_domains_template: 'Allow templates in allowed domains',
        ext_key_usage: 'Extended key usage',
        ext_key_usage_oids: 'Extended key usage OIDs',
        allow_ip_sans: 'Allow IP SANs',
        allowed_uri_sans: 'URI Subject Alternative Names (URI SANs)',
        allowed_uri_sans_template: 'Allow URI SANs template',
        allowed_other_sans: 'Other SANs',
        require_cn: 'Require common name',
        use_csr_common_name: 'Use CSR common name',
        use_csr_sans: 'Use CSR SANs',
        ou: 'Organizational units (OU)',
        locality: 'Locality/City',
        province: 'Province/State',
      }[field] || label
    );
  };

  isArrayField = (field: string) => ['key_usage', 'ext_key_usage', 'ext_key_usage_oids'].includes(field);

  defaultShown = (field: string) => {
    if (this.isArrayField(field)) {
      return 'None';
    } else if (field === 'max_ttl') {
      return 'System default';
    }
    return undefined;
  };

  get displayGroups() {
    const defaultArray = [
      'name',
      'issuer_ref',
      'custom_ttl',
      'not_before_duration',
      'max_ttl',
      'generate_lease',
      'no_store',
      'basic_constraints_valid_for_non_ca',
    ];
    if (this.version.isEnterprise) {
      // insert no_store_metadata after no_store for Enterprise versions
      defaultArray.splice(defaultArray.length - 1, 0, 'no_store_metadata');
    }
    return [
      { default: defaultArray },
      {
        'Domain handling': [
          'allowed_domains',
          'allowed_domains_template',
          'allow_bare_domains',
          'allow_subdomains',
          'allow_glob_domains',
          'allow_wildcard_certificates',
          'allow_localhost',
          'allow_any_name',
          'enforce_hostnames',
        ],
      },
      {
        'Key parameters': ['key_type', 'key_bits', 'signature_bits'],
      },
      {
        'Key usage': ['key_usage', 'ext_key_usage', 'ext_key_usage_oids'],
      },
      { 'Policy identifiers': ['policy_identifiers'] },
      {
        'Subject Alternative Name (SAN) Options': [
          'allow_ip_sans',
          'allowed_uri_sans',
          'allowed_uri_sans_template',
          'allowed_other_sans',
        ],
      },
      {
        'Additional subject fields': [
          'allowed_user_ids',
          'allowed_serial_numbers',
          'serial_number_source',
          'require_cn',
          'use_csr_common_name',
          'use_csr_sans',
          'ou',
          'organization',
          'country',
          'locality',
          'province',
          'street_address',
          'postal_code',
        ],
      },
    ];
  }

  @action
  async deleteRole() {
    try {
      await this.api.secrets.pkiDeleteRole(this.args.role.name, this.secretMountPath.currentPath);
      this.flashMessages.success('Role deleted successfully');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
