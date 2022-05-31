// add all handlers here
// individual lookup done in mirage config
import base from './base';
import mfaLogin from './mfa-login';
import activity from './activity';
import clients from './clients';
import db from './db';
import kms from './kms';
import mfaConfig from './mfa-config';

export { base, activity, mfaLogin, mfaConfig, clients, db, kms };
