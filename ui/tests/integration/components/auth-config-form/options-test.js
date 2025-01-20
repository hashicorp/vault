/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import authConfigForm from 'vault/tests/pages/components/auth-config-form/options';

const component = create(authConfigForm);

module('Integration | Component | auth-config-form options', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.router = this.owner.lookup('service:router');
    this.router.reopen({
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
  });

  test('it submits data correctly', async function (assert) {
    assert.expect(1);
    const model = EmberObject.create({
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
    await render(hbs`<AuthConfigForm::Options @model={{this.model}} />`);
    component.save();
    return settled().then(() => {
      assert.strictEqual(model.config.serialize.callCount, 1, 'config serialize was called once');
    });
  });
});
