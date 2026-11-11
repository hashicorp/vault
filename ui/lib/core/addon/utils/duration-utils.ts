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

// Vault parses every duration into Go's time.Duration (int64 nanoseconds) and rejects anything longer.
// It's a language limit, not a server setting, so it only changes if Vault changes its duration type.
export const MAX_DURATION_SECONDS = 9223372036;

// Matches a number with an optional s/m/h/d unit, including the exponent notation large numbers get when
// interpolated (e.g. '1e+21h') and 'Infinity', which TtlPicker produces when a long digit input overflows.
// Duration.parse can't read either form, so they are handled here.
const NUMBER_WITH_UNIT = /^(\d+(?:\.\d+)?(?:e\+?\d+)?|Infinity)([smhd])?$/;

/**
 * Returns true when a duration is longer than `maxSeconds` (Vault's maximum by default).
 * Empty values and values that cannot be parsed return false so the API can report them.
 */
export const isDurationOverMax = (duration: unknown, maxSeconds = MAX_DURATION_SECONDS): boolean => {
  let seconds: number | null = null;
  if (typeof duration === 'number') {
    seconds = duration;
  } else if (typeof duration === 'string' && duration.trim()) {
    const match = duration.trim().match(NUMBER_WITH_UNIT);
    seconds = match
      ? convertToSeconds(Number(match[1]), match[2] ?? 's')
      : durationToSeconds(duration.trim());
  }
  return seconds !== null && seconds > maxSeconds;
};

const MESSAGE_UNITS: Array<[string, number]> = [
  ['day', secondsMap.d],
  ['hour', secondsMap.h],
  ['minute', secondsMap.m],
  ['second', secondsMap.s],
];

// Writes the maximum in the largest unit it fills, rounded down so the stated limit is always accepted.
const formatMaxDuration = (maxSeconds: number) => {
  const [unit, size] = MESSAGE_UNITS.find(([, unitSeconds]) => maxSeconds >= unitSeconds) ?? ['second', 1];
  const count = Math.floor(maxSeconds / size);
  return `${new Intl.NumberFormat('en-US').format(count)} ${unit}${count === 1 ? '' : 's'}`;
};

export const maxDurationMessage = (label: string, maxSeconds = MAX_DURATION_SECONDS) =>
  `${label} is too large. Maximum allowed value is ${formatMaxDuration(maxSeconds)}.`;
