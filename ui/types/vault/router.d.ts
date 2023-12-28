/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';

/* 
Get the resolved type of an item.
https://docs.ember-cli-typescript.com/cookbook/working-with-route-models

- If the item is a promise, the result will be the resolved value type
- If the item is not a promise, the result will just be the type of the item
*/
export type Resolved<P> = P extends Promise<infer T> ? T : P;

/* 
Get the resolved model value from a route. 
Example use:

import type { ModelFrom } from 'vault/vault/router';
export default class MyRoute extends Route {
  redirect(model: ModelFrom<MyRoute>) {}
}
*/
export type ModelFrom<R extends Route> = Resolved<ReturnType<R['model']>>;
