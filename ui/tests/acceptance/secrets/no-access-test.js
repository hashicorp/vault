import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import listPage from 'vault/tests/pages/secrets/backend/list';
import { startMirage } from 'vault/initializers/ember-cli-mirage';
import Ember from 'ember';

let adapterException;
let loggerError;

moduleForAcceptance('Acceptance | secrets/secret/secret error', {
  beforeEach() {
    this.server = startMirage();
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
    this.server.shutdown();
  },
});

test('it shows a warning if dont have access to the secrets list', function(assert) {
  listPage.visitRoot({ backend: 'secret' });
  andThen(() => {
    assert.ok(find('[data-test-sys-mounts-warning]').length, 'shows the warning for sys/mounts');
  });
});
