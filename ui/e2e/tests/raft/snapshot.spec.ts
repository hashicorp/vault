/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test, holdForVideo, highlightElement } from '../../fixtures/demo';
import fs from 'fs';

import type { Page } from '@playwright/test';

// The raft persona boots real Raft storage; snapshots are not available for inmem storage.
// Tests run in order and share one snapshot file taken before the marker mount is created.
test.describe.configure({ mode: 'serial' });

// Blocking service workers forces the in-memory fetch fallback for downloads instead of the service worker path.
test.use({ serviceWorkers: 'block' });

const MARKER_MOUNT = 'snapshot-marker';
let snapshotPath = '';

const runConsoleCommand = async (page: Page, command: string, success: string) => {
  await page.goto('dashboard');
  await page.getByRole('button', { name: 'Console toggle' }).click();
  const repl = page.getByRole('textbox', { name: 'Web R.E.P.L.' });
  await repl.fill(command);
  await repl.press('Enter');
  await expect(page.locator('.console-ui-panel').getByText(success)).toBeVisible();
  await highlightElement(page, page.locator('.console-ui-panel').getByText(success));
};

const enginePaths = async (page: Page) => {
  await page.goto('secrets-engines');
  await holdForVideo(page);
  return page.getByRole('gridcell', { name: `${MARKER_MOUNT}/`, exact: true });
};

const goToRestore = async (page: Page) => {
  await page.goto('storage/raft');
  await page.getByRole('button', { name: 'Snapshots' }).click();
  await page.getByRole('link', { name: 'Restore' }).click();
  await expect(page.getByRole('heading', { name: 'Restore Snapshot' })).toBeVisible();
};

test('downloads a raft snapshot', async ({ page }, testInfo) => {
  await page.goto('storage/raft');
  await page.getByRole('button', { name: 'Snapshots' }).click();

  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Download' }).click();
  await expect(page.getByText('The snapshot download has begun.')).toBeVisible();
  await highlightElement(page, page.getByText('The snapshot download has begun.'));
  const download = await downloadPromise;

  expect(download.suggestedFilename()).toBe('snapshot.gz');
  snapshotPath = testInfo.outputPath('snapshot.gz');
  await download.saveAs(snapshotPath);
  const bytes = fs.readFileSync(snapshotPath);
  expect(bytes.length).toBeGreaterThan(0);
  // gzip magic number
  expect([bytes[0], bytes[1]]).toEqual([0x1f, 0x8b]);
});

test('restores a snapshot and discards data written after it was taken', async ({ page }) => {
  await test.step('create a marker mount after the snapshot was taken', async () => {
    await runConsoleCommand(page, `write sys/mounts/${MARKER_MOUNT} type=kv`, 'Success! Data written to');
    await expect(await enginePaths(page)).toBeVisible();
  });

  await test.step('restore the snapshot through the UI', async () => {
    await goToRestore(page);
    await page.locator('input[type="file"]').setInputFiles(snapshotPath);
    const responsePromise = page.waitForResponse(
      (response) =>
        response.url().endsWith('/v1/sys/storage/raft/snapshot') && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Restore' }).click();
    const response = await responsePromise;
    expect(response.ok()).toBe(true);
    await expect(page.getByText('The snapshot was successfully uploaded!')).toBeVisible();
    await highlightElement(page, page.getByText('The snapshot was successfully uploaded!'));
  });

  await test.step('the marker mount no longer exists in the backend', async () => {
    await expect(await enginePaths(page)).toHaveCount(0);
    await holdForVideo(page);
  });
});

test('force restores a snapshot through the force endpoint', async ({ page }) => {
  await goToRestore(page);
  await page.locator('input[type="file"]').setInputFiles(snapshotPath);
  await page.getByText('Force restore', { exact: true }).click();
  await expect(page.getByRole('checkbox', { name: 'Force restore' })).toBeChecked();
  await highlightElement(page, page.getByRole('checkbox', { name: 'Force restore' }));
  const responsePromise = page.waitForResponse(
    (response) =>
      response.url().endsWith('/v1/sys/storage/raft/snapshot-force') && response.request().method() === 'POST'
  );
  await page.getByRole('button', { name: 'Restore' }).click();
  const response = await responsePromise;
  expect(response.ok()).toBe(true);
  await expect(page.getByText('The snapshot was successfully uploaded!')).toBeVisible();
  await highlightElement(page, page.getByText('The snapshot was successfully uploaded!'));
});

test('cancelling an upload aborts the restore without applying it', async ({ page }) => {
  await test.step('create a marker mount that must survive the aborted restore', async () => {
    await runConsoleCommand(page, `write sys/mounts/${MARKER_MOUNT} type=kv`, 'Success! Data written to');
    await expect(await enginePaths(page)).toBeVisible();
  });

  await test.step('hold the upload open, then cancel it', async () => {
    // Hold the restore request so the upload is still in flight when it is cancelled.
    await page.route('**/v1/sys/storage/raft/snapshot', () => new Promise(() => undefined));
    await goToRestore(page);
    await page.locator('input[type="file"]').setInputFiles(snapshotPath);
    await page.getByRole('button', { name: 'Restore' }).click();
    await expect(page.getByText('Uploading your file...')).toBeVisible();
    await holdForVideo(page);
    await page.getByRole('button', { name: 'Cancel upload' }).click();
    await expect(page.getByText('Uploading your file...')).toBeHidden();
    await expect(page.getByText('The snapshot was successfully uploaded!')).toBeHidden();
    await holdForVideo(page);
    await page.unroute('**/v1/sys/storage/raft/snapshot');
  });

  await test.step('the marker mount still exists, then clean it up', async () => {
    await expect(await enginePaths(page)).toBeVisible();
    await highlightElement(page, page.getByRole('gridcell', { name: `${MARKER_MOUNT}/`, exact: true }));
    await runConsoleCommand(page, `delete sys/mounts/${MARKER_MOUNT}`, 'Success! Data deleted');
  });
});
