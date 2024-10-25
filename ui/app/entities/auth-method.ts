import { tracked } from '@glimmer/tracking';
import { withModelValidations } from 'vault/decorators/model-validations';
import { WHITESPACE_WARNING } from 'vault/utils/model-helpers/validators';
import { set } from '@ember/object';

import type { components } from 'vault/vault-openapi-schema';

export type AuthMethodData = components['schemas']['AuthReadConfigurationResponse'];
export type AuthEnableMethodRequest = components['schemas']['AuthEnableMethodRequest'];

const validations = {
  path: [
    { type: 'presence', message: "Path can't be blank." },
    {
      type: 'containsWhiteSpace',
      message: WHITESPACE_WARNING('path'),
      level: 'warn',
    },
  ],
};

// eslint-disable-next-line
// @ts-ignore
@withModelValidations(validations)
export default class AuthMethod {
  @tracked type = '';
  @tracked path = '';

  data: AuthMethodData = {
    config: {},
  };

  constructor(data?: AuthMethodData) {
    if (data) {
      this.data = data;
    }
  }

  // shim this for now but get away from old Ember patterns!
  set(key: string, val: unknown) {
    set(this, key, val);
  }

  fieldGroups = [
    {
      default: [{ name: 'path', type: 'string' }],
    },
    {
      'Method Options': [
        {
          name: 'data.description',
          options: {
            label: 'Description',
            editType: 'textarea',
          },
        },
        {
          name: 'data.config.listingVisibility',
          type: 'mountVisibility',
          options: {
            editType: 'boolean',
            label: 'List method when unauthenticated',
            defaultValue: false,
          },
        },
        {
          name: 'data.local',
          type: 'boolean',
          options: {
            label: 'Local',
            helpText:
              'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
          },
        },
        {
          name: 'data.sealWrap',
          type: 'boolean',
          options: {
            label: 'Seal Wrap',
            helpText:
              'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
          },
        },
        {
          name: 'data.config.defaultLeaseTtl',
          options: {
            label: 'Default Lease TTL',
            editType: 'ttl',
          },
        },
        {
          name: 'data.config.maxLeaseTtl',
          options: {
            label: 'Max Lease TTL',
            editType: 'ttl',
          },
        },
        {
          name: 'data.config.tokenType',
          options: {
            label: 'Token type',
            helpText:
              'The type of token that should be generated via this role. For `default-service` and `default-batch` service and batch tokens will be issued respectively, unless the auth method explicitly requests a different type.',
            possibleValues: ['default-service', 'default-batch', 'batch', 'service'],
            noDefault: true,
          },
        },
        {
          name: 'data.config.auditNonHmacRequestKeys',
          options: {
            label: 'Request keys excluded from HMACing in audit',
            editType: 'stringArray',
            helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
          },
        },
        {
          name: 'data.config.auditNonHmacResponseKeys',
          options: {
            label: 'Response keys excluded from HMACing in audit',
            editType: 'stringArray',
            helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
          },
        },
        {
          name: 'data.config.passthroughRequestHeaders',
          options: {
            label: 'Allowed passthrough request headers',
            helpText: 'Headers to allow and pass from the request to the backend',
            editType: 'stringArray',
          },
        },
        {
          name: 'data.config.allowedResponseHeaders',
          options: {
            label: 'Allowed response headers',
            helpText: 'Headers to allow, allowing a plugin to include them in the response.',
            editType: 'stringArray',
          },
        },
        {
          name: 'data.config.pluginVersion',
          type: 'string',
          options: {
            label: 'Plugin version',
            subText:
              'Specifies the semantic version of the plugin to use, e.g. "v1.0.0". If unspecified, the server will select any matching un-versioned plugin that may have been registered, the latest versioned plugin registered, or a built-in plugin in that order of precedence.',
          },
        },
      ],
    },
  ];

  // quick workaround for decorator not being ts
  validate() {
    // eslint-disable-next-line
    // @ts-ignore
    return super.validate();
  }
}
