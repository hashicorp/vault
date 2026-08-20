/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

export type OnboardingCardSize = 'small' | 'large';

interface OnboardingCardSignature {
  Args: {
    /** Card size variant: 'small' or 'large'. Defaults to 'large'. */
    size?: OnboardingCardSize;
  };
  Blocks: {
    /** Main content area */
    default: [];
  };
  Element: HTMLDivElement;
}

/**
 * @module OnboardingCard
 * Reusable onboarding card component with gradient background and decorative elements.
 * Displays content with optional size variants.
 *
 * @example
 * <OnboardingCard @size="large">
 *   {{yield}}
 * </OnboardingCard>
 */
export default class OnboardingCard extends Component<OnboardingCardSignature> {
  get sizeClass(): string {
    return `onboarding-card--${this.args.size || 'large'}`;
  }
}
