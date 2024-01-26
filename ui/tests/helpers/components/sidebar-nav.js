/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { allFeatures } from 'vault/helpers/all-features';
import sinon from 'sinon';

export const stubFeaturesAndPermissions = (owner, isEnterprise = false, setCluster = false) => {
  const permissions = owner.lookup('service:permissions');
  const hasNavPermission = sinon.stub(permissions, 'hasNavPermission');
  hasNavPermission.returns(true);
  sinon.stub(permissions, 'navPathParams');

  const version = owner.lookup('service:version');
  const features = sinon.stub(version, 'features');
  features.value(allFeatures());
  sinon.stub(version, 'isEnterprise').value(isEnterprise);

  const auth = owner.lookup('service:auth');
  sinon.stub(auth, 'authData').value({});

  if (setCluster) {
    owner.lookup('service:currentCluster').setCluster({
      id: 'foo',
      anyReplicationEnabled: true,
      usingRaft: true,
    });
  }

  return { hasNavPermission, features };
};
