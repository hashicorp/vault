import { create, collection, visitable, clickable, isPresent } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

export default create({
  visit: visitable('/vault/secrets/:backend/list/:id'),
  visitRoot: visitable('/vault/secrets/:backend/list'),
  create: clickable('[data-test-secret-create]'),
  createIsPresent: isPresent('[data-test-secret-create]'),
  configure: clickable('[data-test-secret-backend-configure]'),
  configureIsPresent: isPresent('[data-test-secret-backend-configure]'),

  tabs: collection('[data-test-tab]'),
  secrets: collection('[data-test-secret-link]', {
    menuToggle: clickable('[data-test-popup-menu-trigger]'),
  }),
  menuItems: collection('.ember-basic-dropdown-content li', {
    testContainer: '#ember-testing',
  }),

  backendIsEmpty: getter(function() {
    return this.secrets.length === 0;
  }),
});
