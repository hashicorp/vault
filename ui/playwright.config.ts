/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { defineConfig, devices } from '@playwright/test';
import path from 'path';
import fs from 'fs';
import { USER_POLICY_MAP } from './e2e/policies';

import type { UserSetupOptions } from './e2e/init.setup';

const userTypes = Object.keys(USER_POLICY_MAP);

// start at port 8204 and increment for each project to allow them to run concurrently
const getURL = (increment: number, server = false) => {
  const port = `820${4 + increment}`;
  return server ? `127.0.0.1:${port}` : `http://localhost:${port}/ui/vault/`;
};

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
  reporter: 'html',
  // shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions.
  use: {
    // collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer
    trace: 'on-first-retry',
  },
  projects: [
    // create setup project for each user type
    ...userTypes.map((userType, index) => ({
      name: `setup:${userType}`,
      testMatch: /init\.setup\.ts/,
      use: {
        userType,
        baseURL: getURL(index),
      },
    })),
    // create browser projects for each user type
    ...userTypes.map((userType, index) => {
      const sessionFile = path.join(tmpDir, `${userType}-session.json`);
      return {
        name: `chrome:${userType}`,
        dependencies: [`setup:${userType}`],
        workers: 1,
        // only run tests for this user type
        testDir: `./e2e/tests/${userType}`,
        use: {
          ...devices['Desktop Chrome'],
          // only use if file has already been created by the setup project
          storageState: fs.existsSync(sessionFile) ? sessionFile : undefined,
          // start at port 8204 and increment for each project to allow them to run concurrently without conflicts
          baseURL: getURL(index),
        },
      };
    }),
  ],
  webServer: [
    // start a vault server for each project on a different port to allow them to run concurrently
    ...userTypes.map((userType, index) => {
      // read base config file
      const config = JSON.parse(fs.readFileSync(path.join(__dirname, '/e2e/vault-config.json'), 'utf-8'));
      // set the listener address with correct port for this project
      config.listener.tcp.address = getURL(index, true);
      // write the config to a new file for this project
      const configPath = path.join(tmpDir, `${userType}-vault-config.json`);
      fs.writeFileSync(configPath, JSON.stringify(config));

      return {
        // start vault server (not dev) with inmem storage
        command: `pnpm run vault:e2e -config=${configPath}`,
        url: getURL(index),
        reuseExistingServer: false,
      };
    }),
  ],
});
