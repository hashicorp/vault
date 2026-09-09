/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export type Entity = {
  aliases: Alias[];
  creation_time: string;
  direct_group_ids: string[];
  disabled: boolean;
  group_ids: string[];
  id: string;
  inherited_group_ids: string[];
  last_update_time: string;
  merged_entity_ids: string[] | null;
  metadata: Record<string, string> | null;
  mfa_secrets: Record<string, string>;
  name: string;
  namespace_id: string;
  policies: string[];
};

export type Alias = {
  canonical_id: string;
  creation_time: string;
  custom_metadata: Record<string, string> | null;
  id: string;
  last_update_time: string;
  local: boolean;
  merged_from_canonical_ids: string[] | null;
  metadata: Record<string, string> | null;
  mount_accessor: string;
  mount_path: string;
  mount_type: string;
  name: string;
  // enriched from the oauth resource server profile (present only when mount_accessor matches the oauth-resource-server pattern)
  issuer_id?: string;
  external_id?: string;
  profile_name?: string;
  profile_id?: string;
  namespace?: string;
};

export type ListEntity = {
  id: string;
  name: string;
  disabled: boolean;
  creation_time: string;
  last_update_time: string;
  group?: string;
  description?: string;
  policies?: string[];
  merged_entity_ids?: string[] | null;
  metadata?: Record<string, string> | null;
  aliases: {
    id: string;
    mount_accessor: string;
    mount_path: string;
    mount_type: string;
    name: string;
    creation_time: string;
    last_update_time: string;
  }[];
};

export type ListAlias = Partial<
  Omit<Alias, 'creation_time' | 'last_update_time' | 'merged_from_canonical_ids'>
>;

export type Group = {
  alias: Alias;
  creation_time: string;
  id: string;
  last_update_time: string;
  member_entity_ids: string[] | null;
  member_group_ids: string[] | null;
  metadata: Record<string, string> | null;
  modify_index: number;
  name: string;
  namespace_id: string;
  parent_group_ids: string[] | null;
  policies: string[] | null;
  type: 'internal' | 'external';
};
