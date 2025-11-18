/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { addToArray } from 'vault/helpers/add-to-array';
import { capitalize } from '@ember/string';

import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type VersionService from 'vault/services/version';
import type PkiConfigAcmeForm from 'vault/forms/secrets/pki/config/acme';
import type PkiConfigClusterForm from 'vault/forms/secrets/pki/config/cluster';
import type PkiConfigCrlForm from 'vault/forms/secrets/pki/config/crl';
import type PkiConfigUrlsForm from 'vault/forms/secrets/pki/config/urls';
import type { FormField, TtlEvent } from 'vault/app-types';
import type ApiService from 'vault/services/api';

interface Args {
  acmeForm: PkiConfigAcmeForm;
  clusterForm: PkiConfigClusterForm;
  crlForm: PkiConfigCrlForm;
  urlsForm: PkiConfigUrlsForm;
  backend: string;
  capabilities: Record<string, boolean>;
}
interface PkiConfigCrlTtls {
  auto_rebuild_grace_period: string;
  expiry: string;
  delta_rebuild_interval: string;
  ocsp_expiry: string;
}
interface PkiConfigCrlBooleans {
  auto_rebuild: boolean;
  enable_delta: boolean;
  disable: boolean;
  ocsp_disable: boolean;
}

interface ErrorObject {
  type: string;
  message: string;
}
export default class PkiConfigurationEditComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly version: VersionService;
  @service declare readonly api: ApiService;

  @tracked invalidFormAlert = '';
  @tracked errors: Array<ErrorObject> = [];

  get isEnterprise() {
    return this.version.isEnterprise;
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      // first clear errors and sticky flash messages
      this.errors = [];
      this.flashMessages.clearMessages();

      for (const type of ['cluster', 'acme', 'urls', 'crl']) {
        // skip saving and continue to next iteration if user does not have permission
        if (!this.args.capabilities[`canSet${capitalize(type)}`]) continue;

        try {
          const formKey = `${type}Form` as 'clusterForm' | 'acmeForm' | 'urlsForm' | 'crlForm';
          const { data } = this.args[formKey].toJSON();
          const apiKey = `pkiConfigure${capitalize(type)}` as
            | 'pkiConfigureAcme'
            | 'pkiConfigureCluster'
            | 'pkiConfigureUrls'
            | 'pkiConfigureCrl';
          await this.api.secrets[apiKey](this.args.backend, data);
          this.flashMessages.success(`Successfully updated config/${type}`);
        } catch (error) {
          const { message } = await this.api.parseError(error);
          this.errors = addToArray(this.errors, { type, message });
          this.flashMessages.danger(`Error updating config/${type}`, { sticky: true });
        }
      }

      if (this.errors.length) {
        this.invalidFormAlert = 'There was an error submitting this form.';
      } else {
        this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
      }
    })
  );

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
  }

  @action
  handleTtl(field: FormField, e: TtlEvent) {
    const { enabled, goSafeTimeString } = e;
    const {
      name,
      options: { mapToBoolean, isOppositeValue },
    } = field;
    const { data } = this.args.crlForm;

    data[name as keyof PkiConfigCrlTtls] = goSafeTimeString;
    // expiry and ocspExpiry both correspond to 'disable' booleans
    // so when ttl is enabled, the booleans are set to false
    data[mapToBoolean as keyof PkiConfigCrlBooleans] = isOppositeValue ? !enabled : enabled;
  }
}
