/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { computed } from '@ember/object';
import Controller from '@ember/controller';

export default Controller.extend({
  flashMessages: service(),

  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  filter: null,
  page: 1,
  pageFilter: null,

  filterFocused: false,

  isLoading: false, // set via the route `loading` action
  policyToDelete: null, // set when clicking 'Delete' from popup menu

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  },

  filterMatchesKey: computed('filter', 'model', 'model.[]', function () {
    var filter = this.filter;
    var content = this.model;
    return !!(content && content.length && content.find((c) => c.id === filter));
  }),

  firstPartialMatch: computed('filter', 'model', 'model.[]', 'filterMatchesKey', function () {
    var filter = this.filter;
    var content = this.model;
    if (!content) {
      return;
    }
    var filterMatchesKey = this.filterMatchesKey;
    var re = new RegExp('^' + filter);
    return filterMatchesKey
      ? null
      : content.find(function (key) {
          return re.test(key.id);
        });
  }),

  actions: {
    setFilter: function (val) {
      this.set('filter', val);
    },
    setFilterFocus: function (bool) {
      this.set('filterFocused', bool);
    },
    deletePolicy(model) {
      const { policyType } = model;
      const name = model.id;
      const flash = this.flashMessages;
      model
        .destroyRecord()
        .then(() => {
          // this will clear the dataset cache on the store
          this.send('reload');
          flash.success(`${policyType.toUpperCase()} policy "${name}" was successfully deleted.`);
        })
        .catch((e) => {
          const errors = e.errors ? e.errors.join('') : e.message;
          flash.danger(
            `There was an error deleting the ${policyType.toUpperCase()} policy "${name}": ${errors}.`
          );
        })
        .finally(() => this.set('policyToDelete', null));
    },
  },
});
