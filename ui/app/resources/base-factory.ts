/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { getOwner, setOwner } from '@ember/owner';

import type Owner from '@ember/owner';

abstract class BaseResource<T> {
  // pass data that the resource should represent (typically from an API response) to constructor
  // object properties will be assigned to class instance
  // extending classes can define getters and additional properties/methods that are required widely across the app
  constructor(data: T, context?: unknown) {
    Object.assign(this, data) as T;
    // pass in context (this) of Ember class (route, component etc.) where the resource is being constructed
    // this will be used to set the owner on the class so that services can be injected (if required)
    if (context) {
      setOwner(this, getOwner(context) as Owner);
    }
  }
}

// factory that allows for the BaseResource class to be cast to the specific type provided
// without this the compiler is not aware of the properties set on the class via Object.assign
// example usage -> export default class SecretsEngineResource extends baseResourceFactory<Mount>() { ... }
export function baseResourceFactory<T>() {
  return BaseResource as new (data: T, context?: unknown) => T;
}
