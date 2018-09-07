import Application from '../../app';
import config from '../../config/environment';
import { merge } from '@ember/polyfills';
import { run } from '@ember/runloop';

import './auth-login';
import './auth-logout';
import './mount-secret-backend';
import './poll-cluster';

import registerClipboardHelpers from '../helpers/ember-cli-clipboard';
registerClipboardHelpers();

export default function startApp(attrs) {
  let attributes = merge({}, config.APP);
  attributes.autoboot = true;
  attributes = merge(attributes, attrs); // use defaults, but you can override;

  return run(() => {
    let application = Application.create(attributes);
    application.setupForTesting();
    application.injectTestHelpers();
    return application;
  });
}
