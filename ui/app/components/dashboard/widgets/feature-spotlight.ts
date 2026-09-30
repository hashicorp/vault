/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import {
  DASHBOARD_FEATURE_SPOTLIGHT_BACK,
  DASHBOARD_FEATURE_SPOTLIGHT_LEARN_MORE,
  DASHBOARD_FEATURE_SPOTLIGHT_NEXT,
} from 'vault/utils/analytic-events';
import { FEATURE_SPOTLIGHT_CARDS } from 'vault/utils/constants/feature-spotlight';

import type AnalyticsService from 'vault/services/analytics';
import type ThemeService from 'vault/services/theme';
import type { FeatureSpotlightCard } from 'vault/utils/constants/feature-spotlight';

interface Args {
  cards?: FeatureSpotlightCard[];
}

/**
 * @module Dashboard::Widgets::FeatureSpotlight
 * Paginated dashboard widget that cycles through feature highlight cards.
 * The starting card is chosen randomly on each page load; Back/Next wrap around.
 *
 * Pass @cards to override the default FEATURE_SPOTLIGHT_CARDS list (useful in tests).
 *
 * @example
 * ```hbs
 * <Dashboard::Widgets::FeatureSpotlight />
 * ```
 */
export default class DashboardWidgetsFeatureSpotlight extends Component<Args> {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly theme: ThemeService;

  @tracked currentIndex = 0;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.currentIndex = Math.floor(Math.random() * this.cards.length);
  }

  get cards(): FeatureSpotlightCard[] {
    return (this.args.cards ?? FEATURE_SPOTLIGHT_CARDS).filter((c) => c.isVisible !== false);
  }

  get currentFeature(): FeatureSpotlightCard {
    return this.cards[this.currentIndex] as FeatureSpotlightCard;
  }

  get currentImageSrc(): string | undefined {
    const { imageSrc, imageSrcDark } = this.currentFeature;
    if (this.theme.isDarkMode && imageSrcDark) {
      return imageSrcDark;
    }
    return imageSrc;
  }

  get linkText(): string {
    return this.currentFeature.linkText ?? 'Learn more';
  }

  get total(): number {
    return this.cards.length;
  }

  get pageLabel(): string {
    return `${this.currentIndex + 1}/${this.total}`;
  }

  get isFirst(): boolean {
    return this.currentIndex === 0;
  }

  get isLast(): boolean {
    return this.currentIndex === this.total - 1;
  }

  get hasPagination(): boolean {
    return this.total > 1;
  }

  @action
  next(): void {
    this.currentIndex = (this.currentIndex + 1) % this.total;
    this.analytics.trackEvent(DASHBOARD_FEATURE_SPOTLIGHT_NEXT, {
      channel: 'webpage',
      location: 'dashboard',
      objectType: 'feature-spotlight-widget',
      uiElement: 'next-button',
      type: 'Button',
      action: 'clicked',
      cardTitle: this.currentFeature.title,
      cardIndex: this.currentIndex,
    });
  }

  @action
  back(): void {
    this.currentIndex = (this.currentIndex - 1 + this.total) % this.total;
    this.analytics.trackEvent(DASHBOARD_FEATURE_SPOTLIGHT_BACK, {
      channel: 'webpage',
      location: 'dashboard',
      objectType: 'feature-spotlight-widget',
      uiElement: 'back-button',
      type: 'Button',
      action: 'clicked',
      cardTitle: this.currentFeature.title,
      cardIndex: this.currentIndex,
    });
  }

  @action
  trackLearnMore(): void {
    this.analytics.trackEvent(DASHBOARD_FEATURE_SPOTLIGHT_LEARN_MORE, {
      CTA: this.linkText,
      channel: 'webpage',
      location: 'dashboard',
      objectType: 'feature-spotlight-widget',
      uiElement: 'learn-more-link',
      type: 'Link',
      action: 'clicked',
      cardTitle: this.currentFeature.title,
      cardIndex: this.currentIndex,
    });
  }
}
