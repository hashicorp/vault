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

// returns a formatted args object to populate snippets with empty values removed (e.g. empty strings, empty objects/arrays, undefined, null)
export const formatArgsFromPayload = (payload: Record<string, unknown> = {}) => {
  return Object.fromEntries(
    Object.entries(payload).filter(([, value]) => {
      const isEmptyValue = value === '' || value === undefined || value === null;
      const isEmptyObject =
        typeof value === 'object' &&
        value !== null &&
        !Array.isArray(value) &&
        Object.keys(value).length === 0;
      const isEmptyArray = Array.isArray(value) && value.length === 0;
      return !isEmptyValue && !isEmptyObject && !isEmptyArray;
    })
  );
};
