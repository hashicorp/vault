import { later, run } from '@ember/runloop';
import EmberObject, { computed } from '@ember/object';
import Evented from '@ember/object/evented';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Pretender from 'pretender';
import { create } from 'ember-cli-page-object';
import form from '../../pages/components/auth-jwt';
import { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS } from 'vault/components/auth-jwt';

const component = create(form);
const fakeWindow = EmberObject.extend(Evented, {
  init() {
    this._super(...arguments);
    this.__proto__.on('close', () => {
      this.set('closed', true);
    });
  },
  screen: computed(function() {
    return {
      height: 600,
      width: 500,
    };
  }),
  localStorage: computed(function() {
    return {
      removeItem: sinon.stub(),
    };
  }),
  closed: false,
});

fakeWindow.reopen({
  open() {
    return fakeWindow.create();
  },

  close() {
    fakeWindow.prototype.trigger('close');
  },
});

const OIDC_AUTH_RESPONSE = {
  auth: {
    client_token: 'token',
  },
};

const routerStub = Service.extend({
  urlFor() {
    return 'http://example.com';
  },
});

const renderIt = async (context, path = 'jwt') => {
  let handler = (data, e) => {
    if (e && e.preventDefault) e.preventDefault();
  };
  let fake = fakeWindow.create();
  sinon.spy(fake, 'open');
  context.set('window', fake);
  context.set('handler', sinon.spy(handler));
  context.set('roleName', '');
  context.set('selectedAuthPath', path);

  await render(hbs`
    <AuthJwt
      @window={{window}}
      @roleName={{roleName}}
      @selectedAuthPath={{selectedAuthPath}}
      @onError={{action (mut error)}}
      @onLoading={{action (mut isLoading)}}
      @onToken={{action (mut token)}}
      @onNamespace={{action (mut namespace)}}
      @onSelectedAuth={{action (mut selectedAuth)}}
      @onSubmit={{action handler}}
      @onRoleName={{action (mut roleName)}}
    />
    `);
};
module('Integration | Component | auth jwt', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.register('service:router', routerStub);
    this.server = new Pretender(function() {
      this.get('/v1/auth/:path/oidc/callback', function() {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify(OIDC_AUTH_RESPONSE)];
      });
      this.post('/v1/auth/:path/oidc/auth_url', request => {
        let body = JSON.parse(request.requestBody);
        if (body.role === 'test') {
          return [
            200,
            { 'Content-Type': 'application/json' },
            JSON.stringify({
              data: {
                auth_url: 'http://example.com',
              },
            }),
          ];
        }
        if (body.role === 'okta') {
          return [
            200,
            { 'Content-Type': 'application/json' },
            JSON.stringify({
              data: {
                auth_url: 'http://okta.com',
              },
            }),
          ];
        }
        return [400, { 'Content-Type': 'application/json' }, JSON.stringify({ errors: ['nope'] })];
      });
    });
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('it renders the yield', async function(assert) {
    await render(hbs`<AuthJwt @onSubmit={{action (mut submit)}}>Hello!</AuthJwt>`);
    assert.equal(component.yieldContent, 'Hello!', 'yields properly');
  });

  test('jwt: it renders', async function(assert) {
    await renderIt(this);
    await settled();
    assert.ok(component.jwtPresent, 'renders jwt field');
    assert.ok(component.rolePresent, 'renders jwt field');
    assert.equal(this.server.handledRequests.length, 1, 'request to the default path is made');
    assert.equal(this.server.handledRequests[0].url, '/v1/auth/jwt/oidc/auth_url');
    this.set('selectedAuthPath', 'foo');
    await settled();
    assert.equal(this.server.handledRequests.length, 2, 'a second request was made');
    assert.equal(
      this.server.handledRequests[1].url,
      '/v1/auth/foo/oidc/auth_url',
      'requests when path is set'
    );
  });

  test('jwt: it calls passed action on login', async function(assert) {
    await renderIt(this);
    await component.login();
    assert.ok(this.handler.calledOnce);
  });

  test('oidc: test role: it renders', async function(assert) {
    await renderIt(this);
    await settled();
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    await settled();
    assert.notOk(component.jwtPresent, 'does not show jwt input for OIDC type login');
    assert.equal(component.loginButtonText, 'Sign in with OIDC Provider');

    await component.role('okta');
    // 1 for initial render, 1 for each time role changed = 3
    assert.equal(this.server.handledRequests.length, 4, 'fetches the auth_url when the path changes');
    assert.equal(component.loginButtonText, 'Sign in with Okta', 'recognizes auth methods with certain urls');
  });

  test('oidc: it calls window.open popup window on login', async function(assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();

    later(async () => {
      run.cancelTimers();
      await settled();
      let call = this.window.open.getCall(0);
      assert.deepEqual(
        call.args,
        [
          'http://example.com',
          'vaultOIDCWindow',
          'width=500,height=600,resizable,scrollbars=yes,top=0,left=0',
        ],
        'called with expected args'
      );
    }, 50);
    await settled();
  });

  test('oidc: it calls error handler when popup is closed', async function(assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();

    later(async () => {
      this.window.close();
      await settled();
      assert.equal(this.error, ERROR_WINDOW_CLOSED, 'calls onError with error string');
    }, 50);
    await settled();
  });

  test('oidc: storage event fires with wrong key', async function(assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    later(async () => {
      run.cancelTimers();
      this.window.trigger('storage', { key: 'wrongThing' });
      assert.equal(this.window.localStorage.removeItem.callCount, 0, 'never calls removeItem');
    }, 50);
    await settled();
  });

  test('oidc: storage event fires with correct key, wrong params', async function(assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    later(async () => {
      this.window.trigger('storage', { key: 'oidcState', newValue: JSON.stringify({}) });
      await settled();
      assert.equal(this.window.localStorage.removeItem.callCount, 1, 'calls removeItem');
      assert.equal(this.error, ERROR_MISSING_PARAMS, 'calls onError with params missing error');
    }, 50);
    await settled();
  });

  test('oidc: storage event fires with correct key, correct params', async function(assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    later(async () => {
      this.window.trigger('storage', {
        key: 'oidcState',
        newValue: JSON.stringify({
          path: 'foo',
          state: 'state',
          code: 'code',
        }),
      });
      await settled();
      assert.equal(this.selectedAuth, 'token', 'calls onSelectedAuth with token');
      assert.equal(this.token, 'token', 'calls onToken with token');
      assert.ok(this.handler.calledOnce, 'calls the onSubmit handler');
    }, 50);
    await settled();
  });
});
