/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { debug } from '@ember/debug';

export const OIDC_BASE_URL = `/vault/access/oidc`;

export const SELECTORS = {
  oidcHeader: '[data-test-oidc-header]',
  oidcClientCreateButton: '[data-test-oidc-configure]',
  oidcRouteTabs: '[data-test-oidc-tabs]',
  oidcLandingImg: '[data-test-oidc-img]',
  confirmActionButton: '[data-test-confirm-button]',
  inlineAlert: '[data-test-inline-alert]',
  // client route
  clientSaveButton: '[data-test-oidc-client-save]',
  clientCancelButton: '[data-test-oidc-client-cancel]',
  clientDeleteButton: '[data-test-oidc-client-delete]',
  clientEditButton: '[data-test-oidc-client-edit]',
  clientDetailsTab: '[data-test-oidc-client-details]',
  clientProvidersTab: '[data-test-oidc-client-providers]',

  // assignment route
  assignmentSaveButton: '[data-test-oidc-assignment-save]',
  assignmentCreateButton: '[data-test-oidc-assignment-create]',
  assignmentEditButton: '[data-test-oidc-assignment-edit]',
  assignmentDeleteButton: '[data-test-oidc-assignment-delete]',
  assignmentCancelButton: '[data-test-oidc-assignment-cancel]',
  assignmentDetailsTab: '[data-test-oidc-assignment-details]',

  // scope routes
  scopeSaveButton: '[data-test-oidc-scope-save]',
  scopeCancelButton: '[data-test-oidc-scope-cancel]',
  scopeDeleteButton: '[data-test-oidc-scope-delete]',
  scopeEditButton: '[data-test-oidc-scope-edit]',
  scopeDetailsTab: '[data-test-oidc-scope-details]',
  scopeEmptyState: '[data-test-oidc-scope-empty-state]',
  scopeCreateButton: '[data-test-oidc-scope-create]',

  // key route
  keySaveButton: '[data-test-oidc-key-save]',
  keyCancelButton: '[data-test-oidc-key-cancel]',
  keyDeleteButton: '[data-test-oidc-key-delete]',
  keyEditButton: '[data-test-oidc-key-edit]',
  keyRotateButton: '[data-test-oidc-key-rotate]',
  keyDetailsTab: '[data-test-oidc-key-details]',
  keyClientsTab: '[data-test-oidc-key-clients]',

  // provider route
  providerSaveButton: '[data-test-oidc-provider-save]',
  providerCancelButton: '[data-test-oidc-provider-cancel]',
  providerDeleteButton: '[data-test-oidc-provider-delete]',
  providerEditButton: '[data-test-oidc-provider-edit]',
  providerDetailsTab: '[data-test-oidc-provider-details]',
  providerClientsTab: '[data-test-oidc-provider-clients]',
};

export async function clearRecord(store, modelType, id) {
  await store
    .findRecord(modelType, id)
    .then((model) => {
      deleteModelRecord(model);
    })
    .catch(() => {
      debug(`Clearing record failed for ${modelType} with id: ${id}`);
      // swallow error
    });
}

const deleteModelRecord = async (model) => {
  await model.destroyRecord();
};

// MOCK RESPONSES:

export const CLIENT_LIST_RESPONSE = {
  keys: ['test-app', 'app-1'],
  key_info: {
    'test-app': {
      assignments: ['allow_all'],
      client_id: 'whaT7KB0C3iBH1l3rXhd5HPf0n6vXU0s',
      client_secret: 'hvo_secret_nkJSTu2NVYqylXwFbFijsTxJHg4Ic4gqSJw7uOZ4FbSXcObngDkKoVsvyndrf2O8',
      client_type: 'confidential',
      id_token_ttl: 0,
      key: 'default',
      redirect_uris: [],
    },
    'app-1': {
      assignments: ['allow_all'],
      client_id: 'HkmsTA4GG17j0Djy4EUAB2VAyzuLVewg',
      client_secret: 'hvo_secret_g3f30MxAJWLXhhrCejbG4zY3O4LEHhEIO24aMy181AYKnfQtWTVV924ZmnlpUFUw',
      client_type: 'confidential',
      id_token_ttl: 0,
      key: 'test-key',
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

export const ASSIGNMENT_LIST_RESPONSE = {
  keys: ['allow_all', 'test-assignment'],
};

export const ASSIGNMENT_DATA_RESPONSE = {
  group_ids: ['262ca5b9-7b69-0a84-446a-303dc7d778af'],
  entity_ids: ['b6094ac6-baf4-6520-b05a-2bd9f07c66da'],
};

export const SCOPE_LIST_RESPONSE = {
  keys: ['test-scope'],
};

export const SCOPE_DATA_RESPONSE = {
  description: 'this is a test',
  template: `{
    "groups": {{identity.entity.groups.names}}
  }`,
};

export const PROVIDER_LIST_RESPONSE = {
  keys: ['test-provider'],
};

export const PROVIDER_DATA_RESPONSE = {
  allowed_client_ids: ['*'],
  issuer: '',
  scopes_supported: ['test-scope'],
};
