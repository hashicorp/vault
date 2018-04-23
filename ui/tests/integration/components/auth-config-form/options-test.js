import { moduleForComponent, test } from 'ember-qunit';
import Ember from 'ember';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

import { create } from 'ember-cli-page-object';
import authConfigForm from 'vault/tests/pages/components/auth-config-form/options';

const component = create(authConfigForm);

moduleForComponent('auth-config-form/options', 'Integration | Component | auth-config-form options', {
  integration: true,
  beforeEach() {
    Ember.getOwner(this).lookup('service:flash-messages').registerTypes(['success']);
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

test('it submits data correctly', function(assert) {
  let model = Ember.Object.create({
    tune() {
      return Ember.RSVP.resolve();
    },
    config: {
      serialize() {
        return {};
      },
    },
  });
  sinon.spy(model.config, 'serialize');
  this.set('model', model);
  this.render(hbs`{{auth-config-form/options model=model}}`);
  component.save();
  wait().then(() => {
    assert.ok(model.config.serialize.calledOnce);
  });
});
