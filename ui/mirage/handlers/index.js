/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// add all handlers here
// individual lookup done in mirage config
import base from './base';
import clients from './clients';
import db from './db';
import kms from './kms';
import mfaConfig from './mfa-config';
import mfaLogin from './mfa-login';
import oidcConfig from './oidc-config';
import hcpLink from './hcp-link';
import kubernetes from './kubernetes';
import ldap from './ldap';

export { base, clients, db, kms, mfaConfig, mfaLogin, oidcConfig, hcpLink, kubernetes, ldap };
