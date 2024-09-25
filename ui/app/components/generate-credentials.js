/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

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
    formFields: (model) => {
      return {
        iam_user: ['credentialType'],
        assumed_role: ['credentialType', 'ttl', 'roleArn'],
        federation_token: ['credentialType', 'ttl'],
        session_token: ['credentialType', 'ttl'],
      }[model.credentialType];
    },
  },
};

export default class GenerateCredentials extends Component {
  @service controlGroup;
  @service store;
  @service router;

  @tracked model;
  @tracked loading = false;
  @tracked hasGenerated = false;
  emptyData = '{\n}';

  constructor() {
    super(...arguments);
    const modelType = this.modelForType();
    this.model = this.generateNewModel(modelType);
  }

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    if (noTeardown && !this.model.isDestroyed && !this.model.isDestroying) {
      this.model.unloadRecord();
    }
    super.willDestroy();
  }

  modelForType() {
    const type = this.options;
    if (type) {
      return type.model;
    }
    // if we don't have a mode for that type then redirect them back to the backend list
    this.router.transitionTo('vault.cluster.secrets.backend.list-root', this.args.backendPath);
  }

  get helpText() {
    if (this.options?.model === 'aws-credential') {
      return 'For Vault roles of credential type iam_user, there are no inputs, just submit the form. Choose a type to change the input options.';
    }
    return '';
  }

  get options() {
    return CREDENTIAL_TYPES[this.args.backendType];
  }

  get formFields() {
    const typeOpts = this.options;
    if (typeof typeOpts.formFields === 'function') {
      return typeOpts.formFields(this.model);
    }
    return typeOpts.formFields;
  }

  get displayFields() {
    return this.options.displayFields;
  }

  generateNewModel(modelType) {
    if (!modelType) {
      return;
    }
    const { roleName, backendPath, awsRoleType } = this.args;
    const attrs = {
      role: {
        backend: backendPath,
        name: roleName,
      },
      id: `${backendPath}-${roleName}`,
    };
    if (awsRoleType) {
      // this is only set from route if backendType = aws
      attrs.credentialType = awsRoleType;
    }
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
  create(evt) {
    evt.preventDefault();
    this.loading = true;
    this.model
      .save()
      .then(() => {
        this.hasGenerated = true;
      })
      .catch((error) => {
        // Handle control group AdapterError
        if (error.message === 'Control Group encountered') {
          this.controlGroup.saveTokenFromError(error);
          const err = this.controlGroup.logFromError(error);
          error.errors = [err.content];
        }
        throw error;
      })
      .finally(() => {
        this.loading = false;
      });
  }

  @action
  codemirrorUpdated(attr, val, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror.state.lint.marked.length > 0;

    if (!hasErrors) {
      this.model[attr] = JSON.parse(val);
    }
  }

  @action
  reset() {
    this.hasGenerated = false;
    this.replaceModel();
  }
}
