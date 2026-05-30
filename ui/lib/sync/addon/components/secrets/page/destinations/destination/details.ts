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
        'external_id',
        'credential_type',
        'role_arn',
        'access_key_id',
        'secret_access_key',
        'identity_token_ttl',
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

    if (field.startsWith('connection_details.')) {
      const connectionDetails = destination.connection_details;
      return connectionDetails?.[fieldName as keyof DestinationConnectionDetails];
    }

    if (field.startsWith('options.')) {
      const options = destination.options;
      return options?.[fieldName as keyof DestinationOptions];
    }

    return destination[fieldName as keyof Destination];
  };

  // remove connection_details or options from the field name
  fieldName(field: string) {
    return field.replace(/(connection_details|options)\./, '');
  }

  fieldLabel = (field: string) => {
    const fieldName = this.fieldName(field);
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
    }[fieldName];

    return customLabel || toLabel([fieldName]);
  };

  isMasked = (field: string) => {
    const { maskedParams = [] } = findDestination(this.args.destination.type);
    return maskedParams.includes(this.fieldName(field));
  };

  credentialValue = (value: string) => {
    // if this value is empty, a destination uses globally set environment variables
    return value ? 'Destination credentials added' : 'Using environment variable';
  };
}
