/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import apiPath from './api-path';

const API_PATHS = {
  kvv2: [apiPath`${'backend'}/data/${'path'}`, apiPath`${'backend'}/metadata/${'path'}`],
  secretsList: [apiPath`${'backend'}/`],
};

// Regex-based route matching - more flexible for parent/child relationships
const ROUTE_PATTERNS: Array<{ pattern: RegExp; paths: ReturnType<typeof apiPath>[] }> = [
  // KV v2 routes - matches any kv secret route
  {
    pattern: /^vault\.cluster\.secrets\.backend\.kv/,
    paths: API_PATHS.kvv2,
  },
  {
    pattern: /^vault\.cluster\.secrets\./,
    paths: API_PATHS.secretsList,
  },
];

export default function mapApiPathToRoute(routeName: string): ReturnType<typeof apiPath>[] {
  // Try pattern matching
  for (const { pattern, paths } of ROUTE_PATTERNS) {
    if (pattern.test(routeName)) {
      return paths;
    }
  }

  return [];
}
