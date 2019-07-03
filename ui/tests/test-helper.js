import Application from '../app';
import config from '../config/environment';
import { setApplication } from '@ember/test-helpers';
import { start } from 'ember-qunit';
import './helpers/flash-message';
import preloadAssets from 'ember-asset-loader/test-support/preload-assets';
import manifest from 'vault/config/asset-manifest';

preloadAssets(manifest).then(() => {
  setApplication(Application.create(config.APP));
  start({
    setupTestIsolationValidation: true,
  });
});
