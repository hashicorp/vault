import { format } from 'd3-format';

// COLOR THEME:
export const LIGHT_AND_DARK_BLUE = ['#BFD4FF', '#1563FF'];
export const BAR_COLOR_HOVER = ['#1563FF', '#0F4FD1'];
export const GREY = '#EBEEF2';

// TRANSLATIONS:
export const TRANSLATE = { left: -11 };
export const SVG_DIMENSIONS = { height: 190, width: 500 };

// Reference for tickFormat https://www.youtube.com/watch?v=c3MCROTNN8g
export function formatNumbers(number) {
  // replace SI prefix of 'G' for billions to 'B'
  return format('.1s')(number).replace('G', 'B');
}
