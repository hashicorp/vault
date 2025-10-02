/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { WHITESPACE_WARNING } from 'vault/utils/forms/validators';

import type { Validations } from 'vault/app-types';

type KvFormData = {
  path: string;
  secretData: { [key: string]: string };
  custom_metadata?: { [key: string]: string };
  max_versions?: number;
  cas_required?: boolean;
  delete_version_after?: string;
  // readonly options when editing an existing secret
  options?: {
    cas: number;
  };
};

export default class KvForm extends Form<KvFormData> {
  fieldProps = ['secretFields', 'metadataFields'];

  validations: Validations = {
    path: [
      { type: 'presence', message: `Path can't be blank.` },
      { type: 'endsInSlash', message: `Path can't end in forward slash '/'.` },
      {
        type: 'containsWhiteSpace',
        message: WHITESPACE_WARNING('path'),
        level: 'warn',
      },
    ],
    secretData: [
      {
        validator: ({ secretData }: KvForm['data']) =>
          secretData !== undefined && typeof secretData !== 'object' ? false : true,
        message: 'Vault expects data to be formatted as an JSON object.',
      },
    ],
    max_versions: [
      { type: 'number', message: 'Maximum versions must be a number.' },
      { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
    ],
  };

  secretFields = [
    new FormField('path', 'string', {
      label: 'Path for this secret',
      subText: 'Names with forward slashes define hierarchical path structures.',
    }),
  ];

  metadataFields = [
    new FormField('custom_metadata', 'object', {
      editType: 'kv',
      isSectionHeader: true,
      subText:
        'An optional set of informational key-value pairs that will be stored with all secret versions.',
    }),
    new FormField('max_versions', 'number', {
      defaultValue: 0,
      label: 'Maximum number of versions',
      subText:
        'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted.',
    }),
    new FormField('cas_required', 'boolean', {
      defaultValue: false,
      label: 'Require Check and Set',
      subText: `Writes will only be allowed if the key's current version matches the version specified in the cas parameter.`,
    }),
    new FormField('delete_version_after', 'boolean', {
      defaultValue: '0s',
      editType: 'ttl',
      label: 'Automate secret deletion',
      helperTextDisabled: `A secret's version must be manually deleted.`,
      helperTextEnabled: 'Delete all new versions of this secret after:',
    }),
  ];
}
