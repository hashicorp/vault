/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { get, set } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import type ApiService from 'vault/services/api';
import type { FieldValue, FormConfig, FormField, VisibilityRule } from './form-config';
import { validateAllFields, validateField } from './form-validator';
import { getFormConfig, type FormConfigKey } from './get-form-config';

/**
 * V2Form manages the state of a form with type-safe property updates and validation.
 * Supports nested property updates via dotted-path notation (e.g., "user.address.street").
 * Automatically validates fields when values change.
 *
 * Usage with registry name:
 * ```typescript
 * const form = new V2Form('mountsEnableSecretsEngine');
 * form.set('path', 'my-path');
 * const errors = form.getErrors('path');
 * await form.submit(api);
 * ```
 *
 * Usage with config object:
 * ```typescript
 * const form = new V2Form(myFormConfig);
 * form.set('path', 'my-path');
 * await form.submit(api);
 * ```
 *
 * @template TPayload - The form payload type
 * @template TResponse - The API response type
 */
export default class V2Form<TPayload extends object = any, TResponse = unknown> {
  @tracked payload: TPayload;
  @tracked validationErrors: Map<string, string[]> = new Map();
  #config: FormConfig<TPayload, TResponse>;

  /**
   * NOTE: Flexible constructor pattern
   *
   * Supports two instantiation methods to enable flexible form usage:
   * 1. Name-based: `new V2Form('mountsEnableSecretsEngine')` - looks up config from registry
   *    - Use for single-step forms with globally registered configs
   *
   * 2. Config-based: `new V2Form<PayloadType, ResponseType>(configObject)` - accepts config directly
   *    - Use for wizard steps with local overrides (e.g., dynamic payloads)
   *    - Explicit generic parameters should match the config's types
   *    - For wizards, use `new V2Form<any, any>(config)` to simplify typing
   */
  constructor(config: FormConfigKey | FormConfig<TPayload, TResponse>) {
    if (typeof config === 'string') {
      // Registry-based: getFormConfig returns FormConfig<never, never> but we need
      // FormConfig<TPayload, TResponse>. Since generics default to `any`, this
      // type assertion is safe and unavoidable due to TypeScript's structural typing.
      const formConfig = getFormConfig(config);
      this.#config = formConfig as unknown as FormConfig<TPayload, TResponse>;
      this.payload = this.#resolvePayload(formConfig.payload as TPayload);
    } else {
      // Config-based: types must match the provided config
      this.#config = config;
      this.payload = this.#resolvePayload(config.payload);
    }

    // Auto-inject required validations for fields marked with isRequired: true
    this.#injectRequiredValidations();
  }

  /**
   * Automatically adds required validation rules to fields marked with isRequired: true
   * if they don't already have a required validation.
   */
  #injectRequiredValidations(): void {
    for (const section of this.#config.sections) {
      for (const field of section.fields) {
        if (field.isRequired) {
          // Check if field already has a required validation
          const hasRequiredValidation = field.validations?.some(
            (rule) => 'type' in rule && rule.type === 'required'
          );

          if (!hasRequiredValidation) {
            // Add required validation
            if (!field.validations) {
              field.validations = [];
            }
            field.validations.unshift({
              type: 'required',
              message: `${field.label} is required`,
            });
          }
        }
      }
    }
  }

  /**
   * Resolves the payload from the config.
   * If the config.payload is a function, it's assumed to be a static payload
   * (wizard payload resolution happens at the ProvisionForm level).
   * This method just extracts the value.
   */
  #resolvePayload(payload: TPayload | ((wizardState: any) => TPayload)): TPayload {
    // For V2Form, we expect the payload to already be resolved
    // (either a static object or pre-resolved by the parent component)
    const resolvePayload = typeof payload === 'function' ? payload({}) : payload;
    return structuredClone(resolvePayload);
  }

  get config(): FormConfig<TPayload, TResponse> {
    return this.#config;
  }

  get isValid(): boolean {
    return this.validationErrors.size === 0;
  }

  /**
   * Updates a property in the payload using dotted-path notation.
   * Creates intermediate objects if they don't exist.
   * Automatically validates the field after updating.
   *
   * @param propPath - Dotted-path to the property (e.g., "user.address.street")
   * @param value - The new value for the property
   */
  set(propPath: string, value: unknown): void {
    const nextPayload = structuredClone(this.payload);
    set(nextPayload, propPath, value);
    this.payload = nextPayload;
    this.#pruneHiddenFieldErrors();
    this.#validateField(propPath);
  }

  /**
   * Get validation errors for a specific field
   */
  getErrors(propPath: string): string[] {
    return this.validationErrors.get(propPath) ?? [];
  }

  /**
   * Validate all fields in the form.
   * Useful before form submission to show all validation errors.
   */
  validateForm(): { isValid: boolean } {
    const errors = validateAllFields(this.#visibleFields, this.payload);
    this.validationErrors = errors;

    return {
      isValid: this.isValid,
    };
  }

  /**
   * Submit the form after validation.
   * Validates the entire form and calls the config's submit handler.
   *
   * @param api - The API service instance
   * @returns Promise resolving to the API response
   * @throws Error if form validation fails
   */
  async submit(api: ApiService): Promise<TResponse> {
    const { isValid } = this.validateForm();

    if (!isValid) {
      throw new Error('Form validation failed');
    }

    return this.#config.submit(api, this.payload);
  }

  get #visibleFields(): FormField[] {
    return this.#config.sections
      .filter((section) => this.#isVisible(section.isVisible))
      .flatMap((section) => section.fields.filter((field) => this.#isVisible(field.isVisible)));
  }

  #isVisible(rule?: VisibilityRule): boolean {
    if (typeof rule === 'function') {
      return rule(this.payload);
    }

    if (typeof rule === 'boolean') {
      return rule;
    }

    return true;
  }

  #pruneHiddenFieldErrors(): void {
    const visibleFieldNames = new Set(this.#visibleFields.map((field) => field.name));
    const nextErrors = new Map(
      [...this.validationErrors].filter(([fieldName]) => visibleFieldNames.has(fieldName))
    );

    if (nextErrors.size !== this.validationErrors.size) {
      this.validationErrors = nextErrors;
    }
  }

  /**
   * Search through the config structure to find a field configuration object by its name.
   */
  #findField(propPath: string): FormField | null {
    return this.#visibleFields.find((field) => field.name === propPath) ?? null;
  }

  #validateField(propPath: string): void {
    const field = this.#findField(propPath);
    if (!field) {
      this.validationErrors.delete(propPath);
      this.validationErrors = new Map(this.validationErrors);
      return;
    }

    const value = get(this.payload, propPath) as FieldValue;
    const errors = validateField(field, value, this.payload);

    if (errors.length > 0) {
      this.validationErrors.set(propPath, errors);
    } else {
      this.validationErrors.delete(propPath);
    }

    // Trigger reactivity by creating a new Map
    this.validationErrors = new Map(this.validationErrors);
  }
}
