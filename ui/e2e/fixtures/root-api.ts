/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import fs from 'fs';
import path from 'path';
import { expect } from '@playwright/test';

import type { APIRequestContext, PlaywrightWorkerArgs, TestInfo } from '@playwright/test';

/**
 * Root-token API context for seeding test data, including data the persona under test cannot
 * create itself. The root token comes from the keys file that init.setup.ts downloads for the
 * persona, so it never appears in the UI or a recording.
 */
export const rootApi = async (
  playwright: PlaywrightWorkerArgs['playwright'],
  testInfo: TestInfo
): Promise<APIRequestContext> => {
  const persona = testInfo.project.name.replace(/^chrome:/, '');
  const keysPath = path.join(__dirname, `../tmp/${persona}-keys.json`);
  const { root_token } = JSON.parse(fs.readFileSync(keysPath, 'utf-8'));
  const baseURL = new URL(testInfo.project.use.baseURL as string).origin;
  return playwright.request.newContext({ baseURL, extraHTTPHeaders: { 'X-Vault-Token': root_token } });
};

/**
 * Enables a secrets engine at `path` unless one is already mounted there. Playwright re-runs
 * beforeAll in a fresh worker after a failure, and Vault rejects re-creating a mount with a 400.
 *
 * @param api - root-token context from {@link rootApi}
 * @param path - mount path without a trailing slash, e.g. "readonly-kv"
 * @param type - secrets engine type, e.g. "kv-v2"
 */
export const ensureMount = async (api: APIRequestContext, path: string, type: string) => {
  if ((await api.get(`/v1/sys/mounts/${path}`)).ok()) return;
  await expect(await api.post(`/v1/sys/mounts/${path}`, { data: { type } })).toBeOK();
};

/**
 * Creates the namespace at `path` unless it already exists, for the same retry reason as
 * {@link ensureMount}.
 *
 * @param api - root-token context from {@link rootApi}
 * @param path - namespace path relative to the root namespace, e.g. "readonly-ns"
 */
export const ensureNamespace = async (api: APIRequestContext, path: string) => {
  if ((await api.get(`/v1/sys/namespaces/${path}`)).ok()) return;
  await expect(await api.post(`/v1/sys/namespaces/${path}`)).toBeOK();
};
