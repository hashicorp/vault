/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';

declare module 'ember-data/types/registries/model' {
  export default interface ModelRegistry {
    // Catchall for any other models
    [key: string]: any;
  }
}
