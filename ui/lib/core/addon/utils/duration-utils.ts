/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * These utils are used for managing Duration type values
 * (eg. '30m', '365d'). Most often used in the context of TTLs
 */
import Duration from '@icholy/duration';

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

// parses duration string ('3m') and returns seconds
export const durationToSeconds = (duration: string) => {
  // we assume numbers are seconds
  if (typeof duration === 'number') return duration;
  try {
    return Duration.parse(duration).seconds();
  } catch (e) {
    // since 0 is falsy, parent should explicitly check for null and decide how to handle parsing error
    return null;
  }
};
