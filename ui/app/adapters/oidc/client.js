import NamedPathAdapter from '../named-path';

export default class OidcClientAdapter extends NamedPathAdapter {
  pathForType() {
    return 'identity/oidc/client';
  }
}
