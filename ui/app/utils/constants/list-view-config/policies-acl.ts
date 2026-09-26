/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const policiesAclListViewConfig: ListViewDisplayConfig = {
  title: 'ACL policies',
  description:
    'Define fine-grained rules to explicitly grant or forbid access to specific paths and operations within your cluster. Because Vault is a "default deny" system, if a permission is not granted in a policy, an entity would not have permission.',
  breadcrumbs: [], // set by route model() before returning
  primaryAction: {
    label: 'Create ACL policy',
    route: 'vault.cluster.policies.acl.create',
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
  noDataTitle: 'No ACL policies yet',
  noDataDescription: 'A list of policies will be listed here. Create your first ACL policy to get started.',
  filteredEmptyTitle: 'No results for',
};

export default policiesAclListViewConfig;
