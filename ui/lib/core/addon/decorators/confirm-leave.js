/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import Ember from 'ember';

/**
 * Confirm that the user wants to discard unsaved changes before leaving the page. This decorator hooks into
 * the willTransition action. If you override setupController, be sure to set 'model' on the controller to
 * store data or this won't work.
 *
 * By default it will check if the route's model is dirty and prompt when leaving. Usage for this is simple:
 *
 * @withConfirmLeave()
 * export default class MyRoute extends Route {
 *   @service store;
 *   model() {
 *     return this.store.createRecord('some-model')
 *   }
 * }
 *
 * If the route has ember-data models at multiple paths, you can pass an array of secondary modelPaths which
 * will rollback on exit after the prompt for the first model is confirmed. In the example below, the window
 * will only prompt on leave if `model.main` is dirty. Either way, `model.secondary` and `model.optional`
 * will be cleaned up from the data store.
 *
 * @withConfirmLeave('model.main', ['model.secondary', 'model.optional'])
 * export default class MyRoute extends Route {
 *   @service store;
 *   model() {
 *     return {
 *       main: this.store.peekRecord('some-model', 'abc1')
 *       secondary: this.store.createRecord('some-other')
 *       optional: this.store.createRecord('optional')
 *     }
 *   }
 * }
 *
 */
export function withConfirmLeave(modelPath = 'model', silentCleanupPaths) {
  return function decorator(SuperClass) {
    if (!Object.prototype.isPrototypeOf.call(Route, SuperClass)) {
      // eslint-disable-next-line
      console.error(
        'withConfirmLeave decorator must be used on instance of ember Route class. Decorator not applied to returned class'
      );
      return SuperClass;
    }
    return class ConfirmLeave extends SuperClass {
      @service store;

      _rollbackModel(modelPath) {
        const model = this.controller.get(modelPath);
        // we only want to complete rollback if the model is dirty and not saving
        if (model && model.hasDirtyAttributes && !model.isSaving) {
          const method = model.isNew ? 'unloadRecord' : 'rollbackAttributes';
          model[method]();
        }
      }

      @action
      willTransition(transition) {
        try {
          super.willTransition(...arguments);
        } catch (e) {
          // if the SuperClass doesn't have willTransition
          // defined calling it will throw an error.
        }
        const model = this.controller.get(modelPath);

        if (model && model.hasDirtyAttributes && !model.isSaving) {
          if (
            Ember.testing ||
            window.confirm(
              'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
            )
          ) {
            this._rollbackModel(modelPath);
          } else {
            transition.abort();
            return false;
          }
        }
        silentCleanupPaths?.forEach((pathToModel) => {
          this._rollbackModel(pathToModel);
        });
        return true;
      }
    };
  };
}
