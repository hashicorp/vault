import { format, parseISO } from 'date-fns';

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
  if (typeof timestamp !== 'string') return;
  const date = parseISO(timestamp.split('T')[0]);
  if (!style) return date;
  return format(date, style);
};

// convert M/yy (format of dates in charts) to 'Month yyyy' (format in tooltip)
export function formatChartDate(date) {
  const array = date.split('/');
  array.splice(1, 0, '01');
  const dateString = array.join('/');
  return format(new Date(dateString), 'MMMM yyyy');
}
