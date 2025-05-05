/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export type ListDestination = {
  id: string;
  name: string;
  type: DestinationType;
  icon?: string;
  typeDisplayName?: string;
};

export type AssociatedSecret = {
  accessor: string;
  secretName: string;
  syncStatus: string;
  updatedAt: Date;
};

export type AssociatedDestination = {
  type: string;
  name: string;
  syncStatus: string;
  updatedAt: Date;
};

export interface SyncStatus {
  destinationType: string;
  destinationName: string;
  syncStatus: string;
  updatedAt: string;
}

export type DestinationMetrics = {
  icon?: string;
  name?: string;
  type?: DestinationType;
  associationCount: number;
  status: string | null;
  lastUpdated?: Date;
};

export type AssociationMetrics = {
  totalAssociations: number;
  totalSecrets: number;
};

export type DestinationType = 'aws-sm' | 'azure-kv' | 'gcp-sm' | 'gh' | 'vercel-project';

export type DestinationName =
  | 'AWS Secrets Manager'
  | 'Azure Key Vault'
  | 'Google Secret Manager'
  | 'Github Actions'
  | 'Vercel Project';
