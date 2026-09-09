/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export function formatDownloadText(
  order: string[],
  steps: string[],
  guidanceMessage: string,
  mode: 'Upgrade' | 'Rollback'
): string {
  const orderLines = order.map((step, index) => `${index + 1}. ${step}`);
  const stepLines = steps.map((step, index) => `${index + 1}. ${step}`);

  return [
    `# Vault ${mode} Steps`,
    '',
    guidanceMessage,
    '',
    `## ${mode} order`,
    ...orderLines,
    '',
    `## Detailed ${mode} Steps`,
    ...stepLines,
  ].join('\n');
}
