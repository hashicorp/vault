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

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type RouterService from '@ember/routing/router-service';
import type { ValidationMap, Breadcrumb } from 'vault/app-types';
import type Owner from '@ember/owner';
import type KubernetesConfigForm from 'vault/forms/secrets/kubernetes/config';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  form: KubernetesConfigForm;
  breadcrumbs: Array<Breadcrumb>;
}

/**
 * @module Configure
 * ConfigurePage component is a child component to configure kubernetes secrets engine.
 *
 * @param {object} form - config form that contains kubernetes configuration
 */
export default class ConfigurePageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked inferredState: 'success' | 'error' | null = null;
  @tracked modelValidations: ValidationMap | null = null;
  @tracked alert = '';
  @tracked error = '';
  @tracked showConfirm = false;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    const { form } = this.args;
    if (!form.isNew && !form.data.disable_local_ca_jwt) {
      this.inferredState = 'success';
    }
  }

  get isDisabled() {
    if (!this.args.form.data.disable_local_ca_jwt && this.inferredState !== 'success') {
      return true;
    }
    return this.save.isRunning || this.fetchInferred.isRunning;
  }

  leave(route: string) {
    this.router.transitionTo(`vault.cluster.secrets.backend.kubernetes.${route}`);
  }

  @action
  onRadioSelect(value: boolean) {
    this.args.form.data.disable_local_ca_jwt = value;
    this.inferredState = null;
  }

  fetchInferred = task(
    waitFor(async () => {
      try {
        await this.api.secrets.kubernetesCheckConfiguration(this.secretMountPath.currentPath);
        this.inferredState = 'success';
      } catch {
        this.inferredState = 'error';
      }
    })
  );

  save = task(
    waitFor(async () => {
      const { form } = this.args;

      if (!form.isNew && !this.showConfirm) {
        this.showConfirm = true;
        return;
      }
      this.showConfirm = false;
      this.modelValidations = null;
      this.alert = '';

      const { isValid, state, invalidFormMessage, data } = form.toJSON();
      if (isValid) {
        try {
          await this.api.secrets.kubernetesConfigure(this.secretMountPath.currentPath, data);
          this.flashMessages.success('Successfully configured Kubernetes engine');
          this.leave('configuration');
        } catch (error) {
          const { message } = await this.api.parseError(
            error,
            'Error saving configuration. Please try again or contact support'
          );
          this.error = message;
        }
      } else {
        this.modelValidations = state;
        this.alert = invalidFormMessage;
      }
    })
  );

  @action
  cancel() {
    const { isNew } = this.args.form;
    const transitionRoute = isNew ? 'overview' : 'configuration';
    this.leave(transitionRoute);
  }
}
