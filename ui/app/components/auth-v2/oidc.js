import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { withAuthForm } from 'vault/decorators/auth-form';
import parseURL from 'core/utils/parse-url';
import { DOMAIN_STRINGS, PROVIDER_WITH_LOGO } from 'vault/models/role-jwt';
import errorMessage from 'vault/utils/error-message';
import { next } from '@ember/runloop';

const WAIT_TIME = 500;
const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
// const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';
@withAuthForm('oidc')
export default class AuthV2OidcComponent extends Component {
  @tracked roleInfo;

  constructor() {
    super(...arguments);
    next(() => {
      // do it after decorator is setup, so mountPath is available
      this.fetchRoles.perform();
    });
  }

  get redirectUrl() {
    let url = `${window.location.origin}/ui/vault/auth/${this.mountPath}/oidc/callback`;
    if (this.args.namespace) {
      url += `?namespace=${this.args.namespace}`;
    }
    return url;
  }
  @task
  *fetchRoles() {
    this.error = '';
    const { namespace } = this.args;
    const url = `/v1/auth/${this.mountPath}/oidc/auth_url`;
    const options = {
      method: 'POST',
      body: JSON.stringify({
        role: this.state.role,
        redirect_uri: this.redirectUrl,
      }),
    };
    if (namespace) {
      options.headers = {
        ['X-Vault-Namespace']: namespace,
      };
    }
    try {
      const response = yield fetch(url, options);
      const body = yield response.json();
      if (response.status === 403) {
        this.roleInfo = null;
        return;
      } else if (response.status !== 200) {
        this.roleInfo = null;
        throw new Error(
          response.httpStatus === 400
            ? 'Invalid role. Please try again.'
            : `Error fetching role: ${body.errors.join(', ')}`
        );
      }
      const roleInfo = this.roleFromUrl(body.data.auth_url);
      this.roleInfo = roleInfo;
    } catch (e) {
      this.roleInfo = null;
      this.error = errorMessage(e);
    }
  }
  roleFromUrl(authUrl) {
    const { hostname } = parseURL(authUrl);
    const firstMatch = Object.keys(DOMAIN_STRINGS).find((name) => hostname.includes(name));
    const providerName = DOMAIN_STRINGS[firstMatch] || null;
    const providerIcon = PROVIDER_WITH_LOGO.includes(providerName) ? providerName.toLowerCase() : null;
    return {
      authUrl,
      providerIcon,
      providerName,
    };
  }
  @action
  async updateRole(evt) {
    this.state.role = evt.target.value;
    await this.fetchRoles.perform();
  }

  @action handleLogin(evt) {
    evt.preventDefault();
    if (!this.roleInfo?.authUrl) {
      this.error =
        'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.';
      return;
    }
    const win = this.authWindow;
    const POPUP_WIDTH = 500;
    const POPUP_HEIGHT = 600;
    const left = win.screen.width / 2 - POPUP_WIDTH / 2;
    const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
    // Initiate popup
    const oidcWindow = win.open(
      this.roleInfo.authUrl,
      'vaultOIDCWindow',
      `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
    );
    // Watch popup & current window
    this.watchPopup.perform(oidcWindow);
    this.watchCurrent.perform(oidcWindow);
    // Kick off OIDC flow
    this.oidcFlow.perform(oidcWindow);
  }

  // OIDC flow
  get authWindow() {
    return this.window || window;
  }

  @task *oidcFlow(oidcWindow) {
    const win = this.authWindow;
    while (true) {
      const event = yield waitForEvent(win, 'message');
      if (event.origin === win.origin && event.isTrusted && event.data.source === 'oidc-callback') {
        return this.exchangeOIDC.perform(event.data, oidcWindow);
      }
      // continue to wait for the correct message
    }
  }

  @task *exchangeOIDC(oidcState, oidcWindow) {
    if (oidcState === null || oidcState === undefined) {
      // TODO: should show an error?
      return;
    }

    let { namespace, path, state, code } = oidcState;
    if (!path || !state || !code) {
      return this.handleOIDCError(oidcWindow, ERROR_MISSING_PARAMS);
    }

    if (namespace === '') {
      const i = state.indexOf(',ns=');
      if (i >= 0) {
        // ",ns=" is 4 characters
        namespace = state.substring(i + 4);
        state = state.substring(0, i);
      }
    }

    try {
      yield this.session.authenticate(
        `authenticator:${this._type}`,
        { state, code, oidcWindow },
        { backend: this.mountPath, namespace: this.namespace }
      );
      // resp = yield adapter.exchangeOIDC(path, state, code);
      this.closeWindow(oidcWindow);
    } catch (e) {
      // If there was an error on Vault's end, close the popup
      // and show the error on the login screen
      return this.handleOIDCError(oidcWindow, errorMessage(e));
    }
    // yield this.onSubmit(null, null, resp.auth.client_token);
  }

  @task *watchPopup(oidcWindow) {
    while (true) {
      yield timeout(WAIT_TIME);
      if (!oidcWindow || oidcWindow.closed) {
        return this.handleOIDCError(oidcWindow, ERROR_WINDOW_CLOSED);
      }
    }
  }
  @task *watchCurrent(oidcWindow) {
    // when user is about to change pages, close the popup window
    yield waitForEvent(this.authWindow, 'beforeunload');
    oidcWindow.close();
  }

  closeWindow(oidcWindow) {
    // Removes listeners and then closes the window
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    oidcWindow?.close();
  }
  handleOIDCError(oidcWindow, errorMessage) {
    // close window
    this.closeWindow(oidcWindow);
    // cancel waiting for response
    this.oidcFlow.cancelAll();
    // show error
    this.error = errorMessage;
  }
}
