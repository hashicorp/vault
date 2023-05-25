/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';
import { formatDuration, intervalToDuration } from 'date-fns';

export function duration([time]) {
  // time must be in seconds
  // 0 does not always mean 0 seconds, i.e. it can representing using system defaults
  if (Number.isInteger(time) && time !== 0) {
    const milliseconds = time * 1000;
    // pass milliseconds to intervalToDuration returns a durationObject: { years: 0, months: 0, days: 0, hours: 1, minutes: 0, seconds: 6 }
    // formatDuration converts to human-readable format: '1 hour 6 seconds'
    return formatDuration(intervalToDuration({ start: 0, end: milliseconds }));
  }
  // to avoid making any assumptions return strings and 0 as-is
  return time;
}

export default helper(duration);
