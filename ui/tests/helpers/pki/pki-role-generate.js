export const SELECTORS = {
  form: '[data-test-pki-generate-cert-form]',
  commonNameField: '[data-test-input="commonName"]',
  optionsToggle: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  generateButton: '[data-test-pki-generate-button]',
  cancelButton: '[data-test-pki-generate-cancel]',
  downloadButton: '[data-test-pki-cert-download-button]',
  revokeButton: '[data-test-pki-cert-revoke-button]',
  serialNumber: '[data-test-value-div="Serial number"]',
  certificate: '[data-test-value-div="Certificate"]',
  inlineAlert: '[data-test-alert]',
  commonNameInlineError: '[data-test-field="commonName"] [data-test-inline-alert]',
};
