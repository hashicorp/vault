/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { findDestination } from 'core/helpers/sync-destinations';
import { toLabel } from 'core/helpers/to-label';
import {
  DestinationType,
  CLOUD_DESTINATION_TYPES,
  ACCOUNT_CREDENTIAL_FIELDS,
  WIF_CREDENTIAL_FIELDS,
  type CloudDestinationType,
} from 'sync/utils/constants';
import type { Destination, DestinationConnectionDetails, DestinationOptions } from 'vault/sync';
import type { CapabilitiesMap } from 'vault/app-types';

interface Args {
  destination: Destination;
  capabilities: CapabilitiesMap;
}

function isCloudDestinationType(type: DestinationType): type is CloudDestinationType {
  return CLOUD_DESTINATION_TYPES.includes(type as CloudDestinationType);
}

export default class DestinationDetailsPage extends Component<Args> {
  private getCredentialType(destination: Destination): string | undefined {
    if (!isCloudDestinationType(destination.type)) {
      return undefined;
    }

    const isWIF = !!destination.connection_details?.identity_token_audience;

    if (isWIF) {
      return 'WIF';
    }

    const accountCredentialTypes: Record<CloudDestinationType, string> = {
      [DestinationType.AwsSm]: 'IAM',
      [DestinationType.AzureKv]: 'Client secret',
      [DestinationType.GcpSm]: 'JSON',
    };

    return accountCredentialTypes[destination.type];
  }
  get connectionDetailsMap() {
    const { destination } = this.args;
    const isWIF = !!destination.connection_details?.identity_token_audience;

    const baseMap = {
      [DestinationType.AwsSm]: [
        'region',
        'kms_key_id',
        'external_id',
        'credential_type',
        'role_arn',
        'access_key_id',
        'secret_access_key',
        'identity_token_ttl',
        'granularity',
        'secret_name_template',
        'custom_tags',
      ],
      [DestinationType.AzureKv]: [
        'key_vault_uri',
        'tenant_id',
        'cloud',
        'client_id',
        'credential_type',
        'client_secret',
        'identity_token_ttl',
      ],
      [DestinationType.GcpSm]: [
        'project_id',
        'kms_key_id',
        'credential_type',
        'credentials',
        'service_account_email',
        'identity_token_ttl',
      ],
      [DestinationType.Gh]: ['repository_owner', 'repository_name', 'access_token'],
      [DestinationType.VercelProject]: ['access_token', 'project_id', 'team_id', 'deployment_environments'],
    };

    if (isCloudDestinationType(destination.type)) {
      const fieldsToRemove = isWIF
        ? ACCOUNT_CREDENTIAL_FIELDS[destination.type]
        : WIF_CREDENTIAL_FIELDS[destination.type];
      baseMap[destination.type] = baseMap[destination.type].filter(
        (field) => !fieldsToRemove.includes(field)
      );
    }

    return baseMap;
  }

  get displayFields() {
    const { destination } = this.args;
    const type = destination.type as keyof typeof this.connectionDetailsMap;

    const availableFields = this.connectionDetailsMap[type] || [];
    const connectionDetails = availableFields.map((field) => `connection_details.${field}`);

    const fields = [
      'name',
      ...connectionDetails,
      'options.granularity_level',
      'options.secret_name_template',
    ];

    if (type === DestinationType.AwsSm || type === DestinationType.GcpSm) {
      fields.push('connection_details.replica_regions');
    }

    if (CLOUD_DESTINATION_TYPES.includes(type as CloudDestinationType)) {
      fields.push('options.custom_tags');
    }

    return fields;
  }

  getFieldValue = (field: string): unknown => {
    const { destination } = this.args;
    const fieldName = this.fieldName(field);

    if (fieldName === 'credential_type') {
      return this.getCredentialType(destination);
    }

    let value: unknown;
    if (field.startsWith('connection_details.')) {
      const connectionDetails = destination.connection_details;
      value = connectionDetails?.[fieldName as keyof DestinationConnectionDetails];
    } else if (field.startsWith('options.')) {
      const options = destination.options;
      value = options?.[fieldName as keyof DestinationOptions];
    } else {
      value = destination[fieldName as keyof Destination];
    }

    // google-managed encryption only stores selected regions (empty KMS key values), so render as a
    // simple comma separated list rather than the key/value row grouping used when KMS keys are also set
    if (fieldName === 'regional_kms_keys' && this.isGcpRegionsOnly) {
      return Object.keys(value as Record<string, string>).join(', ');
    }

    // google-managed encryption only stores selected regions (empty KMS key values), so render as a
    // simple comma separated list rather than the key/value row grouping used when KMS keys are also set
    if (fieldName === 'replica_regions' && this.isGcpRegionsOnly) {
      return Object.keys(value as Record<string, string>).join(', ');
    }

    return value;
  };

  // remove connection_details or options from the field name
  fieldName(field: string) {
    return field.replace(/(connection_details|options)\./, '');
  }

  fieldLabel = (field: string) => {
    const fieldName = this.fieldName(field);

    if (fieldName === 'replica_regions' && this.args.destination.type === DestinationType.GcpSm) {
      return this.isGcpRegionsOnly ? 'Replica regions' : 'Replica regions and KMS keys';
    }

    // some fields have a specific label that cannot be converted from key name
    const customLabel = {
      granularity_level: 'Secret sync granularity',
      access_key_id: 'Access key ID',
      role_arn: 'Role ARN',
      external_id: 'External ID',
      key_vault_uri: 'Key Vault URI',
      client_id: 'Client ID',
      tenant_id: 'Tenant ID',
      project_id: 'Project ID',
      credentials: 'JSON credentials',
      team_id: 'Team ID',
      identity_token_ttl: 'Identity token time to live',
      kms_key_id: 'KMS key ID',
      replica_regions: 'Replica regions and KMS keys',
    }[fieldName];

    return customLabel || toLabel([fieldName]);
  };

  isMasked = (field: string) => {
    const { maskedParams = [] } = findDestination(this.args.destination.type);
    return maskedParams.includes(this.fieldName(field));
  };

  // object values render as a labeled group of key/value rows instead of a single row
  isKeyValueField = (value: unknown): boolean => {
    return typeof value === 'object' && value !== null && !Array.isArray(value) && !(value instanceof Date);
  };

  // true when a replica_regions object only has region keys selected with no KMS key values populated
  private isRegionsOnly(value: unknown): boolean {
    if (!value || typeof value !== 'object') return false;
    const entries = Object.entries(value as Record<string, string>);
    return entries.length > 0 && entries.every(([, kmsKey]) => !kmsKey);
  }

  private get isGcpRegionsOnly(): boolean {
    const { destination } = this.args;
    return (
      destination.type === DestinationType.GcpSm &&
      this.isRegionsOnly(destination.connection_details?.replica_regions)
    );
  }

  credentialValue = (value: string) => {
    // if this value is empty, a destination uses globally set environment variables
    return value ? 'Destination credentials added' : 'Using environment variable';
  };
}
