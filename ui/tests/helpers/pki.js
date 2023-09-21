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
