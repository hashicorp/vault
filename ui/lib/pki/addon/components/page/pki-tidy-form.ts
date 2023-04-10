/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';

import PkiTidyModel from 'vault/models/pki/tidy';
import RouterService from '@ember/routing/router-service';

interface Args {
  tidy: PkiTidyModel;
}

export default class PkiTidyForm extends Component<Args> {
  @service declare readonly router: RouterService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  returnToConfiguration() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
  }

  @action
  updateSafetyBuffer({ goSafeTimeString }: { goSafeTimeString: string }) {
    this.args.tidy.safetyBuffer = goSafeTimeString;
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    try {
      yield this.args.tidy.save();
      this.returnToConfiguration();
    } catch (e) {
      this.errorBanner = errorMessage();
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  cancel() {
    this.returnToConfiguration();
  }
}
