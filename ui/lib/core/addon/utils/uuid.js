/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { v4 as uuidv4 } from 'uuid';

/**
 * Use this instead of uuidv4() so that generated UUIDs can be stubbed in tests.
 */

export default function uuid() {
  return uuidv4();
}
