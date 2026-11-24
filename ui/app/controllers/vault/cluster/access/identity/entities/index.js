/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import ListController from 'core/mixins/list-controller';

export default Controller.extend(ListController, {
  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
    pageSize: 'pageSize',
  },
  pageSize: 10,

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  },

  actions: {
    reloadRecord() {
      this.send('reload');
    },
  },
});
