/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import type RouterService from '@ember/routing/router-service';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { HTMLElementEvent } from 'vault/forms';

// TODO: improve typing
interface KnownIssue {
  found: string;
  fixed: string;
  workaround: string;
  edition: string;
  issue: string;
  link: string;
}

interface BreakingChange {
  edition: string;
  recommendations: boolean;
  introduced: string;
  change: string;
  link: string;
}

interface NewBehavior {
  edition: string;
  recommendations: boolean;
  introduced: string;
  change: string;
  link: string;
}

interface VersionInfo {
  version: string;
  known_issues?: KnownIssue[];
  breaking_changes?: BreakingChange[];
  new_behavior?: NewBehavior[];
  rollback_steps?: string[];
}

interface UpgradeInfoArgs {
  breadcrumbs: Array<Breadcrumb>;
  upgradeInfo?: VersionInfo[];
  tab?: string;
}

export default class UpgradeInfoComponent extends Component<UpgradeInfoArgs> {
  @service declare readonly router: RouterService;

  get selectedTabIndex() {
    const index = parseInt(this.args.tab ?? '', 10);
    return !isNaN(index) && index >= 0 && index < this.tabs.length ? index : 0;
  }

  @action
  onClickTab(_event: HTMLElementEvent<HTMLInputElement>, index: number) {
    this.router.replaceWith({ queryParams: { tab: String(index) } });
  }

  private formatKnownIssues(issues: KnownIssue[]) {
    return issues.map((data) => {
      const isFixed = data.fixed !== 'No';
      const hasWorkaround = data.workaround === 'Yes';

      return {
        badges: [
          { text: data.edition, color: 'neutral' },
          {
            text: isFixed ? `Fixed in ${data.fixed}` : 'Not fixed',
            color: isFixed ? 'success' : 'critical',
          },
          {
            text: hasWorkaround ? 'Workaround available' : 'No workaround available',
            color: hasWorkaround ? 'neutral' : 'critical',
          },
        ],
        title: data.issue,
        description: `Found in ${data.found}`,
        link: data.link,
      };
    });
  }

  private formatChangeOrBehavior(items: BreakingChange[] | NewBehavior[]) {
    return items.map((data) => ({
      badges: [
        { text: data.edition, color: 'neutral' },
        {
          text: data.recommendations ? 'Recommendations available' : 'No recommendation',
          color: data.recommendations ? 'highlight' : 'critical',
        },
      ],
      title: data.change,
      description: `Introduced in ${data.introduced}`,
      link: data.link,
    }));
  }

  get panels() {
    // TODO: this is flattening upgradeInfo as it is currently in the format of
    // [{version: '1.21', known_issues: []..}, {version: '1.20', known_issues: []..}]
    // since we're directly pulling from the test data
    // the shape may change once we wire up the endpoint and we could flatten it or format it before it gets here

    const knownIssues = this.upgradeInfo?.flatMap((item) => item.known_issues ?? []) ?? [];
    const breakingChanges = this.upgradeInfo?.flatMap((item) => item.breaking_changes ?? []) ?? [];
    const newBehavior = this.upgradeInfo?.flatMap((item) => item.new_behavior ?? []) ?? [];

    const formattedPanels = [
      this.formatKnownIssues(knownIssues),
      this.formatChangeOrBehavior(breakingChanges),
      this.formatChangeOrBehavior(newBehavior),
    ];

    return formattedPanels;
  }

  get tabs() {
    return [
      { text: 'Known issues', icon: 'shield-alert', count: this.panels[0]?.length },
      { text: 'Breaking changes', icon: 'alert-triangle', count: this.panels[1]?.length },
      { text: 'New behavior', icon: 'alert-circle', count: this.panels[2]?.length },
      { text: 'Rollback steps', icon: 'rewind' },
    ];
  }

  get upgradeInfo() {
    return this.args.upgradeInfo;
  }
}
