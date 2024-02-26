/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Clears pki-related data and capabilities so that admin
// capabilities from setup don't rollover
export function clearRecords(store) {
  store.unloadAll('pki/action');
  store.unloadAll('pki/issuer');
  store.unloadAll('pki/key');
  store.unloadAll('pki/role');
  store.unloadAll('pki/sign-intermediate');
  store.unloadAll('pki/tidy');
  store.unloadAll('pki/config/urls');
  store.unloadAll('pki/config/crl');
  store.unloadAll('pki/config/cluster');
  store.unloadAll('pki/config/acme');
  store.unloadAll('pki/certificate/generate');
  store.unloadAll('pki/certificate/sign');
  store.unloadAll('capabilities');
}
