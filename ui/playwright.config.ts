/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { defineConfig, devices } from '@playwright/test';
import path from 'path';
import fs from 'fs';
import { TEST_USERS } from './e2e/test-users';
import { isVideoEnabled, videoOptionsFor, videoTimeout } from './e2e/video-config';

import type { UserSetupOptions } from './e2e/init.setup';

const userTypes = Object.values(TEST_USERS);

// start at port 8204 and increment for each project to allow them to run concurrently
const getURL = (increment: number, server = false) => {
  const port = 8204 + increment;
  return server ? `127.0.0.1:${port}` : `http://localhost:${port}/ui/vault/`;
};

// Default cluster listeners use API port + 1, which would collide with the next persona.
const getClusterAddr = (increment: number) => `127.0.0.1:${8304 + increment}`;

// create tmp dir if it doesn't exist for storing session, keys and vault config files
const tmpDir = path.join(__dirname, '/e2e/tmp');
fs.mkdirSync(tmpDir, { recursive: true });

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig<UserSetupOptions>({
  testDir: './e2e',
  // opt out of parallel execution with a test file - by default tests will run in the order they are defined
  fullyParallel: false,
  // fail the build on CI if you accidentally left test.only in the source code.
  forbidOnly: !!process.env.CI,
  // retry on CI only
  retries: process.env.CI ? 2 : 0,
  // use a worker for each project so they run concurrently
  workers: userTypes.length,
  // reporter to use. See https://playwright.dev/docs/test-reporters
  // the html reporter opens a browser at the end, which interrupts a recording session
  reporter: isVideoEnabled ? 'list' : 'html',
  // recorded runs are paced for viewing and need more than the 30s default
  ...(isVideoEnabled ? { timeout: videoTimeout } : {}),
  // shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions.
  use: {
    // collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer
    trace: 'on-first-retry',
  },
  projects: [
    // create setup project for each user type
    ...userTypes.map((userType, index) => ({
      name: `setup:${userType.name}`,
      testMatch: /init\.setup\.ts/,
      use: {
        userType,
        baseURL: getURL(index),
      },
    })),
    // create browser projects for each user type
    ...userTypes.map((userType, index) => {
      const sessionFile = path.join(tmpDir, `${userType.name}-session.json`);
      const name = `chrome:${userType.name}`;
      return {
        name,
        dependencies: [`setup:${userType.name}`],
        workers: 1,
        // only run tests for this user type
        testDir: `./e2e/tests/${userType.name}`,
        use: {
          ...devices['Desktop Chrome'],
          // The setup dependency creates this file before the browser project starts.
          storageState: sessionFile,
          // start at port 8204 and increment for each project to allow them to run concurrently without conflicts
          baseURL: getURL(index),
          permissions: ['clipboard-read', 'clipboard-write'],
          // no-op unless PW_VIDEO is set; never applied to setup projects
          ...videoOptionsFor(name),
        },
      };
    }),
  ],
  webServer: [
    // start a vault server for each project on a different port to allow them to run concurrently
    ...userTypes.map((userType, index) => {
      const isRaft = userType.storage === 'raft';
      // Read the base config for the persona's storage backend.
      const config = JSON.parse(
        fs.readFileSync(
          path.join(__dirname, isRaft ? '/e2e/vault-config-raft.json' : '/e2e/vault-config.json'),
          'utf-8'
        )
      );
      // set the listener address with correct port for this project
      config.listener.tcp.address = getURL(index, true);
      config.listener.tcp.cluster_address = getClusterAddr(index);
      if (isRaft) {
        // raft storage persists to disk (unlike inmem) — clear prior run data for a fresh cluster
        const raftDataDir = path.join(tmpDir, `${userType.name}-raft-data`);
        fs.rmSync(raftDataDir, { recursive: true, force: true });
        fs.mkdirSync(raftDataDir, { recursive: true });
        config.storage.raft.path = raftDataDir;
        config.api_addr = `http://${getURL(index, true)}`;
        config.cluster_addr = `http://${getClusterAddr(index)}`;
      }
      // write the config to a new file for this project
      const configPath = path.join(tmpDir, `${userType.name}-vault-config.json`);
      fs.writeFileSync(configPath, JSON.stringify(config));

      return {
        // Start a non-dev server with the persona's configured storage.
        command: `pnpm run vault:e2e -config=${configPath}`,
        url: getURL(index),
        reuseExistingServer: false,
      };
    }),
  ],
});
