/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const accessIdentityGroupsListViewConfig: ListViewDisplayConfig = {
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
    placeholder: 'Search by id or name',
    ariaLabel: 'Search groups by id or name',
    filterKey: 'name',
  },
  columns: [
    {
      key: 'name',
      label: 'Group name',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'id',
      label: 'Group ID',
      valueType: 'copy',
    },
    {
      key: 'popupMenu',
      label: 'Actions',
      width: '10%',
    },
  ],
  noDataTitle: 'No groups yet',
  noDataDescription:
    'A list of groups in this namespace will be listed here. Create your first group to get started.',
  filteredEmptyTitle: 'No results for',
};

export default accessIdentityGroupsListViewConfig;
