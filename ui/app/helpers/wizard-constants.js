import TutorialMachineConfig from 'vault/machines/tutorial-machine';
import SecretsMachineConfig from 'vault/machines/secrets-machine';
import PoliciesMachineConfig from 'vault/machines/policies-machine';
import ReplicationMachineConfig from 'vault/machines/replication-machine';
import ToolsMachineConfig from 'vault/machines/tools-machine';
import AuthMachineConfig from 'vault/machines/auth-machine';

export const STORAGE_KEYS = {
  TUTORIAL_STATE: 'vault:ui-tutorial-state',
  FEATURE_LIST: 'vault:ui-feature-list',
  FEATURE_STATE: 'vault:ui-feature-state',
  COMPLETED_FEATURES: 'vault:ui-completed-list',
  COMPONENT_STATE: 'vault:ui-component-state',
  RESUME_URL: 'vault:ui-tutorial-resume-url',
  RESUME_ROUTE: 'vault:ui-tutorial-resume-route',
};

export const MACHINES = {
  tutorial: TutorialMachineConfig,
  secrets: SecretsMachineConfig,
  policies: PoliciesMachineConfig,
  replication: ReplicationMachineConfig,
  tools: ToolsMachineConfig,
  authentication: AuthMachineConfig,
};

export const DEFAULTS = {
  currentState: null,
  featureList: null,
  featureState: null,
  currentMachine: null,
  tutorialComponent: null,
  featureComponent: null,
  stepComponent: null,
  detailsComponent: null,
  componentState: null,
  nextFeature: null,
  nextStep: null,
};
