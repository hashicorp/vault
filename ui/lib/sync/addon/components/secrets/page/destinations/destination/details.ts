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
    'aws-sm': ['region', 'accessKeyId', 'secretAccessKey', 'roleArn', 'externalId'],
    'azure-kv': ['keyVaultUri', 'tenantId', 'cloud', 'clientId', 'clientSecret'],
    'gcp-sm': ['projectId', 'credentials'],
    gh: ['repositoryOwner', 'repositoryName', 'accessToken'],
    'vercel-project': ['accessToken', 'projectId', 'teamId', 'deploymentEnvironments'],
  };

  get displayFields() {
    const { destination } = this.args;
    const type = destination.type as keyof typeof this.connectionDetailsMap;
    const connectionDetails = this.connectionDetailsMap[type].map((field) => `connectionDetails.${field}`);
    const fields = ['name', ...connectionDetails, 'options.granularityLevel', 'options.secretNameTemplate'];

    if (!['gh', 'vercel-project'].includes(type)) {
      fields.push('options.customTags');
    }

    return fields;
  }

  // remove connectionDetails or options from the field name
  fieldName(field: string) {
    return field.replace(/(connectionDetails|options)\./, '');
  }

  fieldLabel = (field: string) => {
    const fieldName = this.fieldName(field);
    // some fields have a specific label that cannot be converted from key name
    const customLabel = {
      granularityLevel: 'Secret sync granularity',
      accessKeyId: 'Access key ID',
      roleArn: 'Role ARN',
      externalId: 'External ID',
      keyVaultUri: 'Key Vault URI',
      clientId: 'Client ID',
      tenantId: 'Tenant ID',
      projectId: 'Project ID',
      credentials: 'JSON credentials',
      teamId: 'Team ID',
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
