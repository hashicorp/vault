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
  customLogin: apiPath`sys/config/ui/login/default-auth/${'id'}`,
  customMessages: apiPath`sys/config/ui/custom-messages/${'id'}`,
  syncActivate: apiPath`sys/activation-flags/secrets-sync/activate`,
  syncDestination: apiPath`sys/sync/destinations/${'type'}/${'name'}`,
  syncSetAssociation: apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/set`,
  syncRemoveAssociation: apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/remove`,
  kvConfig: apiPath`${'path'}/config`,
  kvMetadata: apiPath`${'backend'}/metadata/${'path'}`,
  authMethodConfig: apiPath`auth/${'path'}/config`,
  authMethodConfigAws: apiPath`auth/${'path'}/config/client`,
  authMethodDelete: apiPath`sys/auth/${'path'}`,
  pkiRevoke: apiPath`${'backend'}/revoke`,
  pkiConfigAcme: apiPath`${'backend'}/config/acme`,
  pkiConfigCluster: apiPath`${'backend'}/config/cluster`,
  pkiConfigCrl: apiPath`${'backend'}/config/crl`,
  pkiConfigUrls: apiPath`${'backend'}/config/urls`,
  pkiIssuersImportBundle: apiPath`${'backend'}/issuers/import/bundle`,
  pkiIssuersGenerateRoot: apiPath`${'backend'}/issuers/generate/root/${'type'}`,
  pkiIssuersGenerateIntermediate: apiPath`${'backend'}/issuers/generate/intermediate/${'type'}`,
  pkiIssuersCrossSign: apiPath`${'backend'}/issuers/cross-sign`,
  pkiIssuer: apiPath`${'backend'}/issuer/${'issuerId'}`,
  pkiIssuerSignIntermediate: apiPath`${'backend'}/issuer/${'issuerId'}/sign-intermediate`,
  pkiRoot: apiPath`${'backend'}/root`,
  pkiRootRotate: apiPath`${'backend'}/root/rotate/${'type'}`,
  pkiIntermediateCrossSign: apiPath`${'backend'}/intermediate/cross-sign`,
  pkiKey: apiPath`${'backend'}/key/${'keyId'}`,
  pkiKeysGenerate: apiPath`${'backend'}/keys/generate`,
  pkiKeysImport: apiPath`${'backend'}/keys/import`,
  pkiRole: apiPath`${'backend'}/roles/${'id'}`,
  pkiIssue: apiPath`${'backend'}/issue/${'id'}`,
  pkiSign: apiPath`${'backend'}/sign/${'id'}`,
  pkiSignVerbatim: apiPath`${'backend'}/sign-verbatim/${'id'}`,
  ldapStaticRole: apiPath`${'backend'}/static-role/${'name'}`,
  ldapDynamicRole: apiPath`${'backend'}/role/${'name'}`,
  ldapRotateStaticRole: apiPath`${'backend'}/rotate-role/${'name'}`,
  ldapStaticRoleCreds: apiPath`${'backend'}/static-cred/${'name'}`,
  ldapDynamicRoleCreds: apiPath`${'backend'}/creds/${'name'}`,
  ldapLibrary: apiPath`${'backend'}/library/${'name'}`,
  ldapLibraryCheckOut: apiPath`${'backend'}/library/${'name'}/check-out`,
  ldapLibraryCheckIn: apiPath`${'backend'}/library/${'name'}/check-in`,
  kubernetesRole: apiPath`${'backend'}/role/${'name'}`,
  kubernetesCreds: apiPath`${'backend'}/creds/${'name'}`,
  kmipScope: apiPath`${'backend'}/scopes/${'name'}`,
  kmipRole: apiPath`${'backend'}/scopes/${'scope'}/roles/${'name'}`,
  kmipCredentialsRevoke: apiPath`${'backend'}/scope/${'scope'}/role/${'role'}/credentials/revoke`,
};
