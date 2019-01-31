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

const renderIt = async context => {
  context.set('window', fakeWindow.create());
  let handler = e => {
    if (e) e.preventDefault();
  };
  context.set('handler', sinon.spy(handler));

  await context.render(hbs`
    <AuthJwt
      @window={{window}}
      @onError={{action (mut context.error)}}
      @onLoading={{action (mut context.isLoading)}}
      @onToken={{action (mut context.token)}}
      @onNamespace={{action (mut context.namespace)}}
      @onSelectedAuth={{action (mut context.selectedAuth)}}
      @onSubmit={{action handler}}
      @onRoleName={{action (mut context.roleName)}}
      @roleName={{this.roleName}}
      @selectedAuthPath={{'foo'}}
      @form={{hash options=(component 'auth-form-options' customPath=context.customPath onPathChange=(action (mut context.customPath))  selectedAuthIsPath=context.selectedAuthIsPath )}}
    />
    `);
};
module('Integration | Component | auth jwt', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {});

  hooks.afterEach(function() {});
  test('', async function() {});

  test('', async function() {});
});
