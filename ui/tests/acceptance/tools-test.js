import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import { toolsActions } from 'vault/helpers/tools-actions';

moduleForAcceptance('Acceptance | tools', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
});

const DATA_TO_WRAP = JSON.stringify({ tools: 'tests' });
const TOOLS_ACTIONS = toolsActions();

/*
data-test-tools-input="wrapping-token"
data-test-tools-input="rewrapped-token"
data-test-tools="token-lookup-row"
data-test-tools-action-link=supportedAction
*/

var createTokenStore = () => {
  let token;
  return {
    set(val) {
      token = val;
    },
    get() {
      return token;
    },
  };
};
test('tools functionality', function(assert) {
  var tokenStore = createTokenStore();
  visit('/vault/tools');
  andThen(function() {
    assert.equal(currentURL(), '/vault/tools/wrap', 'forwards to the first action');
    TOOLS_ACTIONS.forEach(action => {
      assert.ok(findWithAssert(`[data-test-tools-action-link="${action}"]`), `${action} link renders`);
    });
    find('.CodeMirror').get(0).CodeMirror.setValue(DATA_TO_WRAP);
  });

  // wrap
  click('[data-test-tools-submit]');
  andThen(function() {
    tokenStore.set(find('[data-test-tools-input="wrapping-token"]').val());
    assert.ok(find('[data-test-tools-input="wrapping-token"]').val(), 'has a wrapping token');
  });

  //lookup
  click('[data-test-tools-action-link="lookup"]');
  // have to wrap this in andThen because tokenStore is sync, but fillIn is async
  andThen(() => {
    fillIn('[data-test-tools-input="wrapping-token"]', tokenStore.get());
  });
  click('[data-test-tools-submit]');
  andThen(() => {
    assert.ok(
      find('[data-test-tools="token-lookup-row"]:eq(0)').text().match(/Creation time/i),
      'show creation time row'
    );
    assert.ok(
      find('[data-test-tools="token-lookup-row"]:eq(1)').text().match(/Creation ttl/i),
      'show creation ttl row'
    );
  });

  //rewrap
  click('[data-test-tools-action-link="rewrap"]');
  andThen(() => {
    fillIn('[data-test-tools-input="wrapping-token"]', tokenStore.get());
  });
  click('[data-test-tools-submit]');
  andThen(() => {
    assert.ok(find('[data-test-tools-input="rewrapped-token"]').val(), 'has a new re-wrapped token');
    assert.notEqual(
      find('[data-test-tools-input="rewrapped-token"]').val(),
      tokenStore.get(),
      're-wrapped token is not the wrapped token'
    );
    tokenStore.set(find('[data-test-tools-input="rewrapped-token"]').val());
  });

  //unwrap
  click('[data-test-tools-action-link="unwrap"]');
  andThen(() => {
    fillIn('[data-test-tools-input="wrapping-token"]', tokenStore.get());
  });
  click('[data-test-tools-submit]');
  andThen(() => {
    assert.deepEqual(
      JSON.parse(find('.CodeMirror').get(0).CodeMirror.getValue()),
      JSON.parse(DATA_TO_WRAP),
      'unwrapped data equals input data'
    );
  });

  //random
  click('[data-test-tools-action-link="random"]');
  andThen(() => {
    assert.equal(find('[data-test-tools-input="bytes"]').val(), 32, 'defaults to 32 bytes');
  });
  click('[data-test-tools-submit]');
  andThen(() => {
    assert.ok(
      find('[data-test-tools-input="random-bytes"]').val(),
      'shows the returned value of random bytes'
    );
  });

  //hash
  click('[data-test-tools-action-link="hash"]');
  fillIn('[data-test-tools-input="hash-input"]', 'foo');
  click('[data-test-tools-b64-toggle="input"]');
  click('[data-test-tools-submit]');
  andThen(() => {
    assert.equal(
      find('[data-test-tools-input="sum"]').val(),
      'LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564=',
      'hashes the data, encodes input'
    );
  });
  click('[data-test-tools-back]');
  fillIn('[data-test-tools-input="hash-input"]', 'e2RhdGE6ImZvbyJ9');

  click('[data-test-tools-submit]');
  andThen(() => {
    assert.equal(
      find('[data-test-tools-input="sum"]').val(),
      'JmSi2Hhbgu2WYOrcOyTqqMdym7KT3sohCwAwaMonVrc=',
      'hashes the data, passes b64 input through'
    );
  });
});
