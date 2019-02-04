import { inject as service } from '@ember/service';
import Component from './outer-html';
import { next, later } from '@ember/runloop';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { computed } from '@ember/object';

export default Component.extend({
  store: service(),
  selectedAuthPath: null,
  roleName: null,
  role: null,

  didReceiveAttrs() {
    next(() => {
      let { oldSelectedAuthPath, selectedAuthPath } = this;
      if (oldSelectedAuthPath !== selectedAuthPath) {
        this.set('role', null);
        this.onRoleName(null);
        this.fetchRole.perform(null, { nodebounce: true });
      }
      this.set('oldSelectedAuthPath', selectedAuthPath);
    });
  },

  // OIDC roles in the JWT/OIDC backend are those with an authUrl,
  // those that are JWT type will 400 when trying to fetch the role
  isOIDC: computed('role', function() {
    return this.role && this.role.authUrl;
  }),

  getWindow() {
    return this.window || window;
  },

  fetchRole: task(function*(roleName, options = {}) {
    if (!options.nodebounce) {
      this.onRoleName(roleName);
      // debounce
      yield timeout(500);
    }
    let path = this.selectedAuthPath || 'jwt';
    let id = JSON.stringify([path, roleName]);
    let role = null;
    try {
      role = yield this.store.findRecord('role-jwt', id);
    } catch (e) {
      if (!e.httpStatus || e.httpStatus !== 400) {
        throw e;
      }
    }
    this.set('role', role);
  }).restartable(),

  handleOIDCError(err) {
    this.onError(err);
  },

  prepareForOIDC: task(function*(oidcWindow) {
    this.waitForClose.perform(oidcWindow);
    let storageEvent = yield waitForEvent(this.getWindow(), 'storage');
    this.exchangeOIDC.perform(storageEvent, oidcWindow);
  }),

  waitForClose: task(function*(oidcWindow) {
    while (true) {
      yield timeout(500);
      if (!oidcWindow || oidcWindow.closed) {
        return this.handleOIDCError('windowClosed');
      }
    }
  }),

  closeWindow(oidcWindow) {
    this.waitForClose.cancelAll();
    oidcWindow.close();
  },

  exchangeOIDC: task(function*(event, oidcWindow) {
    if (event.key !== 'oidcState') {
      return;
    }
    let { namespace, path, state, code } = JSON.parse(event.newValue);
    this.getWindow().localStorage.removeItem('oidcState');
    later(() => {
      this.closeWindow(oidcWindow);
    }, 500);
    if (!path || !state || !code) {
      return this.handleOIDCError('missingParams');
    }
    let adapter = this.store.adapterFor('auth-method');
    // this might be bad to mutate the outer state
    this.onNamespace(namespace);
    let resp = yield adapter.exchangeOIDC(path, state, code);
    let token = resp.auth.client_token;
    this.onSelectedAuth('token');
    this.onToken(token);
    yield this.onSubmit();
  }),

  actions: {
    startOIDCAuth(e) {
      if (e && e.preventDefault) {
        e.preventDefault();
      }
      if (!this.isOIDC) {
        return;
      }
      let win = this.getWindow();

      let left = win.screen.width / 2 - 250;
      let top = win.screen.height / 2 - 300;
      let oidcWindow = win.open(
        this.role.authUrl,
        'vaultOIDCWindow',
        `width=500,height=600,resizable,scrollbars=yes,top=${top},left=${left}`
      );
      this.prepareForOIDC.perform(oidcWindow);
    },
  },
});
