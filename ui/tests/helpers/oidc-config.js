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
};

export function overrideMirageResponse(httpStatus, data) {
  if (httpStatus === 403) {
    return new Response(
      403,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] })
    );
  }
  if (httpStatus === 204) {
    return new Response(204, { 'Content-Type': 'application/json' });
  }
  return new Response(200, { 'Content-Type': 'application/json' }, JSON.stringify(data));
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
