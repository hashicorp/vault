import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import enginesPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | console', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test("refresh reloads the current route's data", async function(assert) {
    let numEngines;
    await enginesPage.visit();
    numEngines = enginesPage.rows.length;
    await enginesPage.consoleToggle();
    let now = Date.now();
    for (let num of [1, 2, 3]) {
      let inputString = `write sys/mounts/${now + num} type=kv`;
      await enginesPage.console.consoleInput(inputString);
      await enginesPage.console.enter();
    }
    await enginesPage.console.consoleInput('refresh');
    await enginesPage.console.enter();
    assert.equal(enginesPage.rows.length, numEngines + 3, 'new engines were added to the page');
  });
});
