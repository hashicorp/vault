/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { isAfter } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

const validations = {
  title: [{ type: 'presence', message: 'Title is required.' }],
  message: [{ type: 'presence', message: 'Message is required.' }],
};

@withFormFields(['authenticated', 'type', 'title', 'message', 'link', 'startTime', 'endTime'])
@withModelValidations(validations)
export default class MessageModel extends Model {
  @attr('boolean') active;
  @attr('string', {
    label: 'Type',
    editType: 'radio',
    subText: 'Display to users after they have successfully logged in to Vault.',
    possibleValues: ['Alert banner', 'Modal'],
  })
  type;
  @attr('boolean', {
    label: 'Where should we display this message?',
    editType: 'radio',
    possibleValues: [true, false],
  })
  authenticated;
  @attr('string', {
    label: 'Title',
    fieldValue: 'title',
    editDisabled: true,
  })
  title;
  @attr('string', {
    label: 'Message',
    fieldValue: 'message',
    editType: 'textarea',
    editDisabled: true,
  })
  message;
  @attr('object') link;
  @attr('string') startTime;
  @attr('string') endTime;

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
