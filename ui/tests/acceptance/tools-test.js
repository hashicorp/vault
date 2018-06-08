import Pretender from 'pretender';
import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import { toolsActions } from 'vault/helpers/tools-actions';

moduleForAcceptance('Acceptance | tools', {
  beforeEach() {
    return authLogin();
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
    let rows = document.querySelectorAll('[data-test-tools="token-lookup-row"]');
    assert.dom(rows[0]).hasText(/Creation path/, 'show creation path row');
    assert.dom(rows[1]).hasText(/Creation time/, 'show creation time row');
    assert.dom(rows[2]).hasText(/Creation TTL/, 'show creation ttl row');
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
    assert.dom('[data-test-tools-input="bytes"]').hasValue('32', 'defaults to 32 bytes');
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
    assert
      .dom('[data-test-tools-input="sum"]')
      .hasValue('LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564=', 'hashes the data, encodes input');
  });
  click('[data-test-tools-back]');
  fillIn('[data-test-tools-input="hash-input"]', 'e2RhdGE6ImZvbyJ9');

  click('[data-test-tools-submit]');
  andThen(() => {
    assert
      .dom('[data-test-tools-input="sum"]')
      .hasValue('JmSi2Hhbgu2WYOrcOyTqqMdym7KT3sohCwAwaMonVrc=', 'hashes the data, passes b64 input through');
  });
});

const AUTH_RESPONSE = {
  request_id: '39802bc4-235c-2f0b-87f3-ccf38503ac3e',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: null,
  wrap_info: null,
  warnings: null,
  auth: {
    client_token: 'ecfc2758-588e-981d-50f4-a25883bbf03c',
    accessor: '6299780b-f2b2-1a3f-7b83-9d3d67629249',
    policies: ['root'],
    metadata: null,
    lease_duration: 0,
    renewable: false,
    entity_id: '',
  },
};

test('ensure unwrap with auth block works properly', function(assert) {
  this.server = new Pretender(function() {
    this.post('/v1/sys/wrapping/unwrap', response => {
      return [response, { 'Content-Type': 'application/json' }, JSON.stringify(AUTH_RESPONSE)];
    });
  });
  visit('/vault/tools');
  //unwrap
  click('[data-test-tools-action-link="unwrap"]');
  andThen(() => {
    fillIn('[data-test-tools-input="wrapping-token"]', 'sometoken');
  });
  click('[data-test-tools-submit]');
  andThen(() => {
    assert.deepEqual(
      JSON.parse(find('.CodeMirror').get(0).CodeMirror.getValue()),
      AUTH_RESPONSE.auth,
      'unwrapped data equals input data'
    );
    this.server.shutdown();
  });
});
