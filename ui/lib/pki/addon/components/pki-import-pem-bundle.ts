/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';

import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPathService from 'vault/services/secret-mount-path';
import type {
  PkiIssuersImportBundleResponse,
  PkiConfigureCaResponse,
} from '@hashicorp/vault-client-typescript';

/**
 * @module PkiImportPemBundle
 * PkiImportPemBundle components are used to import PKI CA certificates and keys via pem_bundle.
 * https://github.com/hashicorp/vault/blob/main/website/content/api-docs/secret/pki.mdx#import-ca-certificates-and-keys
 *
 */

interface Args {
  onSave?: CallableFunction;
  onCancel: CallableFunction;
  onComplete: CallableFunction;
  useIssuer: boolean;
}

export default class PkiImportPemBundle extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPathService;

  @tracked pemBundle = '';
  @tracked errorBanner = '';
  @tracked importedIssuerKeyMap: Array<{ issuer: string; key: string }> | null = null;

  mapResponse(response: PkiIssuersImportBundleResponse | PkiConfigureCaResponse) {
    const { mapping, imported_issuers, imported_keys } = response;
    // Even if there are no imported items, mapping will be an empty object from API response
    if (undefined === mapping) return null;

    const importList = (imported_issuers || []).map((issuer: string) => {
      const key = mapping[issuer as keyof typeof mapping] as string;
      return { issuer, key };
    });

    // Check each imported key and make sure it's in the list
    (imported_keys || []).forEach((key) => {
      const matchIdx = importList.findIndex((item) => item.key === key);
      // If key isn't accounted for, add it without a matching issuer
      if (matchIdx === -1) {
        importList.push({ issuer: '', key });
      }
    });

    if (importList.length === 0) {
      // If no new items were imported but the import call was successful, the UI will show accordingly
      return [{ issuer: '', key: '' }];
    }
    return importList;
  }

  submitForm = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      this.errorBanner = '';
      if (!this.pemBundle) {
        this.errorBanner = 'please upload your PEM bundle';
        return;
      }
      try {
        const { currentPath } = this.secretMountPath;
        const payload = { pem_bundle: this.pemBundle };
        let response: PkiIssuersImportBundleResponse | PkiConfigureCaResponse;

        if (this.args.useIssuer) {
          response = await this.api.secrets.pkiIssuersImportBundle(currentPath, payload);
        } else {
          response = await this.api.secrets.pkiConfigureCa(currentPath, payload);
        }

        this.importedIssuerKeyMap = this.mapResponse(response);
        this.flashMessages.success('Successfully imported data.');
        // This component shows the results, but call `onSave` for any side effects on parent
        if (this.args.onSave) {
          this.args.onSave();
        }
        window?.scrollTo(0, 0);
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
      }
    })
  );

  @action
  onFileUploaded({ value }: { value: string }) {
    this.pemBundle = value;
  }
}
