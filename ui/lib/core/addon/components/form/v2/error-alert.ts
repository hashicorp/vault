/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  /** Error message to display */
  error?: string;
  /** Optional custom title (defaults to "Submission error") */
  title?: string;
}

/**
 * Form::V2::ErrorAlert displays submission errors in a consistent format.
 * Used by both Form::V2 and Form::V2::Wizard for standardized error display.
 *
 * @example
 * ```handlebars
 * <Form::V2::ErrorAlert @error={{this.submissionError}} />
 * <Form::V2::ErrorAlert @error={{this.error}} @title="Configuration Error" />
 * ```
 */
export default class FormV2ErrorAlert extends Component<Args> {
  get title(): string {
    return this.args.title || 'Submission error';
  }
}
