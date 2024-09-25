/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// add all handlers here
// individual lookup done in mirage config
import base from './base';
import chrootNamespace from './chroot-namespace';
import customMessages from './custom-messages';
import clients from './clients';
import database from './database';
import hcpLink from './hcp-link';
import kms from './kms';
import kubernetes from './kubernetes';
import ldap from './ldap';
import mfaConfig from './mfa-config';
import mfaLogin from './mfa-login';
import oidcConfig from './oidc-config';
import reducedDisclosure from './reduced-disclosure';
import sync from './sync';
import replication from './replication';

export {
  base,
  chrootNamespace,
  clients,
  database,
  hcpLink,
  kms,
  kubernetes,
  ldap,
  mfaConfig,
  mfaLogin,
  oidcConfig,
  reducedDisclosure,
  customMessages,
  sync,
  replication,
};
