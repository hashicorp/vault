/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { pluralize } from 'ember-inflector';
import { capitalize } from '@ember/string';
import { ROUTES } from 'vault/utils/routes';

const TABS_FOR_SETTINGS = {
  aws: [
    {
      label: 'Client',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['client'],
    },
    {
      label: 'Identity Allow List Tidy',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['identity-accesslist'],
    },
    {
      label: 'Role Tag Deny List Tidy',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['roletag-denylist'],
    },
  ],
  azure: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  github: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  gcp: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  jwt: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  oidc: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  kubernetes: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  ldap: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  okta: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
  radius: [
    {
      label: 'Configuration',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: ['configuration'],
    },
  ],
};

const TABS_FOR_SHOW = {};

export function tabsForAuthSection([authMethodModel, sectionType = 'authSettings', paths]) {
  let tabs;
  if (sectionType === 'authSettings') {
    tabs = (TABS_FOR_SETTINGS[authMethodModel.type] || []).slice();
    tabs.push({
      label: 'Method Options',
      route: ROUTES.VAULT_CLUSTER_SETTINGS_AUTH_CONFIGURE_SECTION,
      routeParams: [authMethodModel.id, 'options'],
    });
    return tabs;
  }
  if (paths || authMethodModel.paths) {
    if (authMethodModel.paths) {
      paths = authMethodModel.paths.paths.filter((path) => path.navigation);
    }

    // TODO: we're unsure if we actually need compact here
    // but are leaving it just in case OpenAPI ever returns an empty thing
    tabs = paths.compact().map((path) => {
      return {
        label: capitalize(pluralize(path.itemName)),
        route: ROUTES.VAULT_CLUSTER_ACCESS_METHOD_ITEM_LIST,
        routeParams: [path.itemType],
      };
    });
  } else {
    tabs = (TABS_FOR_SHOW[authMethodModel.type] || []).slice();
  }
  tabs.push({
    label: 'Configuration',
    route: ROUTES.VAULT_CLUSTER_ACCESS_METHOD_SECTION,
    routeParams: ['configuration'],
  });

  return tabs.map((tab) => ({
    label: tab.label,
    route: tab.route,
    routeParams: [authMethodModel.id, ...tab.routeParams],
  }));
}

export default buildHelper(tabsForAuthSection);
