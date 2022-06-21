import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from '../../utils/field-to-attrs';
import ArrayProxy from '@ember/array/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';

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

  @attr('string', {
    label: 'Signing key',
    subText: 'Add a key to sign and verify the JSON web tokens (JWT). This cannot be edited later.',
    defaultValue: 'default-key',
    possibleValues: [],
  })
  key; // possibleValues are fetched below, in getter

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
  // if limited, expose assignment dropdown and create or add via search-select + modal
  @attr('array', {
    label: 'Assign access',
    editType: 'searchSelect',
  })
  assignments;

  @attr('string', {
    label: 'Client ID',
  })
  clientId;

  @attr('string', {
    label: 'Client Secret',
  })
  clientSecret;

  // API WIP
  @attr('string', {
    label: 'Providers',
  })
  providers;

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

  get fieldGroups() {
    const groups = [
      { default: ['name', 'clientType', 'redirectUris'] },
      { 'More options': ['key', 'idTokenTtl', 'accessTokenTtl'] },
    ];
    this.attrs.findBy('name', 'key').options.possibleValues = ArrayProxy.extend(PromiseProxyMixin).create({
      promise: this.getKeys(),
    });
    // this.attrs.key.options.possibleValues = ['default-key', ...this.getKeys()];
    return fieldToAttrs(this, groups);
  }

  get attrs() {
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

  async getKeys() {
    let keys = this.store.peekAll('oidc/key').toArray();
    if (!keys.length) {
      keys = await this.store.query('oidc/key', {});
    }
    console.log(keys, 'KEYS');
    return ['default-key', 'hello'];
  }
}
