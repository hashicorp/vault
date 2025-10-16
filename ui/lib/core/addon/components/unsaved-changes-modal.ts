/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type RouterService from '@ember/routing/router-service';
import type UnsavedChangesService from 'vault/services/unsaved-changes';
import type { TaskGenerator, Task } from 'ember-concurrency';

/**
 * @module UnsavedChangesModal handles displaying the unsaved changes modal.
 * 
 * @example
 * <SecretEngine::UnsavedChangesModal
    @model={{this.model}}
    />
 *
 * @param {object} model - A model contains a secret engine resource, lease config from the sys/internal endpoint. 
 * */

interface Args {
  model: { secretsEngine: SecretsEngineResource };
  showUnsavedChanges: boolean;
  onSave: Task<TaskGenerator<[string]>, []>;
  onDiscard: () => void;
}

export default class UnsavedChangesModal extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly unsavedChanges: UnsavedChangesService;

  @action
  closeModal() {
    this.unsavedChanges.showModal = false;
  }

  @action
  closeAndHandle(close: () => void, action: 'save' | 'discard') {
    close();

    if (action === 'save') {
      this.args.onSave.perform();
    }

    if (action === 'discard') {
      this.unsavedChanges.changedFields = [];
      this.args.onDiscard();
    }
  }
}
