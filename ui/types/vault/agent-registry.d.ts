/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { RegistrationUpdateByIdRequest } from '@hashicorp/vault-client-typescript';
import type { ListEntity } from 'vault/identity';

export type Agent = {
  id: string;
  display_name: string;
  entity_id: string;
  description: string;
  owner: string;
  ceiling_policy: string[];
  no_default_ceiling_policy: boolean;
  creation_time: string;
  last_updated_time: string;
};

export type ListAgent = {
  id: string;
  display_name: string;
  entity_id: string;
  creation_time?: string;
  last_updated_time?: string;
  no_default_ceiling_policy: boolean;
  entity?: ListEntity;
};

export type RegistryCsvRow = {
  agentName: string;
  agenticEntityInVault: string;
  entityOrAliasId: string;
  entityStatus: string;
  entityCreatedAt: string;
  entityUpdatedAt: string;
};
