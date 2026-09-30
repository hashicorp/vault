/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { SecretsApiKvV2ListListEnum } from '@hashicorp/vault-client-typescript';
import { didCancel, restartableTask, task, timeout } from 'ember-concurrency';
import { isTransientKvMountError } from 'vault/utils/kv-mount-transient-error';

import type ApiService from 'vault/services/api';
import type { Task } from 'ember-concurrency';

const DEFAULT_MAX_ATTEMPTS = 10;
const DEFAULT_RETRY_DELAY_MS = 500;

interface RetryOptions {
  maxAttempts?: number;
  delayMs?: number;
}

type ParseError = (error: unknown) => Promise<{ status?: number; message?: string }>;
type Request = () => Promise<unknown>;
type RequestTask = Task<unknown, [Request, ParseError, RetryOptions?]>;

// Retries known KV responses while a newly created mount propagates to the current node.
// This is a narrowly scoped workaround, not a general API retry mechanism; use it only after mount
// creation succeeds and retrying the follow-up read is safe.
export default class KvMountRetryService extends Service {
  @service declare api: ApiService;

  @task
  *request(
    requestFn: Request,
    parseError: ParseError,
    options: RetryOptions = {}
  ): Generator<unknown, unknown, unknown> {
    const maxAttempts = options.maxAttempts ?? DEFAULT_MAX_ATTEMPTS;
    const delayMs = options.delayMs ?? DEFAULT_RETRY_DELAY_MS;

    for (let attempt = 1; attempt < maxAttempts; attempt++) {
      try {
        return yield requestFn();
      } catch (error) {
        if (didCancel(error)) throw error;

        const { status, message } = (yield parseError(error)) as Awaited<ReturnType<ParseError>>;

        if (!isTransientKvMountError(status, message)) {
          throw error;
        }

        yield timeout(delayMs);
      }
    }

    return yield requestFn();
  }

  loadMountInfo = restartableTask(async (backend: string) => {
    return (this.request as unknown as RequestTask).perform(
      () => this.api.sys.internalUiReadMountInformation(backend),
      (error) => this.api.parseError(error)
    );
  });

  listKvSecrets = restartableTask(async (pathToSecret: string, backend: string) => {
    return (this.request as unknown as RequestTask).perform(
      () => this.api.secrets.kvV2List(pathToSecret, backend, SecretsApiKvV2ListListEnum.TRUE),
      (error) => this.api.parseError(error)
    );
  });
}
