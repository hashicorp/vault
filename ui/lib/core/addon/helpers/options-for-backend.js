/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';

// TODO move all pki related logic to its ember engine
// then we can remove use of SecretListHeader there, and use self-managed component like TabPageHeader in k8

const DEFAULT_DISPLAY = {
  searchPlaceholder: 'Filter secrets',
  item: 'secret',
  create: 'Create secret',
  navigateTree: true,
  editComponent: 'secret-edit',
  listItemPartial: 'secret-list/item',
};

const SECRET_BACKENDS = {
  aws: {
    displayName: 'AWS',
    searchPlaceholder: 'Filter roles',
    item: 'role',
    create: 'Create role',
    navigateTree: false,
    editComponent: 'role-aws-edit',
    listItemPartial: 'secret-list/aws-role-item',
  },
  ssh: {
    displayName: 'SSH',
    searchPlaceholder: 'Filter roles',
    item: 'role',
    create: 'Create role',
    navigateTree: false,
    editComponent: 'role-ssh-edit',
    listItemPartial: 'secret-list/ssh-role-item',
  },
  database: {
    displayName: 'Database',
    navigateTree: false,
    listItemPartial: 'secret-list/database-list-item',
    hasOverview: true,
    tabs: [
      {
        name: 'connection',
        label: 'Connections',
        searchPlaceholder: 'Filter connections',
        item: 'connection',
        create: 'Create connection',
        editComponent: 'database-connection',
        checkCapabilitiesPath: 'config',
      },
      {
        name: 'role',
        modelPrefix: 'role/',
        label: 'Roles',
        searchPlaceholder: 'Filter roles',
        item: 'role',
        create: 'Create role',
        tab: 'role',
        editComponent: 'database-role-edit',
        checkCapabilitiesPath: 'roles',
      },
    ],
  },
  keymgmt: {
    displayName: 'Key Management',
    navigateTree: false,
    listItemPartial: 'secret-list/item',
    tabs: [
      {
        name: 'key',
        label: 'Keys',
        searchPlaceholder: 'Filter keys',
        item: 'key',
        create: 'Create key',
        editComponent: 'keymgmt/key-edit',
      },
      {
        name: 'provider',
        modelPrefix: 'provider/',
        label: 'Providers',
        searchPlaceholder: 'Filter providers',
        item: 'provider',
        create: 'Create provider',
        tab: 'provider',
        editComponent: 'keymgmt/provider-edit',
      },
    ],
  },
  transform: {
    displayName: 'Transformation',
    navigateTree: false,
    listItemPartial: 'secret-list/transform-list-item',
    firstStep: `To use transform, you'll need to create a transformation and a role.`,
    tabs: [
      {
        name: 'transformations',
        label: 'Transformations',
        searchPlaceholder: 'Filter transformations',
        item: 'transformation',
        create: 'Create transformation',
        editComponent: 'transformation-edit',
        listItemPartial: 'secret-list/transform-transformation-item',
      },
      {
        name: 'role',
        modelPrefix: 'role/',
        label: 'Roles',
        searchPlaceholder: 'Filter roles',
        item: 'role',
        create: 'Create role',
        tab: 'role',
        editComponent: 'transform-role-edit',
      },
      {
        name: 'template',
        modelPrefix: 'template/',
        label: 'Templates',
        searchPlaceholder: 'Filter templates',
        item: 'template',
        create: 'Create template',
        tab: 'template',
        editComponent: 'transform-template-edit',
      },
      {
        name: 'alphabet',
        modelPrefix: 'alphabet/',
        label: 'Alphabets',
        searchPlaceholder: 'Filter alphabets',
        item: 'alphabet',
        create: 'Create alphabet',
        tab: 'alphabet',
        editComponent: 'alphabet-edit',
      },
    ],
  },
  transit: {
    searchPlaceholder: 'Filter keys',
    item: 'key',
    create: 'Create key',
    navigateTree: false,
    editComponent: 'transit-edit',
    listItemPartial: 'secret-list/item',
    firstStep: `To use transit, you'll need to create a key`,
  },
};

export function optionsForBackend(backend, tab) {
  const selected = SECRET_BACKENDS[backend];
  let backendOptions;
  if (selected && selected.tabs) {
    const tabData = selected.tabs.find((t) => t.name === tab || t.modelPrefix === tab) || selected.tabs[0];
    backendOptions = { ...selected, ...tabData };
  } else if (selected) {
    backendOptions = selected;
  } else {
    backendOptions = { ...DEFAULT_DISPLAY, displayName: backend === 'kv' ? 'KV' : capitalize(backend) };
  }
  return backendOptions;
}

export default buildHelper(function ([backend, tab]) {
  return optionsForBackend(backend, tab);
});
