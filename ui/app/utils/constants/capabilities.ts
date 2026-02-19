/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import apiPath from 'vault/utils/api-path';

export const SUDO_PATHS = [
  'sys/seal',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
  'sys/replication/reindex',
  'sys/leases/lookup/',
];

export const SUDO_PATH_PREFIXES = ['sys/leases/revoke-prefix', 'sys/leases/revoke-force'];

export const PATH_MAP = {
  authMethodConfig: apiPath`auth/${'path'}/config`,
  authMethodConfigAws: apiPath`auth/${'path'}/config/client`,
  authMethodDelete: apiPath`sys/auth/${'path'}`,
  clientsActivityExport: apiPath`${'namespace'}/sys/internal/counters/activity/export`,
  clientsConfig: apiPath`sys/internal/counters/config`,
  customLogin: apiPath`sys/config/ui/login/default-auth/${'id'}`,
  customMessages: apiPath`sys/config/ui/custom-messages/${'id'}`,
  kmipCredentialsRevoke: apiPath`${'backend'}/scope/${'scope'}/role/${'role'}/credentials/revoke`,
  kmipRole: apiPath`${'backend'}/scopes/${'scope'}/roles/${'name'}`,
  kmipScope: apiPath`${'backend'}/scopes/${'name'}`,
  kubernetesCreds: apiPath`${'backend'}/creds/${'name'}`,
  kubernetesRole: apiPath`${'backend'}/role/${'name'}`,
  kvConfig: apiPath`${'path'}/config`,
  kvMetadata: apiPath`${'backend'}/metadata/${'path'}`,
  ldapDynamicRole: apiPath`${'backend'}/role/${'name'}`,
  ldapDynamicRoleCreds: apiPath`${'backend'}/creds/${'name'}`,
  ldapLibrary: apiPath`${'backend'}/library/${'name'}`,
  ldapLibraryCheckIn: apiPath`${'backend'}/library/${'name'}/check-in`,
  ldapLibraryCheckOut: apiPath`${'backend'}/library/${'name'}/check-out`,
  ldapRotateStaticRole: apiPath`${'backend'}/rotate-role/${'name'}`,
  ldapStaticRole: apiPath`${'backend'}/static-role/${'name'}`,
  ldapStaticRoleCreds: apiPath`${'backend'}/static-cred/${'name'}`,
  pkiCertificates: apiPath`${'backend'}/certificates`,
  pkiConfigAcme: apiPath`${'backend'}/config/acme`,
  pkiConfigAutoTidy: apiPath`${'backend'}/config/auto-tidy`,
  pkiConfigCluster: apiPath`${'backend'}/config/cluster`,
  pkiConfigCrl: apiPath`${'backend'}/config/crl`,
  pkiConfigUrls: apiPath`${'backend'}/config/urls`,
  pkiIntermediateCrossSign: apiPath`${'backend'}/intermediate/cross-sign`,
  pkiIssue: apiPath`${'backend'}/issue/${'id'}`,
  pkiIssuer: apiPath`${'backend'}/issuer/${'issuerId'}`,
  pkiIssuersCrossSign: apiPath`${'backend'}/issuers/cross-sign`,
  pkiIssuersGenerateIntermediate: apiPath`${'backend'}/issuers/generate/intermediate/${'type'}`,
  pkiIssuersGenerateRoot: apiPath`${'backend'}/issuers/generate/root/${'type'}`,
  pkiIssuerSignIntermediate: apiPath`${'backend'}/issuer/${'issuerId'}/sign-intermediate`,
  pkiIssuersImportBundle: apiPath`${'backend'}/issuers/import/bundle`,
  pkiKey: apiPath`${'backend'}/key/${'keyId'}`,
  pkiKeysGenerate: apiPath`${'backend'}/keys/generate`,
  pkiKeysImport: apiPath`${'backend'}/keys/import`,
  pkiRevoke: apiPath`${'backend'}/revoke`,
  pkiRole: apiPath`${'backend'}/roles/${'id'}`,
  pkiRoles: apiPath`${'backend'}/roles`,
  pkiRoot: apiPath`${'backend'}/root`,
  pkiRootRotate: apiPath`${'backend'}/root/rotate/${'type'}`,
  pkiSign: apiPath`${'backend'}/sign/${'id'}`,
  pkiSignVerbatim: apiPath`${'backend'}/sign-verbatim/${'id'}`,
  pkiTidy: apiPath`${'backend'}/tidy`,
  pkiTidyStatus: apiPath`${'backend'}/tidy/status`,
  syncActivate: apiPath`sys/activation-flags/secrets-sync/activate`,
  syncDestination: apiPath`sys/sync/destinations/${'type'}/${'name'}`,
  syncRemoveAssociation: apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/remove`,
  syncSetAssociation: apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/set`,
};
