/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface FeatureSpotlightCard {
  title: string;
  description: string;
  /** Docs path (e.g. '/vault/docs/...'). Mutually exclusive with `route`. */
  link?: string;
  /** Ember route name (e.g. 'vault.cluster.preferences'). Mutually exclusive with `link`. */
  route?: string;
  /** Custom call-to-action link text. Defaults to 'Learn more'. */
  linkText?: string;
  imageSrc?: string;
  /** Dark-mode variant of the image. Falls back to imageSrc when not provided. */
  imageSrcDark?: string;
  /** When false, the card is hidden (e.g. enterprise-only cards on CE clusters). Omit to show unconditionally. */
  isVisible?: boolean;
}

/**
 * Fixed-order array of "What's New" cards shown in the Feature Spotlight widget.
 * The starting card is randomized on each dashboard load (see feature-spotlight.ts),
 * but Back/Next always pages through the full list in this order with wrap-around.
 *
 * To update cards for a new release, edit this array.
 */
export const FEATURE_SPOTLIGHT_CARDS: FeatureSpotlightCard[] = [
  {
    title: 'New agent registry in Vault',
    description: 'Centralized view of registered AI agents and their owners in Vault.',
    link: '/vault/docs/concepts/native-ai-agent-support',
    imageSrc: '~/agent-registry-dashboard.png',
    imageSrcDark: '~/agent-registry-dashboard-dark.png',
  },
  {
    title: 'Post-quantum certificate issuance with ML-DSA',
    description:
      "Vault's PKI secrets engine now supports ML-DSA for X.509 certificate issuance, giving you a low-friction path to begin migrating your CA hierarchy to post-quantum cryptography.",
    link: '/vault/docs/secrets/pki/ml-dsa',
  },
  {
    title: 'Eliminate secret zero with TPM auth',
    description:
      'Authenticate Linux VMs to Vault using their vTPM — no bootstrap secret required. TPM-based attestation eliminates Secret Zero for on-premises and private cloud environments.',
    link: '/vault/docs/auth/tpm',
  },
  {
    title: 'Dark mode is here',
    description: 'Easier on your eyes, day or night. Head to Preferences to switch modes anytime.',
    route: 'vault.cluster.preferences',
    linkText: 'Try it out',
    imageSrc: '~/preferences-whats-new.png',
    imageSrcDark: '~/preferences-whats-new-dark.png',
  },
];
