/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class PkiTidyForm extends Component {
  @service secretMountPath;
  @service router;

  @tracked tidyCertificateStore = false;
  @tracked tidyCertificateRevocationQueue = false;
  @tracked safetyBuffer;

  returnToConfiguration() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
  }

  @action
  updateSafetyBuffer({ goSafeTimeString }) {
    this.args.tidy.safetyBuffer = goSafeTimeString;
  }

  @action
  async performTidy() {
    try {
      await this.args.tidy.save();
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
