/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// add all handlers here
// individual lookup done in mirage config
import base from './base';
import chrootNamespace from './chroot-namespace';
import clients from './clients';
import customMessages from './custom-messages';
import customLogin from './custom-login';
import database from './database';
import hcpLink from './hcp-link';
import kms from './kms';
import kubernetes from './kubernetes';
import ldap from './ldap';
import mfaConfig from './mfa-config';
import mfaLogin from './mfa-login';
import oidcConfig from './oidc-config';
import reducedDisclosure from './reduced-disclosure';
import replication from './replication';
import sync from './sync';

export {
  base,
  chrootNamespace,
  clients,
  customLogin,
  customMessages,
  database,
  hcpLink,
  kms,
  kubernetes,
  ldap,
  mfaConfig,
  mfaLogin,
  oidcConfig,
  reducedDisclosure,
  replication,
  sync,
};
