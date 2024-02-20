/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { isAfter, addDays, startOfDay, parseISO, isBefore } from 'date-fns';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

const validations = {
  title: [{ type: 'presence', message: 'Title is required.' }],
  message: [{ type: 'presence', message: 'Message is required.' }],
  link: [
    {
      validator(model) {
        if (!model?.link) return true;
        const [title] = Object.keys(model.link);
        const [href] = Object.values(model.link);
        return title || href ? !!(title && href) : true;
      },
      message: 'Link title and url are required.',
    },
  ],
  startTime: [
    {
      validator(model) {
        if (!model.endTime) return true;
        const start = new Date(model.startTime);
        const end = new Date(model.endTime);
        return isBefore(start, end);
      },
      message: 'Start time is after end time.',
    },
  ],
  endTime: [
    {
      validator(model) {
        if (!model.endTime) return true;
        const start = new Date(model.startTime);
        const end = new Date(model.endTime);
        return isAfter(end, start);
      },
      message: 'End time is before start time.',
    },
  ],
};

@withModelValidations(validations)
@withFormFields(['authenticated', 'type', 'title', 'message', 'link', 'startTime', 'endTime'])
export default class MessageModel extends Model {
  @attr('boolean') active;
  @attr('string', {
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
        subText: 'A pop-up window used to bring immediate attention for important notifications or actions.',
        value: 'modal',
      },
    ],
    defaultValue: 'banner',
  })
  type;
  // The authenticated attr is a boolean. The authenticatedString getter and setter is used only in forms to get and set the boolean via
  // strings values. The server and query params expects the attr to be boolean values.
  @attr({
    label: 'Where should we display this message?',
    editType: 'radio',
    fieldValue: 'authenticatedString',
    possibleValues: [
      {
        label: 'After the user logs in',
        subText: 'Display to users after they have successfully logged in to Vault.',
        value: 'authenticated',
      },
      {
        label: 'On the login page',
        subText: 'Display to users on the login page before they have authenticated.',
        value: 'unauthenticated',
      },
    ],
    defaultValue: true,
  })
  authenticated;

  get authenticatedString() {
    return this.authenticated ? 'authenticated' : 'unauthenticated';
  }

  set authenticatedString(value) {
    this.authenticated = value === 'authenticated' ? true : false;
  }

  @attr('string')
  title;
  @attr('string', {
    editType: 'textarea',
  })
  message;
  @attr('dateTimeLocal', {
    editType: 'dateTimeLocal',
    label: 'Message starts',
    subText: 'Defaults to 12:00 a.m. the following day (local timezone).',
    defaultValue: addDays(startOfDay(new Date()), 1).toISOString(),
  })
  startTime;
  @attr('dateTimeLocal', { editType: 'yield', label: 'Message expires' }) endTime;

  @attr('object', {
    editType: 'kv',
    keyPlaceholder: 'Display text (e.g. Learn more)',
    valuePlaceholder: 'Link URL (e.g. https://www.hashicorp.com/)',
    label: 'Link (optional)',
    isSingleRow: true,
    allowWhiteSpace: true,
  })
  link;

  // date helpers
  get isStartTimeAfterToday() {
    return isAfter(parseISO(this.startTime), new Date());
  }

  // capabilities
  @lazyCapabilities(apiPath`sys/config/ui/custom-messages`) customMessagesPath;

  get canCreateCustomMessages() {
    return this.customMessagesPath.get('canCreate') !== false;
  }
  get canReadCustomMessages() {
    return this.customMessagesPath.get('canRead') !== false;
  }
  get canEditCustomMessages() {
    return this.customMessagesPath.get('canUpdate') !== false;
  }
  get canDeleteCustomMessages() {
    return this.customMessagesPath.get('canDelete') !== false;
  }
}
