/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

import type localStorage from 'vault/lib/local-storage';
import type memoryStorage from 'vault/lib/memory-storage';
import type { ApiResponse, WrapInfo } from 'vault/api';
import type Transition from '@ember/routing/transition';

export interface ControlGroupErrorLog {
  type: string;
  content: string;
  href: string;
  token: string;
  accessor: string;
  creation_path: string;
}

export default class ControlGroupService extends Service {
  tokenToUnwrap: WrapInfo | null;
  storage(): localStorage | memoryStorage;
  keyFromAccessor(accessor: string): string | null;
  storeControlGroupToken(info: WrapInfo): void;
  deleteControlGroupToken(accessor: string): void;
  deleteTokens(): void;
  wrapInfoForAccessor(accessor: string): string | null;
  markTokenForUnwrap(accessor: string): void;
  unmarkTokenForUnwrap(): void;
  tokenForUrl(url: string): { token: string; accessor: string; creationTime: string } | null;
  checkForControlGroup(
    callbackArgs: unknown,
    response: ApiResponse,
    wasWrapTTLRequested: boolean
  ): Promise<unknown>;
  saveTokenFromError(error: WrapInfo): void;
  logFromError(error: WrapInfo): ControlGroupErrorLog;
}
