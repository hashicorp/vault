/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type ControlGroupService from 'vault/services/control-group';

interface ControlGroupModel {
  id: string;
  requestPath: string;
  approved: boolean;
  requestEntity?: { id: string; name: string };
  reload(): void;
}

interface ControlGroupResponse {
  token?: string;
  uiParams?: {
    url: string;
  };
}

interface UnwrapResponse {
  auth?: unknown;
  data?: unknown;
}

interface ParseErrorResponse {
  message: string;
}

interface Args {
  model: ControlGroupModel;
  controlGroupResponse: ControlGroupResponse;
}

/**
 * @module ControlGroupSuccess
 * ControlGroupSuccess component handles the success state of control group authorization.
 * It allows users to unwrap tokens and navigate to authorized resources.
 *
 * @example
 * <ControlGroupSuccess @model={{this.model}} @controlGroupResponse={{this.controlGroupResponse}} />
 *
 * @param {ControlGroupModel} model - The control group model containing request details
 * @param {ControlGroupResponse} controlGroupResponse - Response object containing token and navigation info
 */
export default class ControlGroupSuccessComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly controlGroup: ControlGroupService;
  @service declare readonly api: ApiService;

  @tracked error: string | null = null;
  @tracked unwrapData: unknown = null;
  @tracked token = '';

  @task({ drop: true })
  @waitFor
  *unwrapTask(token: string) {
    this.error = null;
    try {
      const response = (yield this.api.sys.unwrap({}, this.api.buildHeaders({ token }))) as UnwrapResponse;
      this.unwrapData = response.auth || response.data;
      this.controlGroup.deleteControlGroupToken(this.args.model.id);
    } catch (e) {
      const { message } = (yield this.api.parseError(e)) as ParseErrorResponse;
      this.error = `Token unwrap failed: ${message}`;
    }
  }

  @task({ drop: true })
  @waitFor
  *markAndNavigateTask() {
    this.controlGroup.markTokenForUnwrap(this.args.model.id);
    const url = this.args.controlGroupResponse.uiParams?.url;
    if (url) {
      yield this.router.transitionTo(url);
    }
  }
}
