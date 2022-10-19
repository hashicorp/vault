// add all handlers here
// individual lookup done in mirage config
import base from './base';
import activity from './activity';
import clients from './clients';
import db from './db';
import kms from './kms';
import mfaConfig from './mfa-config';
import mfaLogin from './mfa-login';
import oidcConfig from './oidc-config';
import hcpLink from './hcp-link';

export { base, activity, clients, db, kms, mfaConfig, mfaLogin, oidcConfig, hcpLink };
