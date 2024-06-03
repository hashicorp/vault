/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type PkiIssuerModel from 'vault/models/pki/issuer';
import { removeFromArray } from 'vault/helpers/remove-from-array';
import { addToArray } from 'vault/helpers/add-to-array';

interface Args {
  model: PkiIssuerModel;
}

export default class PkiIssuerEditComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked usageValues: Array<string> = [];
  @tracked error = null;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.usageValues = (this.args.model.usage || '').split(',');
  }

  toDetails() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.issuers.issuer.details');
  }

  @action
  setUsage(value: string) {
    if (this.usageValues.includes(value)) {
      this.usageValues = removeFromArray(this.usageValues, value);
    } else {
      this.usageValues = addToArray(this.usageValues, value);
    }
    this.args.model.usage = this.usageValues.join(',');
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    try {
      yield this.args.model.save();
      this.flashMessages.success('Successfully updated issuer');
      this.toDetails();
    } catch (error) {
      this.error = errorMessage(error);
    }
  }

  @action
  cancel() {
    this.args.model.rollbackAttributes();
    this.toDetails();
  }
}
