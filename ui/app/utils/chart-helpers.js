/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { format } from 'd3-format';
import { mean } from 'd3-array';

// COLOR THEME:
export const BAR_PALETTE = ['#CCE3FE', '#1060FF', '#C2C5CB', '#656A76'];
export const UPGRADE_WARNING = '#FDEEBA';
export const GREY = '#EBEEF2';

// TRANSLATIONS:
export const TRANSLATE = { left: -11 };
export const SVG_DIMENSIONS = { height: 190, width: 500 };

export const BAR_WIDTH = 17; // data bar width is 17 pixels

// Reference for tickFormat https://www.youtube.com/watch?v=c3MCROTNN8g
export function numericalAxisLabel(number) {
  if (number < 1000) return number;
  if (number < 1100) return format('.1s')(number);
  if (number < 2000) return format('.2s')(number); // between 1k and 2k, show 2 decimals
  if (number < 10000) return format('.1s')(number);
  // replace SI prefix of 'G' for billions to 'B'
  return format('.2s')(number).replace('G', 'B');
}

export function calculateAverage(dataset, objectKey) {
  // before mapping for values, check that the objectKey exists at least once in the dataset because
  // map returns 0 when dataset[objectKey] is undefined in order to calculate average
  if (!Array.isArray(dataset) || !objectKey || !dataset.some((d) => Object.keys(d).includes(objectKey))) {
    return null;
  }

  const integers = dataset.map((d) => (d[objectKey] ? d[objectKey] : 0));
  const checkIntegers = integers.every((n) => Number.isInteger(n)); // decimals will be false
  return checkIntegers ? Math.round(mean(integers)) : null;
}

/**
 * Calculates the sum of an array of numbers with optional decimal precision.
 * This function fixes floating-point arithmetic errors by rounding to a specified
 * number of decimal places. For example, 48.7888 + 0.0112 = 48.800000000000004
 * in JavaScript, but with fixedDecimalPlaces=4, it returns 48.8.
 *
 * @param {number[]} integerArray - Array of numbers to sum
 * @param {number} [fixedDecimalPlaces] - Optional number of decimal places for precision.
 *                                        If provided, the sum is rounded to this precision.
 * @returns {number|null} Returns the sum as a number, or null if invalid input.
 *                        When fixedDecimalPlaces is provided, returns a number rounded
 *                        to that precision (e.g., 48.8 instead of 48.800000000000004).
 *
 * @example
 * calculateSum([2, 3])                    // Returns 5
 * calculateSum([48.7888, 0.0112], 4)      // Returns 48.8 (fixes floating-point error)
 * calculateSum([73.1832, 0.0168], 4)      // Returns 73.2
 * calculateSum([10, 20, 30], 4)           // Returns 60
 * calculateSum(['one', 2])                // Returns null (invalid input)
 */
export function calculateSum(integerArray, fixedDecimalPlaces) {
  if (!Array.isArray(integerArray) || integerArray.some((n) => typeof n !== 'number')) {
    return null;
  }

  const sum = integerArray.reduce((a, b) => a + b, 0);

  if (typeof fixedDecimalPlaces === 'number' && fixedDecimalPlaces >= 0) {
    return parseFloat(sum.toFixed(fixedDecimalPlaces));
  }

  return sum;
}

/**
 * Formats a number for display with a fixed number of decimal places.
 * This helper function converts numbers to strings with trailing zeros preserved,
 * which is useful for displaying values like "48.8000" instead of "48.8".
 *
 * @param {number} number - The number to format
 * @param {number} decimalPlaces - The number of decimal places to display
 * @returns {number|string} Returns the original number if invalid inputs,
 *                          returns 0 as-is for zero values,
 *                          otherwise returns a string with fixed decimal places
 *
 * @example
 * toFixedDisplay(48.8, 4)      // Returns "48.8000"
 * toFixedDisplay(73.2, 4)      // Returns "73.2000"
 * toFixedDisplay(0, 4)         // Returns 0 (not "0.0000")
 * toFixedDisplay(100, 2)       // Returns "100.00"
 */
export function toFixedDisplay(number, decimalPlaces) {
  if (typeof number !== 'number' || typeof decimalPlaces !== 'number' || decimalPlaces < 0) {
    return number;
  }

  if (number === 0) {
    return 0;
  }

  return number.toFixed(decimalPlaces);
}
