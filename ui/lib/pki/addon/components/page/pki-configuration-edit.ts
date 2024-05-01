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
import type VersionService from 'vault/services/version';
import type PkiConfigAcmeModel from 'vault/models/pki/config/acme';
import type PkiConfigClusterModel from 'vault/models/pki/config/cluster';
import type PkiConfigCrlModel from 'vault/models/pki/config/crl';
import type PkiConfigUrlsModel from 'vault/models/pki/config/urls';
import type { FormField, TtlEvent } from 'vault/app-types';
import { addToArray } from 'vault/helpers/add-to-array';

interface Args {
  acme: PkiConfigAcmeModel;
  cluster: PkiConfigClusterModel;
  crl: PkiConfigCrlModel;
  urls: PkiConfigUrlsModel;
}
interface PkiConfigCrlTtls {
  autoRebuildGracePeriod: string;
  expiry: string;
  deltaRebuildInterval: string;
  ocspExpiry: string;
}
interface PkiConfigCrlBooleans {
  autoRebuild: boolean;
  enableDelta: boolean;
  disable: boolean;
  ocspDisable: boolean;
}

interface ErrorObject {
  modelName: string;
  message: string;
}
export default class PkiConfigurationEditComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly version: VersionService;

  @tracked invalidFormAlert = '';
  @tracked errors: Array<ErrorObject> = [];

  get isEnterprise() {
    return this.version.isEnterprise;
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    // first clear errors and sticky flash messages
    this.errors = [];
    this.flashMessages.clearMessages();

    // modelName is also the API endpoint (i.e. pki/config/cluster)
    for (const modelName of ['cluster', 'acme', 'urls', 'crl']) {
      const model = this.args[modelName as keyof Args];
      // skip saving and continue to next iteration if user does not have permission
      if (!model.canSet) continue;
      try {
        yield model.save();
        this.flashMessages.success(`Successfully updated config/${modelName}`);
      } catch (error) {
        const errorObject: ErrorObject = {
          modelName,
          message: errorMessage(error),
        };
        this.flashMessages.danger(`Error updating config/${modelName}`, { sticky: true });
        this.errors = addToArray(this.errors, errorObject);
      }
    }

    if (this.errors.length) {
      this.invalidFormAlert = 'There was an error submitting this form.';
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
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
    this.args.crl[ttlAttr as keyof PkiConfigCrlTtls] = goSafeTimeString;
    // expiry and ocspExpiry both correspond to 'disable' booleans
    // so when ttl is enabled, the booleans are set to false
    this.args.crl[attr.options.mapToBoolean as keyof PkiConfigCrlBooleans] = attr.options.isOppositeValue
      ? !enabled
      : enabled;
  }
}
