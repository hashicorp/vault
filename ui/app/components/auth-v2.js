import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task, timeout } from 'ember-concurrency';
import { capitalize } from '@ember/string';

class MfaState {
  @tracked enforcement;
  @tracked successData;
  @tracked error = '';

  reset() {
    this.enforcement = null;
    this.successData = null;
    this.error = '';
  }
}

export default class AuthV2Component extends Component {
  @service flashMessages;
  @tracked namespace;
  @tracked authType;
  @tracked mountPath;
  @tracked mfa = new MfaState();

  constructor() {
    super(...arguments);
    this.namespace = this.args.namespace || '';
    this.mountPath = this.args.mountPath || '';
    if (this.args.wrappedToken) {
      // Only the token component handles wrapped tokens
      this.authType = 'token';
    } else {
      this.authType = this.args.authType || 'token';
    }
  }

  get authMethods() {
    return ['token', 'userpass', 'oidc'];
    // return ['token', 'userpass', 'ldap', 'okta', 'jwt', 'oidc', 'radius', 'github'];
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
    if (this.args.onUpdate) {
      // Do parent side effects like update query params
      this.args.onUpdate(name, value);
    }
  }

  @action onSuccess() {
    if (this.args.onSuccess) {
      // Do parent side effects like show flash message for root token
      this.args.onSuccess();
    }
  }

  /* Multi Factor Authentication */
  @task *waitForMfa() {
    while (true) {
      yield timeout(500);
      if (this.mfa.successData) {
        return this.mfa.successData;
      }
      if (this.mfa.error) {
        throw this.mfa.error;
      }
    }
  }

  @action onMfaSuccess(response) {
    this.mfa.successData = response;
  }

  @action
  async handleData(payload) {
    if (payload.warnings) {
      payload.warnings.forEach((message) => {
        this.flashMessages.info(message);
      });
    }
    if (payload.auth.mfa_requirement) {
      // Show MFA form
      const { mfa_requirement } = this._parseMfaResponse(payload.auth.mfa_requirement);
      // this.mfa.enforcement = { mfa_requirement, backend, data };
      this.mfa.enforcement = { mfa_requirement };
      // Listen for updates from MFA form
      const authed = await this.waitForMfa.perform();
      return authed;
    }
    // TODO: wait for okta number challenge
    return payload;
  }

  @action cancelMfa() {
    this.mfa.reset();
    this.waitForMfa.cancelAll();
  }

  _parseMfaResponse(mfa_requirement) {
    // mfa_requirement response comes back in a shape that is not easy to work with
    // convert to array of objects and add necessary properties to satisfy the view
    if (mfa_requirement) {
      const { mfa_request_id, mfa_constraints } = mfa_requirement;
      const constraints = [];
      for (const key in mfa_constraints) {
        const methods = mfa_constraints[key].any;
        const isMulti = methods.length > 1;

        // friendly label for display in MfaForm
        methods.forEach((m) => {
          const typeFormatted = m.type === 'totp' ? m.type.toUpperCase() : capitalize(m.type);
          m.label = `${typeFormatted} ${m.uses_passcode ? 'passcode' : 'push notification'}`;
        });
        constraints.push({
          name: key,
          methods,
          selectedMethod: isMulti ? null : methods[0],
        });
      }

      return {
        mfa_requirement: { mfa_request_id, mfa_constraints: constraints },
      };
    }
    return {};
  }
}
