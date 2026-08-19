/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type { FlyoutData } from '../page';

export type Tab = 'agent' | 'policies' | 'entity' | 'alias';

type TabMap = {
  key: Tab;
  label: string;
};

interface Args {
  data: FlyoutData;
  onClose: () => void;
}

export default class AgentRegistryFlyoutComponent extends Component<Args> {
  @service declare readonly api: ApiService;

  @tracked selectedTab: Tab = this.args.data.aliasId ? 'alias' : 'agent';

  allTabs: TabMap[] = [
    {
      key: 'agent',
      label: 'Agent details',
    },
    { key: 'policies', label: 'Policies' },
    {
      key: 'entity',
      label: 'Entity details',
    },
    { key: 'alias', label: 'Alias details' },
  ];

  get tabs() {
    const { entity, aggregatePolicy } = this.args.data;
    return this.allTabs.filter((tab) => {
      if (tab.key === 'alias') {
        return entity?.aliases?.length;
      }

      if (tab.key === 'policies') {
        return aggregatePolicy.policyString;
      }

      return true;
    });
  }

  get selectedTabIndex() {
    const index = this.tabs.findIndex((tab) => tab.key === this.selectedTab);
    return index < 0 ? 0 : index;
  }

  @action
  onClickTab(_event: Event, index: number) {
    this.selectedTab = this.tabs[index]?.key || this.tabs[0]?.key || 'agent';
  }
}
