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
// TYPES
import RouterService from '@ember/routing/router-service';
import FlashMessageService from 'vault/services/flash-messages';
import VersionService from 'vault/services/version';
import PkiCrlModel from 'vault/models/pki/crl';
import PkiUrlsModel from 'vault/models/pki/urls';
import { FormField, TtlEvent } from 'vault/app-types';

interface Args {
  crl: PkiCrlModel;
  urls: PkiUrlsModel;
}
interface PkiCrlTtls {
  autoRebuildGracePeriod: string;
  expiry: string;
  deltaRebuildInterval: string;
  ocspExpiry: string;
}
interface PkiCrlBooleans {
  autoRebuild: boolean;
  enableDelta: boolean;
  disable: boolean;
  ocspDisable: boolean;
}
export default class PkiConfigurationEditComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly version: VersionService;

  @tracked invalidFormAlert = '';
  @tracked errorBanner = '';

  get isEnterprise() {
    return this.version.isEnterprise;
  }

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
    const ttlAttr = attr.name;
    this.args.crl[ttlAttr as keyof PkiCrlTtls] = goSafeTimeString;
    // expiry and ocspExpiry both correspond to 'disable' booleans
    // so when ttl is enabled, the booleans are set to false
    this.args.crl[attr.options.mapToBoolean as keyof PkiCrlBooleans] = attr.options.isOppositeValue
      ? !enabled
      : enabled;
  }
}
