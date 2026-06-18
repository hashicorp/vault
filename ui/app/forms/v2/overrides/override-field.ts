/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { FormConfig, FormField, FormSection } from '../form-config';
import { deepCopyValue } from '../utils/deep-copy';

/**
 * Creates a section builder for easily modifying form configurations.
 *
 * @param generatedConfig - The generated form configuration to build from
 * @returns A builder object with methods for modifying sections
 *
 * @example
 * import generatedConfig from '../generated/mounts-enable-secrets-engine-config';
 * import { configBuilder } from './override-field';
 *
 * const builder = configBuilder(generatedConfig);
 *
 * export default builder
 *   .addSection({
 *     name: 'engine_selection',
 *     title: 'Engine Type',
 *     fields: [
 *       { name: 'engineType', type: 'Radio', label: 'Type', options: [...] }
 *     ]
 *   })
 *   .updateField('default', 'path', { label: 'Mount Path', helperText: 'Custom help' })
 *   .removeField('default', 'MountsEnableSecretsEngineRequest.seal_wrap')
 *   .moveField('path', 'params', 'basic_config')
 *   .build();
 */
export interface SectionBuilder<Request extends object = object, Response = unknown> {
  addSection(section: FormSection, position?: number): SectionBuilder<Request, Response>;
  removeSection(sectionName: string): SectionBuilder<Request, Response>;
  updateSection(
    sectionName: string,
    updates: Partial<Omit<FormSection, 'name' | 'fields'>>
  ): SectionBuilder<Request, Response>;
  addField(sectionName: string, fieldInfo: FormField): SectionBuilder<Request, Response>;
  updateField(
    sectionName: string,
    fieldName: string,
    overrides: Partial<Omit<FormField, 'name'>>
  ): SectionBuilder<Request, Response>;
  removeField(sectionName: string, fieldName: string): SectionBuilder<Request, Response>;
  moveField(
    fieldName: string,
    fromSection: string,
    toSection: string,
    position?: number
  ): SectionBuilder<Request, Response>;
  reorderFields(sectionName: string, fieldNames: string[]): SectionBuilder<Request, Response>;
  build(): FormConfig<Request, Response>;
  getSections(): FormSection[];
}

