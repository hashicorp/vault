/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { format, parseISO } from 'date-fns';

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
export const parseAPITimestamp = (timestamp, style) => {
  if (typeof timestamp !== 'string') return timestamp;
  const date = parseISO(timestamp.split('T')[0]);
  if (!style) return date;
  return format(date, style);
};
