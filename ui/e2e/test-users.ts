/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { USER_POLICY_MAP } from './policies';

/** A test identity selects its environment independently of its ACL policy. */
export interface TestUserType {
  name: string;
  storage: 'inmem' | 'raft';
  policy: keyof typeof USER_POLICY_MAP;
}

// Names also identify projects, test directories, and persisted setup artifacts.
export const TEST_USERS = {
  superuser: { name: 'superuser', storage: 'inmem', policy: 'superuser' },
  raft: { name: 'raft', storage: 'raft', policy: 'superuser' },
} satisfies Record<string, TestUserType>;
