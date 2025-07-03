/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type SecretsEngineResource from 'vault/resources/secrets/engine';

/**
 * @module DashboardSecretsEnginesCard
 * DashboardSecretsEnginesCard component are used to display 5 secrets engines to the user.
 *
 * @example
 * ```js
 * <DashboardSecretsEnginesCard @secretsEngines={{@model.secretsEngines}} />
 * ```
 * @param {array} secretsEngines - list of secrets engines
 */

interface Args {
  secretsEngines: SecretsEngineResource[];
}

export default class DashboardSecretsEnginesCard extends Component<Args> {
  @tracked favoriteEngines: Array<string> = [];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.loadFavorites();
  }

  loadFavorites() {
    try {
      const stored = localStorage.getItem('vault-favorite-engines');
      if (stored) {
        this.favoriteEngines = JSON.parse(stored);
      }
    } catch (e) {
      this.favoriteEngines = [];
    }
  }

  get filteredSecretsEngines() {
    const engines = this.args.secretsEngines?.filter((secretEngine) => secretEngine.shouldIncludeInList);

    // sort favorites first, then alphabetically
    return engines?.sort((a, b) => {
      const aIsFavorite = this.favoriteEngines.includes(a.id);
      const bIsFavorite = this.favoriteEngines.includes(b.id);

      // sort by favorites first
      if (aIsFavorite && !bIsFavorite) return -1;
      if (!aIsFavorite && bIsFavorite) return 1;

      // else sort alphabetically by path
      return (a.path || '').localeCompare(b.path || '');
    });
  }

  get firstFiveSecretsEngines() {
    return this.filteredSecretsEngines?.slice(0, 5);
  }
}
