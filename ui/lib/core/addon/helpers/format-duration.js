/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';
import { formatDuration, intervalToDuration } from 'date-fns';

export function duration([time], { nullable = false }) {
  // intervalToDuration creates a durationObject that turns the seconds (ex 3600) to respective:
  // { years: 0, months: 0, days: 0, hours: 1, minutes: 0, seconds: 0 }
  // then formatDuration returns the filled in keys of the durationObject
  // nullable if you don't want a value to be returned instead of 0s

  if (nullable && (time === '0' || time === 0)) {
    return null;
  }

  // time must be in seconds
  const duration = Number.parseInt(time, 10);
  if (isNaN(duration)) {
    return time;
  }

  return formatDuration(intervalToDuration({ start: 0, end: duration * 1000 }));
}

export default helper(duration);
