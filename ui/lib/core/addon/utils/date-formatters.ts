/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { formatInTimeZone } from 'date-fns-tz';
import isValid from 'date-fns/isValid';

export const datetimeLocalStringFormat = "yyyy-MM-dd'T'HH:mm";

export const ARRAY_OF_MONTHS = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
];

// datetime may be returned from the API client as either a Date object or an ISO string
// strings will be converted from an ISO string ('2021-03-21T00:00:00Z') to date object and optionally formatted
// the timezone of the formatted output will be in UTC and not the local timezone of the user
export function parseAPITimestamp(timestamp: string | Date, style: string): string;
export function parseAPITimestamp(timestamp: string | Date): Date | null;
export function parseAPITimestamp(timestamp: string | Date, style?: string) {
  if (timestamp) {
    if (timestamp instanceof Date && isValid(timestamp)) {
      // if no style (format) is provided return the Date object as is since there is nothing more to parse
      return style ? formatInTimeZone(timestamp, 'UTC', style) : timestamp;
    } else if (typeof timestamp === 'string') {
      // if no style return a date object
      if (!style) {
        const date = new Date(timestamp);
        return isValid(date) ? date : null;
      }
      // otherwise format it as a calendar date that is in UTC.
      return formatInTimeZone(timestamp, 'UTC', style);
    }
  }
  return null;
}

export const buildISOTimestamp = (args: { monthIdx: number; year: number; isEndDate: boolean }) => {
  const { monthIdx, year, isEndDate } = args;
  // passing `0` for the "day" arg to Date.UTC() returns the last day of the previous month
  // which is why the monthIdx is increased by one for end dates.
  // Date.UTC() also returns December if -1 is passed (which happens when January is selected)
  const utc = isEndDate
    ? new Date(Date.UTC(year, monthIdx + 1, 0, 23, 59, 59))
    : new Date(Date.UTC(year, monthIdx, 1));

  // remove milliseconds to return a UTC timestamp that matches the API
  // e.g. "2025-05-01T00:00:00Z" or "2025-09-30T23:59:59Z"
  return utc.toISOString().replace('.000', '');
};

export const isSameMonthUTC = (timestampA: string, timestampB: string): boolean => {
  const dateA = parseAPITimestamp(timestampA) as Date;
  const dateB = parseAPITimestamp(timestampB) as Date;
  if (isValid(dateA) && isValid(dateB)) {
    // Compare in UTC as any date-fns comparisons will be in localized timezones!
    return dateA.getUTCFullYear() === dateB.getUTCFullYear() && dateA.getUTCMonth() === dateB.getUTCMonth();
  }
  return false;
};
