/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';

import type RouterService from '@ember/routing/router';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type DownloadService from 'vault/services/download';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type PkiCertificateForm from 'vault/forms/secrets/pki/certificate';
import type CapabilitiesService from 'vault/services/capabilities';
import type {
  PkiIssueWithRoleRequest,
  PkiIssueWithRoleResponse,
  PkiSignWithRoleRequest,
  PkiSignWithRoleResponse,
} from '@hashicorp/vault-client-typescript';

interface Args {
  role: string;
  form: PkiCertificateForm;
  onSuccess: CallableFunction;
  mode: 'generate' | 'sign';
}

export default class PkiRoleGenerate extends Component<Args> {
  @service declare readonly download: DownloadService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked declare certData: PkiIssueWithRoleResponse | PkiSignWithRoleResponse;
  @tracked canRevoke = false;

  save = task(
    waitFor(async (evt: Event) => {
      evt.preventDefault();
      this.errorBanner = '';
      const { role, form, mode, onSuccess } = this.args;

      try {
        const { data } = form.toJSON();
        if (mode === 'generate') {
          this.certData = await this.api.secrets.pkiIssueWithRole(
            role,
            this.secretMountPath.currentPath,
            data as PkiIssueWithRoleRequest
          );
        } else {
          this.certData = await this.api.secrets.pkiSignWithRole(
            role,
            this.secretMountPath.currentPath,
            data as PkiSignWithRoleRequest
          );
        }
        // check for revoke capabilities for certificate details component
        const { canCreate } = await this.capabilities.for('pkiRevoke', {
          backend: this.secretMountPath.currentPath,
        });
        this.canRevoke = canCreate;
        onSuccess();
        // since we are staying on the same page to display cert details scroll to top
        window.scrollTo(0, 0);
      } catch (err) {
        const { message } = await this.api.parseError(
          err,
          `Could not ${mode} certificate. See Vault logs for details.`
        );
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );

  @action cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
  }
}
