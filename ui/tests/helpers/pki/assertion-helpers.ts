/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { GENERAL } from '../general-selectors';

export const EXTERNAL_TABS = ['Overview', 'Roles', 'Recent orders', 'DNS providers', 'ACME accounts'];

export const assertTabState = (assert: Assert, activeTab: string, tabList: string[]) => {
  const inactive = tabList.filter((t) => t !== activeTab);
  inactive.forEach((t) => {
    assert.dom(GENERAL.linkTo(t)).exists().doesNotHaveClass('active', `${t} is inactive`);
  });
  assert.dom(GENERAL.linkTo(activeTab)).exists().hasClass('active', `${activeTab} is active`);
};
