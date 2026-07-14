/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type { RoleActiveOrdersRouteModel } from 'pki/routes/external/roles/role/active-orders';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  model: RoleActiveOrdersRouteModel;
}

export default class ExternalPkiPageActiveOrdersComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked searchInput = '';

  tableColumns = [
    {
      key: 'order_id',
      label: 'Order ID',
      isSortable: true,
      customTableItem: true,
    },
  ];

  get ordersList() {
    const filteredOrders = this.args.model.activeOrders.filter((id) => id.includes(this.searchInput));
    return filteredOrders?.map((o) => ({ order_id: o }));
  }

  @action
  handleSearch(e: HTMLElementEvent<HTMLInputElement>) {
    this.searchInput = e.target.value;
  }

  @action
  refresh() {
    this.router.refresh('vault.cluster.secrets.backend.pki.external.roles.role.active-orders');
  }
}
