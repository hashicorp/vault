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
import RouterService from '@ember/routing/router-service';
import FlashMessageService from 'vault/services/flash-messages';
import { FormField, TtlEvent } from 'vault/app-types';
import PkiCrlModel from 'vault/models/pki/crl';
import PkiUrlsModel from 'vault/models/pki/urls';
import errorMessage from 'vault/utils/error-message';

interface Args {
  crl: PkiCrlModel;
  urls: PkiUrlsModel;
}
interface PkiCrlAttrs {
  autoRebuildData: object;
  deltaCrlBuildingData: object;
  crlExpiryData: object;
  ocspExpiryData: object;
}
export default class PkiConfigurationEditComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked invalidFormAlert = '';
  @tracked errorBanner = '';

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    try {
      yield this.args.urls.save();
      yield this.args.crl.save();
      this.flashMessages.success('Successfully updated configuration');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
    } catch (error) {
      this.invalidFormAlert = 'There was an error submitting this form.';
      this.errorBanner = errorMessage(error);
    }
  }

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
  }

  @action
  handleTtl(attr: FormField, e: TtlEvent) {
    const { enabled, goSafeTimeString } = e;
    const modelAttr = attr.name;
    this.args.crl[modelAttr as keyof PkiCrlAttrs] = { enabled, duration: goSafeTimeString };
  }
}
