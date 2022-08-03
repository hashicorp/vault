import { Response } from 'miragejs';

export const BASE_URL = `/vault/access/oidc`;

export const SELECTORS = {
  oidcHeader: '[data-test-oidc-header]',
  oidcClientCreateButton: '[data-test-oidc-configure]',
  oidcRouteTabs: '[data-test-oidc-tabs]',
  oidcLandingImg: '[data-test-oidc-img]',
  clientSaveButton: '[data-test-oidc-client-save]',
  clientDeleteButton: '[data-test-client-delete] button',
  confirmDeleteButton: '[data-test-confirm-button="true"]',
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
