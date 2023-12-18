/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import Route from '@ember/routing/route';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { inject as service } from '@ember/service';
import { Response } from 'miragejs';

module('Unit | Decorators | fetch-secrets-engine-config', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy(console, 'error');
    this.store = this.owner.lookup('service:store');
    this.backend = 'test-path';
    this.owner.lookup('service:secretMountPath').update(this.backend);

    this.createClass = () => {
      @withConfig('ldap/config')
      class Foo extends Route {
        @service store;
        @service secretMountPath;
      }
      // service injection will fail if class is not instantiated with an owner
      return new Foo(this.owner);
    };
  });
  hooks.afterEach(function () {
    this.spy.restore();
  });

  test('it should warn when applying decorator to class that does not extend Route', function (assert) {
    @withConfig()
    class Foo {} // eslint-disable-line
    const message =
      'withConfig decorator must be used on an instance of Ember Route class. Decorator not applied to returned class';
    assert.ok(this.spy.calledWith(message), 'Error is printed to console');
  });

  test('it should return cached record from store if it exists', async function (assert) {
    this.store.pushPayload('ldap/config', {
      modelName: 'ldap/config',
      backend: this.backend,
    });
    const peekSpy = sinon.spy(this.store, 'peekRecord');
    const route = this.createClass();

    await route.beforeModel();
    assert.true(peekSpy.calledWith('ldap/config', this.backend), 'peekRecord called for config model');
    assert.strictEqual(route.configModel.backend, this.backend, 'config model set on class');
    assert.strictEqual(route.configError, null, 'error is unset when model is found');
    assert.false(route.promptConfig, 'promptConfig is false when model is found');
  });

  test('it should fetch record when not in the store', async function (assert) {
    assert.expect(4);

    this.server.get('/test-path/config', () => {
      assert.ok(true, 'fetch request is made');
      return {};
    });

    const route = this.createClass();
    await route.beforeModel();

    assert.strictEqual(route.configModel.backend, this.backend, 'config model set on class');
    assert.strictEqual(route.configError, null, 'error is unset when model is found');
    assert.false(route.promptConfig, 'promptConfig is false when model is found');
  });

  test('it should set prompt value when fetch returns a 404', async function (assert) {
    assert.expect(4);

    this.server.get('/test-path/config', () => {
      assert.ok(true, 'fetch request is made');
      return new Response(404, {}, { errors: [] });
    });

    const route = this.createClass();
    await route.beforeModel();

    assert.strictEqual(route.configModel, null, 'config is not set when error is returned');
    assert.strictEqual(route.configError, null, 'error is unset when 404 is returned');
    assert.true(route.promptConfig, 'promptConfig is true when 404 is returned');
  });

  test('it should set error value when fetch returns error other than 404', async function (assert) {
    assert.expect(4);

    const error = { errors: ['Permission denied'] };
    this.server.get('/test-path/config', () => {
      assert.ok(true, 'fetch request is made');
      return new Response(403, {}, error);
    });

    const route = this.createClass();
    await route.beforeModel();

    assert.strictEqual(route.configModel, null, 'config is not set when error is returned');
    assert.deepEqual(
      route.configError.errors,
      error.errors,
      'error is set when error other than 404 is returned'
    );
    assert.false(route.promptConfig, 'promptConfig is false when error other than 404 is returned');
  });
});
