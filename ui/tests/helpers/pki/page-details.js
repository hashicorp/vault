export const SELECTORS = {
  breadcrumbContainer: '[data-test-breadcrumbs]',
  breadcrumbs: '[data-test-breadcrumbs] li',
  confirmDelete: '[data-test-confirm-action-trigger]',
};
export const ROLE_SELECTORS = {
  title: '[data-test-role-details-title]',
  issuerLabel: '[data-test-row-label="Issuer"]',
  noStoreValue: '[data-test-value-div="Store in storage backend"]',
  keyUsageValue: '[data-test-value-div="Key usage"]',
  extKeyUsageValue: '[data-test-value-div="Ext key usage"]',
};
export const KEY_SELECTORS = {
  title: '[data-test-key-details-title]',
  keyIdValue: '[data-test-value-div="Key ID"]',
  keyNameValue: '[data-test-value-div="Key name"]',
  keyTypeValue: '[data-test-value-div="Key type"]',
  keyBitsValue: '[data-test-value-div="Key bits"]',
  keyDeleteButton: '[data-test-pki-key-delete]',
};
