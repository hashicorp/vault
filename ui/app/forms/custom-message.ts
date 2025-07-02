/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Form from './form';
import FormField from 'vault/utils/forms/field';
import { isBefore, isAfter } from 'date-fns';
import { encodeString } from 'core/utils/b64';

import type { CreateCustomMessageRequest } from '@hashicorp/vault-client-typescript';
import { Validations } from 'vault/vault/app-types';

type CustomMessageFormData = Partial<CreateCustomMessageRequest>;

export default class CustomMessageForm extends Form<CustomMessageFormData> {
  formFields = [
    new FormField('authenticated', undefined, {
      label: 'Where should we display this message?',
      editType: 'radio',
      possibleValues: [
        {
          label: 'After the user logs in',
          subText: 'Display to users after they have successfully logged in to Vault.',
          value: true,
          id: 'authenticated',
        },
        {
          label: 'On the login page',
          subText: 'Display to users on the login page before they have authenticated.',
          value: false,
          id: 'unauthenticated',
        },
      ],
    }),

    new FormField('type', 'string', {
      label: 'Type',
      editType: 'radio',
      possibleValues: [
        {
          label: 'Alert message',
          subText:
            'A banner that appears on the top of every page to display brief but high-signal messages like an update or system alert.',
          value: 'banner',
        },
        {
          label: 'Modal',
          subText:
            'A pop-up window used to bring immediate attention for important notifications or actions.',
          value: 'modal',
        },
      ],
    }),

    new FormField('title', 'string'),

    new FormField('message', 'string', { editType: 'textarea' }),

    new FormField('link', 'string', {
      editType: 'kv',
      keyPlaceholder: 'Display text (e.g. Learn more)',
      valuePlaceholder: 'Link URL (e.g. https://www.hashicorp.com/)',
      label: 'Link (optional)',
      isSingleRow: true,
      allowWhiteSpace: true,
    }),

    new FormField('startTime', 'dateTimeLocal', {
      editType: 'dateTimeLocal',
      label: 'Message starts',
      subText: 'Defaults to 12:00 a.m. the following day (local timezone).',
    }),

    new FormField('endTime', 'dateTimeLocal', {
      editType: 'yield',
      label: 'Message expires',
    }),
  ];

  validations: Validations = {
    title: [{ type: 'presence', message: 'Title is required.' }],
    message: [{ type: 'presence', message: 'Message is required.' }],
    link: [
      {
        validator({ link }: CustomMessageFormData) {
          if (!link) return true;
          const [title] = Object.keys(link);
          const [href] = Object.values(link);
          return title || href ? !!(title && href) : true;
        },
        message: 'Link title and url are required.',
      },
    ],
    startTime: [
      {
        validator({ startTime, endTime }: CustomMessageFormData) {
          if (!startTime || !endTime) return true;
          return isBefore(new Date(startTime), new Date(endTime));
        },
        message: 'Start time is after end time.',
      },
    ],
    endTime: [
      {
        validator({ startTime, endTime }: CustomMessageFormData) {
          if (!startTime || !endTime) return true;
          return isAfter(new Date(endTime), new Date(startTime));
        },
        message: 'End time is before start time.',
      },
    ],
  };

  toJSON() {
    // overriding to do some date serialization
    // form sets dates as strings but client expects Date objects
    const startTime = this.data.startTime ? new Date(this.data.startTime as unknown as string) : undefined;
    const endTime = this.data.endTime ? new Date(this.data.endTime as unknown as string) : undefined;
    // encode message to base64
    const message = this.data.message ? encodeString(this.data.message) : undefined;
    return super.toJSON({ ...this.data, startTime, endTime, message });
  }
}
