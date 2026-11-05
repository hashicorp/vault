/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  id: () => (Math.random() + 1).toString(36).substring(2),
  use_signature: true,
  idp_url: 'https://foobar.pingidentity.com/pingid',
  admin_url: 'https://foobar.pingidentity.com/pingid',
  authenticator_url: 'https://authenticator.pingone.com/pingid/ppm',
  org_alias: 'foobarbaz',
  type: 'pingid',
  username_template: '',
  namespace_path: 'admin/',
});
