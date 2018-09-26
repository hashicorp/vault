// meant for use mixed-in to a Route file
//
// When a route is deactivated, this mixin clears the Ember Data store of
// models of type specified by the required param `modelType`.
//
// example:
// Using this as with a modelType of `datacenter` on the infrastructure
// route will cause all `datacenter` models to get unloaded when the
// infrastructure route is navigated away from.

import Route from '@ember/routing/route';

import { isPresent } from '@ember/utils';
import { warn } from '@ember/debug';
import { on } from '@ember/object/evented';
import Mixin from '@ember/object/mixin';

export default Mixin.create({
  modelType: null,
  modelTypes: null,

  verifyProps: on('init', function() {
    var modelType = this.get('modelType');
    var modelTypes = this.get('modelTypes');
    warn(
      'No `modelType` or `modelTypes` specified for `' +
        this.toString() +
        '`. Check to make sure you still need to use the `model-boundary-route` mixin.',
      isPresent(modelType) || isPresent(modelTypes),
      { id: 'model-boundary-init' }
    );

    warn(
      'Expected `model-boundary-route` to be used on an Ember.Route, not `' + this.toString() + '`.',
      this instanceof Route,
      { id: 'mode-boundary-is-route' }
    );
  }),

  clearModelCache: on('deactivate', function() {
    var modelType = this.get('modelType');
    var modelTypes = this.get('modelTypes');

    if (!modelType && !modelTypes) {
      warn(
        'Attempted to clear store clear store cache when leaving `' +
          this.routeName +
          '`, but no `modelType` or `modelTypes` was specified.',
        isPresent(modelType),
        { id: 'model-boundary-clear' }
      );
      return;
    }
    if (modelType) {
      this.store.unloadAll(modelType);
    }
    if (modelTypes) {
      modelTypes.forEach(type => {
        this.store.unloadAll(type);
      });
    }
  }),
});
