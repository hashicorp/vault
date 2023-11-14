/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { isAfter } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';

export default class MessageModel extends Model {
  @attr('boolean') active;
  @attr('string') type;
  @attr('boolean') authenticated;
  @attr('string') title;
  @attr('string') message;
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
