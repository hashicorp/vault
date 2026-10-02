/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import ThemeService from 'vault/services/theme';

// Returns a fake MediaQueryList whose `matches` can be updated, and whose
// `change` event can be fired synchronously by calling `triggerChange()`.
function makeFakeMediaQuery(initialMatches = false) {
  const listeners = [];
  const mql = {
    matches: initialMatches,
    addEventListener: sinon.spy((event, fn) => {
      if (event === 'change') listeners.push(fn);
    }),
    removeEventListener: sinon.spy((event, fn) => {
      const idx = listeners.indexOf(fn);
      if (idx !== -1) listeners.splice(idx, 1);
    }),
    triggerChange(newMatches) {
      this.matches = newMatches;
      listeners.forEach((fn) => fn({ matches: newMatches }));
    },
  };
  return mql;
}

module('Unit | Service | theme', function (hooks) {
  setupTest(hooks);

  const STORAGE_KEY = 'vault:prefs:theme';

  hooks.beforeEach(function () {
    window.localStorage.removeItem(STORAGE_KEY);
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
    window.localStorage.removeItem(STORAGE_KEY);
  });

  // --- initial theme resolution ---

  test('defaults to system when no localStorage value', function (assert) {
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(false));
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'system', 'theme is system');
    assert.false(service.isDarkMode, 'isDarkMode reflects system light preference');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
  });

  test('system theme follows dark media query', function (assert) {
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(true));
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'system', 'theme is system');
    assert.true(service.isDarkMode, 'isDarkMode reflects system dark preference');
    assert.dom(document.documentElement).hasAttribute('data-theme', 'dark');
  });

  test('stored "dark" value overrides system preference', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('dark'));
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(false));
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'dark', 'theme is dark');
    assert.true(service.isDarkMode, 'isDarkMode is true');
    assert.dom(document.documentElement).hasAttribute('data-theme', 'dark');
  });

  test('stored "light" value overrides system dark preference', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('light'));
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(true));
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'light', 'theme is light');
    assert.false(service.isDarkMode, 'isDarkMode is false');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
  });

  test('falls back to system when localStorage has invalid value', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, 'invalid-theme');
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(false));
    const service = this.owner.lookup('service:theme');

    assert.strictEqual(service.theme, 'system', 'theme defaults to system');
    assert.false(service.isDarkMode, 'isDarkMode is false');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
  });

  // --- setTheme action ---

  test('setTheme("dark") sets dark theme and persists to localStorage', function (assert) {
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(false));
    const service = this.owner.lookup('service:theme');

    service.setTheme('dark');

    assert.strictEqual(service.theme, 'dark', 'theme is dark');
    assert.true(service.isDarkMode, 'isDarkMode is true');
    assert.dom(document.documentElement).hasAttribute('data-theme', 'dark');
    assert.strictEqual(
      window.localStorage.getItem(STORAGE_KEY),
      JSON.stringify('dark'),
      'persisted to localStorage'
    );
  });

  test('setTheme("light") sets light theme and persists to localStorage', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('dark'));
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(false));
    const service = this.owner.lookup('service:theme');

    service.setTheme('light');

    assert.strictEqual(service.theme, 'light', 'theme is light');
    assert.false(service.isDarkMode, 'isDarkMode is false');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
    assert.strictEqual(
      window.localStorage.getItem(STORAGE_KEY),
      JSON.stringify('light'),
      'persisted to localStorage'
    );
  });

  test('setTheme("system") defers to media query', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('dark'));
    sinon.stub(window, 'matchMedia').returns(makeFakeMediaQuery(false));
    const service = this.owner.lookup('service:theme');

    service.setTheme('system');

    assert.strictEqual(service.theme, 'system', 'theme is system');
    assert.false(service.isDarkMode, 'isDarkMode follows light media query');
    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme');
    assert.strictEqual(
      window.localStorage.getItem(STORAGE_KEY),
      JSON.stringify('system'),
      'persisted to localStorage'
    );
  });

  // --- OS-level media query listener (the bug fix) ---

  test('OS switching to dark mode while theme is "system" immediately applies dark theme', function (assert) {
    const fmq = makeFakeMediaQuery(false); // starts light
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    assert.dom(document.documentElement).doesNotHaveAttribute('data-theme', 'dark is not applied initially');

    // Simulate the OS switching to dark mode.
    fmq.triggerChange(true);

    assert
      .dom(document.documentElement)
      .hasAttribute('data-theme', 'dark', 'data-theme is set to dark immediately when OS switches to dark');
    assert.true(service.isDarkMode, 'isDarkMode reflects new OS dark preference');
  });

  test('OS switching to light mode while theme is "system" immediately removes dark theme', function (assert) {
    const fmq = makeFakeMediaQuery(true); // starts dark
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    assert
      .dom(document.documentElement)
      .hasAttribute('data-theme', 'dark', 'dark theme is applied initially');

    // Simulate the OS switching to light mode.
    fmq.triggerChange(false);

    assert
      .dom(document.documentElement)
      .doesNotHaveAttribute('data-theme', 'data-theme is removed when OS switches to light');
    assert.false(service.isDarkMode, 'isDarkMode reflects new OS light preference');
  });

  test('OS change has no effect when theme is explicitly "dark"', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('dark'));
    const fmq = makeFakeMediaQuery(true);
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    assert.false(fmq.addEventListener.called, 'no listener registered for explicit dark theme');

    // OS flips to light — should be ignored.
    fmq.triggerChange(false);

    assert
      .dom(document.documentElement)
      .hasAttribute(
        'data-theme',
        'dark',
        'dark theme is unchanged when OS switches to light under explicit dark setting'
      );
    assert.true(service.isDarkMode, 'isDarkMode is still true');
  });

  test('OS change has no effect when theme is explicitly "light"', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('light'));
    const fmq = makeFakeMediaQuery(false);
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    assert.false(fmq.addEventListener.called, 'no listener registered for explicit light theme');

    // OS flips to dark — should be ignored.
    fmq.triggerChange(true);

    assert
      .dom(document.documentElement)
      .doesNotHaveAttribute(
        'data-theme',
        'light theme is unchanged when OS switches to dark under explicit light setting'
      );
    assert.false(service.isDarkMode, 'isDarkMode is still false');
  });

  test('switching from "system" to "dark" removes the OS listener', function (assert) {
    const fmq = makeFakeMediaQuery(false);
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    // While on 'system' a listener should be registered.
    assert.true(fmq.addEventListener.calledOnce, 'listener added when theme is system');

    service.setTheme('dark');

    // The listener should be removed after switching away from 'system'.
    assert.true(fmq.removeEventListener.calledOnce, 'listener removed when switching away from system');

    // Trigger an OS change — the listener is gone so the theme should be unaffected.
    fmq.triggerChange(false);

    assert
      .dom(document.documentElement)
      .hasAttribute('data-theme', 'dark', 'dark theme is unaffected by OS change after listener was removed');
  });

  test('switching from "dark" to "system" registers the OS listener', function (assert) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify('dark'));
    const fmq = makeFakeMediaQuery(false); // OS is light
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    // Started as 'dark' — no listener should have been registered.
    assert.false(fmq.addEventListener.called, 'no listener when starting as explicit dark');

    service.setTheme('system');

    // Switching to 'system' must register the listener.
    assert.true(fmq.addEventListener.calledOnce, 'listener added when switching to system');

    // Now an OS change should take effect.
    fmq.triggerChange(true);

    assert
      .dom(document.documentElement)
      .hasAttribute(
        'data-theme',
        'dark',
        'dark theme applied when OS switches to dark after switching to system'
      );
  });

  test('willDestroy removes the OS listener to prevent memory leaks', function (assert) {
    const fmq = makeFakeMediaQuery(false);
    sinon.stub(window, 'matchMedia').returns(fmq);
    const service = this.owner.lookup('service:theme');

    assert.true(fmq.addEventListener.calledOnce, 'listener was registered');

    // Call willDestroy directly — it is a public lifecycle method. This is
    // synchronous and does not involve Ember's destroyable scheduling, so we
    // can assert on the spy counts immediately.
    service.willDestroy();

    assert.true(fmq.removeEventListener.calledOnce, 'listener removed on willDestroy');

    // After willDestroy, OS changes must not touch the DOM.
    fmq.triggerChange(true);

    assert
      .dom(document.documentElement)
      .doesNotHaveAttribute('data-theme', 'DOM is not updated after willDestroy is called');
  });
});
