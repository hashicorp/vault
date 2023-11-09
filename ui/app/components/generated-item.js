/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module GeneratedItem
 * The `GeneratedItem` component is the form to configure generated items related to mounts (e.g. groups, roles, users)
 *
 * @example
 * <GeneratedItem @model={{model}} @mode={{mode}} @itemType={{itemType/>
 *
 *
 * @property model=null {DS.Model} - The corresponding item model that is being configured.
 * @property mode {String} - which config mode to use. either `show`, `edit`, or `create`
 * @property itemType {String} - the type of item displayed
 *
 */

export default Component.extend({
  model: null,
  itemType: null,
  flashMessages: service(),
  router: service(),
  modelValidations: null,
  isFormInvalid: false,
  props: computed('model', function () {
    return this.model.serialize();
  }),
  saveModel: task(
    waitFor(function* () {
      try {
        yield this.model.save();
      } catch (err) {
        // AdapterErrors are handled by the error-message component
        // in the form
        if (err instanceof AdapterError === false) {
          throw err;
        }
        return;
      }
      this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
      this.flashMessages.success(`Successfully saved ${this.itemType} ${this.model.id}.`);
    })
  ),
  init() {
    this._super(...arguments);
    if (this.mode === 'edit') {
      // For validation to work in edit mode,
      // reconstruct the model values from field group
      this.model.fieldGroups.forEach((element) => {
        if (element.default) {
          element.default.forEach((attr) => {
            const fieldValue = attr.options && attr.options.fieldValue;
            if (fieldValue) {
              this.model[attr.name] = this.model[fieldValue];
            }
          });
        }
      });
    }
  },
  actions: {
    onKeyUp(name, value) {
      this.model.set(name, value);
      if (this.model.validate) {
        // Set validation error message for updated attribute
        const { isValid, state } = this.model.validate();
        this.setProperties({
          modelValidations: state,
          isFormInvalid: !isValid,
        });
      } else {
        this.set('isFormInvalid', false);
      }
    },
    deleteItem() {
      this.model.destroyRecord().then(() => {
        this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
        this.flashMessages.success(`Successfully deleted ${this.itemType} ${this.model.id}.`);
      });
    },
  },
});
