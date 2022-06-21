import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class OidcScopeModel extends Model {
  @attr('string') name;
  @attr('string') description;
  @attr('string', {
    label: 'JSON Template',
    editType: 'json',
  })
  template;

  @lazyCapabilities(apiPath`identity/oidc/scope/${'name'}`, 'name') scopePath;
  @lazyCapabilities(apiPath`identity/oidc/scope`) scopesPath;
  get canCreate() {
    return this.scopePath.get('canCreate');
  }
  get canRead() {
    return this.scopePath.get('canRead');
  }
  get canEdit() {
    return this.scopePath.get('canUpdate');
  }
  get canDelete() {
    return this.scopePath.get('canDelete');
  }
  get canList() {
    return this.scopesPath.get('canList');
  }
}
