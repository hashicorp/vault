/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { isBlank, isEmpty, typeOf } from '@ember/utils';

import type { FormMetadata } from 'vault/services/path-help';

type GeneratedItemData = Record<string, unknown>;

// Form metadata carries descriptions as inline subText; these pages have always shown them as tooltips.
const withTooltipDescriptions = ({ formFields, formFieldGroups }: FormMetadata) => {
  const fields = new Map(
    formFields.map((field) => {
      const { subText, ...options } = field.options;
      const tooltipField = subText
        ? new FormField(field.name, field.type, { ...options, helpText: subText })
        : field;
      return [field.name, tooltipField];
    })
  );
  const toField = (field: FormField) => fields.get(field.name) as FormField;
  return {
    formFields: formFields.map(toField),
    formFieldGroups: formFieldGroups.flatMap((group) =>
      Object.entries(group).map(([name, groupFields]) => new FormFieldGroup(name, groupFields.map(toField)))
    ),
  };
};

// Applies the attribute transforms the generated Ember Data model used, so payload types are unchanged.
const serializeValue = (type: string | undefined, value: unknown) => {
  switch (type) {
    case 'number': {
      if (isEmpty(value)) return null;
      const number = Number(value);
      return Number.isFinite(number) ? number : null;
    }
    case 'boolean':
      return Boolean(value);
    case 'string':
      return value === null || value === undefined ? null : String(value);
    case 'object':
      return typeOf(value) === 'object' ? value : {};
    default:
      return value;
  }
};

/**
 * Form for an auth method item (user, group, role, certificate) whose fields come from the mount's
 * OpenAPI help at runtime. The path identifier is sent in the URL, not the body, and is read-only on edit.
 */
export default class GeneratedItemForm extends Form<GeneratedItemData> {
  declare formFields: FormField[];
  declare formFieldGroups: FormFieldGroup[];
  declare identifierField: string | undefined;
  declare initialData: GeneratedItemData;

  constructor(metadata: FormMetadata, data: GeneratedItemData = {}, options: { isNew?: boolean } = {}) {
    const { formFields, formFieldGroups } = withTooltipDescriptions(metadata);
    const defaults: GeneratedItemData = {};
    formFields.forEach(({ name, options }) => {
      if (options.defaultValue !== undefined && data[name] === undefined) {
        defaults[name] = options.defaultValue;
      }
    });
    super({ ...defaults, ...data }, options);
    this.formFields = formFields;
    this.formFieldGroups = formFieldGroups;
    this.identifierField = metadata.identifier;
    this.initialData = { ...this.data };
  }

  get id() {
    return this.identifierField ? (this.data[this.identifierField] as string) : undefined;
  }

  // Field groups for read-only views, which show the identifier as the page title instead of a row.
  get fieldGroups() {
    return this.formFieldGroups.map((group) =>
      Object.fromEntries(
        Object.entries(group).map(([name, fields]) => [
          name,
          fields.filter((field) => field.name !== this.identifierField),
        ])
      )
    );
  }

  // Omits the identifier and any blank value the user did not change, as the Ember Data serializer did.
  toJSON() {
    const payload: GeneratedItemData = {};
    this.formFields.forEach(({ name, type, options }) => {
      const value = this.data[name];
      if (options.readOnly || (isBlank(value) && value === this.initialData[name])) {
        return;
      }
      payload[name] = serializeValue(type, value);
    });
    return super.toJSON(payload);
  }
}
