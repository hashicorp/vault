import NamedPathAdapter from '../named-path';

export default class OidcScopeAdapter extends NamedPathAdapter {
  pathForType() {
    return 'identity/oidc/scope';
  }
}
