import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller, { inject as controller } from '@ember/controller';
import { task, timeout } from 'ember-concurrency';

export default Controller.extend({
  vaultController: controller('vault'),
  clusterController: controller('vault.cluster'),
  namespaceService: service('namespace'),
  featureFlagService: service('featureFlag'),
  auth: service(),
  router: service(),

  queryParams: [{ authMethod: 'with', oidcProvider: 'o' }],

  namespaceQueryParam: alias('clusterController.namespaceQueryParam'),
  wrappedToken: alias('vaultController.wrappedToken'),
  redirectTo: alias('vaultController.redirectTo'),
  managedNamespaceRoot: alias('featureFlagService.managedNamespaceRoot'),

  authMethod: '',
  oidcProvider: '',

  get managedNamespaceChild() {
    let fullParam = this.namespaceQueryParam;
    let split = fullParam.split('/');
    if (split.length > 1) {
      split.shift();
      return `/${split.join('/')}`;
    }
    return '';
  },

  updateManagedNamespace: task(function* (value) {
    // debounce
    yield timeout(500);
    // TODO: Move this to shared fn
    const newNamespace = `${this.managedNamespaceRoot}${value}`;
    this.namespaceService.setNamespace(newNamespace, true);
    this.set('namespaceQueryParam', newNamespace);
  }).restartable(),

  updateNamespace: task(function* (value) {
    // debounce
    yield timeout(500);
    this.namespaceService.setNamespace(value, true);
    this.set('namespaceQueryParam', value);
  }).restartable(),

  authSuccess({ isRoot, namespace }) {
    let transition;
    if (this.redirectTo) {
      // here we don't need the namespace because it will be encoded in redirectTo
      transition = this.router.transitionTo(this.redirectTo);
      // reset the value on the controller because it's bound here
      this.set('redirectTo', '');
    } else {
      transition = this.router.transitionTo('vault.cluster', { queryParams: { namespace } });
    }
    transition.followRedirects().then(() => {
      if (isRoot) {
        this.flashMessages.warning(
          'You have logged in with a root token. As a security precaution, this root token will not be stored by your browser and you will need to re-authenticate after the window is closed or refreshed.'
        );
      }
    });
  },

  actions: {
    onAuthResponse(authResponse, backend, data) {
      const { mfa_requirement } = authResponse;
      // if an mfa requirement exists further action is required
      if (mfa_requirement) {
        this.set('mfaAuthData', { mfa_requirement, backend, data });
      } else {
        this.authSuccess(authResponse);
      }
    },
    onMfaSuccess(authResponse) {
      this.authSuccess(authResponse);
    },
    onMfaErrorDismiss() {
      this.setProperties({
        mfaAuthData: null,
        mfaErrors: null,
      });
    },
  },
});
