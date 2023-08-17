/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { durationToSeconds } from 'core/utils/duration-utils';
import { formatDuration, intervalToDuration } from 'date-fns';

export function duration([time]) {
  // 0 (integer) does not necessarily mean 0 seconds, i.e. it can represent using system ttl defaults
  if (time === 0) return time;

  // assume numbers are seconds or parses duration strings into seconds
  const seconds = durationToSeconds(time);

  if (Number.isInteger(seconds)) {
    // durationObject: { years: 0, months: 0, days: 0, hours: 1, minutes: 0, seconds: 6 }
    const durationObject = intervalToDuration({ start: 0, end: seconds * 1000 });

    if (Object.values(durationObject).every((v) => v === 0)) {
      // formatDuration returns an empty string if every value is 0
      return '0 seconds';
    }
    // convert to human-readable format: '1 hour 6 seconds'
    return formatDuration(durationObject);
  }
  return time;
}

export default helper(duration);
