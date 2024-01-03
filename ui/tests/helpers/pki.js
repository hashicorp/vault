/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  caChain: '[data-test-value-div="CA chain"] [data-test-certificate-card]',
  certificate: '[data-test-value-div="Certificate"] [data-test-certificate-card]',
  commonName: '[data-test-row-value="Common name"]',
  csr: '[data-test-value-div="CSR"] [data-test-certificate-card]',
  expiryDate: '[data-test-row-value="Expiration date"]',
  issueDate: '[data-test-row-value="Issue date"]',
  issuingCa: '[data-test-value-div="Issuing CA"] [data-test-certificate-card]',
  privateKey: '[data-test-value-div="Private key"] [data-test-certificate-card]',
  revocationTime: '[data-test-row-value="Revocation time"]',
  serialNumber: '[data-test-row-value="Serial number"]',
};

export const STANDARD_META = {
  total: 2,
  currentPage: 1,
  pageSize: 100,
};

export const clearPkiRecords = (store) => {
  // Clears pki-related data and capabilities so that admin
  // capabilities from setup don't rollover in permissions tests
  store.unloadAll('pki/issuer');
  store.unloadAll('pki/action');
  store.unloadAll('pki/config/acme');
  store.unloadAll('pki/certificate/generate');
  store.unloadAll('pki/certificate/sign');
  store.unloadAll('pki/config/cluster');
  store.unloadAll('pki/key');
  store.unloadAll('pki/role');
  store.unloadAll('pki/sign-intermediate');
  store.unloadAll('pki/tidy');
  store.unloadAll('pki/config/urls');
  store.unloadAll('capabilities');
};

export function arbitraryWait() {
  // this is a temporary fix to resolve an issue where about 50% of the tests
  // returned a 404 of the issuers list endpoint even after configuring the mount
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, 1000);
  });
}
