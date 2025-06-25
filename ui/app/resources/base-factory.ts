/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

abstract class BaseResource<T> {
  // pass data that the resource should represent (typically from an API response) to constructor
  // object properties will be assigned to class instance
  // extending classes can define getters and additional properties/methods that are required widely across the app
  constructor(readonly data: T) {
    Object.assign(this, data) as T;
  }
}

// factory that allows for the BaseResource class to be casted to the specific type provided
// without this the compiler is not aware of the properties set on the class via Object.assign
// example usage -> export default class SecretsEngineResource extends baseResourceFactory<SecretsEngine>() { ... }
export function baseResourceFactory<T>() {
  return BaseResource as new (data: T) => T;
}
