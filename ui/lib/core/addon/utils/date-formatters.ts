/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { format, parse, parseISO } from 'date-fns';
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

// convert API timestamp ( '2021-03-21T00:00:00Z' ) to date object, optionally format
export const parseAPITimestamp = (timestamp: string, style?: string): Date | string | null => {
  if (!timestamp || typeof timestamp !== 'string') return null;

  if (!style) {
    // If no style, return a date object in UTC
    const parsed = parseISO(timestamp) as Date;
    return isValid(parsed) ? parsed : null;
  }

  // Otherwise format it as a calendar date that is timezone agnostic.
  const yearMonthDay = timestamp.split('T')[0] ?? '';
  const date = parse(yearMonthDay, 'yyyy-MM-dd', new Date()); // 'yyyy-MM-dd' lets parse() know the format of yearMonthDay
  return format(date, style) as string;
};

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
