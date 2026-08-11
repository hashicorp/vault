/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

// Aggregates multiple Vault ACL policy strings into a single set of unique capabilities per path.

interface AggregatedPolicy {
  [path: string]: string[];
}

export interface AggregatePolicy {
  policy: AggregatedPolicy;
  policyString: string;
}

/**
 * Parses a single Vault policy string and extracts paths with their capabilities
 * @param policyString - HCL formatted policy string
 * @returns Object mapping paths to their capabilities
 */
function parsePolicyString(policyString: string): AggregatedPolicy {
  const policies: AggregatedPolicy = {};

  // Match path blocks: path "..." { capabilities = [...] }
  // This regex handles multi-line blocks and captures the path and capabilities
  const pathRegex = /path\s+"([^"]+)"\s*\{[^}]*capabilities\s*=\s*\[([^\]]*)\]/g;

  let match: RegExpExecArray | null;
  while ((match = pathRegex.exec(policyString)) !== null) {
    const path = match[1];
    const capabilitiesStr = match[2] ?? '';

    if (!path) {
      continue;
    }

    // Parse capabilities array, removing quotes and whitespace
    const capabilities = capabilitiesStr
      .split(',')
      .map((cap) => cap.trim().replace(/['"]/g, ''))
      .filter((cap) => cap.length > 0);

    if (!policies[path]) {
      policies[path] = [];
    }

    // Add capabilities, avoiding duplicates
    capabilities.forEach((cap) => {
      const pathCaps = policies[path];
      if (pathCaps && !pathCaps.includes(cap)) {
        pathCaps.push(cap);
      }
    });
  }

  return policies;
}

/**
 * Aggregates multiple Vault policies, combining capabilities for duplicate paths
 * @param policyStrings - Array of HCL formatted policy strings
 * @returns Object containing aggregated policies and formatted string output
 */
export function aggregatePolicies(policyStrings: string[]): AggregatePolicy {
  const aggregated: AggregatedPolicy = {};

  // Parse each policy string and merge results
  policyStrings.forEach((policyString) => {
    const parsed = parsePolicyString(policyString);

    Object.entries(parsed).forEach(([path, capabilities]) => {
      if (!aggregated[path]) {
        aggregated[path] = [];
      }

      // Add capabilities, avoiding duplicates
      capabilities.forEach((cap) => {
        const pathCaps = aggregated[path];
        if (pathCaps && !pathCaps.includes(cap)) {
          pathCaps.push(cap);
        }
      });
    });
  });

  // Sort capabilities for consistent output
  Object.keys(aggregated).forEach((path) => {
    aggregated[path]?.sort();
  });

  // Generate policy string in HCL format
  const policyString = Object.entries(aggregated)
    .map(([path, capabilities]) => {
      const capsStr = capabilities.map((cap) => `"${cap}"`).join(', ');
      return `path "${path}" {\n    capabilities = [${capsStr}]\n}`;
    })
    .join('\n');

  return {
    policy: aggregated,
    policyString,
  };
}
