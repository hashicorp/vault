/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import AwsSmForm from './aws-sm';
import AzureKvForm from './azure-kv';
import GcpSmForm from './gcp-sm';
import GhForm from './gh';
import VercelProjectForm from './vercel-project';
import { DestinationType, CredentialType } from 'sync/utils/constants';

import type { DestinationConnectionDetails } from 'vault/sync';
import type { FormOptions } from '../form';
import type { Validations } from 'vault/app-types';

// given the differences in form fields across destination types, each type has a specific form class
// to make it easier in routes, this resolver will instantiate a new instance of the correct form class for a given type
export default function destinationFormResolver(type: DestinationType, data = {}, options?: FormOptions) {
  const validations: Validations = {
    name: [
      { type: 'presence', message: 'Name is required.' },
      { type: 'containsWhiteSpace', message: 'Name cannot contain whitespace.' },
    ],
    role_arn: [
      {
        validator({ role_arn, credential_type }: DestinationConnectionDetails) {
          if (type === DestinationType.AwsSm && credential_type === CredentialType.WIF) {
            return !!role_arn;
          }
          return true;
        },
        message: 'Role ARN is required.',
      },
    ],
    identity_token_audience: [
      {
        validator({ identity_token_audience, credential_type }: DestinationConnectionDetails) {
          if (credential_type === CredentialType.WIF) {
            return !!identity_token_audience;
          }
          return true;
        },
        message: 'Identity token audience is required.',
      },
    ],
    key_vault_uri: [
      {
        validator({ key_vault_uri }: DestinationConnectionDetails) {
          if (type === DestinationType.AzureKv) {
            return !!key_vault_uri;
          }
          return true;
        },
        message: 'Key Vault URI is required.',
      },
    ],
    tenant_id: [
      {
        validator({ tenant_id }: DestinationConnectionDetails) {
          if (type === DestinationType.AzureKv) {
            return !!tenant_id;
          }
          return true;
        },
        message: 'Tenant ID is required.',
      },
    ],
    client_id: [
      {
        validator({ client_id }: DestinationConnectionDetails) {
          if (type === DestinationType.AzureKv) {
            return !!client_id;
          }
          return true;
        },
        message: 'Client ID is required.',
      },
    ],
    project_id: [
      {
        validator({ project_id, credential_type }: DestinationConnectionDetails) {
          if (type === DestinationType.GcpSm && credential_type === CredentialType.WIF) {
            return !!project_id;
          }
          return true;
        },
        message: 'Project ID is required.',
      },
    ],
    service_account_email: [
      {
        validator({ service_account_email, credential_type }: DestinationConnectionDetails) {
          if (type === DestinationType.GcpSm && credential_type === CredentialType.WIF) {
            return !!service_account_email;
          }
          return true;
        },
        message: 'Service account email is required.',
      },
    ],
  };

  if (type === DestinationType.AwsSm) {
    return new AwsSmForm(data, options, validations);
  }
  if (type === DestinationType.AzureKv) {
    return new AzureKvForm(data, options, validations);
  }
  if (type === DestinationType.GcpSm) {
    return new GcpSmForm(data, options, validations);
  }
  if (type === DestinationType.Gh) {
    return new GhForm(data, options, validations);
  }
  if (type === DestinationType.VercelProject) {
    const teamId = (data as VercelProjectForm['data'])['team_id'];
    validations['team_id'] = [
      {
        validator: (formData: VercelProjectForm['data']) => {
          if (!options?.isNew && formData['team_id'] !== teamId) {
            return false;
          }
          return true;
        },
        message: 'Team ID should only be updated if the project was transferred to another account.',
        level: 'warn',
      },
    ];
    validations['deployment_environments'] = [
      { type: 'presence', message: 'At least one environment is required.' },
    ];
    return new VercelProjectForm(data, options, validations);
  }

  throw new Error(`Unknown destination type: ${type}`);
}
