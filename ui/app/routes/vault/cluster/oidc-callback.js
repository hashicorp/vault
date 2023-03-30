/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

export function getParamsForCallback(qp, searchString) {
  const queryString = decodeURIComponent(searchString);
  let { path, code, state, namespace } = qp;
  // namespace from state takes precedence over the cluster's ns
  if (state?.includes(',ns=')) {
    [state, namespace] = state.split(',ns=');
  }
  // some SSO providers do not return a url-encoded state param
  // check for namespace using URLSearchParams instead of paramsFor
  const urlParams = new URLSearchParams(queryString);
  const checkState = urlParams.get('state');
  if (checkState?.includes(',ns=')) {
    [state, namespace] = checkState.split(',ns=');
  }
  path = window.decodeURIComponent(path);
  const payload = { source: 'oidc-callback', path: path || '', code: code || '', state: state || '' };
  if (namespace) {
    payload.namespace = namespace;
  }
  return payload;
}

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  afterModel() {
    const { auth_path: path, code, state } = this.paramsFor(this.routeName);
    const { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    const queryString = window.location.search;
    const payload = getParamsForCallback({ path, code, state, namespace }, queryString);
    window.opener.postMessage(payload, window.origin);
  },
  setupController(controller) {
    this._super(...arguments);
    controller.set('pageContainer', document.querySelector('.page-container'));
  },
});
