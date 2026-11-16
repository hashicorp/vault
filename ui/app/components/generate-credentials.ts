/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import AwsCredentialForm from 'vault/forms/aws-credential';

import type ApiService from 'vault/services/api';
import type ControlGroupService from 'vault/vault/services/control-group';
import AwsCredentialResource from 'vault/resources/aws-credential';
import type { AwsCredentialResponse } from 'vault/resources/aws-credential';

const CREDENTIAL_TYPES = {
  aws: {
    title: 'Generate AWS Credentials',
    backIsListLink: true,
    formFields: (credentialType?: string) => {
      const fields: Record<string, string[]> = {
        iam_user: ['credential_type'],
        assumed_role: ['credential_type', 'ttl', 'role_arn'],
        federation_token: ['credential_type', 'ttl'],
        session_token: ['credential_type', 'ttl'],
      };
      return fields[credentialType || 'iam_user'];
    },
  },
};

interface Args {
  awsRoleType: string | undefined;
  backendPath: string;
  backendType: 'aws';
  roleName: string;
}

export default class GenerateCredentials extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly controlGroup: ControlGroupService;

  @tracked model: AwsCredentialResource | undefined;
  @tracked credentialForm: AwsCredentialForm;
  @tracked errorMessage: string | null = null;
  @tracked loading = false;
  @tracked hasGenerated = false;

  cannotReadAwsRole = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.credentialForm = new AwsCredentialForm({ credential_type: this.args.awsRoleType || 'iam_user' });
    this.cannotReadAwsRole = this.args.backendType === 'aws' && !this.args.awsRoleType;
  }

  get breadcrumbs() {
    const { backendPath, roleName } = this.args;
    return [
      { label: backendPath, route: 'vault.cluster.secrets.backend', model: backendPath },
      { label: 'Credentials', route: 'vault.cluster.secrets.backend', model: backendPath },
      { label: roleName, route: 'vault.cluster.secrets.backend.show', model: roleName },
      { label: this.options.title },
    ];
  }

  get helpText(): string {
    let message = '';
    if (this.cannotReadAwsRole) {
      message =
        'You do not have permissions to read this role so Vault cannot infer the credential type. Select the credential type you want to generate. ';
    }
    if (this.credentialForm.data.credential_type === 'iam_user')
      message += 'For Vault roles of credential type iam_user, there are no inputs, just submit the form.';
    return message;
  }

  get options() {
    return CREDENTIAL_TYPES[this.args.backendType];
  }

  get formFields() {
    const fields = this.options.formFields(this.credentialForm.data.credential_type) ?? [];
    return this.credentialForm.formFields.filter(
      (field) => fields.includes(field.name) && (this.cannotReadAwsRole || field.name !== 'credential_type')
    );
  }

  @action
  @waitFor
  async create(evt: Event) {
    evt.preventDefault();
    if (this.loading) return;
    this.loading = true;
    this.errorMessage = null;
    const { roleName, backendPath } = this.args;
    const params = { name: roleName, aws_mount_path: backendPath };
    const { data } = this.credentialForm.toJSON();
    try {
      const result =
        this.credentialForm.data.credential_type === 'iam_user'
          ? await this.api.secrets.awsGenerateCredentialsRaw(params)
          : await this.api.secrets.awsGenerateCredentialsWithParametersRaw({
              ...params,
              AwsGenerateCredentialsWithParametersRequest: data,
            });
      // The generated value() strips top-level lease metadata from Vault's envelope.
      const response: AwsCredentialResponse = await result.raw.json();
      if (!response.data?.access_key) {
        throw new Error('AWS did not return credentials. Please try again.');
      }
      this.model = new AwsCredentialResource(
        {
          ...response.data,
          lease_id: response.lease_id,
          lease_duration: response.lease_duration,
          renewable: response.renewable,
        },
        this
      );
      this.hasGenerated = true;
    } catch (error) {
      const { message, response } = await this.api.parseError(error);
      if (response?.isControlGroupError) {
        this.controlGroup.saveTokenFromError(response);
        this.errorMessage = this.controlGroup.logFromError(response).content;
      } else {
        this.errorMessage = message;
      }
    } finally {
      this.loading = false;
    }
  }

  @action
  reset() {
    this.hasGenerated = false;
    this.model = undefined;
    this.errorMessage = null;
    this.credentialForm = new AwsCredentialForm({ credential_type: this.args.awsRoleType || 'iam_user' });
  }
}
