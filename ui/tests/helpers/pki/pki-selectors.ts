/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_OVERVIEW = {
  issuersCardTitle: '[data-test-overview-card-title="Issuers"]',
  issuersCardSubtitle: '[data-test-overview-card-subtitle="Issuers"]',
  issuersCardLink: '[data-test-overview-card-container="Issuers"] a',
  issuersCardOverviewNum: '[data-test-overview-card-container="Issuers"] h2',
  rolesCardTitle: '[data-test-overview-card-title="Roles"]',
  rolesCardSubtitle: '[data-test-overview-card-subtitle="Roles"]',
  rolesCardLink: '[data-test-overview-card-container="Roles"] a',
  rolesCardOverviewNum: '[data-test-overview-card-container="Roles"] h2',
  issueCertificate: '[data-test-overview-card-title="Issue certificate"]',
  issueCertificateInput: '[data-test-issue-certificate-input]',
  issueCertificatePowerSearch: '[data-test-issue-certificate-input] span',
  issueCertificateButton: '[data-test-issue-certificate-button]',
  viewCertificate: '[data-test-overview-card-title="View certificate"]',
  viewCertificateInput: '[data-test-view-certificate-input]',
  viewCertificatePowerSearch: '[data-test-view-certificate-input] span',
  viewCertificateButton: '[data-test-view-certificate-button]',
  viewIssuerInput: '[data-test-issue-issuer-input]',
  viewIssuerPowerSearch: '[data-test-issue-issuer-input] span',
  viewIssuerButton: '[data-test-view-issuer-button]',
  firstPowerSelectOption: '[data-option-index="0"]',
};

// COMPONENTS
export const PKI_CONFIGURE_CREATE = {
  // page::pki-configure-create
  nextStepsBanner: '[data-test-config-next-steps]',
  option: '[data-test-pki-config-option]',
  optionByKey: (key: string) => `[data-test-pki-config-option="${key}"]`,
  doneButton: '[data-test-done]',
  configureButton: '[data-test-configure-pki-button]',
  // pki-generate-root
  generateRootOption: '[data-test-pki-config-option="generate-root"]',
  // pki-ca-cert-import
  importForm: '[data-test-pki-import-pem-bundle-form]',
  importSubmit: '[data-test-pki-import-pem-bundle]',
  importSectionLabel: '[data-test-import-section-label]',
  importMapping: '[data-test-imported-bundle-mapping]',
  importedIssuer: '[data-test-imported-issuer]',
  importedKey: '[data-test-imported-key]',
  // generate-intermediate
  csrDetails: '[data-test-generate-csr-result]',
};

export const PKI_DELETE_ALL_ISSUERS = {
  issuerLink: '[data-test-delete-all-issuers-link]',
  deleteAllIssuerModal: '#confirmation-modal',
  deleteAllIssuerInput: '[data-test-confirmation-modal-input="Delete all issuers?"]',
  deleteAllIssuerButton: '[data-test-confirm-button="Delete all issuers?"]',
};

export const PKI_GENERATE_ROOT = {
  mainSectionTitle: '[data-test-generate-root-title="Root parameters"]',
  urlSectionTitle: '[data-test-generate-root-title="Issuer URLs"]',
  keyParamsGroupToggle: '[data-test-toggle-group="Key parameters"]',
  sanGroupToggle: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalGroupToggle: '[data-test-toggle-group="Additional subject fields"]',
  toggleGroupDescription: '[data-test-toggle-group-description]',
  groupFields: (group: string) => `[data-test-group="${group}"] [data-test-field]`,
  formInvalidError: '[data-test-pki-generate-root-validation-error]',
  urlsSection: '[data-test-urls-section]',
  urlField: '[data-test-urls-section] [data-test-input]',
  // Shown values after save
  saved: {
    certificate: '[data-test-value-div="Certificate"] [data-test-certificate-card]',
    commonName: '[data-test-row-value="Common name"]',
    issuerName: '[data-test-row-value="Issuer name"]',
    issuerLink: '[data-test-value-div="Issuer ID"] a',
    keyName: '[data-test-row-value="Key name"]',
    keyLink: '[data-test-value-div="Key ID"] a',
    privateKey: '[data-test-value-div="Private key"] [data-test-certificate-card]',
    serialNumber: '[data-test-row-value="Serial number"]',
  },
};

