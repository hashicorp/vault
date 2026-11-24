/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { CsvColumn } from 'vault/utils/generate-csv';
import type { RegistryCsvRow } from 'vault/agent-registry';

export const AGENT_REGISTRY_PATH = 'agent-registry';
export const AGENTS_REGISTRY_CSV_FILENAME = 'agent_registry';

export const CSV_COLUMNS: CsvColumn<RegistryCsvRow>[] = [
  { header: 'Agent name', key: 'agentName' },
  { header: 'Agentic entity in Vault', key: 'agenticEntityInVault' },
  { header: 'Entity / Alias ID', key: 'entityOrAliasId' },
  { header: 'Entity status', key: 'entityStatus' },
  { header: 'Entity created at', key: 'entityCreatedAt' },
  { header: 'Entity updated at', key: 'entityUpdatedAt' },
];