export const configBuilder = <Request extends object = object, Response = unknown>(
  generatedConfig: FormConfig<Request, Response>
): SectionBuilder<Request, Response> => {
  let sections: FormSection[] = deepCopyValue(generatedConfig.sections) as FormSection[];

  const builder: SectionBuilder<Request, Response> = {
    /**
     * Add a new section to the configuration
     */
    addSection(section: FormSection, position?: number): SectionBuilder<Request, Response> {
      if (position !== undefined) {
        sections.splice(position, 0, section);
      } else {
        sections.push(section);
      }
      return builder;
    },

    /**
     * Remove a section by name
     */
    removeSection(sectionName: string): SectionBuilder<Request, Response> {
      sections = sections.filter((s) => s.name !== sectionName);
      return builder;
    },

    /**
     * Update section properties (title, description, isVisible)
     */
    updateSection(
      sectionName: string,
      updates: Partial<Omit<FormSection, 'name' | 'fields'>>
    ): SectionBuilder<Request, Response> {
      const section = sections.find((s) => s.name === sectionName);
      if (!section) {
        throw new Error(`Section "${sectionName}" not found`);
      }
      Object.assign(section, updates);
      return builder;
    },

    /**
     * Add a field to a section
     */
    addField(sectionName: string, fieldInfo: FormField): SectionBuilder<Request, Response> {
      const section = sections.find((s) => s.name === sectionName);
      if (!section) {
        throw new Error(`Section "${sectionName}" not found`);
      }

      section.fields.push(fieldInfo);
      return builder;
    },

    /**
     * Update a field within a section
     * Merges overrides with the existing field definition
     */
    updateField(
      sectionName: string,
      fieldName: string,
      overrides: Partial<Omit<FormField, 'name'>>
    ): SectionBuilder<Request, Response> {
      const section = sections.find((s) => s.name === sectionName);
      if (!section) {
        throw new Error(`Section "${sectionName}" not found`);
      }

      const fieldIndex = section.fields.findIndex((f) => f.name === fieldName);
      const existingField = section.fields[fieldIndex];
      if (!existingField) {
        throw new Error(`Field "${fieldName}" not found in section "${sectionName}"`);
      }

      // Merge overrides with existing field
      section.fields[fieldIndex] = {
        ...existingField,
        ...overrides,
      };
      return builder;
    },

    /**
     * Remove a field from a section
     */
    removeField(sectionName: string, fieldName: string): SectionBuilder<Request, Response> {
      const section = sections.find((s) => s.name === sectionName);
      if (!section) {
        throw new Error(`Section "${sectionName}" not found`);
      }

      section.fields = section.fields.filter((f) => f.name !== fieldName);
      return builder;
    },

    /**
     * Move a field from one section to another
     */
    moveField(
      fieldName: string,
      fromSection: string,
      toSection: string,
      position?: number
    ): SectionBuilder<Request, Response> {
      const fromIndex = sections.findIndex((s) => s.name === fromSection);
      const toIndex = sections.findIndex((s) => s.name === toSection);

      const from = sections[fromIndex];
      const to = sections[toIndex];

      if (!from) {
        throw new Error(`Source section "${fromSection}" not found`);
      }
      if (!to) {
        throw new Error(`Target section "${toSection}" not found`);
      }

      const fieldIndex = from.fields.findIndex((f) => f.name === fieldName);
      if (fieldIndex === -1) {
        throw new Error(`Field "${fieldName}" not found in section "${fromSection}"`);
      }

      // Reuse the already-cloned field object. Avoid JSON serialization as it strips functions.
      const movedField = from.fields[fieldIndex];
      if (!movedField) {
        throw new Error(`Field "${fieldName}" not found in section "${fromSection}"`);
      }

      // Remove from source section
      from.fields = from.fields.filter((f) => f.name !== fieldName);

      // Add to target section at specified position
      if (position !== undefined) {
        to.fields.splice(position, 0, movedField);
      } else {
        to.fields.push(movedField);
      }

      return builder;
    },

    /**
     * Reorder fields within a section
     */
    reorderFields(sectionName: string, fieldNames: string[]): SectionBuilder<Request, Response> {
      const section = sections.find((s) => s.name === sectionName);
      if (!section) {
        throw new Error(`Section "${sectionName}" not found`);
      }

      const fieldMap = new Map(section.fields.map((f) => [f.name, f]));
      const reorderedFields: FormField[] = [];

      for (const name of fieldNames) {
        const field = fieldMap.get(name);
        if (!field) {
          throw new Error(`Field "${name}" not found in section "${sectionName}"`);
        }
        reorderedFields.push(field);
        fieldMap.delete(name);
      }

      // Add any remaining fields that weren't in the reorder list
      reorderedFields.push(...fieldMap.values());
      section.fields = reorderedFields;
      return builder;
    },

    /**
     * Build and return the final configuration
     * Preserves all non-section properties (including functions) from the original config
     */
    build(): FormConfig<Request, Response> {
      return {
        ...generatedConfig,
        sections,
      };
    },

    /**
     * Get the current sections (useful for debugging)
     */
    getSections(): FormSection[] {
      return sections;
    },
  };

  return builder;
};

/**
 * Simplified helper for overriding multiple fields in a section.
 * This is the recommended pattern for most override use cases.
 *
 * @param generatedConfig - The generated form configuration to override
 * @param sectionName - The name of the section containing the fields to override
 * @param fieldOverrides - Object mapping field names to their override properties
 * @returns A new FormConfig with the overrides applied
 *
 * @example
 * import generatedConfig from '../generated/mounts-enable-secrets-engine-config';
 * import { overrideFieldsInSection } from '../overrides/override-field';
 *
 * export default overrideFieldsInSection(
 *   generatedConfig,
 *   'default',
 *   {
 *     'path': {
 *       label: 'Mount Path',
 *       helperText: 'Where to mount this secrets engine'
 *     },
 *     'description': {
 *       type: 'TextArea'
 *     }
 *   }
 * );
 */
export const overrideFieldsInSection = <Request extends object = object, Response = unknown>(
  generatedConfig: FormConfig<Request, Response>,
  sectionName: string,
  fieldOverrides: Record<string, Partial<Omit<FormField, 'name'>>>
): FormConfig<Request, Response> => {
  const builder = configBuilder(generatedConfig);

  // Apply each field override
  for (const [fieldName, overrides] of Object.entries(fieldOverrides)) {
    builder.updateField(sectionName, fieldName, overrides);
  }

  return builder.build();
};
