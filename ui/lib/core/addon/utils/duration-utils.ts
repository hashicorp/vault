/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * These utils are used for managing Duration type values
 * (eg. '30m', '365d'). Most often used in the context of TTLs
 */
interface SecondsMap {
  s: 1;
  m: 60;
  h: 3600;
  d: 86400;
  [key: string]: number;
}
export const secondsMap: SecondsMap = {
  s: 1,
  m: 60,
  h: 3600,
  d: 86400,
};
export const convertToSeconds = (time: number, unit: string) => {
  return time * (secondsMap[unit] || 1);
};
export const convertFromSeconds = (seconds: number, unit: string) => {
  return seconds / (secondsMap[unit] || 1);
};
export const goSafeConvertFromSeconds = (seconds: number, unit: string) => {
  // Go only accepts s, m, or h units
  const u = unit === 'd' ? 'h' : unit;
  return convertFromSeconds(seconds, u) + u;
};
export const largestUnitFromSeconds = (seconds: number) => {
  let unit = 's';
  if (seconds === 0) return unit;
  // get largest unit with no remainder
  if (seconds % secondsMap.d === 0) {
    unit = 'd';
  } else if (seconds % secondsMap.h === 0) {
    unit = 'h';
  } else if (seconds % secondsMap.m === 0) {
    unit = 'm';
  }
  return unit;
};
