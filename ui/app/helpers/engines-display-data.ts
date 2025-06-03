/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

export default function engineDisplayData(methodType: string) {
  const engine = ALL_ENGINES?.find((t) => t.type === methodType);
  return engine;
}
