/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Wraps content in a heredoc (EOT - End of Text) block.
 * @param content - The content to wrap in the heredoc block
 * @returns A formatted heredoc string with EOT delimiters
 */
export const formatEot = (content = '') => {
  return `<<EOT
${content}
EOT`;
};
