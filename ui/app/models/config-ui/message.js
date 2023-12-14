/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { isAfter, format, addDays, startOfDay } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

export const localDateTimeString = "yyyy-MM-dd'T'HH:mm";

const validations = {
  title: [{ type: 'presence', message: 'Title is required.' }],
  message: [{ type: 'presence', message: 'Message is required.' }],
};

@withFormFields(['authenticated', 'type', 'title', 'message', 'linkTitle', 'startTime', 'endTime'])
@withModelValidations(validations)
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

  @attr('string', {
    label: 'Title',
  })
  title;
  @attr('string', {
    label: 'Message',
    editType: 'textarea',
  })
  message;
  @attr('date', {
    editType: 'dateTimeLocal',
    label: 'Message starts',
    subText: 'Defaults to 12:00 a.m. the following day (local timezone).',
    defaultValue: format(addDays(startOfDay(new Date() || this.startTime), 1), localDateTimeString),
  })
  startTime;
  @attr('date', { editType: 'yield', label: 'Message expires' }) endTime;

  // the api returns link as an object with title and href as keys, but we separate the link key/values into
  // different attributes to easily show link title and href fields on the create form. In our serializer,
  // we send the link attribute in to the correct format (as an object) to the server.
  @attr('string', { fieldValue: 'linkTitle' }) linkTitle;
  @attr('string', { fieldValue: 'linkHref' }) linkHref;

  // date helpers
  get isStartTimeAfterToday() {
    return isAfter(parseAPITimestamp(this.startTime), new Date());
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
