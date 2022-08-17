import { Response } from 'miragejs';

export const OIDC_BASE_URL = `/vault/access/oidc`;

export const SELECTORS = {
  oidcHeader: '[data-test-oidc-header]',
  oidcClientCreateButton: '[data-test-oidc-configure]',
  oidcRouteTabs: '[data-test-oidc-tabs]',
  oidcLandingImg: '[data-test-oidc-img]',
  confirmDeleteButton: '[data-test-confirm-button="true"]',
  // client route
  clientSaveButton: '[data-test-oidc-client-save]',
  clientCancelButton: '[data-test-oidc-client-cancel]',
  clientDeleteButton: '[data-test-oidc-client-delete] button',
  clientEditButton: '[data-test-oidc-client-edit]',
  clientDetailsTab: '[data-test-oidc-client-details]',
  clientProvidersTab: '[data-test-oidc-client-providers]',

  // assignment route
  assignSaveButton: '[data-test-oidc-assignment-save]',

  // scope routes
  scopeSaveButton: '[data-test-oidc-scope-save]',
  scopeCancelButton: '[data-test-oidc-scope-cancel]',
  scopeDeleteButton: '[data-test-oidc-scope-delete] button',
  scopeEditButton: '[data-test-oidc-scope-edit]',
  scopeDetailsTab: '[data-test-oidc-scope-details]',
  scopeEmptyState: '[data-test-oidc-scope-empty-state]',
  scopeCreateButtonEmptyState: '[data-test-oidc-scope-create-empty-state]',
  scopeCreateButton: '[data-test-oidc-scope-create]',
};

export function overrideMirageResponse(httpStatus, data) {
  if (httpStatus === 403) {
    return new Response(
      403,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] })
    );
  }
  if (httpStatus === 404) {
    return new Response(404, { 'Content-Type': 'application/json' });
  }
  if (httpStatus === 200) {
    return new Response(200, { 'Content-Type': 'application/json' }, JSON.stringify(data));
  }
  return {
    request_id: crypto.randomUUID(),
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    wrap_info: null,
    warnings: null,
    auth: null,
    data: { ...data },
  };
}

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

// MOCK RESPONSES:

export const CLIENT_LIST_RESPONSE = {
  keys: ['some-app'],
  key_info: {
    'some-app': {
      assignments: ['allow_all'],
      client_id: 'whaT7KB0C3iBH1l3rXhd5HPf0n6vXU0s',
      client_secret: 'hvo_secret_nkJSTu2NVYqylXwFbFijsTxJHg4Ic4gqSJw7uOZ4FbSXcObngDkKoVsvyndrf2O8',
      client_type: 'confidential',
      id_token_ttl: 0,
      key: 'default',
      redirect_uris: [],
    },
  },
};

export const CLIENT_DATA_RESPONSE = {
  access_token_ttl: 0,
  assignments: ['allow_all'],
  client_id: 'whaT7KB0C3iBH1l3rXhd5HPf0n6vXU0s',
  client_secret: 'hvo_secret_nkJSTu2NVYqylXwFbFijsTxJHg4Ic4gqSJw7uOZ4FbSXcObngDkKoVsvyndrf2O8',
  client_type: 'confidential',
  id_token_ttl: 0,
  key: 'default',
  redirect_uris: [],
};

export const SCOPE_LIST_RESPONSE = {
  keys: ['test-scope'],
};

export const SCOPE_DATA_RESPONSE = {
  description: 'this is a test',
  template: '{ test }',
};

export const PROVIDER_LIST_RESPONSE = {
  keys: ['test-provider'],
};

export const PROVIDER_DATA_RESPONSE = {
  allowed_client_ids: ['*'],
  issuer: '',
  scopes_supported: ['test-scope'],
};