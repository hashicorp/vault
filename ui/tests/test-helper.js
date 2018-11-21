import Application from '../app';
import config from '../config/environment';
import { setApplication } from '@ember/test-helpers';
import { start } from 'ember-qunit';
import { useNativeEvents } from 'ember-cli-page-object/extend';

useNativeEvents();

setApplication(Application.create(config.APP));

start();
