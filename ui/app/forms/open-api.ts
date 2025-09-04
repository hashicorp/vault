/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { propsForSchema } from 'vault/utils/openapi-helpers';

import type { OpenApiHelpResponse } from 'vault/utils/openapi-helpers';

export default class OpenApiForm<T extends object> extends Form<T> {
  declare formFieldGroups: FormFieldGroup[];

  constructor(helpResponse: OpenApiHelpResponse, ...formArgs: ConstructorParameters<typeof Form>) {
    super(...formArgs);
    // create formFieldGroups from the OpenAPI properties
    const props = propsForSchema(helpResponse);
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
      groups[group]?.[arrMethod](new FormField(name, type, options));
      // set the default value on the data object
      if (defaultValue && this.data[name as keyof typeof this.data] === undefined) {
        this.data = { ...this.data, [name]: defaultValue };
      }
    }

    // ensure default group is the first item in the formFieldGroups

    // create formFieldGroups from the expanded groups
    this.formFieldGroups = Object.entries(groups).reduce<FormFieldGroup[]>(
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
  }
}
