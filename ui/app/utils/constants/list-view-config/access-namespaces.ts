/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const accessNamespacesListViewConfig: ListViewDisplayConfig = {
  title: 'Namespaces',
  description:
    'Create logically separated, multi-tenant environments so teams can manage secrets, policies, and authentication methods independently.',
  breadcrumbs: [], // set by route model() before returning
  columns: [
    { key: 'id', label: 'Namespace' },
    { key: 'popupMenu', label: 'Action', width: '10%' },
  ],
  badge: '', // set by route model() before returning
  badgeIcon: 'org',
  noDataTitle: 'No namespaces yet',
  noDataDescription: 'Your namespaces will be listed here. Add a namespace to get started.',
  filteredEmptyTitle: 'No results for',
  primaryAction: {
    label: 'Create namespace',
    route: 'vault.cluster.access.namespaces.create',
    icon: 'plus',
  },
  filter: {
    type: 'text',
    placeholder: 'Filter by namespace path',
    ariaLabel: 'Filter by namespace path',
    filterKey: 'id',
  },
};

export default accessNamespacesListViewConfig;
