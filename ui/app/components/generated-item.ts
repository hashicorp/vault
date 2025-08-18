/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module GeneratedItem
 * The `GeneratedItem` component is the form to configure generated items related to mounts (e.g. groups, roles, users)
 *
 * @example
 * ```js
 * <GeneratedItem @model={{model}} @mode={{mode}} @itemType={{itemType/>
 * ```
 *
 * @property model=null {DS.Model} - The corresponding item model that is being configured.
 * @property mode {String} - which config mode to use. either `show`, `edit`, or `create`
 * @property itemType {String} - the type of item displayed
 *
 */

interface GeneratedItemArgs {
  model: any;
  itemType: string;
  mode: string;
  mountPath: string;
}

export default class GeneratedItem extends Component<GeneratedItemArgs> {
  @service declare flashMessages: any;
  @service declare router: any;
  
  @tracked modelValidations: any = null;
  @tracked isFormInvalid = false;

  get props() {
    return this.args.model.serialize();
  }

  constructor(owner: any, args: GeneratedItemArgs) {
    super(owner, args);
    
    if (this.args.mode === 'edit') {
      // For validation to work in edit mode,
      // reconstruct the model values from field group
      this.args.model.fieldGroups.forEach((element: any) => {
        if (element.default) {
          element.default.forEach((attr: any) => {
            const fieldValue = attr.options && attr.options.fieldValue;
            if (fieldValue) {
              this.args.model[attr.name] = this.args.model[fieldValue];
            }
          });
        }
      });
    }
  }

  validateForm() {
    // Only validate on new models because blank passwords will not be updated
    // in practice this only happens for userpass users
    if (this.args.model.validate && this.args.model.isNew) {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = state;
      this.isFormInvalid = !isValid;
      return isValid;
    } else {
      this.isFormInvalid = false;
      return true;
    }
  }

  saveModel = task(
    waitFor(function* (this: GeneratedItem) {
      const isValid = this.validateForm();
      if (!isValid) {
        return;
      }
      try {
        yield this.args.model.save();
      } catch (err) {
        // AdapterErrors are handled by the error-message component
        // in the form
        if (err instanceof AdapterError === false) {
          throw err;
        }
        return;
      }
      this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
      this.flashMessages.success(`Successfully saved ${this.args.itemType} ${this.args.model.id}.`);
    })
  );

  onKeyUp = (name: string, value: any) => {
    this.args.model.set(name, value);
  };

  deleteItem = () => {
    this.args.model.destroyRecord().then(() => {
      this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
      this.flashMessages.success(`Successfully deleted ${this.args.itemType} ${this.args.model.id}.`);
    });
  };
}
