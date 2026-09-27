/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const accessIdentityGroupAliasesListViewConfig: ListViewDisplayConfig = {
  title: 'Groups',
  description:
    'Create and name logical collections of entities to simplify policy management and permission scaling across your organization.',
  breadcrumbs: [], // set by route model() before returning
  primaryAction: {
    label: 'Create group',
    route: 'vault.cluster.access.identity.groups.create',
    icon: 'plus',
  },
  filter: {
    type: 'text',
    placeholder: 'Search by alias id or name',
    ariaLabel: 'Search group aliases by id or name',
    filterKey: 'name',
  },
  columns: [
    {
      key: 'name',
      label: 'Aliases name',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'id',
      label: 'Aliases ID',
      valueType: 'copy',
    },
    {
      key: 'mount_type',
      label: 'Mount type',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'popupMenu',
      label: 'Actions',
      width: '10%',
    },
  ],
  noDataTitle: 'No group aliases yet',
  noDataDescription:
    'A list of group aliases in this namespace will be listed here. Choose one of the groups and click "Add alias" to get started.',
  filteredEmptyTitle: 'No results for',
};

export default accessIdentityGroupAliasesListViewConfig;
