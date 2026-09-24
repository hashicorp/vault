/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import ThemeService from 'vault/services/theme';

module('Unit | Service | theme', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    window.localStorage.removeItem('vault:theme');
    document.documentElement.removeAttribute('data-theme');
    // Destroy the cached singleton and evict it from the container cache so
    // each test's lookup() constructs a brand-new instance. unregister() alone
    // only clears the registry — the live instance stays in container.cache.
    const container = this.owner.__container__;
    const cached = container.cache['service:theme'];
    if (cached) {
      if (typeof cached.destroy === 'function') cached.destroy();
      delete container.cache['service:theme'];
      delete container.factoryManagerCache['service:theme'];
    }
    this.owner.unregister('service:theme');
    this.owner.register('service:theme', ThemeService);
  });

  hooks.afterEach(function () {
    sinon.restore();
    document.documentElement.removeAttribute('data-theme');
    window.localStorage.removeItem('vault:theme');
  });

  // --- initial theme resolution ---

  test('defaults to system when no localStorage value', function (assert) {
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: false });
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'system', 'theme is system');
    assert.false(service.isDarkMode, 'isDarkMode reflects system light preference');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');

    stub.restore();
  });

  test('system theme follows dark media query', function (assert) {
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: true });
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'system', 'theme is system');
    assert.true(service.isDarkMode, 'isDarkMode reflects system dark preference');
    assert.dom(document.documentElement).hasAttribute('data-theme', 'dark');

    stub.restore();
  });

  test('stored "dark" value overrides system preference', function (assert) {
    window.localStorage.setItem('vault:theme', 'dark');
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: false });
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'dark', 'theme is dark');
    assert.true(service.isDarkMode, 'isDarkMode is true');
    assert.dom(document.documentElement).hasAttribute('data-theme', 'dark');

    stub.restore();
  });

  test('stored "light" value overrides system dark preference', function (assert) {
    window.localStorage.setItem('vault:theme', 'light');
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: true });
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'light', 'theme is light');
    assert.false(service.isDarkMode, 'isDarkMode is false');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');

    stub.restore();
  });

  test('falls back to system when localStorage has invalid value', function (assert) {
    window.localStorage.setItem('vault:theme', 'invalid-theme');
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: false });
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'system', 'theme defaults to system');
    assert.false(service.isDarkMode, 'isDarkMode is false');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');

    stub.restore();
  });

  // --- setTheme action ---

  test('setTheme("dark") sets dark theme and persists to localStorage', function (assert) {
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: false });
    const service = this.owner.lookup('service:theme');

    service.setTheme('dark');

    assert.strictEqual(service.theme, 'dark', 'theme is dark');
    assert.true(service.isDarkMode, 'isDarkMode is true');
    assert.dom(document.documentElement).hasAttribute('data-theme', 'dark');
    assert.strictEqual(window.localStorage.getItem('vault:theme'), 'dark', 'persisted to localStorage');

    stub.restore();
  });

  test('setTheme("light") sets light theme and persists to localStorage', function (assert) {
    window.localStorage.setItem('vault:theme', 'dark');
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: false });
    const service = this.owner.lookup('service:theme');

    service.setTheme('light');

    assert.strictEqual(service.theme, 'light', 'theme is light');
    assert.false(service.isDarkMode, 'isDarkMode is false');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
    assert.strictEqual(window.localStorage.getItem('vault:theme'), 'light', 'persisted to localStorage');

    stub.restore();
  });

  test('setTheme("system") defers to media query', function (assert) {
    window.localStorage.setItem('vault:theme', 'dark');
    const stub = sinon.stub(window, 'matchMedia').returns({ matches: false });
    const service = this.owner.lookup('service:theme');

    service.setTheme('system');

    assert.strictEqual(service.theme, 'system', 'theme is system');
    assert.false(service.isDarkMode, 'isDarkMode follows light media query');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
    assert.strictEqual(window.localStorage.getItem('vault:theme'), 'system', 'persisted to localStorage');

    stub.restore();
  });
});
