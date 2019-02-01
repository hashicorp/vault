import { later, run } from '@ember/runloop';
import EmberObject, { computed } from '@ember/object';
import Evented from '@ember/object/evented';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Pretender from 'pretender';
import { create } from 'ember-cli-page-object';
import form from '../../pages/components/auth-jwt';

const component = create(form);
const fakeWindow = EmberObject.extend(Evented, {
  screen: computed(function() {
    return {
      height: 200,
      width: 200,
    };
  }),
  localStorage: computed(function() {
    return {
      removeItem: sinon.stub(),
    };
  }),
  open: sinon.stub(),
  closed: false,
  close() {
    this.set('closed', true);
  },
});

let response = {
  auth: {
    client_token: 'token',
  },
};

let adapter = {
  exchangeOIDC() {
    return resolve(response);
  },
};

const storeStub = Service.extend({
  adapterFor() {
    return adapter;
  },
});

const routerStub = Service.extend({
  urlFor() {
    return 'http://example.com';
  },
});

const renderIt = async (context, path) => {
  let handler = e => {
    if (e) e.preventDefault();
  };
  context.set('window', fakeWindow.create());
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
    this.owner.register('service:store', storeStub);
    this.server = new Pretender(function() {
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
        return [400, { 'Content-Type': 'application/json' }, JSON.stringify({ errors: ['nope'] })];
      });
    });
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });
  test('jwt: it renders', async function(assert) {
    await renderIt(this);
    assert.ok(component.jwtPresent, 'renders jwt field');
    assert.ok(component.rolePresent, 'renders jwt field');
    assert.equal(this.server.handledRequests.length, 0, 'no requests made when there is no path set');

    this.set('selectedAuthPath', 'foo');

    await settled();
    assert.equal(
      this.server.handledRequests[0].url,
      '/v1/auth/foo/oidc/auth_url',
      'requests when path is set'
    );
  });

  test('', async function() {});
});
