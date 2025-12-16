/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { add } from 'date-fns';
import timestamp from 'vault/utils/timestamp';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { HTMLElementEvent } from 'vault/forms';
import type { KubernetesCredentials } from 'vault/secrets/kubernetes';

/**
 * @module Credentials
 * CredentialsPage component is a child component to show the generate and view
 * credentials form.
 *
 * @param {string} roleName - role name as a string
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */

interface Args {
  roleName: string;
}

export default class CredentialsPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked ttl = '';
  @tracked clusterRoleBinding = false;
  @tracked kubernetesNamespace = '';
  @tracked error = '';

  @tracked declare credentials: KubernetesCredentials;

  get leaseExpiry() {
    return add(timestamp.now(), { seconds: this.credentials.lease_duration });
  }

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles.role.details');
  }

  @action
  setKubernetesNamespace(event: HTMLElementEvent<HTMLInputElement>) {
    this.kubernetesNamespace = event.target.value;
  }

  @action
  updateTtl({ goSafeTimeString }: { goSafeTimeString: string }) {
    this.ttl = goSafeTimeString;
  }

  fetchCredentials = task(
    waitFor(async (event) => {
      event.preventDefault();
      this.error = '';

      try {
        const { currentPath } = this.secretMountPath;
        const payload = {
          kubernetes_namespace: this.kubernetesNamespace,
          cluster_role_binding: this.clusterRoleBinding,
          ttl: this.ttl,
        };
        const { lease_duration, lease_id, data } = await this.api.secrets.kubernetesGenerateCredentials(
          this.args.roleName,
          currentPath,
          payload
        );
        this.credentials = { lease_duration, lease_id, ...(data as object) } as KubernetesCredentials;
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.error = message;
      }
    })
  );
}
