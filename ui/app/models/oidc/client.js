import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from '../../utils/field-to-attrs';

export default class OidcClientModel extends Model {
  @attr('string', {
    label: 'Application name',
  })
  name;

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

  // might be a good candidate for using @belongsTo relationship?
  @attr('string', {
    label: 'Signing key',
    subText: 'Add a key to sign and verify the JSON web tokens (JWT). This cannot be edited later.',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['oidc/key'],
  })
  key;

  @attr({
    label: 'Access Token TTL',
    editType: 'ttl',
    helperTextDisabled: 'Vault will use the default lease duration',
  })
  accessTokenTtl;

  @attr({
    label: 'ID Token TTL',
    editType: 'ttl',
    helperTextDisabled: 'Vault will use the default lease duration',
  })
  idTokenTtl;

  // >> END MORE OPTIONS TOGGLE <<

  // form has a radio option to allow_all (default), or limit access to selected 'assignment'
  // if limited, expose search-select to select or create new and add via modal
  @attr('array', {
    label: 'Assign access',
    editType: 'searchSelect',
  })
  assignments; // might be a good candidate for @hasMany relationship instead of @attr

  @attr('string', {
    label: 'Client ID',
  })
  clientId;

  @attr('string', {
    label: 'Client Secret',
  })
  clientSecret;

  // API WIP - param TBD
  @attr('string', {
    label: 'Providers',
  })
  provider_ds; // might be a good candidate for @hasMany relationship instead of @attr

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
