/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ConfigureAwsSecretComponent
 *
 * @example
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
 *
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
export default class ConfigureAwsSecretComponent extends Component {
  @action
  saveRootCreds(data, event) {
    event.preventDefault();
    this.args.saveAWSRoot(data);
  }

  @action
  saveLease(data, event) {
    event.preventDefault();
    this.args.saveAWSLease(data);
  }

  @action
  handleTtlChange(name, ttlObj) {
    // lease values cannot be undefined, set to 0 to use default
    const valueToSet = ttlObj.enabled ? ttlObj.goSafeTimeString : 0;
    this.args.model.set(name, valueToSet);
  }
}
