/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/raft-cluster';

// Unseal and login enter secrets; do not retain them in trace/video/screenshot artifacts.
test.use({
  storageState: { cookies: [], origins: [] },
  trace: 'off',
  video: 'off',
  screenshot: 'off',
});

test('joins and removes a real Raft peer with persisted membership', async ({ page, raftCluster }) => {
  test.setTimeout(120_000);
  const { leader, peer, token, unsealKey, readServers } = raftCluster;

  await test.step('starts with one leader and an uninitialized peer', async () => {
    const servers = await readServers();
    expect(servers).toEqual([
      expect.objectContaining({
        node_id: leader.id,
        address: leader.clusterAddress,
        leader: true,
        voter: true,
      }),
    ]);
    const status = await page.request.get(`${peer.apiURL}/v1/sys/seal-status`);
    await expect(status).toBeOK();
    expect(await status.json()).toMatchObject({ initialized: false, sealed: true });
    await test.info().attach('initial-membership', {
      body: JSON.stringify(servers),
      contentType: 'application/json',
    });
  });

  await test.step('joins through the UI and sends the leader address to the real backend', async () => {
    await page.goto(`${peer.apiURL}/ui/vault`);
    await page.getByRole('radio', { name: 'Join an existing Raft cluster' }).check();
    await page.getByRole('button', { name: 'Next' }).click();
    await page.getByRole('textbox', { name: 'Leader API Address' }).fill(leader.apiURL);
    await page.getByRole('checkbox', { name: 'Keep retrying to join in case of failures' }).uncheck();
    const responsePromise = page.waitForResponse(
      (response) =>
        response.url() === `${peer.apiURL}/v1/sys/storage/raft/join` && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Join', exact: true }).click();
    const response = await responsePromise;
    expect(response.ok()).toBe(true);
    expect(response.request().postDataJSON()).toMatchObject({ leader_api_addr: leader.apiURL });
    expect(response.request().postDataJSON().retry ?? false).toBe(false);
    await expect(page.getByRole('textbox', { name: 'Unseal Key Portion' })).toBeVisible();
    await page.getByRole('textbox', { name: 'Unseal Key Portion' }).fill(unsealKey);
    await page.getByRole('button', { name: 'Unseal', exact: true }).click();
    await expect
      .poll(readServers, { timeout: 60_000, message: 'The joined peer becomes a Raft voter' })
      .toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            node_id: leader.id,
            address: leader.clusterAddress,
            leader: true,
            voter: true,
          }),
          expect.objectContaining({
            node_id: peer.id,
            address: peer.clusterAddress,
            leader: false,
            voter: true,
          }),
        ])
      );
    expect(await readServers()).toHaveLength(2);
    const status = await page.request.get(`${peer.apiURL}/v1/sys/seal-status`);
    await expect(status).toBeOK();
    expect(await status.json()).toMatchObject({ initialized: true, sealed: false });
  });

  await test.step('reads the joined membership in the UI and after reload', async () => {
    await page.goto(`${leader.apiURL}/ui/vault/auth`);
    await page.getByRole('textbox', { name: 'Token' }).fill(token);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await page.getByRole('link', { name: 'Raft storage' }).click();
    const rows = page.locator('[data-raft-row]');
    await expect(rows).toHaveCount(2);
    const peerRow = rows.filter({ hasText: peer.clusterAddress });
    await expect(peerRow.getByRole('cell').nth(1)).toHaveText('No');
    await expect(peerRow.getByRole('cell').nth(2)).toHaveText('Yes');
    await page.reload();
    await expect(peerRow).toHaveCount(1);
    await expect(rows).toHaveCount(2);
    await test.info().attach('joined-membership-after-reload', {
      body: JSON.stringify({ servers: await readServers(), rows: await rows.allTextContents() }),
      contentType: 'application/json',
    });
  });

  await test.step('removes only the peer through the UI and verifies backend read-back', async () => {
    const peerRow = page.locator('[data-raft-row]').filter({ hasText: peer.clusterAddress });
    await peerRow.getByRole('button', { name: 'Raft server actions' }).click();
    await page.getByRole('button', { name: 'Remove Peer', exact: true }).click();
    const dialog = page.getByRole('dialog');
    await expect(dialog).toContainText(`Remove ${peer.id}?`);
    const responsePromise = page.waitForResponse(
      (response) =>
        response.url() === `${leader.apiURL}/v1/sys/storage/raft/remove-peer` &&
        response.request().method() === 'POST'
    );
    await dialog.getByRole('button', { name: 'Confirm', exact: true }).click();
    const response = await responsePromise;
    expect(response.ok()).toBe(true);
    expect(response.request().postDataJSON()).toEqual({ server_id: peer.id });
    await expect.poll(readServers).toEqual([
      expect.objectContaining({
        node_id: leader.id,
        address: leader.clusterAddress,
        leader: true,
        voter: true,
      }),
    ]);
    await expect(page.locator('[data-raft-row]')).toHaveCount(1);
    await expect(peerRow).toHaveCount(0);
    await page.reload();
    await expect(page.locator('[data-raft-row]')).toHaveCount(1);
    await expect(page.locator('[data-raft-row]')).toContainText(leader.clusterAddress);
    const servers = await readServers();
    expect(servers).toEqual([
      expect.objectContaining({
        node_id: leader.id,
        address: leader.clusterAddress,
        leader: true,
        voter: true,
      }),
    ]);
    await test.info().attach('removed-membership-after-reload', {
      body: JSON.stringify({
        servers,
        rows: await page.locator('[data-raft-row]').allTextContents(),
      }),
      contentType: 'application/json',
    });
  });
});
