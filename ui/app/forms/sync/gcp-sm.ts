/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { commonFields, getPayload } from './shared';

import type { SystemWriteSyncDestinationsGcpSmNameRequest } from '@hashicorp/vault-client-typescript';

type GcpSmFormData = SystemWriteSyncDestinationsGcpSmNameRequest & {
  name: string;
};

export default class GcpSmForm extends Form<GcpSmFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      commonFields.name,
      new FormField('projectId', 'string', {
        label: 'Project ID',
        subText:
          'The target project to manage secrets in. If set, overrides the project derived from the service account JSON credentials or application default credentials.',
      }),
    ]),
    new FormFieldGroup('Credentials', [
      new FormField('credentials', 'string', {
        label: 'JSON credentials',
        subText:
          'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
        editType: 'file',
        docLink: '/vault/docs/secrets/gcp#authentication',
      }),
    ]),
    new FormFieldGroup('Advanced configuration', [
      commonFields.granularity,
      commonFields.secretNameTemplate,
      commonFields.customTags,
    ]),
  ];

  toJSON() {
    const formState = super.toJSON();
    const data = getPayload<GcpSmFormData>('gcp-sm', this.data, this.isNew);
    return { ...formState, data };
  }
}
