/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { capitalize } from '@ember/string';

/**
 * Converts a snake_case, kebab-case, or plain lowercase string to sentence case.
 * Following HDS guidelines, only the first character is capitalised; all other
 * characters are lower-cased.
 *
 * Examples:
 *   "some_example_text" → "Some example text"
 *   "client-id"         → "Client id"
 *   "kubernetes"        → "Kubernetes"
 */
const ACRONYMS = new Set([
  'aws',
  'gcp',
  'jwt',
  'kmip',
  'kv',
  'ldap',
  'mfa',
  'oidc',
  'pki',
  'saml',
  'ssh',
  'tls',
  'totp',
]);

const WORD_OVERRIDES: Record<string, string> = {
  alicloud: 'AliCloud',
  approle: 'AppRole',
  github: 'GitHub',
  rabbitmq: 'RabbitMQ',
};

function formatWord(word: string, capitalizeWord = false): string {
  const override = WORD_OVERRIDES[word];
  if (override) return override;

  if (ACRONYMS.has(word)) return word.toUpperCase();

  return capitalizeWord ? capitalize(word) : word;
}

interface ToSentenceCaseOptions {
  acronymsOnly?: boolean;
}

export function toSentenceCase(str: string, options: ToSentenceCaseOptions = {}): string {
  if (!str) return '';

  const words = str.replace(/[_-]/g, ' ').toLowerCase().split(/\s+/).filter(Boolean);

  if (!words.length) return '';

  if (options.acronymsOnly) {
    return words
      .map((word) => {
        const override = WORD_OVERRIDES[word];
        if (override) return override;
        return ACRONYMS.has(word) ? word.toUpperCase() : word;
      })
      .join(' ');
  }

  const firstWord = words[0] ?? '';
  const remainingWords = words.slice(1);
  const formattedFirstWord = formatWord(firstWord, true);

  return [formattedFirstWord, ...remainingWords.map((word) => formatWord(word))].join(' ');
}
