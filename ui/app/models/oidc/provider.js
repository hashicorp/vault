import Model, { attr, hasMany } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class OidcProviderModel extends Model {
  @hasMany('identity/oidc/scope') scopes_supported;
  @attr('string') name;
  @attr('string', {
    label: 'Issuer URL',
  })
  issuer;

  @attr('array', {
    editType: 'searchSelect',
  })
  scopesSupported;

  // form has a radio option to allow_all, or limit access to selected 'application'
  // if limited, expose assignment dropdown and add via search-select
  @attr('array', {
    editType: 'searchSelect',
  })
  allowedClientIds;

  @lazyCapabilities(apiPath`identity/oidc/provider/${'name'}`, 'name') providerPath;
  @lazyCapabilities(apiPath`identity/oidc/provider`) providersPath;
  get canCreate() {
    return this.providerPath.get('canCreate');
  }
  get canRead() {
    return this.providerPath.get('canRead');
  }
  get canEdit() {
    return this.providerPath.get('canUpdate');
  }
  get canDelete() {
    return this.providerPath.get('canDelete');
  }
  get canList() {
    return this.providersPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/oidc/client`) clientsPath;
  get canListClients() {
    return this.clientsPath.get('canList');
  }
  @lazyCapabilities(apiPath`identity/oidc/scope`) scopesPath;
  get canListScopes() {
    return this.scopesPath.get('canList');
  }
}
