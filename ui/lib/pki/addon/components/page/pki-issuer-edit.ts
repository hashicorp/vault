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
import { removeFromArray } from 'vault/helpers/remove-from-array';
import { addToArray } from 'vault/helpers/add-to-array';

import type PkiIssuerForm from 'vault/forms/secrets/pki/issuers/issuer';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPathService from 'vault/services/secret-mount-path';

interface Args {
  form: PkiIssuerForm;
  issuerRef: string;
}

export default class PkiIssuerEditComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPathService;

  @tracked usageValues: Array<string> = [];
  @tracked error = null;
  declare initialIssuerName: string;

  fields = [
    'issuer_name',
    'leaf_not_after_behavior',
    'usage',
    'manual_chain',
    'revocation_signature_algorithm',
    'issuing_certificates',
    'crl_distribution_points',
    'ocsp_servers',
  ];

  usageOptions = [
    { label: 'Issuing certificates', value: 'issuing-certificates' },
    { label: 'Signing CRLs', value: 'crl-signing' },
    { label: 'Signing OCSPs', value: 'ocsp-signing' },
  ];

  notAfterOptions = ['always_enforce_err', 'err', 'truncate', 'permit'];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // comma separated strings appear to be represented as string[] in the spec
    // this will need to be addressed globally but for now cast usage to string and then back to string[]
    this.usageValues = ((this.args.form.data.usage as unknown as string) || '').split(',');
    // in the case where an issuer_name was not provided it will be considered 'default' and returned as an empty string
    // sending an empty string back to the API will result in an error
    // cache the initial value so we can check against it later before saving
    this.initialIssuerName = this.args.form.data.issuer_name || '';
  }

  @action
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
    this.args.form.data.usage = this.usageValues.join(',') as unknown as string[];
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      try {
        const { data } = this.args.form.toJSON();
        // if the issuer_name was originally default and unchanged then remove it from the payload
        if (!this.initialIssuerName && !data.issuer_name) {
          delete data.issuer_name;
        }
        await this.api.secrets.pkiWriteIssuer(this.args.issuerRef, this.secretMountPath.currentPath, data);
        this.flashMessages.success('Successfully updated issuer');
        this.toDetails();
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.error = message;
      }
    })
  );
}
