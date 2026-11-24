/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import type { FieldValue, FormField, FormSection, VisibilityRule } from 'vault/forms/v2/form-config';
import type V2Form from 'vault/forms/v2/v2-form';

interface Args {
  /** The V2Form instance containing payload, config, and validation state */
  form: V2Form;
  /** Optional error message to display at the top of the form */
  error?: string;
  /** Whether to render form fields (used by wizard's apply step) */
  renderFields?: boolean;
}

/**
 * Form::V2::Renderer is a shared form rendering component that encapsulates
 * the common form structure (Hds::Form, error alerts, sections, and fields).
 *
 * Usage:
 * ```handlebars
 * <Form::V2::Renderer @form={{this.form}} @error={{this.submissionError}} as |Form|>
 *   <Form.Section>
 *     <Hds::Button @text="Submit" type="submit" ... />
 *   </Form.Section>
 * </Form::V2::Renderer>
 * ```
 *
 * The component yields the Hds::Form context for consumers to define
 * their own submit/navigation UI.
 */
export default class FormV2Renderer extends Component<Args> {
  /**
   * Get validation errors for a specific field.
   * Arrow function to preserve `this` context when called from template.
   */
  getFieldErrors = (fieldName: string): string[] => {
    return this.args.form.getErrors(fieldName);
  };

  /**
   * Handle field value changes.
   * Arrow function to preserve `this` context when called from template.
   */
  handleFieldChange = (name: string, value: FieldValue): void => {
    this.args.form.set(name, value);
  };

  isSectionVisible = (section: FormSection): boolean => {
    return this.#isVisible(section.isVisible);
  };

  isFieldVisible = (field: FormField): boolean => {
    return this.#isVisible(field.isVisible);
  };

  #isVisible(rule?: VisibilityRule): boolean {
    if (typeof rule === 'function') {
      return rule(this.args.form.payload);
    }

    if (typeof rule === 'boolean') {
      return rule;
    }

    return true;
  }
}
