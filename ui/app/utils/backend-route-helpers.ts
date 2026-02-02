/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

/**
 * Utility functions for backend-related route operations.
 * Replaces the deprecated backend-helpers mixin.
 */

/**
 * Get the effective engine type for a given route's backend.
 * This handles external plugin mapping to builtin types.
 *
 * @param route - The Ember route instance
 * @returns The effective engine type
 */
export function getBackendEffectiveType(route: Route): string {
  const backendModel = route.modelFor('vault.cluster.secrets.backend') as { engineType: string };
  return getEffectiveEngineType(backendModel?.engineType);
}

/**
 * Get the current backend path parameter from a route.
 *
 * @param route - The Ember route instance
 * @returns The backend path
 */
export function getEnginePathParam(route: Route): string {
  const params = route.paramsFor('vault.cluster.secrets.backend') as { backend: string };
  return params?.backend;
}
