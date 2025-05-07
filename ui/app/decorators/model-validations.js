/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { validate } from 'vault/utils/forms/validate';

// see documentation at ui/docs/model-validations.md for detailed usage information
export function withModelValidations(validations) {
  return function decorator(SuperClass) {
    return class ModelValidations extends SuperClass {
      static _validations;

      constructor() {
        super(...arguments);
        if (!validations || typeof validations !== 'object') {
          throw new Error('Validations object must be provided to constructor for setup');
        }
        this._validations = validations;
      }

      validate() {
        return validate(this, this._validations);
      }
    };
  };
}
