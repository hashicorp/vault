/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

// ⚠️ AUTO-GENERATED FILE - DO NOT EDIT
// This file is generated from openapi.json
// To customize this form, create an override in
// forms/v2/overrides/

import type ApiService from 'vault/services/api';
import type { FormConfig } from '../form-config';
import type { SystemApiMountsEnableSecretsEngineOperationRequest } from '@hashicorp/vault-client-typescript';

/**
 * Form configuration for mountsEnableSecretsEngine
 * Auto-generated from OpenAPI specification
 */
const mountsEnableSecretsEngineConfig: FormConfig<
  SystemApiMountsEnableSecretsEngineOperationRequest,
  unknown
> = {
  name: 'mountsEnableSecretsEngine',
  description: 'Mount a new backend at a new path.',
  submit: async (api: ApiService, payload: SystemApiMountsEnableSecretsEngineOperationRequest) => {
    return await api.sys.mountsEnableSecretsEngineRaw(payload);
  },
  payload: {
    path: '',
    MountsEnableSecretsEngineRequest: {
      config: {},
      description: '',
      external_entropy_access: false,
      local: false,
      options: {},
      plugin_name: '',
      plugin_version: '',
      seal_wrap: false,
      type: '',
    },
  },
  sections: [
    {
      name: 'params',
      fields: [
        {
          name: 'path',
          type: 'TextInput',
          label: 'Path',
          helperText: 'The path to mount to. Example: "aws/east"',
        },
      ],
    },
    {
      name: 'default',
      fields: [
        {
          name: 'MountsEnableSecretsEngineRequest.config',
          type: 'TextInput',
          label: 'Config',
          helperText: 'Configuration for this mount, such as default_lease_ttl and max_lease_ttl.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.description',
          type: 'TextInput',
          label: 'Description',
          helperText: 'User-friendly description for this mount.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.external_entropy_access',
          type: 'TextInput',
          label: 'External Entropy Access',
          helperText: "Whether to give the mount access to Vault's external entropy.",
        },
        {
          name: 'MountsEnableSecretsEngineRequest.local',
          type: 'TextInput',
          label: 'Local',
          helperText:
            'Mark the mount as a local mount, which is not replicated and is unaffected by replication.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.options',
          type: 'TextInput',
          label: 'Options',
          helperText:
            'The options to pass into the backend. Should be a json object with string keys and values.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.plugin_name',
          type: 'TextInput',
          label: 'Plugin Name',
          helperText: 'Name of the plugin to mount based from the name registered in the plugin catalog.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.plugin_version',
          type: 'TextInput',
          label: 'Plugin Version',
          helperText: 'The semantic version of the plugin to use, or image tag if oci_image is provided.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.seal_wrap',
          type: 'TextInput',
          label: 'Seal Wrap',
          helperText: 'Whether to turn on seal wrapping for the mount.',
        },
        {
          name: 'MountsEnableSecretsEngineRequest.type',
          type: 'TextInput',
          label: 'Type',
          helperText: 'The type of the backend. Example: "passthrough"',
        },
      ],
    },
  ],
};

export default mountsEnableSecretsEngineConfig;
