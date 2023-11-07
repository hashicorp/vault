import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { allSupportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import errorMessage from 'vault/utils/error-message';

class AuthState {
  @tracked type = '';
  @tracked token = '';
  @tracked username = '';
  @tracked password = '';
  @tracked role = '';
  @tracked jwt = '';

  resetFields() {
    this.token = '';
    this.userame = '';
    this.password = '';
    this.role = '';
    this.jwt = '';
  }

  constructor(type) {
    this.type = type || 'token';
  }
}
export default class AuthV2Component extends Component {
  @service permissions;
  @service session;

  @tracked namespace = '';
  @tracked authType = 'token';
  @tracked mountPath = '';
  @tracked error = '';
  // Two ways of doing things
  @tracked form = new AuthState();
  @tracked fields = {};

  get authMethods() {
    return ['token', 'userpass', 'oidc'];
    // return ['token', 'userpass', 'ldap', 'okta', 'jwt', 'oidc', 'radius', 'github'];
  }
  get showFields() {
    const backend = allSupportedAuthBackends().findBy('type', this.authType);
    return backend.formAttributes;
  }

  maybeMask = (field) => {
    if (field === 'token' || field === 'password') {
      return 'password';
    }
    return 'text';
  };

  @action
  handleFormChange(evt) {
    this.form[evt.target.name] = evt.target.value;
  }

  @action
  handleChange(evt) {
    // For changing values in this backing class, not on form
    const { name, value } = evt.target;
    this[name] = value;
    if (name === 'authType') {
      // if the authType changes, reset the form and mount path
      // this.form.resetFields();
      this.mountPath = '';
    }
  }

  @action
  async authenticate(fields) {
    const authenticator = `authenticator:${this.authType}`;
    try {
      await this.session.authenticate(authenticator, fields, {
        backend: this.mountPath,
        namespace: this.namespace,
      });
    } catch (e) {
      this.error = errorMessage(e);
    }

    if (this.session.isAuthenticated) {
      // TODO: Show root warning flash message
      this.permissions.getPaths.perform();
    }
  }

  @action
  async handleFormLogin(evt) {
    evt.preventDefault();
    const authenticator = `authenticator:${this.authType}`;
    const fields = this.showFields.reduce((obj, field) => {
      console.log({ field, val: this.form[field] });
      obj[field] = this.form[field];
      return obj;
    }, {});
    console.log({ fields });
    try {
      await this.session.authenticate(authenticator, fields, {
        backend: this.mountPath,
        namespace: this.namespace,
      });
    } catch (e) {
      this.error = errorMessage(e);
    }

    if (this.session.isAuthenticated) {
      // TODO: Show root warning flash message
      this.permissions.getPaths.perform();
    }
  }

  @action onSuccess() {
    this.permissions.getPaths.perform();
    // TODO: show flash message if root
  }
}
