import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class OidcProviderModel extends Model {
  @attr('string') name;
  @attr('string', { label: 'Issuer UR' }) issuer;
  @attr('array', { label: 'Supported scopes', editType: 'searchSelect' }) scopesSupported;
  @attr('array', { label: 'Allowed applications' }) allowedClientIds; // no editType because does not use form-field component

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
