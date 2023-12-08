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

@withFormFields(['authenticated', 'type', 'title', 'message', 'linkTitle', 'startTime', 'endTime'])
@withModelValidations(validations)
export default class MessageModel extends Model {
  @attr('boolean') active;
  @attr('string', {
    defaultValue: 'banner',
  })
  type;
  @attr('boolean', {
    defaultValue: true,
  })
  authenticated;
  @attr('string', {
    label: 'Title',
    fieldValue: 'title',
  })
  title;
  @attr('string', {
    label: 'Message',
    fieldValue: 'message',
    editType: 'textarea',
  })
  message;
  @attr('string') startTime;
  @attr('string', { defaultValue: '' }) endTime;

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
