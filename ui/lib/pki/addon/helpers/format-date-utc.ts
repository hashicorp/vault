/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseAPITimestamp } from 'core/utils/date-formatters';

export default function formatDateUTC(isoString: string): string | null {
  // parseAPITimestamp formats to UTC so we include it as part of the display string
  return parseAPITimestamp(isoString, "MM/dd/yyyy, HH:mm 'UTC'");
}
