/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { create } from 'jsondiffpatch';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type Transition from '@ember/routing/transition';
import type RouterService from '@ember/routing/router-service';
import FlagsService from 'vault/services/flags';

// this service tracks the unsaved changes modal state.
export default class UnsavedChangesService extends Service {
  @service declare readonly router: RouterService;
  @service declare readonly flags: FlagsService;

  @tracked showModal = false;

  @tracked initialState: Record<string, unknown> | undefined;
  @tracked currentState: Record<string, unknown> | undefined;
  @tracked intendedTransition: Transition | undefined; // saved transition from willTransition hook before exiting with unsaved changes

  setup(state: Record<string, unknown> | undefined) {
    // ensure unsaved-changes intendedTransition is initially set to undefined each time the user transition
    this.intendedTransition = undefined;
    // set up unsaved-changes service state
    this.currentState = state;
  }

  get changedFields() {
    const diffpatcher = create({});
    const delta = diffpatcher.diff(this.initialState, this.currentState);

    return delta ? Object.keys(delta) : [];
  }

  get hasChanges() {
    return this.changedFields.length > 0;
  }

  get transitionInfo() {
    return {
      routeName: this.intendedTransition?.to?.name,
      params: this.intendedTransition?.to?.params,
    };
  }

  show(transition: Transition) {
    this.intendedTransition = transition;
    this.showModal = true;
  }

  // This method is to update the initial state so it can be called after a successful
  // save or if the user has decided to discard changes
  @action
  resetUnsavedState() {
    this.initialState = this.currentState;
  }

  @action
  transition(route: string) {
    const { routeName: intendedRoute } = this.transitionInfo || {};
    if (intendedRoute) {
      this.resetUnsavedState();
      this.router.transitionTo(intendedRoute);
    } else {
      // due to an issue with the routing model not reloading when a route has query params (ie. is in a namespace)
      // we need to refresh the route to ensure the model is updated
      this.router.refresh(route);
    }
  }
}
