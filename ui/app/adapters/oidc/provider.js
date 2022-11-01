import NamedPathAdapter from '../named-path';

export default class OidcProviderAdapter extends NamedPathAdapter {
  pathForType() {
    return 'identity/oidc/provider';
  }
}
