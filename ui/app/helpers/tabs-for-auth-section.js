/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { pluralize } from 'ember-inflector';
import { capitalize } from '@ember/string';

const TABS_FOR_SETTINGS = {
  aws: [
    {
      label: 'Client',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['client'],
    },
    {
      label: 'Identity Allow List Tidy',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['identity-accesslist'],
    },
    {
      label: 'Role Tag Deny List Tidy',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['roletag-denylist'],
    },
  ],
  azure: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  github: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  gcp: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  jwt: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  oidc: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  kubernetes: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  ldap: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  okta: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
      routeParams: ['configuration'],
    },
  ],
  radius: [
    {
      label: 'Configuration',
      route: 'vault.cluster.settings.auth.configure.section',
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
      route: 'vault.cluster.settings.auth.configure.section',
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
        route: 'vault.cluster.access.method.item.list',
        routeParams: [path.itemType],
      };
    });
  } else {
    tabs = (TABS_FOR_SHOW[authMethodModel.type] || []).slice();
  }
  tabs.push({
    label: 'Configuration',
    route: 'vault.cluster.access.method.section',
    routeParams: ['configuration'],
  });

  return tabs.map((tab) => ({
    label: tab.label,
    route: tab.route,
    routeParams: [authMethodModel.id, ...tab.routeParams],
  }));
}

export default buildHelper(tabsForAuthSection);
