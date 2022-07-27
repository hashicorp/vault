import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from '../../utils/field-to-attrs';

export default class OidcClientModel extends Model {
  @attr('string', { label: 'Application name' }) name;
  @attr('string', {
    label: 'Type',
    subText: 'Specify whether the application type is confidential or public. The public type must use PKCE.',
    editType: 'radio',
    possibleValues: ['confidential', 'public'],
  })
  clientType;

  @attr('array', {
    label: 'Redirect URIs',
    subText:
      'One of these values must exactly match the redirect_uri parameter value used in each authentication request.',
    editType: 'stringArray',
  })
  redirectUris;

  // >> MORE OPTIONS TOGGLE <<

  @attr('string', {
    label: 'Signing key',
    subText: 'Add a key to sign and verify the JSON web tokens (JWT). This cannot be edited later.',
    editType: 'searchSelect',
    editDisabled: true,
    disallowNewItems: true,
    defaultValue() {
      return ['default'];
    },
    fallbackComponent: 'input-search',
    selectLimit: 1,
    models: ['oidc/key'],
  })
  key;

  @attr({
    label: 'Access Token TTL',
    editType: 'ttl',
    helperTextDisabled: 'Vault will use the default lease duration',
    setDefault: false,
  })
  accessTokenTtl;

  @attr({
    label: 'ID Token TTL',
    editType: 'ttl',
    helperTextDisabled: 'Vault will use the default lease duration',
    setDefault: false,
  })
  idTokenTtl;

  // >> END MORE OPTIONS TOGGLE <<

  @attr('array', { label: 'Assign access' }) assignments; // no editType because does not use form-field component
  @attr('string', { label: 'Client ID' }) clientId;
  @attr('string') clientSecret;

  // TODO API WIP - attr TBD, current PR proposes the following:
  // $ curl -H "X-Vault-Token: ..." -X LIST http://127.0.0.1:8200/v1/identity/oidc/provider?allowed_client_id="<client_id>"
  // add to model on route? or query from tab in UI?
  @attr('string', { label: 'Providers' }) providers;

  // CAPABILITIES //
  @lazyCapabilities(apiPath`identity/oidc/client/${'name'}`, 'name') clientPath;
  @lazyCapabilities(apiPath`identity/oidc/client`) clientsPath;
  get canCreate() {
    return this.clientPath.get('canCreate');
  }
  get canRead() {
    return this.clientPath.get('canRead');
  }
  get canEdit() {
    return this.clientPath.get('canUpdate');
  }
  get canDelete() {
    return this.clientPath.get('canDelete');
  }
  get canList() {
    return this.clientsPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/oidc/key`) keysPath;
  get canListKeys() {
    return this.keysPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/oidc/assignment/${'name'}`, 'name') assignmentPath;
  @lazyCapabilities(apiPath`identity/oidc/assignment`) assignmentsPath;
  get canCreateAssignments() {
    return this.assignmentPath.get('canCreate');
  }
  get canListAssignments() {
    return this.assignmentsPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/oidc/${'name'}/provider`, 'backend', 'name') clientProvidersPath; // API is WIP
  // API WIP
  get canListProviders() {
    return this.clientProvidersPath.get('canList');
  }

  // fieldGroups was behaving buggy so may not use
  get fieldGroups() {
    const groups = [
      { default: ['name', 'clientType', 'redirectUris'] },
      { 'More options': ['key', 'idTokenTtl', 'accessTokenTtl'] },
    ];
    return fieldToAttrs(this, groups);
  }

  // WIP
  get fieldAttrs() {
    return expandAttributeMeta(this, [
      'name',
      'clientType',
      'redirectUris',
      'key',
      'idTokenTtl',
      'accessTokenTtl',
      'assignments',
    ]);
  }
}
