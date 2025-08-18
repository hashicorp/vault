/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import MountForm from 'vault/forms/mount';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { AuthMethodFormData } from 'vault/auth/methods';

export default class AuthMethodForm extends MountForm<AuthMethodFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [this.fields.path]),
    new FormFieldGroup('Method Options', [
      this.fields.description,
      this.fields.listingVisibility,
      this.fields.local,
      this.fields.sealWrap,
      this.fields.defaultLeaseTtl,
      this.fields.maxLeaseTtl,
      new FormField('config.token_type', 'string', {
        label: 'Token type',
        helpText:
          'The type of token that should be generated via this role. For `default-service` and `default-batch` service and batch tokens will be issued respectively, unless the auth method explicitly requests a different type.',
        possibleValues: ['default-service', 'default-batch', 'batch', 'service'],
        noDefault: true,
      }),
      this.fields.auditNonHmacRequestKeys,
      this.fields.auditNonHmacResponseKeys,
      this.fields.passthroughRequestHeaders,
      this.fields.allowedResponseHeaders,
      this.fields.pluginVersion,
    ]),
  ];
}
