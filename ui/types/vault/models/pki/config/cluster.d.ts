import type { Model } from 'vault/app-types';
import type CapabilitiesModel from 'vault/models/capabilities';

type PkiConfigClusterModel = Model & {
  path: boolean;
  aiaPath: string;
  clusterPath: CapabilitiesModel;
  get canSet(): boolean;
};

export default PkiConfigClusterModel;
