import { format } from 'd3-format';
import { mean } from 'd3-array';

// COLOR THEME:
export const LIGHT_AND_DARK_BLUE = ['#BFD4FF', '#1563FF'];
export const UPGRADE_WARNING = '#FDEEBA';
export const BAR_COLOR_HOVER = ['#1563FF', '#0F4FD1'];
export const GREY = '#EBEEF2';

// TRANSLATIONS:
export const TRANSLATE = { left: -11 };
export const SVG_DIMENSIONS = { height: 190, width: 500 };

// Reference for tickFormat https://www.youtube.com/watch?v=c3MCROTNN8g
export function formatNumbers(number) {
  if (number < 1000) return number;
  if (number < 10000) return format('.1s')(number);
  // replace SI prefix of 'G' for billions to 'B'
  return format('.2s')(number).replace('G', 'B');
}

export function formatTooltipNumber(value) {
  if (typeof value !== 'number') {
    return value;
  }
  // formats a number according to the locale
  return new Intl.NumberFormat().format(value);
}

export function calculateAverage(dataset, objectKey) {
  if (!Array.isArray(dataset) || dataset?.length === 0) return null;
  // if an array of objects, objectKey of the integer we want to calculate, ex: 'entity_clients'
  // if d[objectKey] is undefined there is no value, so return 0
  const getIntegers = objectKey ? dataset?.map((d) => (d[objectKey] ? d[objectKey] : 0)) : dataset;
  const checkIntegers = getIntegers.every((n) => Number.isInteger(n)); // decimals will be false
  return checkIntegers ? Math.round(mean(getIntegers)) : null;
}
