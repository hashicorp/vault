/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, select, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setPreference, getStringPreference, setStringPreference } from 'vault/utils/preferences';

const SELECTOR = {
  section: '[data-test-persona-section]',
  select: '[data-test-persona-select]',
  otherInput: '[data-test-persona-other-input]',
  otherInputError: '[data-test-persona-other-input-error]',
};

module('Integration | Component | user-preferences/persona', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.analytics = this.owner.lookup('service:analytics');
    this.trackEvent = sinon.stub(this.analytics, 'trackEvent');
    window.localStorage.clear();
  });

  hooks.afterEach(function () {
    sinon.restore();
    window.localStorage.clear();
  });

  test('it renders the Your role section with a select field', async function (assert) {
    // Verifies the section heading and select control render on mount.
    await render(hbs`<UserPreferences::Persona />`);

    assert.dom(SELECTOR.section).exists('the Persona section renders');
    assert.dom(SELECTOR.select).exists('the select control renders');
  });

  test('the select defaults to the placeholder option when no persona is stored', async function (assert) {
    // Verifies that an absent localStorage key leaves the placeholder selected.
    await render(hbs`<UserPreferences::Persona />`);

    assert.dom(`${SELECTOR.select} option:checked`).hasValue('', 'placeholder option is selected by default');
  });

  test('the select reflects a previously stored persona on mount', async function (assert) {
    // Verifies that a value already in localStorage is pre-selected when the
    // component renders, so users see their current choice rather than the placeholder.
    // setStringPreference must be called before render so the @tracked property
    // initialises from the already-stored value.
    setStringPreference('persona', JSON.stringify({ value: 'developer' }));

    await render(hbs`<UserPreferences::Persona />`);

    assert
      .dom(`${SELECTOR.select} option:checked`)
      .hasValue('developer', 'the previously stored persona is pre-selected');
  });

  test('selecting a persona persists it to localStorage as structured JSON', async function (assert) {
    // Verifies the storage write path when the user makes a selection.
    await render(hbs`<UserPreferences::Persona />`);

    await select(SELECTOR.select, 'security-analyst');

    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'security-analyst' },
      'persona is persisted to localStorage via the registry as structured JSON'
    );
  });

  test('all four persona options are available', async function (assert) {
    // Verifies the complete option list matches the spec.
    await render(hbs`<UserPreferences::Persona />`);

    const options = [...document.querySelectorAll(`${SELECTOR.select} option`)].map((o) => o.value);

    assert.deepEqual(
      options,
      ['', 'developer', 'platform-engineer', 'security-analyst', 'other'],
      'placeholder + four persona options are present'
    );
  });

  test('selecting a persona fires trackEvent when telemetry consent is on', async function (assert) {
    // Verifies the Segment event is emitted only when the user has opted in.
    setPreference('telemetryConsent', true);

    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'developer');

    assert.true(
      this.trackEvent.calledOnceWith('UI Interaction', {
        namespace: 'user-preferences',
        action: 'persona_set',
        elementId: 'persona-select',
        channel: 'webpage',
        location: 'user-preferences',
        objectType: 'persona',
        object: 'developer',
        resultValue: 'developer',
      }),
      'trackEvent is called with UI Interaction when consent is granted'
    );
  });

  test('selecting a persona does NOT fire trackEvent when telemetry consent is off', async function (assert) {
    // Verifies the persona is stored locally but Segment is not called when the
    // user has not opted into telemetry.
    setPreference('telemetryConsent', false);

    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'platform-engineer');

    assert.true(this.trackEvent.notCalled, 'trackEvent is not called when telemetry consent is off');
    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'platform-engineer' },
      'persona is still stored locally even without telemetry consent'
    );
  });

  test('selecting a persona does NOT fire trackEvent when consent is absent (never decided)', async function (assert) {
    // Verifies that an unrecorded consent (empty localStorage) is treated as off.
    // The user has never seen the banner so they have not opted in.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'platform-engineer');

    assert.true(this.trackEvent.notCalled, 'trackEvent is not called when consent has never been recorded');
  });

  // ---------------------------------------------------------------------------
  // "Other" free-text field
  // ---------------------------------------------------------------------------

  test('selecting Other shows the free-text input, persists "other", and hides the field for other options', async function (assert) {
    // Verifies the text field only appears when "Other" is chosen, and that
    // selecting Other immediately persists { value: "other" } to localStorage before the
    // user has typed anything into the free-text field.
    await render(hbs`<UserPreferences::Persona />`);

    assert.dom(SELECTOR.otherInput).doesNotExist('text input is hidden before Other is selected');

    await select(SELECTOR.select, 'other');
    assert.dom(SELECTOR.otherInput).exists('text input appears when Other is selected');
    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'other' },
      '"other" is persisted immediately on selection'
    );

    await select(SELECTOR.select, 'developer');
    assert.dom(SELECTOR.otherInput).doesNotExist('text input is hidden again when a fixed option is chosen');
  });

  test('typing a valid custom role in the Other field stores it in localStorage', async function (assert) {
    // Verifies valid free text is preserved with original casing in localStorage.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');

    await fillIn(SELECTOR.otherInput, 'DevOps SRE');

    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'other', customRole: 'DevOps SRE' },
      'valid custom role is stored as { value: "other", customRole: "DevOps SRE" }'
    );
  });

  test('typing an invalid custom role in the Other field does not update localStorage', async function (assert) {
    // Verifies that invalid input (containing characters outside letters, digits, spaces)
    // is not persisted — localStorage retains only { value: "other" } until a valid
    // string is entered.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');

    await fillIn(SELECTOR.otherInput, 'my ~DREAM~ role');

    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'other' },
      'invalid input is not written to localStorage'
    );
  });

  test('correcting an invalid custom role persists the new valid value', async function (assert) {
    // Verifies that after typing an invalid value and then correcting it,
    // the valid value is stored and the error is gone.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');

    await fillIn(SELECTOR.otherInput, 'bad ~value~');
    await fillIn(SELECTOR.otherInput, 'Platform Engineer');

    assert.dom(SELECTOR.otherInputError).doesNotExist('error is cleared after correction');
    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'other', customRole: 'Platform Engineer' },
      'corrected valid value is persisted to localStorage'
    );
  });

  test('clearing the Other text field stores only { value: "other" } without a customRole key', async function (assert) {
    // Verifies that an empty custom role is not serialised into localStorage so the
    // stored object stays clean rather than carrying a redundant empty-string key.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');
    await fillIn(SELECTOR.otherInput, 'designer');
    await fillIn(SELECTOR.otherInput, '');

    assert.deepEqual(
      JSON.parse(getStringPreference('persona')),
      { value: 'other' },
      'clearing the field stores only { value: "other" } with no customRole key'
    );
  });

  test('Other text input is pre-filled verbatim when a stored object is loaded on mount', async function (assert) {
    // Verifies that returning to the page displays the exact typed string without modification.
    setStringPreference('persona', JSON.stringify({ value: 'other', customRole: 'DevOps / SRE (Cloud)' }));

    await render(hbs`<UserPreferences::Persona />`);

    assert.dom(SELECTOR.otherInput).exists('text input is shown for a stored other object');
    assert.dom(SELECTOR.otherInput).hasValue('DevOps / SRE (Cloud)', 'exact customRole is displayed');
    assert
      .dom(`${SELECTOR.select} option:checked`)
      .hasValue('other', 'Other option is selected in the dropdown');
  });

  test('typing invalid characters in the Other field shows an inline error', async function (assert) {
    // Verifies that characters outside the allowed set (letters, digits, spaces) trigger
    // an inline error message on the field without blocking the input.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');

    await fillIn(SELECTOR.otherInput, 'my ~**DREAM**~ role');

    assert.dom(SELECTOR.otherInputError).exists('error message appears for invalid characters');
    assert
      .dom(SELECTOR.otherInputError)
      .hasText(
        'Role description can only contain letters, numbers, and spaces.',
        'error message describes the constraint'
      );
    assert
      .dom(SELECTOR.otherInput)
      .hasClass('hds-form-text-input--is-invalid', 'field has the HDS invalid CSS class');
  });

  test('correcting invalid input in the Other field clears the error', async function (assert) {
    // Verifies that once the value is brought back to a valid string the error
    // disappears, giving users clear, live feedback.
    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');

    await fillIn(SELECTOR.otherInput, 'bad ~input~');
    assert.dom(SELECTOR.otherInputError).exists('error is shown for invalid input');

    await fillIn(SELECTOR.otherInput, 'Platform Engineer');
    assert.dom(SELECTOR.otherInputError).doesNotExist('error clears once input is valid');
  });

  test('Other trackEvent fires with the full other-<role> value when consent is on', async function (assert) {
    // Verifies the Segment event carries the full stored value, not the raw input.
    setPreference('telemetryConsent', true);

    await render(hbs`<UserPreferences::Persona />`);
    await select(SELECTOR.select, 'other');
    await fillIn(SELECTOR.otherInput, 'admin');

    assert.true(
      this.trackEvent.calledWith('UI Interaction', {
        namespace: 'user-preferences',
        action: 'persona_set',
        elementId: 'persona-select',
        channel: 'webpage',
        location: 'user-preferences',
        objectType: 'persona',
        object: 'other-admin',
        resultValue: 'other-admin',
      }),
      'trackEvent is called with the full other-admin value'
    );
  });
});
