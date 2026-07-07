/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';

import type { PATH_MAP } from 'vault/utils/constants/capabilities';
import type { Capabilities } from 'vault/vault/app-types';
import type CapabilitiesService from 'vault/services/capabilities';

const DEFAULT_CAPABILITIES = {
  canCreate: true,
  canDelete: true,
  canList: true,
  canPatch: true,
  canRead: true,
  canSudo: true,
  canUpdate: true,
};
/**
 * Stubs `capabilities.for()` for a specific PATH_MAP key.
 * Returns the specified capability for the target path and
 * DEFAULT_CAPABILITIES (all true) for all other paths.
 *
 * @example
 * stubCapabilitiesFor(
 *   this.owner.lookup('service:capabilities'),
 *   'pkiExternalConfigAcmeAccount',
 *   { canList: true }
 * );
 *
 * @example
 * stubCapabilitiesFor(
 *   this.owner.lookup('service:capabilities'),
 *   'pkiExternalLookupOrders',
 *   { canRead: true, canList: false }
 * );
 */
export const stubCapabilitiesFor = (
  capabilities: CapabilitiesService,
  pathKey: keyof typeof PATH_MAP,
  capabilityOverrides: Capabilities
) => {
  const stub = sinon.stub(capabilities, 'for');

  // Return specific capabilities for the target path
  stub.withArgs(pathKey).returns({ ...DEFAULT_CAPABILITIES, ...capabilityOverrides });

  // Return default capabilities for all other paths
  stub.returns(DEFAULT_CAPABILITIES);

  return stub;
};

/**
 * Stubs `capabilities.fetch()` to return specific capabilities for specified paths
 * and DEFAULT_CAPABILITIES (all true) for all other paths.
 *
 * @example
 * stubCapabilitiesFetch(
 *   this.owner.lookup('service:capabilities'),
 *   {
 *     'pki-external-ca/config/acme-account': { canList: true, canRead: false },
 *     'pki-external-ca/config/dns': { canList: false, canRead: true }
 *   }
 * );
 */
export const stubCapabilitiesFetch = (
  capabilities: CapabilitiesService,
  capabilityOverrides: Record<string, Capabilities>
) => {
  const stub = sinon.stub(capabilities, 'fetch');

  // Stub returns a function that processes the paths array
  stub.callsFake(async (paths: string[]) => {
    return paths.reduce(
      (obj, path) => {
        // Check if this path has an override
        if (path in capabilityOverrides) {
          // spread defaults and then override specified capabilities
          obj[path] = { ...DEFAULT_CAPABILITIES, ...capabilityOverrides[path] };
        } else {
          // Return default capabilities for paths not in overrides
          obj[path] = DEFAULT_CAPABILITIES;
        }
        return obj;
      },
      {} as Record<string, Capabilities>
    );
  });

  return stub;
};

/**
 * Stubs `capabilities.fetch()` using PATH_MAP keys instead of explicit paths.
 * Automatically resolves PATH_MAP keys to their API paths using the provided params.
 * Returns DEFAULT_CAPABILITIES for paths not specified in pathCapabilities.
 *
 * @example - in acceptance test with mount path
 * stubCapabilitiesForPaths(
 *   this.owner.lookup('service:capabilities'),
 *   {
 *     pkiExternalConfigAcmeAccount: { canList: true, canRead: false },
 *     pkiExternalConfigDns: { canList: false, canRead: true }
 *   },
 *   { backend: 'my-pki-mount/' }
 * );
 *
 * @example - in integration test with role name
 * stubCapabilitiesForPaths(
 *   this.owner.lookup('service:capabilities'),
 *   {
 *     pkiExternalRole: { canRead: true }
 *   },
 *   { backend: 'pki/', roleName: 'my-role' }
 * );
 */
export const stubCapabilitiesForPaths = (
  capabilities: CapabilitiesService,
  pathCapabilities: Partial<Record<keyof typeof PATH_MAP, Capabilities>>,
  params: Record<string, string> = {}
) => {
  const paths = (Object.keys(pathCapabilities) as Array<keyof typeof PATH_MAP>).reduce(
    (obj, key) => {
      const capability = pathCapabilities[key];
      if (capability) {
        obj[capabilities.pathFor(key, params)] = capability;
      }
      return obj;
    },
    {} as Record<string, Capabilities>
  );
  return stubCapabilitiesFetch(capabilities, paths);
};
