/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ListViewDisplayConfig } from 'core/types/list-view-config';

export const secretsBackendsListViewConfig: ListViewDisplayConfig = {
  title: 'Secrets engines',
  description:
    'View and manage your configured secrets engines in the current cluster, ranging from key value store (kv) to dynamic database credentials and more.',
  breadcrumbs: [], // filled by route setupController
  primaryAction: {
    label: 'Enable new engine',
    route: 'vault.cluster.secrets.enable',
    icon: 'plus',
  },
  columns: [
    {
      key: 'path',
      label: 'Engine path',
      isSortable: true,
      valueType: 'link',
      routeKey: 'backendLink',
      unavailableTooltip:
        'The UI only supports configuration views for these secret engines. The CLI must be used to manage other engine resources.',
    },
    {
      key: 'engineType',
      label: 'Engine Type',
      isSortable: true,
      valueType: 'icon-text',
      iconKey: 'icon',
      textKey: 'displayName',
    },
    {
      key: 'accessor',
      label: 'Accessor',
      valueType: 'header-tooltip',
      headerTooltip:
        'An accessor is a stable, system-generated identifier used to represent a secret engine path.',
    },
    {
      key: 'description',
      label: 'Description',
      valueType: 'text',
    },
    {
      key: 'version',
      label: 'Version',
      isSortable: true,
      valueType: 'text',
    },
    {
      key: 'popupMenu',
      label: 'Actions',
      width: '10%',
    },
  ],
  rowActions: [
    {
      label: 'View configuration',
      dataTest: 'view-configuration',
      kind: 'route',
      routeKey: 'backendConfigurationLink',
    },
    {
      label: 'Delete engine path',
      dataTest: 'delete-engine-path',
      kind: 'modal',
      color: 'critical',
      modal: {
        title: 'Delete secrets engine',
        titleItemDisplayKey: 'path',
        body: 'The following data will be permanently deleted:',
        bodyItems: ['secrets engine', 'Engine configuration and third-party integrations'],
        itemDisplayKey: 'path',
        color: 'critical',
        confirmButtonText: 'Delete engine',
        cancelButtonText: 'Cancel',
        confirmText: 'delete-engine',
        confirmLabel: 'Confirm deletion',
        apiCall: { service: 'sys', method: 'mountsDisableSecretsEngine', argKey: 'id' },
      },
    },
  ],
  noDataTitle: 'No secrets engines',
  noDataDescription: 'No secrets engines have been configured yet.',
  filteredEmptyTitle: 'No results for',
};

export default secretsBackendsListViewConfig;
