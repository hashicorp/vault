/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import type { EntityFlyoutData, EntityFlyoutTab } from 'vault/components/identity/page/entities';

type TabMap = {
  key: EntityFlyoutTab;
  label: string;
};

interface Args {
  data: EntityFlyoutData | null | undefined;
  onClose: () => void;
  initialTab?: EntityFlyoutTab;
}

const ALL_TABS: TabMap[] = [
  { key: 'policies', label: 'Policies' },
  { key: 'entity', label: 'Entity details' },
  { key: 'alias', label: 'Alias details' },
];

export default class IdentityEntitiesFlyoutComponent extends Component<Args> {
  @tracked _userSelectedTab: EntityFlyoutTab | null = null;

  get tabs(): TabMap[] {
    const { data } = this.args;
    // Only show the alias tab when the entity has aliases; policies tab when the entity or its groups reference policies
    return ALL_TABS.filter((tab) => {
      if (tab.key === 'alias') {
        return data?.entity?.aliases?.length;
      }

      if (tab.key === 'policies') {
        // Check the attached policy names rather than the aggregated policy, since deleted policies stay
        // attached (and regain access if recreated) but can't be read to build the aggregate
        return data?.entity?.policies?.length || data?.groups?.some((group) => group.policies?.length);
      }

      return true;
    });
  }

  get selectedTab(): EntityFlyoutTab {
    const preferred = this._userSelectedTab ?? this.args.initialTab;
    const isVisible = preferred && this.tabs.some((t) => t.key === preferred);
    return isVisible ? preferred : this.tabs[0]?.key ?? 'entity';
  }

  get selectedTabIndex() {
    const index = this.tabs.findIndex((tab) => tab.key === this.selectedTab);
    return index < 0 ? 0 : index;
  }

  @action
  onClickTab(_event: Event, index: number): void {
    this._userSelectedTab = this.tabs[index]?.key ?? this.tabs[0]?.key ?? 'entity';
  }
}
