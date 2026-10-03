/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import { action } from '@ember/object';

import type FlashMessageService from 'vault/services/flash-messages';
import type PkiKeyForm from 'vault/forms/secrets/pki/key';
import type { Capabilities, ValidationMap } from 'vault/app-types';
import type ApiService from 'vault/services/api';
import type SecretMountPathService from 'vault/services/secret-mount-path';
import type RouterService from '@ember/routing/router-service';
import type CapabilitiesService from 'vault/services/capabilities';
import type {
  PkiGenerateInternalKeyRequest,
  PkiGenerateExportedKeyRequest,
  PkiGenerateInternalKeyResponse,
  PkiGenerateExportedKeyResponse,
  PkiWriteKeyRequest,
} from '@hashicorp/vault-client-typescript';

type PkiGenerateKeyResponse = PkiGenerateInternalKeyResponse | PkiGenerateExportedKeyResponse;

/**
 * @module PkiKeyForm
 * PkiKeyForm components are used to create and update PKI keys.
 *
 * @param {Form} form - pki key form.
 */

interface Args {
  form: PkiKeyForm;
  canUpdate?: boolean;
  canDelete?: boolean;
}

export default class PkiKeyFormComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPathService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly capabilities: CapabilitiesService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;
  @tracked declare generatedKey: PkiGenerateKeyResponse;
  @tracked declare generatedKeyCapabilities: Capabilities;

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      try {
        const { currentPath } = this.secretMountPath;
        const { form } = this.args;
        const { isValid, state, invalidFormMessage, data } = form.toJSON();

        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;

        if (isValid) {
          const { type, key_id, ...payload } = data;
          if (!form.isNew) {
            this.generatedKey = await this.api.secrets.pkiWriteKey(
              key_id as string,
              currentPath,
              payload as PkiWriteKeyRequest
            );
          } else if (data.type === 'internal') {
            this.generatedKey = await this.api.secrets.pkiGenerateInternalKey(
              currentPath,
              payload as PkiGenerateInternalKeyRequest
            );
          } else {
            this.generatedKey = await this.api.secrets.pkiGenerateExportedKey(
              currentPath,
              payload as PkiGenerateExportedKeyRequest
            );
          }

          this.flashMessages.success(
            `Successfully ${form.isNew ? 'generated' : 'updated'} key${
              data.key_name ? ` ${data.key_name}.` : '.'
            }`
          );

          const { private_key, key_id: keyId } = this.generatedKey;
          // only transition to details if there is no private_key data to display
          if (!private_key) {
            this.router.transitionTo('vault.cluster.secrets.backend.pki.keys.key.details', keyId);
          } else {
            // check capabilities on newly generated key
            this.generatedKeyCapabilities = await this.capabilities.for('pkiKey', {
              backend: currentPath,
              keyId,
            });
          }
        }
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );

  @action
  onCancel() {
    if (this.args.form.isNew) {
      this.router.transitionTo('vault.cluster.secrets.backend.pki.keys.index');
    } else {
      this.router.transitionTo(
        'vault.cluster.secrets.backend.pki.keys.key.details',
        this.args.form.data.key_id
      );
    }
  }
}
