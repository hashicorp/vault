/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module AuthForm
 * The `PkiPaginatedList` is used to handle a list page layout with lazyPagination response.
 * It is specific to PKI so we can make certain assumptions about routing.
 * The toolbar has no filtering since users can go directly to an item from the overview page.
 *
 * @example ```js
 * <PkiPaginatedList @list={{this.model.roles}} @hasConfig={{this.model.hasConfig}} @listRoute="roles.index">
 *   <:list as |items|>
 *     {{#each items as |item}}
 *       <div>for each thing</div>
 *     {{/each}}
 *   </:list>
 * </PkiPaginatedList>
 * ```
 */

interface Args {
  list: unknown[];
  listRoute: string;
  hasConfig?: boolean;
  backend: string;
}
export default class PkiPaginatedListComponent extends Component<Args> {
  get paginationQueryParams() {
    return (page: number) => ({ page });
  }
  get hasConfig() {
    if (typeof this.args.hasConfig === 'boolean') return this.args.hasConfig;
    return true;
  }
}
