/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import type { FieldValue, FormField } from 'vault/forms/v2/form-config';

/**
 * FormV2Field component renders a single form field based on its configuration.
 * Supports multiple value types (string, number, boolean, arrays) and delegates
 * change events to the parent form component.
 * Displays validation errors passed from the parent FormState.
 */
interface Args {
  field: FormField;
  value?: FieldValue;
  errors?: string[];
  onChange: (name: string, value: FieldValue) => void;
}

export default class FormV2Field extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // Log warning for unsupported field types
    this.checkUnsupportedType();
  }

  get name(): string {
    return this.args.field.name;
  }

  get label(): string {
    return this.args.field.label;
  }

  get helperText(): string | undefined {
    return this.args.field.helperText;
  }

  get inputType(): string {
    return this.args.field.inputType || 'text';
  }

  get value(): FieldValue {
    return this.args.value;
  }

  get errors(): string[] {
    return this.args.errors || [];
  }

  get isInvalid(): boolean {
    return this.errors.length > 0;
  }

  /**
   * Checks if the field type is supported and logs a warning for unsupported types.
   * Called in constructor to ensure warning is logged once when component is created.
   */
  private checkUnsupportedType(): void {
    const supportedTypes = [
      'TextInput',
      'TextArea',
      'Select',
      'Toggle',
      'Checkbox',
      'Radio',
      'RadioCard',
      'MaskedInput',
    ];
    const isUnsupported = !supportedTypes.includes(this.args.field.type);

    if (isUnsupported) {
      console.warn(
        `[Form::V2::Field] Unsupported field type "${this.args.field.type}" for field "${this.args.field.label}". Falling back to text input.`
      );
    }
  }

  /**
   * Handles field value changes, supporting multiple value types
   */
  handleChange = (name: string, value: FieldValue): void => {
    this.args.onChange(name, value);
  };
}
