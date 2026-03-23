/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { PkiWriteRoleRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { Validations } from 'vault/vault/app-types';

type PkiRoleFormData = PkiWriteRoleRequest & { name: string };

export default class PkiRoleForm extends OpenApiForm<PkiRoleFormData> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('PkiWriteRoleRequest', ...args);

    this.formFields.push(
      // add name field since it's not part of the OpenAPI spec
      new FormField('name', 'string', {
        label: 'Role Name',
        editDisabled: true,
      }),
      // add customTtl which is a convenience field that sets ttl and notAfter via one input <PkiNotValidAfterForm>
      new FormField('customTtl', undefined, {
        label: 'Not valid after',
        subText:
          'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
        editType: 'yield',
      })
    );
    // setup form field groups
    this.formFieldGroups = [];
    for (const group in this.fieldGroupKeys) {
      const fieldKeys = this.fieldGroupKeys[group as keyof typeof this.fieldGroupKeys];
      this.formFieldGroups.push(
        new FormFieldGroup(
          group,
          fieldKeys.map((key) => this.findAndTransform(key))
        )
      );
    }
  }

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
  };

  fieldGroupKeys = {
    default: [
      'name',
      'issuer_ref',
      'customTtl',
      'not_before_duration',
      'max_ttl',
      'generate_lease',
      'no_store',
      'no_store_metadata',
      'basic_constraints_valid_for_non_ca',
    ],
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
    'Key parameters': ['key_type', 'key_bits', 'signature_bits'],
    'Key usage': ['key_usage', 'ext_key_usage', 'ext_key_usage_oids'],
    'Policy identifiers': ['policy_identifiers'],
    'Subject Alternative Name (SAN) Options': [
      'allow_ip_sans',
      'allowed_uri_sans',
      'allowed_uri_sans_template',
      'allowed_other_sans',
    ],
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
  };

  fieldGroupsInfo = {
    'Domain handling': {
      footer: {
        text: 'These options can interact intricately with one another. For more information,',
        docText: 'learn more here.',
        docLink: '/vault/api-docs/secret/pki#allowed_domains',
      },
    },
    'Key parameters': {
      header: {
        text: `These are the parameters for generating or validating the certificate's key material.`,
      },
    },
    'Subject Alternative Name (SAN) Options': {
      header: {
        text: `Subject Alternative Names (SANs) are identities (domains, IP addresses, and URIs) Vault attaches to the requested certificates.`,
      },
    },
    'Additional subject fields': {
      header: {
        text: `Additional identity metadata Vault can attach to the requested certificates.`,
      },
    },
  };

  findAndTransform(key: string) {
    const field = this.formFields.find((field) => field.name === key) as FormField;
    if (key === 'key_usage' && !Array.isArray(this.data.key_usage)) {
      // default value for key_usage needs to be array
      // in the spec there is a default param that has the correct array but also a value in the x-vault-displayAttrs which is taking precedence
      // these should ideally align but perhaps we need to look at the default value over the value in x-vault-displayAttrs
      const keyUsage = (this.data.key_usage as unknown as string) || '';
      this.data.key_usage = keyUsage?.split(',');
    } else {
      const label = {
        not_before_duration: 'Backdate validity',
        no_store: 'Do not store certificates in storage backend',
        no_store_metadata: 'Do not store certificate metadata in storage backend',
      }[key];

      if (label) {
        field.options.label = label;
      }
      if (key === 'not_before_duration') {
        field.options.subText =
          'Also called the not_before_duration property. Allows certificates to be valid for a certain time period before now. This is useful to correct clock misalignment on various systems when setting up your CA.';
      }
    }

    return field;
  }
}
