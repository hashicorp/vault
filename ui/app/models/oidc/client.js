import Model, { attr, hasMany } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from '../../utils/field-to-attrs';

export default class OidcClientModel extends Model {
  @hasMany('oidc/assignment') assignments;
  @attr('string', {
    label: 'Application name',
  })
  name;

  @attr('string', {
    label: 'Type',
    subText: 'Specify whether the application type is confidential or public. The public type must use PKCE.',
    editType: 'radio',
    defaultValue: 'confidential',
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
  providerIds; // might be a good candidate for @hasMany relationship instead of @attr

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

  // default form fields
  get formFields() {
    return expandAttributeMeta(this, ['name', 'clientType', 'redirectUris']);
  }

  // more options fields
  get fieldGroups() {
    return fieldToAttrs(this, [{ 'More options': ['key', 'idTokenTtl', 'accessTokenTtl'] }]);
  }
}
