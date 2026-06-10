/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { sanitizePath } from '../sanitize-path';
import { formatArgsFromPayload } from './formatters';

/**
 * Replaces OpenAPI style tokens {token} with values from an object.
 * eg. - path /identity/oidc/client/{name} and params { name: 'root' } returns /identity/oidc/client/root
 */
export const formatDynamicApiPath = <T extends object>(path: string, params: T) => {
  const formattedPath = path.replace(/{(\w+)}/g, (match, key) => {
    // Type guard: check if 'key' is actually a property of 'params'
    if (key in params) {
      const value = params[key as keyof T];
      return value !== undefined ? encodeURIComponent(String(value)) : match;
    }
    return match;
  });
  return sanitizePath(formattedPath);
};

// returns formatted CURL command for given API path and payload
export const generateCurlCommand = (path: string, payload: Record<string, unknown>, namespace?: string) => {
  return `curl \\
  --header "X-Vault-Token: $VAULT_TOKEN"${
    namespace ? `\\\n  --header "X-Vault-Namespace: ${namespace}"\\` : ''
  }
  --request POST \\
  --data '${JSON.stringify(formatArgsFromPayload(payload))}' \\
  $VAULT_ADDR/v1/${formatDynamicApiPath(path, payload)}
`;
};
