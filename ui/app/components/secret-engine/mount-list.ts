/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { allEngines, mountableEngines } from 'vault/helpers/mountable-secret-engines';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import type Router from '@ember/routing/router';
import type VersionService from 'vault/services/version';
import type SecretEngineModel from 'vault/models/secret-engine';

/**
 todo: we will add to this by comparing the plugin/catelog enpoint to what is shown and grey out blah blah
 // ideallly this might become a template only component?/
 */

interface Args {
  secretEngineModel: SecretEngineModel;
  isMountTypeSaved: boolean;
}

export default class mountList extends Component<Args> {
  @service declare readonly version: VersionService;
  @service declare readonly router: Router;

  @tracked mountType = '';

  get secretEngineMounts() {
    return this.version.isEnterprise ? allEngines() : mountableEngines();
  }

  @action
  setMountType(type: string) {
    // modifying the model here... I think this will work?
    this.mountType = type;
  }

  @action
  saveMountType() {
    this.args.secretEngineModel.type = this.mountType;
    this.router.transitionTo('vault.cluster.secrets.mount.create', this.mountType);
  }

  @action
  onCancel() {
    this.args.secretEngineModel.unloadRecord();
    this.router.transitionTo('vault.cluster.secrets.backends');
  }
}
