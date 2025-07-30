/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { tracked } from '@glimmer/tracking';
import { WHITESPACE_WARNING } from 'vault/utils/forms/validators';

import type { Validations } from 'vault/app-types';
import type { SecretsEngineFormData } from 'vault/secrets/engine';
import type { AuthMethodFormData } from 'vault/vault/auth/methods';

// common fields and validations shared between secrets engine and auth methods (mounts)
// used in form classes for consistency and to avoid duplication
export default class MountForm<T extends SecretsEngineFormData | AuthMethodFormData> extends Form<T> {
  @tracked declare type: string;

  validations: Validations = {
    path: [
      { type: 'presence', message: "Path can't be blank." },
      {
        type: 'containsWhiteSpace',
        message: WHITESPACE_WARNING('path'),
        level: 'warn',
      },
    ],
  };

  fields = {
    path: new FormField('path', 'string'),
    description: new FormField('description', 'string', { editType: 'textarea' }),
    listingVisibility: new FormField('config.listing_visibility', 'boolean', {
      label: 'Use as preferred UI login method',
      editType: 'toggleButton',
      helperTextEnabled:
        'This mount will be included in the unauthenticated UI login endpoint and display as a preferred login method.',
      helperTextDisabled:
        'Turn on the toggle to use this auth mount as a preferred login method during UI login.',
    }),
    local: new FormField('local', 'boolean', {
      helpText:
        'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
    }),
    sealWrap: new FormField('seal_wrap', 'boolean', {
      helpText:
        'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
    }),
    defaultLeaseTtl: new FormField('config.default_lease_ttl', 'string', {
      label: 'Default Lease TTL',
      editType: 'ttl',
    }),
    maxLeaseTtl: new FormField('config.max_lease_ttl', 'string', { label: 'Max Lease TTL', editType: 'ttl' }),
    auditNonHmacRequestKeys: new FormField('config.audit_non_hmac_request_keys', 'string', {
      label: 'Request keys excluded from HMACing in audit',
      editType: 'stringArray',
      helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
    }),
    auditNonHmacResponseKeys: new FormField('config.audit_non_hmac_response_keys', 'string', {
      label: 'Response keys excluded from HMACing in audit',
      editType: 'stringArray',
      helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
    }),
    passthroughRequestHeaders: new FormField('config.passthrough_request_headers', 'string', {
      label: 'Allowed passthrough request headers',
      helpText: 'Headers to allow and pass from the request to the backend',
      editType: 'stringArray',
    }),
    allowedResponseHeaders: new FormField('config.allowed_response_headers', 'string', {
      label: 'Allowed response headers',
      helpText: 'Headers to allow, allowing a plugin to include them in the response.',
      editType: 'stringArray',
    }),
    pluginVersion: new FormField('plugin_version', 'string', {
      label: 'Plugin version',
      subText:
        'Specifies the semantic version of the plugin to use, e.g. "v1.0.0". If unspecified, the server will select any matching un-versioned plugin that may have been registered, the latest versioned plugin registered, or a built-in plugin in that order of precedence.',
    }),
  };

  // namespaces introduced types with a `ns_` prefix for built-in engines so we will strip that out for consistency
  get normalizedType() {
    return (this.type || '').replace(/^ns_/, '');
  }

  toJSON() {
    const { config } = this.data;
    const data = {
      type: this.type,
      ...this.data,
      config: {
        ...(config || {}),
        force_no_cache: config?.force_no_cache ?? false,
        listing_visibility: config?.listing_visibility ? 'unauth' : 'hidden',
      },
    };
    // options are only relevant for kv/generic engines
    if (!['kv', 'generic'].includes(this.type)) {
      delete data.options;
    }

    return super.toJSON(data);
  }
}
