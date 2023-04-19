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

interface Args {
  crl: PkiCrlModel;
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

  @tracked invalidFormAlert = null;
  @tracked errorMessage = null;

  @task
  @waitFor
  *save(event: Event) {
    yield event.preventDefault(); // remove yield
    // do something
  }

  @action
  cancel() {
    // handle cancel
  }

  @action
  handleTtl(attr: FormField, e: TtlEvent) {
    const { enabled, goSafeTimeString } = e;
    const ttlAttr = attr.name;
    this.args.crl[ttlAttr as keyof PkiCrlTtls] = goSafeTimeString;
    this.args.crl[attr.options.booleanBuddy as keyof PkiCrlBooleans] = enabled;
  }
}
