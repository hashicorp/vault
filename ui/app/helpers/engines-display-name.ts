/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ALL_ENGINES } from 'vault/utils/engines-display-data';

export default function engineDisplayName(methodType: string) {
  const displayName = ALL_ENGINES?.find((t) => t.type === methodType)?.displayName;
  return displayName || methodType;
}
