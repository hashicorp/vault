/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import RouterService from '@ember/routing/router-service';
import FlashMessages from 'vault/services/flash-messages';
import PkiIssuerModel from 'vault/models/pki/issuer';

interface Args {
  model: PkiIssuerModel;
}

export default class PkiIssuerEditComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessages;

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
    const method = this.usageValues.includes(value) ? 'removeObject' : 'addObject';
    this.usageValues[method](value);
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
