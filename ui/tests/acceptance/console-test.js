import { module, test } from 'qunit';
import { create } from 'ember-cli-page-object';
import { later } from '@ember/runloop';
import { setupApplicationTest } from 'ember-qunit';
import enginesPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

module('Acceptance | console', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test("refresh reloads the current route's data", async function(assert) {
    await enginesPage.visit();
    let numEngines = enginesPage.rows.length;
    await consoleComponent.toggle();
    let now = Date.now();
    for (let num of [1, 2, 3]) {
      let inputString = `write sys/mounts/${now + num} type=kv`;
      await consoleComponent.runCommands(inputString);
    }
    await consoleComponent.runCommands('refresh');
    assert.equal(enginesPage.rows.length, numEngines + 3, 'new engines were added to the page');
  });

  test('fullscreen command expands the cli panel', async function(assert) {
    await consoleComponent.toggle();
    await consoleComponent.runCommands('fullscreen');

    // have to wrap in a later so that we can wait for the CSS transition to finish
    await later(() => {
      let consoleEle = document.querySelector('[data-test-component="console/ui-panel"]');

      assert.equal(
        consoleEle.offsetHeight,
        window.innerHeight,
        'fullscreen is the same height as the window'
      );

      assert.equal(consoleEle.offsetTop, 0, 'fullscreen is aligned to the top of window');
    }, 300);
  });

  test('array output is correctly formatted', async function(assert) {
    await consoleComponent.toggle();
    await consoleComponent.runCommands('read -field=policies /auth/token/lookup-self');

    // have to wrap in a later so that we can wait for the CSS transition to finish
    await later(() => {
      let consoleOut = document.querySelector('.console-ui-output>pre').innerText;

      assert.notOk(consoleOut.includes('function(){'));
      assert.equal(consoleOut, '["root"]');
    }, 300);
  });

  test('number output is correctly formatted', async function(assert) {
    await consoleComponent.toggle();
    await consoleComponent.runCommands('read -field=creation_time /auth/token/lookup-self');

    // have to wrap in a later so that we can wait for the CSS transition to finish
    await later(() => {
      let consoleOut = document.querySelector('.console-ui-output>pre').innerText;
      assert.ok(consoleOut.match(/^\d+$/).length == 1);
    }, 300);
  });

  test('boolean output is correctly formatted', async function(assert) {
    await consoleComponent.toggle();
    await consoleComponent.runCommands('read -field=orphan /auth/token/lookup-self');

    // have to wrap in a later so that we can wait for the CSS transition to finish
    await later(() => {
      let consoleOut = document.querySelector('.console-ui-output>pre').innerText;
      assert.ok(consoleOut.match(/^(true|false)$/g).length == 1);
    }, 300);
  });
});
