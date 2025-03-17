/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
// import { tracked } from '@glimmer/tracking';

export default class ShowcaseController extends Controller {
  // ###########################################
  // FORMFIELD
  // ###########################################

  jsonExample = JSON.stringify(
    { type: 'JSON', text: 'Lorem ipsum', value: 1234, object: { a: 'b', c: [9, 8, 7] } },
    null,
    2
  );

  @action
  dynamicFormFieldModel(editType, variant) {
    const model = {
      // emulate model's setter
      set: () => {},
    };
    if (editType === 'boolean') {
      if (variant === 'checked' || variant === 'disabled') {
        model.boolean = true;
      }
    } else if (editType === 'checkboxList') {
      model.checkboxList = 'Ipsum';
    } else if (editType === 'dateTimeLocal') {
      if (variant === 'with value' || variant === 'with validation error') {
        model.dateTimeLocal = new Date();
      }
    } else if (editType === 'input') {
      if (variant === 'with value' || variant === 'with validation error') {
        model.input = 'Lorem ipsum dolor sit amet';
      } else if (variant === 'with character limit') {
        model.input = '123456789';
      } else if (variant === 'readonly' || variant === 'disabled') {
        model.input = 'Lorem ipsum';
      }
    } else if (editType === 'json') {
      model.json = '{}';
      if (variant === 'with value' || variant === 'with value + restore') {
        model.json = this.jsonExample;
      }
    } else if (editType === 'kv') {
      if (variant === 'with value' || variant === 'with validation error') {
        model.kv = {
          'my-key': 'This is the value for `my-key`',
        };
      } else if (variant === 'with whitespace in key') {
        model.kv = {
          'my key with space': 'You need to set `@allowWhiteSpace` to avoid this warning',
        };
      }
    } else if (editType === 'mountAccessor') {
      if (variant === 'with value') {
        // TODO! not sure what to use here, it's too connected with the `authMethods` task (see the `mount-accessor-select` controller)
        model.mountAccessor = '???';
      }
    } else if (editType === 'object') {
      model.object = '{}';
      if (variant === 'with value' || variant === 'with value + restore') {
        model.object = { type: 'object', text: 'Hello world', number: 4567, array: ['a', 'b', 'c'] };
      }
    } else if (editType === 'optionalText') {
      if (variant === 'with value' || variant === 'with value + subText + docLink') {
        model.optionalText = 'Lorem ipsum';
      }
    } else if (editType === 'password') {
      if (variant === 'with value') {
        model.password = '123abc';
      }
    } else if (editType === 'radio') {
      model.radio = 2;
    } else if (editType === 'regex') {
      if (variant === 'with value') {
        model.regex = '/^lorem .*$/i';
      }
    } else if (editType === 'stringArray') {
      if (variant === 'with value') {
        model.stringArray = ['Lorem', 'Ipsum'];
      }
    } else if (editType === 'searchSelect') {
      if (variant === 'with value') {
        model.searchSelect = ['Lorem', 'Ipsum'];
      }
    } else if (editType === 'select') {
      if (variant === 'with value') {
        model.select = 'Ipsum';
      }
    } else if (editType === 'sensitive') {
      if (variant === 'with value' || variant === 'with copy') {
        model.sensitive = 'Lorem ipsum dolor';
      }
    } else if (editType === 'textarea') {
      if (variant === 'with value' || variant === 'with validation error') {
        model.textarea = 'Lorem\nipsum\ndolor';
      }
    } else if (editType === 'ttl') {
      if (variant === 'with value' || variant === 'with validation error') {
        model.ttl = 123;
      } else if (variant === 'with value 0s') {
        model.ttl = '0s';
      } else if (variant === 'with value 1h') {
        model.ttl = '1h';
      }
    }
    return model;
  }

