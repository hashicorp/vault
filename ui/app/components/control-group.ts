/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type AuthService from 'vault/services/auth';
import type ControlGroupService from 'vault/services/control-group';

interface ControlGroupModel {
  id: string;
  requestPath: string;
  approved: boolean;
  canAuthorize: boolean;
  requestEntity?: {
    id: string;
    name: string;
    canRead: boolean;
  };
  authorizations: Array<{
    id: string;
    name: string;
    canRead: boolean;
  }>;
  reload(): Promise<void>;
  save(): Promise<void>;
}

interface ControlGroupResponse {
  token?: string;
  uiParams?: {
    url: string;
  };
}

interface Args {
  model: ControlGroupModel;
}

/**
 * @module ControlGroup
 * ControlGroup component handles the display and authorization of control group requests.
 * It shows request details, authorization status, and provides actions for authorization or refresh.
 *
 * @example
 * <ControlGroup @model={{this.model}} />
 *
 * @param {ControlGroupModel} model - The control group model containing request and authorization details
 */
export default class ControlGroupComponent extends Component<Args> {
  @service declare readonly auth: AuthService;
  @service declare readonly controlGroup: ControlGroupService;

  @tracked errors: unknown = null;
  @tracked controlGroupResponse: ControlGroupResponse = {};

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.updateControlGroupResponse();
  }

  private updateControlGroupResponse() {
    const accessor = this.args.model.id;
    const data = this.controlGroup.wrapInfoForAccessor(accessor);
    this.controlGroupResponse = typeof data === 'object' && data !== null ? data : {};
  }

  get currentUserEntityId(): string | undefined {
    return this.auth.authData?.entityId;
  }

  get currentUserIsRequesting(): boolean {
    if (!this.args.model.requestEntity) return false;
    return this.currentUserEntityId === this.args.model.requestEntity.id;
  }

  get currentUserHasAuthorized(): boolean {
    const authorizations = this.args.model.authorizations || [];
    return Boolean(authorizations.find((authz) => authz.id === this.currentUserEntityId));
  }

  get isSuccess(): boolean {
    return this.currentUserHasAuthorized || this.args.model.approved;
  }

  get requestorName(): string {
    const entity = this.args.model.requestEntity;

    if (this.currentUserIsRequesting) {
      return 'You';
    }
    if (entity && entity.name) {
      return entity.name;
    }
    return 'Someone';
  }

  get bannerPrefix(): string {
    if (this.currentUserHasAuthorized) {
      return 'Thanks!';
    }
    if (this.args.model.approved) {
      return 'Success!';
    }
    return 'Locked';
  }

  get bannerText(): string {
    const isApproved = this.args.model.approved;
    const { currentUserHasAuthorized, currentUserIsRequesting } = this;

    if (currentUserHasAuthorized) {
      return 'You have given authorization';
    }
    if (currentUserIsRequesting && isApproved) {
      return 'You have been given authorization';
    }
    if (isApproved) {
      return 'This Control Group has been authorized';
    }
    if (currentUserIsRequesting) {
      return 'The path you requested is locked by a Control Group';
    }
    return 'Someone is requesting access to a path locked by a Control Group';
  }

  @task({ drop: true })
  @waitFor
  *refreshTask() {
    try {
      yield this.args.model.reload();
      this.updateControlGroupResponse();
    } catch (e) {
      this.errors = e;
    }
  }

  @task({ drop: true })
  @waitFor
  *authorizeTask() {
    try {
      yield this.args.model.save();
      yield this.args.model.reload();
      this.updateControlGroupResponse();
    } catch (e) {
      this.errors = e;
    }
  }
}
