import { format } from 'd3-format';

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
