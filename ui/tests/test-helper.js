import Application from '../app';
import config from '../config/environment';
import { setApplication } from '@ember/test-helpers';
import { start } from 'ember-qunit';

import './helpers/auth-login';
import './helpers/auth-logout';
import './helpers/mount-secret-backend';
import './helpers/poll-cluster';
import { useNativeEvents } from 'ember-cli-page-object/extend';
import './helpers/flash-message';

import registerClipboardHelpers from './helpers/ember-cli-clipboard';
registerClipboardHelpers();

setApplication(Application.create(config.APP));
//application.setupForTesting();
//application.injectTestHelpers();
useNativeEvents();

start();
