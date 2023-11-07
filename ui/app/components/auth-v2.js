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

  // calcRedirect(mountPath, namespace) {
  //   let url = `${window.location.origin}/ui/auth/${mountPath}/oidc/callback`;
  //   console.log({ url });
  //   if (namespace) {
  //     url += `?namespace=${namespace}`;
  //   }
  //   return url;
  // }

  // @task
  // *fetchRoles(mountPath, namespace) {
  //   const path = this.selectedAuthPath || this.selectedAuthType;
  //   // const id = JSON.stringify([path, this.roleName]);
  //   const url = `/v1/auth/${mountPath}/oidc/auth_url`;
  //   const redirect_uri = calcRedirect(mountPath, namespace);
  //   const options = {
  //     method: 'POST',
  //     body: JSON.stringify({
  //       role: this.form.role,
  //       redirect_uri,
  //     }),
  //   };
  //   try {
  //     const role = yield fetch(url, options);
  //     // this.set('role', role);
  //   } catch (e) {
  //     console.log('error fetching roles', e);
  //     this.error = 'Could not fetch roles';
  //   }
  // }

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
      this.form.resetFields();
      this.mountPath = '';
    }
  }

  @action
  async handleLogin(evt) {
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
}
