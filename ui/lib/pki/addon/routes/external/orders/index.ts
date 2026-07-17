/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SecretsApiPkiExternalCaListLookupOrdersRecentListEnum } from '@hashicorp/vault-client-typescript';

import type { Breadcrumb } from 'vault/app-types';
import type { OrderStatusName } from 'pki/helpers/map-order-status';
import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export interface RecentOrderListItem {
  order_id: string;
  creation_date: string;
  identifiers: string[];
  last_update: Date;
  order_status: OrderStatusName;
  role_name: string;
}

export interface OrdersIndexRouteParams {
  within: string;
}

export default class PkiExternalOrdersIndexRoute extends Route {
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;

  queryParams = {
    within: {
      refreshModel: true,
    },
  };

  // If we set the default query param on the controller, Ember will clear
  // that param from the URL whenever the dropdown selection matches the
  // default value. We want query changes signaled to the user, so set the
  // default in the route and send it explicitly to the backend.
  defaultWithinQuery = '1h';

  async model(params: { within: string }) {
    let recentOrders: RecentOrderListItem[] = [];
    // The backend returns an error if the query param is invalid, so no need to validate it client-side
    const query = params?.within ? { within: params.within } : { within: this.defaultWithinQuery };
    try {
      const resp = await this.api.secrets.pkiExternalCaListLookupOrdersRecent(
        this.secretMountPath.currentPath,
        SecretsApiPkiExternalCaListLookupOrdersRecentListEnum.TRUE,
        (context) => this.api.addQueryParams(context, query)
      );
      recentOrders = this.api.keyInfoToArray<RecentOrderListItem>(resp, 'order_id');
    } catch (e) {
      // Catch 404s and render empty state instead; throw all other errors.
      const error = await this.api.parseError(e);
      if (error.status !== 404) {
        throw e;
      }
    }

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      recentOrders,
      // Annoyingly we cannot specify a default query param that displays in the URL without
      // adding transitions that cause the parent routes to refire repeatedly.
      // Manually pass the query as source of truth so the UI always displays what was sent to the API.
      query,
    };
  }

  setupController(controller: RouteController, resolvedModel: SecretsEngineResource) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'Recent orders' },
    ];
  }
}
