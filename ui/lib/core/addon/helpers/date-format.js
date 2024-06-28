/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { format, parseISO } from 'date-fns';

function checkType(value) {
  if (typeof value === 'string') {
    // if it's a number when multiplied by 1, it's just a number in quotes
    return isNaN(value * 1) ? 'string' : 'number';
  }
  if (typeof value === 'object') {
    // Dates are technically an object
    try {
      value.toUTCString();
      return 'date';
    } catch (e) {
      return 'object';
    }
  }
  return typeof value;
}

function dateFromNumber(number) {
  if (number.toString().length === 10) {
    // is seconds, convert to millis
    return new Date(number * 1000);
  }
  // Multiply by 1 in case it's a number in quotes
  return new Date(number * 1);
}

function dateFromString(str) {
  // Check ISO format first
  let val = parseISO(str);
  if (val.toString() !== 'Invalid Date') return val;

  val = new Date(str);
  if (val.toString() !== 'Invalid Date') return val;

  return null;
}

export function dateFormat([value, style], { withTimeZone = false }) {
  // see format breaking in upgrade to date-fns 2.x https://github.com/date-fns/date-fns/blob/master/CHANGELOG.md#changed-5
  let date;
  switch (checkType(value)) {
    case 'string':
      date = dateFromString(value);
      break;
    case 'number':
      date = dateFromNumber(value);
      break;
    case 'date':
      date = value;
      break;
    default:
      // date is not a recognized format
      break;
  }

  // at this point, date is either falsey or a Date object
  if (!date) {
    return value || '';
  }

  const zone = withTimeZone ? formatTimeZone(date) : '';
  return format(date, style) + zone;
}

// separate function for testing
export function formatTimeZone(date) {
  let zone; // local timezone ex: 'PST'
  try {
    // passing undefined means default to the browser's locale
    zone = date.toLocaleTimeString(undefined, { timeZoneName: 'short' }).split(' ')[2];
  } catch (e) {
    zone = '';
  }

  return zone ? ` ${zone}` : '';
}

export default helper(dateFormat);
