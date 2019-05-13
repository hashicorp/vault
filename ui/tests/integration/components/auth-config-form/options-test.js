import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Service from '@ember/service';

import { create } from 'ember-cli-page-object';
import authConfigForm from 'vault/tests/pages/components/auth-config-form/options';

const component = create(authConfigForm);
const routerService = Service.extend({
  transitionTo() {
    return {
      followRedirects() {
        return resolve();
      },
    };
  },
  replaceWith() {
    return resolve();
  },
});

module('Integration | Component | auth-config-form options', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.owner.register('service:router', routerService);
    this.router = this.owner.lookup('service:router');

    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it submits data correctly', async function(assert) {
    let model = EmberObject.create({
      tune() {
        return resolve();
      },
      config: {
        serialize() {
          return {};
        },
      },
    });
    sinon.spy(model.config, 'serialize');
    this.set('model', model);
    await render(hbs`{{auth-config-form/options model=model}}`);
    component.save();
    return settled().then(() => {
      assert.ok(model.config.serialize.calledOnce);
    });
  });
});
