/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type AdapterError from 'vault/@ember-data/adapter/error';
import type ApiService from 'vault/services/api';
import type AwsCredential from 'vault/models/aws-credential';
import type ControlGroupService from 'vault/vault/services/control-group';
import type RouterService from '@ember/routing/router-service';
import type Store from '@ember-data/store';

const CREDENTIAL_TYPES = {
  ssh: {
    model: 'ssh-otp-credential',
    title: 'Generate SSH Credentials',
    formFields: ['username', 'ip'],
    displayFields: ['username', 'ip', 'key', 'keyType', 'port'],
  },
  aws: {
    model: 'aws-credential',
    title: 'Generate AWS Credentials',
    backIsListLink: true,
    displayFields: ['accessKey', 'secretKey', 'securityToken', 'leaseId', 'renewable', 'leaseDuration'],
    // aws form fields are dynamic
    formFields: (model: AwsCredential) => {
      return {
        iam_user: ['credentialType'],
        assumed_role: ['credentialType', 'ttl', 'roleArn'],
        federation_token: ['credentialType', 'ttl'],
        session_token: ['credentialType', 'ttl'],
      }[model.credentialType as string];
    },
  },
};

interface Args {
  awsRoleType: string | undefined;
  backendPath: string;
  backendType: 'ssh' | 'aws';
  roleName: string;
}

export default class GenerateCredentials extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly controlGroup: ControlGroupService;
  @service declare readonly store: Store;
  @service declare readonly router: RouterService;

  @tracked model;
  @tracked loading = false;
  @tracked hasGenerated = false;

  cannotReadAwsRole = false;
  emptyData = '{\n}';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const modelType = this.modelForType();
    this.model = this.generateNewModel(modelType);

    // if user lacks role read permissions, awsRoleType will be undefined
    // the role type dictates which form inputs are available, so this case
    // will need special handling when generating credentials
    this.cannotReadAwsRole = this.args.backendType == 'aws' && !this.args.awsRoleType;
  }

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    if (noTeardown && !this.model.isDestroyed && !this.model.isDestroying) {
      this.model.unloadRecord();
    }
    super.willDestroy();
  }

  modelForType(): string | undefined {
    const type = this.options;
    if (type) {
      return type.model;
    }
    // if we don't have a model for that type then redirect them back to the backend list
    this.router.transitionTo('vault.cluster.secrets.backend.list-root', this.args.backendPath);
    return undefined;
  }

  get breadcrumbs() {
    return [
      {
        label: this.args.backendPath,
        route: 'vault.cluster.secrets.backend',
        model: this.args.backendPath,
      },
      { label: 'Credentials', route: 'vault.cluster.secrets.backend', model: this.args.backendPath },
      { label: this.args.roleName, route: 'vault.cluster.secrets.backend.show', model: this.args.roleName },
      { label: this.options.title },
    ];
  }

  get helpText(): string {
    let message = '';
    if (this.cannotReadAwsRole) {
      message =
        'You do not have permissions to read this role so Vault cannot infer the credential type. Select the credential type you want to generate. ';
    }
    if (this.options?.model === 'aws-credential' && this.model.credentialType === 'iam_user')
      message += 'For Vault roles of credential type iam_user, there are no inputs, just submit the form.';
    return message;
  }

  get options() {
    return CREDENTIAL_TYPES[this.args.backendType];
  }

  get formFields(): string[] | undefined {
    const typeOpts = this.options;

    if (typeof typeOpts.formFields === 'function') {
      // without read access to the role, awsRoleType will be undefined and will default to iam_user
      // so we will need to show credentialType input for user selection
      // otherwise, we can omit that input
      const fields = typeOpts.formFields(this.model) ?? [];

      if (!this.cannotReadAwsRole) {
        return fields.filter((f) => f !== 'credentialType');
      }

      return fields;
    }
    return typeOpts.formFields;
  }

  get displayFields() {
    return this.options.displayFields;
  }

  generateNewModel(modelType?: string) {
    if (!modelType) {
      return;
    }
    const { roleName, backendPath, awsRoleType } = this.args;
    // conditionally add credentialType so that if not present, it will default to iam_user
    const attrs = {
      role: {
        backend: backendPath,
        name: roleName,
      },
      id: `${backendPath}-${roleName}`,
      ...(awsRoleType ? { credentialType: awsRoleType } : {}),
    };

    return this.store.createRecord(modelType, attrs);
  }

  replaceModel() {
    const modelType = this.modelForType();
    if (!modelType) {
      return;
    }
    if (this.model) {
      this.model.unloadRecord();
    }
    this.model = this.generateNewModel(modelType);
  }

  @action
  create(evt: Event) {
    evt.preventDefault();
    this.loading = true;
    this.model
      .save()
      .then(() => {
        this.hasGenerated = true;
      })
      .catch(async (error: AdapterError) => {
        const { response } = await this.api.parseError(error);
        // Handle control group AdapterError
        if (response?.isControlGroupError) {
          this.controlGroup.saveTokenFromError(response);
          const err = this.controlGroup.logFromError(response);
          error.errors = [err.content];
        }
        throw error;
      })
      .finally(() => {
        this.loading = false;
      });
  }

  @action
  editorUpdated(attr: string, val: string) {
    // wont set invalid JSON to the model
    try {
      this.model[attr] = JSON.parse(val);
    } catch {
      // linting is handled by the component
    }
  }

  @action
  reset() {
    this.hasGenerated = false;
    this.replaceModel();
  }
}