export const PKI_CROSS_SIGN = {
  objectListInput: (key: string, row = 0) => `[data-test-object-list-input="${key}-${row}"]`,
  addRow: '[data-test-object-list-add-button',
  statusCount: '[data-test-cross-sign-status-count]',
  signedIssuerRow: (row = 0) => `[data-test-info-table-row="${row}"]`,
  signedIssuerCol: (attr: string) => `[data-test-info-table-column="${attr}"]`,
  rowValue: (attr: string) => `[data-test-row-value="${attr}"]`,
  copyButton: (attr: string) => `[data-test-value-div="${attr}"] [data-test-copy-button]`,
};

export const PKI_ISSUER_DETAILS = {
  configure: '[data-test-pki-issuer-configure]',
  copyButtonByName: (name: string) => `[data-test-value-div="${name}"] [data-test-copy-button]`,
  crossSign: '[data-test-pki-issuer-cross-sign]',
  defaultGroup: '[data-test-details-group="default"]',
  download: '[data-test-issuer-download]',
  groupTitle: '[data-test-group-title]',
  parsingAlertBanner: '[data-test-parsing-error-alert-banner]',
  rotateModal: '#pki-rotate-root-modal',
  rotateModalGenerate: '[data-test-root-rotate-step-one]',
  rotateRoot: '[data-test-pki-issuer-rotate-root]',
  row: '[data-test-component="info-table-row"]',
  signIntermediate: '[data-test-pki-issuer-sign-int]',
  urlsGroup: '[data-test-details-group="Issuer URLs"]',
  valueByName: (name: string) => `[data-test-value-div="${name}"]`,
};

export const PKI_KEY_FORM = {
  validationError: '[data-test-pki-key-validation-error]',
};

export const PKI_NOT_VALID_AFTER = {
  radioTtl: '[data-test-radio-button="ttl"]',
  radioTtlLabel: '[data-test-radio-label="ttl"]',
  radioDate: '[data-test-radio-button="not_after"]',
  radioDateLabel: '[data-test-radio-label="specificDate"]',
  ttlForm: '[data-test-ttl-inputs]',
  ttlTimeInput: '[data-test-ttl-value="TTL"]',
  ttlUnitInput: '[data-test-select="ttl-unit"]',
  dateInput: '[data-test-input="not_after"]',
};

export const PKI_ROLE_FORM = {
  domainHandling: '[data-test-toggle-group="Domain handling"]',
  keyParams: '[data-test-toggle-group="Key parameters"]',
  keyUsage: '[data-test-toggle-group="Key usage"]',
  digitalSignature: '[data-test-checkbox="DigitalSignature"]',
  keyAgreement: '[data-test-checkbox="KeyAgreement"]',
  keyEncipherment: '[data-test-checkbox="KeyEncipherment"]',
  any: '[data-test-checkbox="Any"]',
  serverAuth: '[data-test-checkbox="ServerAuth"]',
  policyIdentifiers: '[data-test-toggle-group="Policy identifiers"]',
  san: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalSubjectFields: '[data-test-toggle-group="Additional subject fields"]',
};

export const PKI_ROLE_GENERATE = {
  form: '[data-test-pki-generate-cert-form]',
  optionsToggle: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  downloadButton: '[data-test-pki-cert-download-button]',
  revokeButton: '[data-test-pki-cert-revoke-button]',
};

// PAGES
export const PKI_CONFIG_EDIT = {
  errorBanner: '[data-test-config-edit-error]',
  acmeEditSection: '[data-test-acme-edit-section]',
  configEditSection: '[data-test-cluster-config-edit-section]',
  configInput: (attr: string) => `[data-test-input="${attr}"]`,
  stringListInput: (attr: string) => `[data-test-input="${attr}"] [data-test-string-list-input="0"]`,
  urlsEditSection: '[data-test-urls-edit-section]',
  urlFieldInput: (attr: string) => `[data-test-input="${attr}"] textarea`,
  urlFieldLabel: (attr: string) => `[data-test-input="${attr}"] label`,
  crlEditSection: '[data-test-crl-edit-section]',
  crlToggleInput: (attr: string) => `[data-test-input="${attr}"] input`,
  crlTtlInput: (attr: string) => `[data-test-ttl-value="${attr}"]`,
  crlFieldLabel: (attr: string) => `[data-test-input="${attr}"] label`,
  saveButton: '[data-test-configuration-edit-save]',
  cancelButton: '[data-test-configuration-edit-cancel]',
  validationAlert: '[data-test-configuration-edit-validation-alert]',
  deleteButton: (attr: string) => `[data-test-input="${attr}"] [data-test-string-list-button="delete"]`,
  groupHeader: (group: string) => `[data-test-crl-header="${group}"]`,
  checkboxInput: (attr: string) => `[data-test-input="${attr}"]`,
};

