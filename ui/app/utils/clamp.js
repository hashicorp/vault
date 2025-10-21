/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (num, min, max) {
  let inRangeNumber;
  if (typeof num !== 'number') {
    inRangeNumber = min;
  } else {
    inRangeNumber = num;
  }
  return Math.min(Math.max(inRangeNumber, min), max);
}
