import NamedPathAdapter from '../named-path';

export default class OidcAssignmentAdapter extends NamedPathAdapter {
  pathForType() {
    return 'identity/oidc/assignment';
  }
}
