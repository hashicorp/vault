/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
const SELECTORS = {
  card: '[data-test-onboarding-card]',
  content: '.onboarding-card__content',
  vaultIcon: '.onboarding-card__vault-icon',
  aurora: '.onboarding-card__aurora',
};

module('Integration | Component | onboarding-card', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders with default props', async function (assert) {
    await render(hbs`<OnboardingCard>Test content</OnboardingCard>`);
    assert.dom(SELECTORS.card).exists('Card renders');
    assert.dom(SELECTORS.card).hasText('Test content', 'Content is displayed');
    assert.dom(SELECTORS.vaultIcon).exists('Vault icon decoration renders');
    assert.dom(SELECTORS.aurora).exists('Aurora decoration renders');
  });

  test('it applies large size class by default', async function (assert) {
    await render(hbs`<OnboardingCard>Content</OnboardingCard>`);
    assert.dom(SELECTORS.card).hasClass('onboarding-card--large', 'Default size is large');
  });

  test('it applies small size class when specified', async function (assert) {
    await render(hbs`<OnboardingCard @size="small">Content</OnboardingCard>`);
    assert.dom(SELECTORS.card).hasClass('onboarding-card--small', 'Small size class applied');
  });

  test('it renders yielded content correctly', async function (assert) {
    await render(hbs`
      <OnboardingCard>
        <h2>Test Heading</h2>
        <p>Test paragraph</p>
      </OnboardingCard>
    `);
    assert.dom(SELECTORS.content).includesText('Test Heading');
    assert.dom(SELECTORS.content).includesText('Test paragraph');
  });

  test('vault icon decorations have proper accessibility attributes', async function (assert) {
    await render(hbs`<OnboardingCard>Content</OnboardingCard>`);
    assert
      .dom(SELECTORS.vaultIcon)
      .hasAttribute('aria-hidden', 'true', 'Vault icon is hidden from screen readers');
    assert.dom(SELECTORS.aurora).hasAttribute('aria-hidden', 'true', 'Aurora is hidden from screen readers');
    assert.dom(SELECTORS.vaultIcon).hasAttribute('alt', '', 'Vault icon has empty alt for decorative image');
    assert.dom(SELECTORS.aurora).hasAttribute('alt', '', 'Aurora has empty alt for decorative image');
  });

  test('it forwards attributes to the card element', async function (assert) {
    await render(hbs`<OnboardingCard class="custom-class">Content</OnboardingCard>`);
    assert.dom(SELECTORS.card).hasClass('custom-class', 'Custom attributes are forwarded');
  });
});
