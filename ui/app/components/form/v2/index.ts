/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import type V2Form from 'vault/forms/v2/v2-form';
import type ApiService from 'vault/services/api';

interface Args {
  form: V2Form;
  hideFields?: boolean;
  onSuccess?: (response: unknown) => void;
  onError?: (errorMessage: string) => void;
}

export default class FormV2 extends Component<Args> {
  @service declare readonly api: ApiService;
  @tracked submissionError?: string;
  @tracked lastResponse?: unknown;

  get form() {
    return this.args.form;
  }

  /**
   * Handles form submission with validation and API call.
   * Uses ember-concurrency to prevent double-submission and provide derived state.
   * Invokes component callbacks before form config callbacks for wizard orchestration.
   */
  submitTask = task({ drop: true }, async () => {
    const { isValid } = this.form.validateForm();
    if (!isValid) return;

    try {
      const response = await this.form.submit(this.api);
      this.lastResponse = response;

      // Call component's onSuccess first (for wizard orchestration)
      this.args.onSuccess?.(response);

      // Then call form config's onSuccess (for custom logic)
      this.form.config.onSuccess?.(response);

      return response;
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.submissionError = message;

      // Call component's onError first (for wizard orchestration)
      this.args.onError?.(message);

      // Then call form config's onError (for custom logic)
      this.form.config.onError?.(message);

      // Re-throw to maintain task error state
      throw error;
    }
  });
}
