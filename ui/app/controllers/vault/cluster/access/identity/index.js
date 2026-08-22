/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import ListController from 'core/mixins/list-controller';

export default Controller.extend(ListController, {
  flashMessages: service(),
  api: service(),

  entityToDisable: null,
  itemToDelete: null,

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  },

  actions: {
    async delete(model) {
      const type = this.identityType;
      const id = model.id;

      try {
        const methodType = type === 'group' ? 'groupDeleteById' : 'entityDeleteById';
        await this.api.identity[methodType](id);
        this.flashMessages.success('Successfully deleted');
      } catch (e) {
        const { message } = await this.api.parseError(e);

        this.flashMessages.danger(`There was a problem deleting ${type}: ${id} - ${message}`);
      } finally {
        this.set('itemToDelete', null);
        this.send('reload');
      }
    },

    toggleDisabled(model) {
      const action = model.disabled ? ['enabled', 'enabling'] : ['disabled', 'disabling'];
      const type = model.identityType;
      const id = model.id;
      model.toggleProperty('disabled');

      model
        .save()
        .then(() => {
          this.flashMessages.success(`Successfully ${action[0]} ${type}: ${id}`);
        })
        .catch((e) => {
          this.flashMessages.success(
            `There was a problem ${action[1]} ${type}: ${id} - ${e.errors.join(' ') || e.message}`
          );
        })
        .finally(() => this.set('entityToDisable', null));
    },
    reloadRecord() {
      this.send('reload');
    },
  },
});
