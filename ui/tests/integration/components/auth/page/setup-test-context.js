/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

export default (context) => {
  context.version = context.owner.lookup('service:version');
  context.api = context.owner.lookup('service:api');
  context.cluster = { id: '1' };
  context.directLinkData = null;
  context.loginSettings = null;
  context.namespaceQueryParam = '';
  context.oidcProviderQueryParam = '';
  // mocking as an object with the `perform()` method because loginAndTransition is a concurrency task
  context.loginAndTransition = { perform: sinon.spy() };
  context.onNamespaceUpdate = sinon.spy();
  context.visibleAuthMounts = null;

  context.renderComponent = () => {
    return render(hbs`<Auth::Page
  @cluster={{this.cluster}}
  @directLinkData={{this.directLinkData}}
  @loginSettings={{this.loginSettings}}
  @namespaceQueryParam={{this.namespaceQueryParam}}
  @oidcProviderQueryParam={{this.oidcProviderQueryParam}}
  @loginAndTransition={{this.loginAndTransition}}
  @onNamespaceUpdate={{this.onNamespaceUpdate}}
  @visibleAuthMounts={{this.visibleAuthMounts}}
/>`);
  };
};
