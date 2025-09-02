/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { findDestination } from 'core/helpers/sync-destinations';
import { toLabel } from 'core/helpers/to-label';

import type { Destination } from 'vault/sync';
import type { CapabilitiesMap } from 'vault/app-types';

interface Args {
  destination: Destination;
  capabilities: CapabilitiesMap;
}

export default class DestinationDetailsPage extends Component<Args> {
  connectionDetailsMap = {
    'aws-sm': ['region', 'access_key_id', 'secret_access_key', 'role_arn', 'external_id'],
    'azure-kv': ['key_vault_uri', 'tenant_id', 'cloud', 'client_id', 'client_secret'],
    'gcp-sm': ['project_id', 'credentials'],
    gh: ['repository_owner', 'repository_name', 'access_token'],
    'vercel-project': ['access_token', 'project_id', 'team_id', 'deployment_environments'],
  };

  get displayFields() {
    const { destination } = this.args;
    const type = destination.type as keyof typeof this.connectionDetailsMap;
    const connectionDetails = this.connectionDetailsMap[type].map((field) => `connection_details.${field}`);
    const fields = [
      'name',
      ...connectionDetails,
      'options.granularity_level',
      'options.secret_name_template',
    ];

    if (!['gh', 'vercel-project'].includes(type)) {
      fields.push('options.custom_tags');
    }

    return fields;
  }

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
