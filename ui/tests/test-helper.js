import Application from '../app';
import config from '../config/environment';
import { setApplication } from '@ember/test-helpers';
import { start } from 'ember-qunit';

import './helpers/auth-login';
//import logout from './helpers/auth-logout';
//import enableSecret from './helpers/mount-secret-backend';
//import pollCluster from './helpers/poll-cluster';
//import { useNativeEvents } from 'ember-cli-page-object/extend';
//import './helpers/flash-message';

//import registerClipboardHelpers from './helpers/ember-cli-clipboard';
//registerClipboardHelpers();

setApplication(Application.create(config.APP));
//useNativeEvents();

start();
