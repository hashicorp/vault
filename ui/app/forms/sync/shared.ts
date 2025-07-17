/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import { findDestination } from 'core/helpers/sync-destinations';

import type { DestinationType } from 'vault/sync';

export const commonFields = {
  name: new FormField('name', 'string', {
    subText: 'Specifies the name for this destination.',
    editDisabled: true,
  }),

  secretNameTemplate: new FormField('secretNameTemplate', 'string', {
    subText:
      'Go-template string that indicates how to format the secret name at the destination. The default template varies by destination type but is generally in the form of "vault-{{ .MountAccessor }}-{{ .SecretPath }}" e.g. "vault-kv_9a8f68ad-my-secret-1". Optional.',
  }),

  granularity: new FormField('granularity', 'string', {
    editType: 'radio',
    label: 'Secret sync granularity',
    possibleValues: [
      {
        label: 'Secret path',
        subText: 'Sync entire secret contents as a single entry at the destination.',
        value: 'secret-path',
      },
      {
        label: 'Secret key',
        subText: 'Sync each key-value pair of secret data as a distinct entry at the destination.',
        helpText:
          'Only top-level keys will be synced and any nested or complex values will be encoded as a JSON string.',
        value: 'secret-key',
      },
    ],
  }),

  customTags: new FormField('customTags', 'object', {
    subText:
      'An optional set of informational key-value pairs added as additional metadata on secrets synced to this destination. Custom tags are merged with built-in tags.',
    editType: 'kv',
  }),
};

export function getPayload<T>(type: DestinationType, data: T, isNew: boolean) {
  const { maskedParams, readonlyParams } = findDestination(type);
  const payload: T = { ...data };

  // the server returns ****** for sensitive fields
  // these are represented as maskedParams in the sync-destinations helper
  // when editing, remove these fields from the payload if they haven't been changed
  if (!isNew) {
    maskedParams.forEach((maskedParam) => {
      const key = maskedParam as keyof T;
      const value = (payload[key] as string) || '';
      // if the value is asterisks, remove it from the payload
      if (value.match(/^\*+$/)) {
        delete payload[key];
      }
    });

    // to preserve the original Ember Data payload structure, remove fields that are not editable
    // since editing is disabled in the form the value will not change so this is mostly to satisfy existing test conditions
    readonlyParams.forEach((readonlyParam) => {
      delete payload[readonlyParam as keyof T];
    });
  }

  return payload;
}
