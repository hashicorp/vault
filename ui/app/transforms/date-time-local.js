/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Transform from '@ember-data/serializer/transform';
import { datetimeLocalStringFormat } from 'core/utils/date-formatters';
import { format } from 'date-fns';

export default class DateTimeLocalTransform extends Transform {
  getISODateFormat(deserializedDate) {
    if (!deserializedDate) return null;

    // if the date is a date object or in local date time format ("yyyy-MM-dd'T'HH:mm"), we want to ensure
    // it gets converted to an ISOString
    if (
      typeof deserializedDate === 'object' ||
      (typeof deserializedDate === 'string' && !deserializedDate.includes('Z'))
    ) {
      return new Date(deserializedDate).toISOString();
    }

    return deserializedDate;
  }

  deserialize(serialized) {
    if (!serialized) return null;
    return format(new Date(serialized), datetimeLocalStringFormat);
  }

  serialize(deserialized) {
    return this.getISODateFormat(deserialized);
  }
}
