import { allFeatures } from 'vault/helpers/all-features';
import sinon from 'sinon';

export const stubFeaturesAndPermissions = (owner, isEnterprise = false, setCluster = false) => {
  const permissions = owner.lookup('service:permissions');
  sinon.stub(permissions, 'hasNavPermission').returns(true);
  sinon.stub(permissions, 'navPathParams');

  const version = owner.lookup('service:version');
  sinon.stub(version, 'features').value(allFeatures());
  sinon.stub(version, 'isEnterprise').value(isEnterprise);

  if (setCluster) {
    owner.lookup('service:currentCluster').setCluster({
      id: 'foo',
      anyReplicationEnabled: true,
      usingRaft: true,
    });
  }
};
