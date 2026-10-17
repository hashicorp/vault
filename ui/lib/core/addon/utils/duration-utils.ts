/**
 * Copyright IBM Corp. 2016, 2026
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

// Vault parses durations into Go's time.Duration (int64 nanoseconds) and rejects anything longer.
export const MAX_DURATION_SECONDS = 9223372036;

// Matches plain seconds, including the exponent notation large numbers get when interpolated (e.g. '1e+21s')
// and 'Infinity', which TtlPicker produces when a long digit input overflows.
const PLAIN_SECONDS = /^(\d+(?:\.\d+)?(?:e\+?\d+)?|Infinity)s?$/;

/**
 * Returns true when a duration is longer than Vault can store.
 * Empty values and values that cannot be parsed return false so the API can report them.
 */
export const isDurationOverMax = (duration: unknown): boolean => {
  let seconds: number | null = null;
  if (typeof duration === 'number') {
    seconds = duration;
  } else if (typeof duration === 'string' && duration.trim()) {
    const plain = duration.trim().match(PLAIN_SECONDS);
    seconds = plain ? Number(plain[1]) : durationToSeconds(duration.trim());
  }
  return seconds !== null && seconds > MAX_DURATION_SECONDS;
};

export const maxDurationMessage = (label: string) =>
  `${label} is too large. Maximum allowed value is 106,751 days.`;
