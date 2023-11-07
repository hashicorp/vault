import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

export default class AuthV2OidcComponent extends Component {
  fields = ['role', 'jwt'];
  @tracked roleName = '';
  @tracked jwt = '';
  @tracked role;

  constructor() {
    super(...arguments);
    this.fetchRoles.perform();
  }

  get mountPath() {
    return this.args.mountPath || 'oidc';
  }

  get redirectUrl() {
    let url = `${window.location.origin}/ui/auth/${this.args.mountPath}/oidc/callback`;
    console.log({ url });
    if (this.args.namespace) {
      url += `?namespace=${this.args.namespace}`;
    }
    return url;
  }

  @task
  *fetchRoles() {
    const { namespace } = this.args;
    console.log('fetching roles', namespace);
    // const path = this.selectedAuthPath || this.selectedAuthType;
    // const id = JSON.stringify([path, this.roleName]);
    const url = `/v1/auth/${this.mountPath}/oidc/auth_url`;
    const options = {
      method: 'POST',
      body: JSON.stringify({
        role: this.roleName,
        redirect_uri: this.redirectUrl,
      }),
    };
    if (namespace) {
      options.headers['X-Vault-Namespace'] = namespace;
    }
    try {
      const role = yield fetch(url, options);
      // this.set('role', role);
    } catch (e) {
      console.log('error fetching roles', e);
      this.error = 'Could not fetch roles';
    }
  }

  @action
  async updateRole(evt) {
    this.roleName = evt.target.value;
    await this.fetchRoles.perform();
  }

  @action handleLogin(evt) {
    evt.preventDefault();
  }
}
