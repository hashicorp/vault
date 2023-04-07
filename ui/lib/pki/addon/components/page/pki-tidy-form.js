/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

export default class PkiTidyForm extends Component {
  @service router;

  returnToConfiguration() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
  }

  @action
  updateSafetyBuffer({ goSafeTimeString }) {
    this.args.tidy.safetyBuffer = goSafeTimeString;
  }

  @task
  @waitFor
  *save(e) {
    e.preventDefault();
    try {
      yield this.args.tidy.save();
      this.returnToConfiguration();
    } catch (e) {
      const message = e.errors ? e.errors.join('. ') : e.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  cancel() {
    this.returnToConfiguration();
  }
}
