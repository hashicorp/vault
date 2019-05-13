import { create, visitable, collection, clickable, text } from 'ember-cli-page-object';
import uiPanel from 'vault/tests/pages/components/console/ui-panel';

export default create({
  console: uiPanel,
  consoleToggle: clickable('[data-test-console-toggle]'),
  visit: visitable('/vault/secrets'),
  rows: collection('[data-test-secret-backend-row]', {
    path: text('[data-test-secret-path]'),
    menu: clickable('[data-test-popup-menu-trigger]'),
  }),
  configLink: clickable('[data-test-engine-config]', {
    testContainer: '#ember-testing',
  }),
  disableButton: clickable('[data-test-confirm-action-trigger]', {
    testContainer: '#ember-testing',
  }),
  confirmDisable: clickable('[data-test-confirm-button]', {
    testContainer: '#ember-testing',
  }),
});
