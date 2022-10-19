export const PKI_BASE_URL = `/vault/cluster/secrets/backend/pki/roles`;

export const SELECTORS = {
  // Pki role
  roleName: '[data-test-input="name"]',
  issuerRef: '[data-test-input="issuerRef"]',
  backdateValidity: '[data-test-ttl-value="Backdate validity"]',
  maxTtl: '[data-test-toggle-label="Max TTL"]',
  generateLease: '[data-test-field="generateLease"]',
  noStore: '[data-test-field="noStore"]',
  addBasicConstraints: '[data-test-input="addBasicConstraints"]',
  domainHandling: '[data-test-toggle-group="Domain handling"]',
  keyParams: '[data-test-toggle-group="Key parameters"]',
  keyUsage: '[data-test-toggle-group="Key usage"]',
  policyIdentifiers: '[data-test-toggle-group="Policy identifiers"]',
  san: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalSubjectFields: '[data-test-toggle-group="Additional subject fields"]',
  roleCreateButton: '[data-test-pki-role-save]',
  roleCancelButton: '[data-test-pki-role-cancel]',
};

export function overrideCapabilities(requestPath, capabilitiesArray) {
  // sample of capabilitiesArray: ['read', 'update']
  return {
    request_id: '40f7e44d-af5c-9b60-bd20-df72eb17e294',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      capabilities: capabilitiesArray,
      [requestPath]: capabilitiesArray,
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
}

export async function clearRecord(store, modelType, id) {
  await store
    .findRecord(modelType, id)
    .then((model) => {
      deleteModelRecord(model);
    })
    .catch(() => {
      // swallow error
    });
}

const deleteModelRecord = async (model) => {
  await model.destroyRecord();
};
