// add all handlers here
// individual lookup done in mirage config
import base from './base';
import mfa from './mfa';
import activity from './activity';
import clients from './clients';
import db from './db';
import kms from './kms';

export { base, activity, mfa, clients, db, kms };
