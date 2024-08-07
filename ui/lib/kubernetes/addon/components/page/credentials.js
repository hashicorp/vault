/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { add } from 'date-fns';
import errorMessage from 'vault/utils/error-message';
import timestamp from 'vault/utils/timestamp';
/**
 * @module Credentials
 * CredentialsPage component is a child component to show the generate and view
 * credentials form.
 *
 * @param {string} roleName - role name as a string
 * @param {string} backend - backend as a string
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */
export default class CredentialsPageComponent extends Component {
  @service store;
  @service router;

  @tracked ttl = '';
  @tracked clusterRoleBinding = false;
  @tracked kubernetesNamespace;
  @tracked error;

  @tracked credentials;

  get leaseExpiry() {
    return add(timestamp.now(), { seconds: this.credentials.lease_duration });
  }

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles.role.details');
  }

  @action
  setKubernetesNamespace({ target }) {
    this.kubernetesNamespace = target.value;
  }

  @action
  updateTtl({ goSafeTimeString }) {
    this.ttl = goSafeTimeString;
  }

  @task
  @waitFor
  *fetchCredentials(event) {
    event.preventDefault();
    try {
      const payload = {
        role: this.args.roleName,
        kubernetes_namespace: this.kubernetesNamespace,
        cluster_role_binding: this.clusterRoleBinding,
        ttl: this.ttl,
      };

      this.credentials = yield this.store
        .adapterFor('kubernetes/role')
        .generateCredentials(this.args.backend, payload);
    } catch (error) {
      this.error = errorMessage(error);
    }
  }
}
