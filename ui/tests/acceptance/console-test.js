import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import enginesPage from 'vault/tests/pages/secrets/backends';

moduleForAcceptance('Acceptance | console', {
  beforeEach() {
    return authLogin();
  },
});

test('refresh reloads the current route\'s data', function(assert) {
  let numEngines;
  enginesPage.visit();
  andThen(() => {
    numEngines = enginesPage.rows().count;
    enginesPage.consoleToggle();
    let now = Date.now();
    [1, 2, 3].forEach(num => {
      let inputString = `write sys/mounts/${now + num} type=kv`;
      enginesPage.console.consoleInput(inputString);
      enginesPage.console.enter();
    });
  });
  andThen(() => {
    enginesPage.console.consoleInput('refresh');
    enginesPage.console.enter();
  });
  andThen(() => {
    assert.equal(enginesPage.rows().count, numEngines + 3, 'new engines were added to the page');
  });
});
