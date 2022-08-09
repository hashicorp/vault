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
  scopeDeleteButton: '[data-test-scope-delete]',
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
  if (httpStatus === 204) {
    return new Response(204, { 'Content-Type': 'application/json' });
  }
  return new Response(200, { 'Content-Type': 'application/json' }, JSON.stringify(data));
}
