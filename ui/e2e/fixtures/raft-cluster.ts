/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test as base, expect } from '@playwright/test';
import { spawn } from 'child_process';
import { mkdtemp, mkdir, rm, writeFile } from 'fs/promises';
import { createServer } from 'net';
import { tmpdir } from 'os';
import path from 'path';
import { USER_POLICY_MAP } from '../policies';
import { TEST_USERS } from '../test-users';

import type { ChildProcess } from 'child_process';
import type { Server } from 'net';

interface RaftServer {
  node_id: string;
  address: string;
  leader: boolean;
  voter: boolean;
}

interface RaftNode {
  id: string;
  apiURL: string;
  clusterAddress: string;
}

interface RaftCluster {
  leader: RaftNode;
  peer: RaftNode;
  unsealKey: string;
  token: string;
  readServers: () => Promise<RaftServer[]>;
}

/** Destructive membership tests get fresh nodes and ports, including on retries. */
export const test = base.extend<{ raftCluster: RaftCluster }>({
  raftCluster: [
    async ({ playwright }, use) => {
      const directory = await mkdtemp(path.join(tmpdir(), 'vault-playwright-raft-'));
      const reservations: Server[] = [];
      const processes: { child: ChildProcess; closed: Promise<void> }[] = [];
      const request = await playwright.request.newContext();

      const reserveAddress = async () => {
        const server = createServer();
        reservations.push(server);
        await new Promise<void>((resolve, reject) => {
          server.once('error', reject);
          server.listen(0, '127.0.0.1', resolve);
        });
        const address = server.address();
        if (!address || typeof address === 'string') {
          throw new Error('Expected a TCP address for the Raft test node');
        }
        return `127.0.0.1:${address.port}`;
      };

      const startNode = async (id: string): Promise<RaftNode> => {
        const address = await reserveAddress();
        const clusterAddress = await reserveAddress();
        const apiURL = `http://${address}`;
        const dataPath = path.join(directory, id);
        await mkdir(dataPath);
        const configPath = path.join(directory, `${id}.json`);
        await writeFile(
          configPath,
          JSON.stringify({
            ui: true,
            disable_mlock: true,
            api_addr: apiURL,
            cluster_addr: `https://${clusterAddress}`,
            storage: { raft: { path: dataPath, node_id: id } },
            listener: { tcp: { address, cluster_address: clusterAddress, tls_disable: true } },
          })
        );
        // Hold both ports until the node is ready to bind; never reuse another test's server.
        for (const server of reservations.splice(0)) {
          await new Promise<void>((resolve, reject) =>
            server.close((error) => (error ? reject(error) : resolve()))
          );
        }
        const child = spawn('vault', ['server', `-config=${configPath}`, '-log-level=error'], {
          stdio: 'ignore',
        });
        let startupError: Error | undefined;
        child.once('error', (error) => (startupError = error));
        const closed = new Promise<void>((resolve) => child.once('close', () => resolve()));
        processes.push({ child, closed });
        await expect
          .poll(
            async () => {
              if (startupError) throw startupError;
              if (child.exitCode !== null || child.signalCode !== null) {
                throw new Error(`Raft test node ${id} exited before becoming ready`);
              }
              try {
                const response = await request.get(`${apiURL}/v1/sys/seal-status`);
                await expect(response).toBeOK();
                return (await response.json()).storage_type;
              } catch (error) {
                if (error instanceof Error && error.message.includes('ECONNREFUSED')) return 'starting';
                throw error;
              }
            },
            { timeout: 30_000, message: `Raft node ${id} starts with real Raft storage` }
          )
          .toBe('raft');
        return { id, apiURL, clusterAddress };
      };

      try {
        const leader = await startNode('leader');
        const peer = await startNode('peer');
        const init = await request.put(`${leader.apiURL}/v1/sys/init`, {
          data: { secret_shares: 1, secret_threshold: 1 },
        });
        await expect(init).toBeOK();
        const { keys, root_token } = await init.json();
        const unseal = await request.put(`${leader.apiURL}/v1/sys/unseal`, {
          data: { key: keys[0] },
        });
        await expect(unseal).toBeOK();
        const headers = { 'X-Vault-Token': root_token };
        await expect
          .poll(async () => (await request.get(`${leader.apiURL}/v1/sys/health`)).status())
          .toBe(200);
        const policy = TEST_USERS.raft.policy;
        const policyResponse = await request.put(`${leader.apiURL}/v1/sys/policies/acl/${policy}`, {
          headers,
          data: { policy: USER_POLICY_MAP[policy] },
        });
        await expect(policyResponse).toBeOK();
        const tokenResponse = await request.post(`${leader.apiURL}/v1/auth/token/create`, {
          headers,
          data: { policies: [policy], ttl: '10m' },
        });
        await expect(tokenResponse).toBeOK();
        const { auth } = await tokenResponse.json();
        expect(auth.policies).toContain(policy);
        expect(auth.policies).not.toContain('root');

        await use({
          leader,
          peer,
          unsealKey: keys[0],
          token: auth.client_token,
          readServers: async () => {
            const response = await request.get(`${leader.apiURL}/v1/sys/storage/raft/configuration`, {
              headers: { 'X-Vault-Token': auth.client_token },
            });
            await expect(response).toBeOK();
            return (await response.json()).data.config.servers;
          },
        });
      } finally {
        try {
          await Promise.all(
            processes.map(async ({ child, closed }) => {
              if (child.exitCode === null && child.signalCode === null) child.kill('SIGTERM');
              const forceStop = setTimeout(() => child.kill('SIGKILL'), 5000);
              try {
                await closed;
              } finally {
                clearTimeout(forceStop);
              }
            })
          );
        } finally {
          await request.dispose();
          for (const server of reservations) server.close();
          await rm(directory, { recursive: true, force: true });
        }
      }
    },
    { timeout: 120_000 },
  ],
});
