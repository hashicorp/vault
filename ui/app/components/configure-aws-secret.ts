/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

import type SecretEngineModel from 'vault/models/secret-engine';
import type { TtlEvent } from 'vault/app-types';

/**
 * @module ConfigureAwsSecretComponent
 *
 * @example
 * ```js
 * <ConfigureAwsSecret
    @model={{model}}
    @tab={{tab}}
    @accessKey={{accessKey}}
    @secretKey={{secretKey}}
    @region={{region}}
    @iamEndpoint={{iamEndpoint}}
    @stsEndpoint={{stsEndpoint}}
    @saveAWSRoot={{action "save" "saveAWSRoot"}}
    @saveAWSLease={{action "save" "saveAWSLease"}} />
 * ```
 *
 * @param {object} model - aws secret engine model
 * @param {string} tab - current tab selection
 * @param {string} accessKey - AWS access key
 * @param {string} secretKey - AWS secret key
 * @param {string} region - AWS region
 * @param {string} iamEndpoint - IAM endpoint
 * @param {string} stsEndpoint - Sts endpoint
 * @param {Function} saveAWSRoot - parent action which saves AWS root credentials
 * @param {Function} saveAWSLease - parent action which updates AWS lease information
 *
 */

type AWSRootCredsFields = {
  access_key: string | null;
  iam_endpoint: string | null;
  sts_endpoint: string | null;
  secret_key: string | null;
  region: string | null;
};

type LeaseFields = { lease: string; lease_max: string };

interface Args {
  model: SecretEngineModel;
  tab?: string;
  accessKey: string;
  secretKey: string;
  region: string;
  iamEndpoint: string;
  stsEndpoint: string;
  saveAWSRoot: (data: AWSRootCredsFields) => void;
  saveAWSLease: (data: LeaseFields) => void;
}

export default class ConfigureAwsSecretComponent extends Component<Args> {
  @action
  saveRootCreds(data: AWSRootCredsFields, event: Event) {
    event.preventDefault();
    this.args.saveAWSRoot(data);
  }

  @action
  saveLease(data: LeaseFields, event: Event) {
    event.preventDefault();
    this.args.saveAWSLease(data);
  }

  @action
  handleTtlChange(name: string, ttlObj: TtlEvent) {
    // lease values cannot be undefined, set to 0 to use default
    const valueToSet = ttlObj.enabled ? ttlObj.goSafeTimeString : 0;
    this.args.model.set(name, valueToSet);
  }
}
