/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import type SecretEngineModel from 'vault/models/secret-engine';

type AWSRootCredsFields = {
  access_key: string | null;
  iam_endpoint: string | null;
  sts_endpoint: string | null;
  secret_key: string | null;
  region: string | null;
};

interface Args {
  model: SecretEngineModel;
  accessKey: string;
  iamEndpoint: string;
  region: string;
  secretKey: string;
  stsEndpoint: string;
  onSubmit: (data: AWSRootCredsFields) => void;
}

export default class ConfigureAwsSecretAccessToAwsFormComponent extends Component<Args> {
  @tracked showOptions = false;

  @action
  saveRootCreds(data: AWSRootCredsFields, event: Event) {
    event.preventDefault();
    this.args.onSubmit(data);
  }
}
