import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import enginesPage from 'vault/tests/pages/secrets/backends';

module('Acceptance | console', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test("refresh reloads the current route's data", function(assert) {
    let numEngines;
    enginesPage.visit();
    numEngines = enginesPage.rows().count;
    enginesPage.consoleToggle();
    let now = Date.now();
    [1, 2, 3].forEach(num => {
      let inputString = `write sys/mounts/${now + num} type=kv`;
      enginesPage.console.consoleInput(inputString);
      enginesPage.console.enter();
    });
    enginesPage.console.consoleInput('refresh');
    enginesPage.console.enter();
    assert.equal(enginesPage.rows().count, numEngines + 3, 'new engines were added to the page');
  });
});
