/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AwsSmForm from './aws-sm';
import AzureKvForm from './azure-kv';
import GcpSmForm from './gcp-sm';
import GhForm from './gh';
import VercelProjectForm from './vercel-project';

import type { DestinationType } from 'vault/sync';
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
  };

  if (type === 'aws-sm') {
    return new AwsSmForm(data, options, validations);
  }
  if (type === 'azure-kv') {
    return new AzureKvForm(data, options, validations);
  }
  if (type === 'gcp-sm') {
    return new GcpSmForm(data, options, validations);
  }
  if (type === 'gh') {
    return new GhForm(data, options, validations);
  }
  if (type === 'vercel-project') {
    const teamId = (data as VercelProjectForm['data'])['teamId'];
    validations['teamId'] = [
      {
        validator: (formData: VercelProjectForm['data']) =>
          !options?.isNew && formData['teamId'] !== teamId ? false : true,
        message: 'Team ID should only be updated if the project was transferred to another account.',
        level: 'warn',
      },
    ];
    validations['deploymentEnvironments'] = [
      { type: 'presence', message: 'At least one environment is required.' },
    ];
    return new VercelProjectForm(data, options, validations);
  }

  throw new Error(`Unknown destination type: ${type}`);
}
