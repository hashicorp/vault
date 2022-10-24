export const PKI_BASE_URL = `/vault/cluster/secrets/backend/pki/roles`;

export const SELECTORS = {
  // Pki role
  roleName: '[data-test-input="name"]',
  issuerRef: '[data-test-input="issuerRef"]',
  customTtl: '[data-test-field="customTtl"]',
  backdateValidity: '[data-test-ttl-value="Backdate validity"]',
  maxTtl: '[data-test-toggle-label="Max TTL"]',
  generateLease: '[data-test-field="generateLease"]',
  noStore: '[data-test-field="noStore"]',
  addBasicConstraints: '[data-test-input="addBasicConstraints"]',
  domainHandling: '[data-test-toggle-group="Domain handling"]',
  keyParams: '[data-test-toggle-group="Key parameters"]',
  keyUsage: '[data-test-toggle-group="Key usage"]',
  extKeyUsageOids: '[data-test-input="extKeyUsageOids"]',
  digitalSignature: '[data-test-input="DigitalSignature"]',
  keyAgreement: '[data-test-input="KeyAgreement"]',
  keyEncipherment: '[data-test-input="KeyEncipherment"]',
  any: '[data-test-input="Any"]',
  serverAuth: '[data-test-input="ServerAuth"]',
  policyIdentifiers: '[data-test-toggle-group="Policy identifiers"]',
  san: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalSubjectFields: '[data-test-toggle-group="Additional subject fields"]',
  roleCreateButton: '[data-test-pki-role-save]',
  roleCancelButton: '[data-test-pki-role-cancel]',
};

export async function clearRecord(store, modelType, id) {
  await store
    .findRecord(modelType, id)
    .then((model) => {
      deleteModelRecord(model);
    })
    .catch(() => {});
}

const deleteModelRecord = async (model) => {
  await model.destroyRecord();
};
