import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class AuthV2Component extends Component {
  @service permissions;

  @tracked namespace = '';
  @tracked authType = 'token';
  @tracked mountPath = '';

  get authMethods() {
    return ['token', 'userpass', 'ldap', 'okta', 'jwt', 'oidc', 'radius', 'github'];
  }

  @action
  handleChange(evt) {
    // For changing values in this backing class, not on form
    const { name, value } = evt.target;
    this[name] = value;
    if (name === 'authType') {
      // if the authType changes, reset the mount path
      this.mountPath = '';
    }
  }

  @action onSuccess() {
    this.permissions.getPaths.perform();
    // TODO: show flash message if root
  }
}
