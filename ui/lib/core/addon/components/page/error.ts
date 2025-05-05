/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ResponseError } from '@hashicorp/vault-client-typescript';

import type Owner from '@ember/owner';
import type ApiService from 'vault/services/api';

interface Args {
  error: ResponseError | Error;
}

type Error = {
  httpStatus?: number;
  path?: string;
  message: string;
  errors: string[];
};

export default class PageErrorComponent extends Component<Args> {
  @service declare readonly api: ApiService;

  @tracked declare error: Error;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    this.unpackError();
  }

  async unpackError() {
    if (this.args.error instanceof ResponseError) {
      const { status, path, response } = await this.api.parseError(this.args.error);
      this.error = {
        httpStatus: status,
        path,
        message: response?.message,
        errors: response?.errors || [],
      };
    } else {
      this.error = this.args.error;
    }
  }
}
