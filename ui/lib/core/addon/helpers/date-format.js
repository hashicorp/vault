/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { format, parseISO } from 'date-fns';

export function dateFormat([date, style], { isFormatted = false, withTimeZone = false }) {
  // see format breaking in upgrade to date-fns 2.x https://github.com/date-fns/date-fns/blob/master/CHANGELOG.md#changed-5
  if (isFormatted) {
    return format(new Date(date), style);
  }
  let number = typeof date === 'string' ? parseISO(date) : date;
  if (!number) {
    return;
  }
  if (number.toString().length === 10) {
    number = new Date(number * 1000);
  }
  let zone; // local timezone ex: 'PST'
  try {
    // passing undefined means default to the browser's locale
    zone = ' ' + number.toLocaleTimeString(undefined, { timeZoneName: 'short' }).split(' ')[2];
  } catch (e) {
    zone = '';
  }
  zone = withTimeZone ? zone : '';
  return format(number, style) + zone;
}

export default helper(dateFormat);
