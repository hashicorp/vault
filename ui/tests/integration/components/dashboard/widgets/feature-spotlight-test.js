/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { FEATURE_SPOTLIGHT_CARDS } from 'vault/utils/constants/feature-spotlight';
import {
  DASHBOARD_FEATURE_SPOTLIGHT_BACK,
  DASHBOARD_FEATURE_SPOTLIGHT_LEARN_MORE,
  DASHBOARD_FEATURE_SPOTLIGHT_NEXT,
} from 'vault/utils/analytic-events';

// Fixed card data used by all tests — decoupled from the production constant so
// tests remain stable when FEATURE_SPOTLIGHT_CARDS is edited between releases.
const CARDS = [
  { title: 'Card One', description: 'Description one', link: '/one' },
  { title: 'Card Two', description: 'Description two', link: '/two' },
  { title: 'Card Three', description: 'Description three', link: '/three' },
];

module('Integration | Component | dashboard/widgets/feature-spotlight', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.cards = CARDS;
    this.analytics = this.owner.lookup('service:analytics');
    this.trackEventStub = sinon.stub(this.analytics, 'trackEvent');

    this.renderComponent = () => render(hbs`<Dashboard::Widgets::FeatureSpotlight @cards={{this.cards}} />`);
  });

  hooks.afterEach(function () {
    this.trackEventStub.restore();
  });

  // ─── Rendering ────────────────────────────────────────────────────────────

  test('it renders the card container', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('feature-spotlight')).exists();
  });

  test('it renders the first card title and description', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;
    assert.dom(GENERAL.textDisplay(CARDS[0].title)).hasText(CARDS[0].title);
    assert.dom(GENERAL.textBody('feature-spotlight-description')).hasText(CARDS[0].description);
  });

  test('it renders the Learn more link', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;
    assert.dom(GENERAL.linkTo('feature-spotlight-learn-more')).exists().hasText('Learn more');
  });

  test('it renders the Learn more link pointing to an internal route when route is provided', async function (assert) {
    this.cards = [
      { title: 'Preferences', description: 'Go to preferences', route: 'vault.cluster.preferences' },
    ];
    await this.renderComponent();
    assert.dom(GENERAL.linkTo('feature-spotlight-learn-more')).exists().hasText('Learn more');
  });

  test('it renders custom linkText when provided', async function (assert) {
    this.cards = [
      {
        title: 'Dark mode is here',
        description: 'Go to preferences',
        route: 'vault.cluster.preferences',
        linkText: 'Try it out',
      },
    ];
    await this.renderComponent();
    assert.dom(GENERAL.linkTo('feature-spotlight-learn-more')).exists().hasText('Try it out');
  });

  test('it renders the page label as "1/3" when starting on the first card', async function (assert) {
    // Math.random = 0 → index 0
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('1/3');
  });

  test('it renders Back and Next buttons when there are multiple cards', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.button('back')).exists();
    assert.dom(GENERAL.button('next')).exists();
  });

  test('Back and Next buttons are hidden when there is only one card', async function (assert) {
    this.cards = [CARDS[0]];
    await this.renderComponent();
    assert.dom(GENERAL.button('back')).doesNotExist();
    assert.dom(GENERAL.button('next')).doesNotExist();
    assert.dom('[data-test-feature-spotlight-page-label]').doesNotExist();
  });

  test('renders nothing inside the card when zero cards are provided', async function (assert) {
    this.cards = [];
    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('feature-spotlight')).exists('card container still renders');
    assert.dom(GENERAL.linkTo('feature-spotlight-learn-more')).doesNotExist();
    assert.dom(GENERAL.button('back')).doesNotExist();
    assert.dom(GENERAL.button('next')).doesNotExist();
  });

  test('two-card layout: Back and Next exist and page label shows "1/2"', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    this.cards = [CARDS[0], CARDS[1]];
    await this.renderComponent();
    Math.random = orig;

    assert.dom(GENERAL.button('back')).exists();
    assert.dom(GENERAL.button('next')).exists();
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('1/2');
  });

  test('two-card layout: Next wraps from card 2 back to card 1', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    this.cards = [CARDS[0], CARDS[1]];
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next'));
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('2/2');

    await click(GENERAL.button('next'));
    assert.dom(GENERAL.textDisplay(CARDS[0].title)).hasText(CARDS[0].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('1/2');
  });

  // ─── Pagination ───────────────────────────────────────────────────────────

  test('clicking Next advances to the second card', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next'));

    assert.dom(GENERAL.textDisplay(CARDS[1].title)).hasText(CARDS[1].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('2/3');
  });

  test('clicking Next twice advances to the third card', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next'));
    await click(GENERAL.button('next'));

    assert.dom(GENERAL.textDisplay(CARDS[2].title)).hasText(CARDS[2].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('3/3');
  });

  test('Next wraps from the last card back to the first', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next'));
    await click(GENERAL.button('next'));
    await click(GENERAL.button('next'));

    assert.dom(GENERAL.textDisplay(CARDS[0].title)).hasText(CARDS[0].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('1/3');
  });

  test('clicking Back from the first card wraps to the last', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('back'));

    assert.dom(GENERAL.textDisplay(CARDS[2].title)).hasText(CARDS[2].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('3/3');
  });

  test('clicking Next then Back returns to the first card', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next'));
    await click(GENERAL.button('back'));

    assert.dom(GENERAL.textDisplay(CARDS[0].title)).hasText(CARDS[0].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('1/3');
  });

  // ─── Random start ─────────────────────────────────────────────────────────

  test('random start: Math.random returning 0.99 starts on the last card', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0.99; // floor(0.99 * 3) = 2 → index 2
    await this.renderComponent();
    Math.random = orig;

    assert.dom(GENERAL.textDisplay(CARDS[2].title)).hasText(CARDS[2].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('3/3');
  });

  test('random start: Math.random returning 0.4 starts on the second card', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0.4; // floor(0.4 * 3) = 1 → index 1
    await this.renderComponent();
    Math.random = orig;

    assert.dom(GENERAL.textDisplay(CARDS[1].title)).hasText(CARDS[1].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('2/3');
  });

  // ─── Analytics ────────────────────────────────────────────────────────────

  test('clicking Next fires the NEXT analytic event with the card landed on', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next'));

    assert.ok(
      this.trackEventStub.calledOnceWith(DASHBOARD_FEATURE_SPOTLIGHT_NEXT, {
        channel: 'webpage',
        location: 'dashboard',
        objectType: 'feature-spotlight-widget',
        uiElement: 'next-button',
        type: 'Button',
        action: 'clicked',
        object: CARDS[1].title,
        elementId: 'feature-spotlight-card-1',
      }),
      'trackEvent called with correct next payload'
    );
  });

  test('clicking Back fires the BACK analytic event with the card landed on', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.button('next')); // move to index 1 first
    this.trackEventStub.reset();

    await click(GENERAL.button('back'));

    assert.ok(
      this.trackEventStub.calledOnceWith(DASHBOARD_FEATURE_SPOTLIGHT_BACK, {
        channel: 'webpage',
        location: 'dashboard',
        objectType: 'feature-spotlight-widget',
        uiElement: 'back-button',
        type: 'Button',
        action: 'clicked',
        object: CARDS[0].title,
        elementId: 'feature-spotlight-card-0',
      }),
      'trackEvent called with correct back payload'
    );
  });

  test('clicking Learn more fires the LEARN_MORE analytic event', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    await this.renderComponent();
    Math.random = orig;

    await click(GENERAL.linkTo('feature-spotlight-learn-more'));

    assert.ok(
      this.trackEventStub.calledWith(DASHBOARD_FEATURE_SPOTLIGHT_LEARN_MORE, {
        CTA: 'Learn more',
        channel: 'webpage',
        location: 'dashboard',
        objectType: 'feature-spotlight-widget',
        uiElement: 'learn-more-link',
        type: 'Link',
        action: 'clicked',
        object: CARDS[0].title,
        elementId: 'feature-spotlight-card-0',
      }),
      'trackEvent called with correct learn more payload'
    );
  });

  // ─── isVisible filtering ──────────────────────────────────────────────────

  test('cards with isVisible: false are excluded from the rendered list', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    this.cards = [{ ...CARDS[0], isVisible: false }, { ...CARDS[1] }, { ...CARDS[2], isVisible: false }];
    await this.renderComponent();
    Math.random = orig;

    // only CARDS[1] should be visible — no pagination, shows card 2 content
    assert.dom(GENERAL.textDisplay(CARDS[1].title)).hasText(CARDS[1].title);
    assert.dom(GENERAL.button('back')).doesNotExist('no pagination for a single visible card');
    assert.dom(GENERAL.button('next')).doesNotExist('no pagination for a single visible card');
  });

  test('isVisible: true is treated the same as omitting isVisible', async function (assert) {
    const orig = Math.random;
    Math.random = () => 0;
    this.cards = [
      { ...CARDS[0], isVisible: true },
      { ...CARDS[1], isVisible: true },
    ];
    await this.renderComponent();
    Math.random = orig;

    assert.dom(GENERAL.textDisplay(CARDS[0].title)).hasText(CARDS[0].title);
    assert.dom('[data-test-feature-spotlight-page-label]').hasText('1/2');
  });

  test('all cards hidden via isVisible: false renders an empty card container', async function (assert) {
    this.cards = [
      { ...CARDS[0], isVisible: false },
      { ...CARDS[1], isVisible: false },
    ];
    await this.renderComponent();

    assert.dom(GENERAL.cardContainer('feature-spotlight')).exists('card container still renders');
    assert.dom(GENERAL.linkTo('feature-spotlight-learn-more')).doesNotExist();
    assert.dom(GENERAL.button('back')).doesNotExist();
    assert.dom(GENERAL.button('next')).doesNotExist();
  });

  // ─── Card config ──────────────────────────────────────────────────────────

  test('each production card has required fields', async function (assert) {
    FEATURE_SPOTLIGHT_CARDS.forEach((card, i) => {
      assert.ok(card.title, `card ${i} has a title`);
      assert.ok(card.description, `card ${i} has a description`);
      const hasDestination = Boolean(card.link || card.route);
      assert.true(hasDestination, `card ${i} has a link or route`);
    });
  });
});