export const PKI_ISSUER_LIST = {
  // Page::PkiIssuerList
  issuerListItem: (id: string) => `[data-test-issuer-list="${id}"]`,
  importIssuerLink: '[data-test-generate-issuer="import"]',
  generateIssuerDropdown: '[data-test-issuer-generate-dropdown]',
  generateIssuerRoot: '[data-test-generate-issuer="root"]',
  generateIssuerIntermediate: '[data-test-generate-issuer="intermediate"]',
  issuerPopupDetails: '[data-test-pki-issuer-details]',
};

export const PKI_KEYS = {
  // key index
  importKey: '[data-test-pki-key-import]',
  generateKey: '[data-test-pki-key-generate]',
  keyId: '[data-test-key="id"]',
  keyName: '[data-test-key="name"]',
  popupMenuDetails: '[data-test-key-menu-link="details"]',
  popupMenuEdit: '[data-test-key-menu-link="edit"]',
  // key details
  keyDeleteButton: '[data-test-pki-key-delete]',
  downloadButton: '[data-test-download-button]',
  keyEditLink: '[data-test-pki-key-edit]',
  nextStepsAlert: '[data-test-pki-key-next-steps]',
};

export const PKI_ROLE_DETAILS = {
  issuerLabel: '[data-test-row-label="Issuer"]',
  noStoreValue: '[data-test-value-div="Store in storage backend"]',
  noStoreMetadataValue: '[data-test-value-div="Store metadata in storage backend"]',
  keyUsageValue: '[data-test-value-div="Key usage"]',
  extKeyUsageValue: '[data-test-value-div="Ext key usage"]',
  customTtlValue: '[data-test-value-div="Issued certificates expire after"]',
  deleteRoleButton: '[data-test-pki-role-delete]',
  generateCertLink: '[data-test-pki-role-generate-cert]',
  signCertLink: '[data-test-pki-role-sign-cert]',
  editRoleLink: '[data-test-pki-role-edit-link]',
  createRoleLink: '[data-test-pki-role-create-link]',
};

export const PKI_TIDY_FORM = {
  tidyFormName: (attr: string) => `[data-test-tidy-form="${attr}"]`,
  inputByAttr: (attr: string) => `[data-test-input="${attr}"]`,
  toggleInput: (attr: string) => `[data-test-input="${attr}"] input`,
  intervalDuration: '[data-test-ttl-value="Automatic tidy enabled"]',
  acmeAccountSafetyBuffer: '[data-test-ttl-value="Tidy ACME enabled"]',
  toggleLabel: (label: string) => `[data-test-toggle-label="${label}"]`,
  tidySectionHeader: (header: string) => `[data-test-tidy-header="${header}"]`,
  tidySave: '[data-test-pki-tidy-button]',
  tidyCancel: '[data-test-pki-tidy-cancel]',
  tidyPauseDuration: '[data-test-ttl-value="Pause duration"]',
  editAutoTidyButton: '[data-test-pki-edit-tidy-auto-link]',
};

export const PKI_TIDY = {
  hdsAlertTitle: '[data-test-tidy-status-alert-title]',
  hdsAlertDescription: '[data-test-tidy-status-alert-description]',
  alertUpdatedAt: '[data-test-tidy-status-alert-updated-at]',
  cancelTidyAction: '[data-test-cancel-tidy-action]',
  hdsAlertButtonText: '[data-test-cancel-tidy-action] .hds-button__text',
  timeStartedRow: '[data-test-value-div="Time started"]',
  timeFinishedRow: '[data-test-value-div="Time finished"]',
  cancelTidyModalBackground: '#pki-cancel-tidy-modal',
  tidyEmptyStateConfigure: '[data-test-tidy-empty-state-configure]',
  manualTidyToolbar: '[data-test-pki-manual-tidy-config]',
  autoTidyToolbar: '[data-test-pki-auto-tidy-config]',
  tidyConfigureModal: {
    configureTidyModal: '#pki-tidy-modal',
    tidyModalAutoButton: '[data-test-tidy-modal-auto-button]',
    tidyModalManualButton: '[data-test-tidy-modal-manual-button]',
    tidyModalCancelButton: '[data-test-tidy-modal-cancel-button]',
    tidyOptionsModal: '[data-test-pki-tidy-options-modal]',
  },
};
