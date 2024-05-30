/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { task } from 'ember-concurrency';
import { humanize } from 'vault/helpers/humanize';
import { waitFor } from '@ember/test-waiters';

export default Component.extend({
  flashMessages: service(),
  store: service(),
  'data-test-component': 'identity-edit-form',
  attributeBindings: ['data-test-component'],
  model: null,

  // 'create', 'edit', 'merge'
  mode: 'create',
  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully saving an entity
   */
  onSave: () => {},

  cancelLink: computed('mode', 'model.identityType', function () {
    const { model, mode } = this;
    const routes = {
      'create-entity': 'vault.cluster.access.identity',
      'edit-entity': 'vault.cluster.access.identity.show',
      'merge-entity-merge': 'vault.cluster.access.identity',
      'create-entity-alias': 'vault.cluster.access.identity.aliases',
      'edit-entity-alias': 'vault.cluster.access.identity.aliases.show',
      'create-group': 'vault.cluster.access.identity',
      'edit-group': 'vault.cluster.access.identity.show',
      'create-group-alias': 'vault.cluster.access.identity.aliases',
      'edit-group-alias': 'vault.cluster.access.identity.aliases.show',
    };
    const key = model ? `${mode}-${model.identityType}` : 'merge-entity-alias';
    return routes[key];
  }),

  getMessage(model, isDelete = false) {
    const mode = this.mode;
    const typeDisplay = humanize([model.identityType]);
    const action = isDelete ? 'deleted' : 'saved';
    if (mode === 'merge') {
      return 'Successfully merged entities';
    }
    if (model.id) {
      return `Successfully ${action} ${typeDisplay} ${model.id}.`;
    }
    return `Successfully ${action} ${typeDisplay}.`;
  },

  save: task(
    waitFor(function* () {
      const model = this.model;
      const message = this.getMessage(model);

      try {
        yield model.save();
      } catch (err) {
        // err will display via model state
        return;
      }
      this.flashMessages.success(message);
      yield this.onSave({ saveType: 'save', model });
    })
  ).drop(),

  willDestroy() {
    // components are torn down after store is disconnected and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    const model = this.model;
    if (noTeardown && model && model.isDirty && !model.isDestroyed && !model.isDestroying) {
      model.rollbackAttributes();
    }
    this._super(...arguments);
  },

  actions: {
    deleteItem(model) {
      const message = this.getMessage(model, true);
      const flash = this.flashMessages;
      model.destroyRecord().then(() => {
        flash.success(message);
        return this.onSave({ saveType: 'delete', model });
      });
    },
  },
});
