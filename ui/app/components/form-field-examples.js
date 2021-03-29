/**
 * @module FormFieldExamples
 * FormFieldExamples components are used to...
 *
 * @example
 * ```js
 * <FormFieldExamples @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

const EXAMPLES = [
  {
    name: 'regularInput',
    type: 'string',
  },
  {
    name: 'input',
    type: 'string',
    options: {
      label: 'Custom input',
      defaultValue: 'dogs',
      helpText: 'Help text is here',
      subText: 'Sub text is here',
    },
  },
  {
    name: 'regularSensitive',
    type: 'string',
    options: {
      sensitive: true,
    },
  },
  {
    name: 'customSensitive',
    type: 'string',
    options: {
      label: 'Some secret value',
      defaultValue: 'dogs',
      helpText: 'Help text is here',
      subText: 'Sub text is here',
      sensitive: true,
    },
  },
  {
    name: 'regularBoolean',
    type: 'boolean',
  },
  {
    name: 'stringBoolean',
    type: 'string',
    options: {
      editType: 'boolean',
      trueValue: 'on',
      falseValue: 'off',
    },
  },
  {
    name: 'timeToLive',
    type: 'string',
    options: {
      editType: 'ttl',
    },
  },
  {
    name: 'someObject',
    type: 'object',
  },
  {
    name: 'kvSecret',
    type: 'string',
    options: {
      editType: 'kv',
    },
  },
];

const reverseOpenApi = attr => {
  // description, items, type, format, isId, deprecated, enum
  // x-dis name, value, group, sensitive, editType
  let returnVal = {
    Type: attr.type,
    DisplayAttrs: {},
  };
  if (attr.options?.helpText) {
    returnVal.Description = attr.options.helpText;
  }
  if (attr.options?.editType) {
    returnVal.DisplayAttrs.EditType = attr.options.editType;
  }
  if (attr.options?.defaultValue) {
    returnVal.DisplayAttrs.Value = attr.options.defaultValue;
  }
  if (attr.options?.label) {
    returnVal.DisplayAttrs.Name = attr.options.label;
  }
  if (attr.options?.sensitive) {
    returnVal.DisplayAttrs.Sensitive = attr.options.sensitive;
  }
  console.log(returnVal);
  if (returnVal.DisplayAttrs === {}) {
    delete returnVal.DisplayAttrs;
  }
  const stringed = JSON.stringify(returnVal);
  console.log(stringed);
  return stringed;
};

export default class FormFieldExamples extends Component {
  @tracked
  model = {
    set: (path, val) => {
      console.log('TODO: set', path, val);
    },
  };

  get examples() {
    return EXAMPLES.map(e => ({
      attr: e,
      string: JSON.stringify(e, null, 2),
    }));
  }
}
