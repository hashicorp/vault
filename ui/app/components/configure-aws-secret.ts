/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

import type LeaseConfigModel from 'vault/models/aws/lease-config';
import type RootConfigModel from 'vault/models/aws/root-config';
import type Router from '@ember/routing/router';
import type Store from '@ember-data/store';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ConfigureAwsSecretComponent is used to configure the AWS secret engine
 * A user can configure the endpoint root/config and/or lease/config. 
 * The fields for these endpoints are on one form.
 *
 * @example
 * ```js
 * <ConfigureAwsSecret
    @rootConfig={{this.model.rootConfig}}
    @leaseConfig={{this.model.leaseConfig}} />
 * ```
 *
 * @param {object} awsRootConfig - AWS secret engine config/root model
 * @param {object} awsLeaseConfig - AWS secret engine config/lease model
 * @param {string} id - name of the secret engine, ex: 'aws-123'
 */

interface Args {
  leaseConfig: LeaseConfigModel;
  rootConfig: RootConfigModel;
  id: string;
}

export default class ConfigureAwsSecretComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessageRoot: string | null = null;
  @tracked invalidFormAlert: string | null = null;

  @task
  *save(event: Event) {
    event.preventDefault();
    this.resetErrors();
    const { id, leaseConfig, rootConfig } = this.args;

    try {
      yield rootConfig.save();
      this.flashMessages.success(`Successfully saved ${id}'s root configuration.`);
    } catch (error) {
      this.errorMessageRoot = errorMessage(error);
    }

    try {
      yield leaseConfig.save();
      this.flashMessages.success(`Successfully saved ${id}'s lease configuration.`);
    } catch (error) {
      // if lease config fails, notify with a flash message but still allow users to save the root config.
      if (!this.errorMessageRoot) {
        this.flashMessages.danger(`Lease configuration was not saved: ${errorMessage(error)}`, {
          sticky: true,
        });
      } else {
        this.flashMessages.danger(
          `Configuration not saved: ${errorMessage(error)}. ${this.errorMessageRoot}`,
          {
            sticky: true,
          }
        );
      }
    }
    // prevent transition if there are errors with root configuration
    if (this.errorMessageRoot) {
      this.invalidFormAlert = 'There was an error submitting this form.';
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.configuration', id);
    }
  }

  @action
  onCancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', this.args.id);
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessageRoot = null;
    this.invalidFormAlert = null;
  }
}
