/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { expandOpenApiProps } from 'vault/utils/openapi-helpers';
import openApiSpec from '@hashicorp/vault-client-typescript/openapi.json';

import type { OpenApiProps } from 'vault/utils/openapi-helpers';

export default class OpenApiForm<T extends object> extends Form<T> {
  formFieldGroups: FormFieldGroup[] = [];
  formFields: FormField[] = [];

  constructor(schemaKey: string, ...formArgs: ConstructorParameters<typeof Form>) {
    const [data = {}, ...restArgs] = formArgs;
    const defaultValues = {} as Partial<T>;
    const formFields: FormField[] = [];
    let formFieldGroups: FormFieldGroup[] = [];
    // find the schema in the OpenAPI spec that contains the properties for dynamic form generation
    const schema = openApiSpec.components.schemas[schemaKey as keyof typeof openApiSpec.components.schemas];
    // there could be an instance where the schema isn't found but an inherited class has defined static fields
    // in that case we will simply bypass dynamic form generation in favor of throwing an error
    if (schema) {
      // create formFieldGroups from the OpenAPI properties
      const props = expandOpenApiProps(schema.properties as OpenApiProps, 'form');
      const groups: { [groupName: string]: FormField[] } = {};
      // iterate over the properties and organize them into groups
      for (const [name, prop] of Object.entries(props)) {
        // disabling lint rule since we need to ignore certain options returned from expandOpenApiProps util
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { fieldGroup, fieldValue, type, defaultValue, ...options } = prop;
        // groupName from groupsMap takes precedence over fieldGroup from the property
        const group = fieldGroup || 'default';
        // organize the form fields so we can create formFieldGroups later
        if (!(group in groups)) {
          groups[group] = [];
        }
        // create a new FormField for the property and associate it with the appropriate fieldGroup
        // props marked as `identifier` are primary fields that should be rendered first in the form
        const arrMethod = options.identifier ? 'unshift' : 'push';
        const field = new FormField(name, type, options);
        groups[group]?.[arrMethod](field);
        // push all fields to formFields array for convenience access and flexibility
        formFields.push(field);
        // set the default value from the schema if it is not already set in the data
        if (defaultValue !== undefined && data[name as keyof typeof data] === undefined) {
          defaultValues[name as keyof T] = defaultValue as T[keyof T];
        }
      }

      // create formFieldGroups from the expanded groups
      formFieldGroups = Object.entries(groups).reduce<FormFieldGroup[]>(
        (formFieldGroups, [groupName, fields]) => {
          const group = new FormFieldGroup(groupName, fields);
          // ensure the default group is the first group to render
          if (groupName === 'default') {
            return [group, ...formFieldGroups];
          }
          return [...formFieldGroups, group];
        },
        []
      );
    } else {
      // to aide in development log out an error to the console if the schema is not found
      console.error(`OpenApiForm: Schema '${schemaKey}' not found in OpenAPI spec.`);
    }

    // call the super constructor with the merged default values and data
    super({ ...defaultValues, ...data }, ...restArgs);
    // add the generated form fields and groups to the instance
    // this allows for a base class to define static form fields/groups and have them merged with the dynamic ones if necessary
    this.formFields.push(...formFields);
    this.formFieldGroups.push(...formFieldGroups);
  }
}