  @action
  dynamicFormFieldModelValidations(editType, variant) {
    const modelValidations = {};
    if (editType === 'checkboxList') {
      if (variant === 'with validation error') {
        modelValidations.checkboxList = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'dateTimeLocal') {
      if (variant === 'with validation error') {
        modelValidations.dateTimeLocal = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'file') {
      if (variant === 'with validation error') {
        // NOTICE! this generates a double error message, it's a bug in the code (error is already output by the `FormField`, see line 374, but is also output by `TextFile` via the argument `@validationError`)
        modelValidations.file = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'kv') {
      if (variant === 'with validation error') {
        modelValidations.kv = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'stringArray') {
      if (variant === 'with validation error') {
        modelValidations.stringArray = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'select') {
      if (variant === 'with validation error') {
        modelValidations.select = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'textarea') {
      if (variant === 'with validation error') {
        modelValidations.textarea = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    } else if (editType === 'password') {
      if (variant === 'with validation errors and warnings') {
        modelValidations.password = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
          warnings: [
            'This is the validation warning message #1',
            'This is the validation warning message #2',
          ],
        };
      } else if (variant === 'with validation errors') {
        modelValidations.password = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      } else if (variant === 'with validation warnings') {
        modelValidations.password = {
          isValid: true,
          warnings: [
            'This is the validation warning message #1',
            'This is the validation warning message #2',
          ],
        };
      }
    } else if (editType === 'ttl') {
      if (variant === 'with validation error') {
        // NOTICE: there is a bug in the CSS for the class "ttl-picker-form-field-error" because the border color is applied to the `input` child, but such element is hidden so the red border is not visible!
        modelValidations.ttl = {
          isValid: false,
          errors: ['This is the validation error message #1', 'This is the validation error message #2'],
        };
      }
    }
    return modelValidations;
  }

  @action
  dynamicFormFieldOptionsModels(editType) {
    let optionsModels = [];
    if (editType === 'searchSelect') {
      // TODO! find which is the right format for this
      optionsModels = [];
    }
    return optionsModels;
  }

  @action
  dynamicFormFieldOptionsPossibleValues(editType, variant) {
    let possibleValues = [];
    if (editType === 'checkboxList') {
      possibleValues = ['Lorem', 'Ipsum', 'Dolor'];
    } else if (editType === 'radio') {
      possibleValues = [
        {
          label: 'One',
          value: 1,
        },
        {
          label: 'Two',
          value: 2,
        },
        {
          label: 'Three',
          value: 3,
        },
      ];
      if (variant === 'with item subText') {
        possibleValues.forEach((i) => (i.subText = `Subtext for ${i.label.toLowerCase()}`));
      }
      if (variant === 'with item helpText') {
        possibleValues.forEach((i) => (i.helpText = `Helptext for ${i.label.toLowerCase()}`));
      }
    } else if (editType === 'select') {
      possibleValues = ['Lorem', 'Ipsum', 'Dolor'];
    }
    return possibleValues;
  }

  @action
  dynamicFormFieldOptionsFieldValue(editType) {
    let fieldValue;
    if (editType === 'checkboxList') {
      fieldValue = 'checkboxList';
    } else if (editType === 'input') {
      fieldValue = 'input';
    } else if (editType === 'optionalText') {
      fieldValue = 'optionalText';
    } else if (editType === 'radio') {
      fieldValue = 'radio';
    } else if (editType === 'regex') {
      fieldValue = 'regex';
    }
    return fieldValue;
  }

  // ###########################################
  // READONLY-FORMFIELD
  // ###########################################

  @action
  dynamicReadonlyFormFieldValue(attrType, variant) {
    let value;
    // if (attrType === 'select') {
    // } else {
    // }
    if (variant === 'with value' || variant === 'with value + helpText + subText') {
      value = 'Lorem ipsum';
    }
    return value;
  }

  // ###########################################
  // OBJECT-LIST-INPUT
  // ###########################################

  @action
  dynamicObjectListInputObjectKeys(variant) {
    // any variant needs the definition of the object keys
    if (variant) {
      return [
        { label: 'Label for input A', key: 'A', placeholder: 'Placeholder for A' },
        { label: 'Label for input B', key: 'B', placeholder: 'Placeholder for B' },
        { label: 'Label for input C', key: 'C' },
      ];
    }
  }

  @action
  dynamicObjectListInputInputValue(variant) {
    let inputValue = [];
    if (variant === 'with single set of values') {
      inputValue = [{ A: 'First value for A', B: 'First value for B', C: '' }];
    } else if (variant === 'with multiple sets of values' || variant === 'with validation error') {
      inputValue = [
        { A: 'First value for A', B: 'First value for B', C: '' },
        { A: 'Second value for A', B: 'Second value for B', C: '' },
      ];
    }
    return inputValue;
  }

  @action
  dynamicObjectListInputValidationErrors(variant) {
    let validationErrors = [];
    if (variant === 'with validation error') {
      validationErrors = [
        { A: { errors: ['Error message for first A'], isValid: false } },
        { B: { errors: ['Error message for second B'], isValid: false } },
      ];
    }
    return validationErrors;
  }

  // ###########################################
  // OTHER
  // ###########################################

  @action
  noop() {}
}
