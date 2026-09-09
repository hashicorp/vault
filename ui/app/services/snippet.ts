/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { CreationMethod } from 'vault/utils/constants/snippet';

export interface SnippetTab {
  key: string;
  label: string;
  snippet: string;
}

export default class SnippetService extends Service {
  @tracked selectedTabIdx = 0;
  @tracked creationMethodChoice: CreationMethod = CreationMethod.TERRAFORM;
  @tracked codeSnippet: string | null = null;

  @action
  setCreationMethod(choice: CreationMethod, tfSnippet: string, customTabs: SnippetTab[]) {
    this.creationMethodChoice = choice;
    this.persistSnippet(tfSnippet, customTabs);
  }

  @action
  setSelectedTab(idx: number, tfSnippet: string, customTabs: SnippetTab[]) {
    this.selectedTabIdx = idx;
    this.persistSnippet(tfSnippet, customTabs);
  }

  @action
  persistSnippet(tfSnippet: string, customTabs: SnippetTab[]) {
    if (this.creationMethodChoice === CreationMethod.TERRAFORM) {
      this.codeSnippet = tfSnippet;
    } else if (this.creationMethodChoice === CreationMethod.APICLI) {
      this.codeSnippet = customTabs[this.selectedTabIdx]?.snippet ?? null;
    } else {
      this.codeSnippet = null;
    }
  }

  reset(initialChoice: CreationMethod = CreationMethod.TERRAFORM) {
    this.selectedTabIdx = 0;
    this.creationMethodChoice = initialChoice;
    this.codeSnippet = null;
  }
}
