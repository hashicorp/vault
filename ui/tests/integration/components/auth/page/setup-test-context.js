/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

export default (context) => {
  context.version = context.owner.lookup('service:version');
  context.cluster = { id: '1' };
  context.directLinkData = null;
  context.loginSettings = null;
  context.namespaceQueryParam = '';
  context.oidcProviderQueryParam = '';
  context.onAuthSuccess = sinon.spy();
  context.onNamespaceUpdate = sinon.spy();
  context.visibleAuthMounts = false;

  context.renderComponent = () => {
    return render(hbs`<Auth::Page
  @cluster={{this.cluster}}
  @directLinkData={{this.directLinkData}}
  @loginSettings={{this.loginSettings}}
  @namespaceQueryParam={{this.namespaceQueryParam}}
  @oidcProviderQueryParam={{this.oidcProviderQueryParam}}
  @onAuthSuccess={{this.onAuthSuccess}}
  @onNamespaceUpdate={{this.onNamespaceUpdate}}
  @visibleAuthMounts={{this.visibleAuthMounts}}
/>`);
  };
};
