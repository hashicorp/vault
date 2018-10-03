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
    await enginesPage.visit();
    let numEngines = enginesPage.rows.length;
    await enginesPage.consoleToggle();
    let now = Date.now();
    for (let num of [1, 2, 3]) {
      let inputString = `write sys/mounts/${now + num} type=kv`;
      await enginesPage.console.runCommands(inputString);
    }
    await enginesPage.console.runCommands('refresh');
    assert.equal(enginesPage.rows.length, numEngines + 3, 'new engines were added to the page');
  });
});
