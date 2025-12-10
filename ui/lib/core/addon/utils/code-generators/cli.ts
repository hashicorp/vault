/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface FormatCliArgs {
  command: string; // The CLI command (e.g., "policy write my-policy")
  content: string; // The content/body to pass to the command
}

export const formatCli = ({ command, content }: FormatCliArgs) => {
  return `vault ${command} ${content}`.trim();
};

export const writePolicy = (name: string) => `policy write ${name}`;
