/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const policiesEgpListViewConfig: ListViewDisplayConfig = {
  title: 'EGP policies',
  badge: 'Sentinel',
  description:
    'Use Sentinel to specify policies as code that apply to discrete API paths and enforce organizational compliance standards.',
  breadcrumbs: [], // set by route model() before returning
  primaryAction: {
    label: 'Create EGP policy',
    route: 'vault.cluster.policies.egp.create',
    icon: 'plus',
  },
  filter: {
    type: 'text',
    placeholder: 'Filter policies',
    ariaLabel: 'Filter policies',
    filterKey: 'name',
  },
  columns: [
    {
      key: 'name',
      label: 'Policy name',
      customTableItem: true,
      isSortable: true,
    },
    {
      key: 'popupMenu',
      label: 'Actions',
      width: '10%',
    },
  ],
  noDataTitle: 'No EGP policies yet',
  noDataDescription: 'A list of policies will be listed here. Create your first EGP policy to get started.',
  filteredEmptyTitle: 'No results for',
};

export default policiesEgpListViewConfig;
