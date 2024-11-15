# Unreleased

# v1.69.0

## Enhancements

* Adds BETA support for a variable set `Parent` relation, which is EXPERIMENTAL, SUBJECT TO CHANGE, and may not be available to all users by @jbonhag [#992](https://github.com/hashicorp/go-tfe/pull/992)
* Add support for adding/updating key/value tags by @brandonc [#991](https://github.com/hashicorp/go-tfe/pull/991)
* Add support for reading a registry module by its unique identifier by @dsa0x [#988](https://github.com/hashicorp/go-tfe/pull/988)
* Add support for enabling Stacks on an organization by @brandonc [#987](https://github.com/hashicorp/go-tfe/pull/987)
* Add support for filtering by key/value tags by @brandonc [#987](https://github.com/hashicorp/go-tfe/pull/987)
* Adds `SpeculativePlanManagementEnabled` field to `Organization` by @lilincmu [#983](https://github.com/hashicorp/go-tfe/pull/983)

# v1.68.0

## Enhancements

* Add support for reading a no-code module's variables by @paladin-devops [#979](https://github.com/hashicorp/go-tfe/pull/979)
* Add Waypoint entitlements (the `waypoint-actions` and `waypoint-templates-and-addons` attributes) to `Entitlements` by @ignatius-j [#984](https://github.com/hashicorp/go-tfe/pull/984)

# v1.67.1

## Bug Fixes

* Fixes a bug in `NewRequest` that did not allow query parameters to be specified in the first parameter, which broke several methods: `RegistryModules ReadVersion`, `VariableSets UpdateWorkspaces`, and `Workspaces Readme` by @brandonc [#982](https://github.com/hashicorp/go-tfe/pull/982)

# v1.67.0

## Enhancements

* `Workspaces`: The `Unlock` method now returns a `ErrWorkspaceLockedStateVersionStillPending` error if the latest state version upload is still pending within the platform. This is a retryable error. by @brandonc [#978](https://github.com/hashicorp/go-tfe/pull/978)

# v1.66.0

## Enhancements

* Adds `billable-rum-count` attribute to `StateVersion` by @shoekstra [#974](https://github.com/hashicorp/go-tfe/pull/974)

## Bug Fixes

* Fixed the incorrect error "workspace already unlocked" being returned when attempting to unlock a workspace that was locked by a Team or different User @ctrombley / @lucasmelin [#975](https://github.com/hashicorp/go-tfe/pull/975)

# v1.65.0

## Enhancements

* Adds support for deleting `Stacks` that still have deployments through `ForceDelete` by @hashimoon [#969](https://github.com/hashicorp/go-tfe/pull/969)

## Bug Fixes

* Fixed `RegistryNoCodeModules` method `UpgradeWorkspace` to return a `WorkspaceUpgrade` type. This resulted in a BREAKING CHANGE, yet the previous type was not properly decoded nor reflective of the actual API result by @paladin-devops [#955](https://github.com/hashicorp/go-tfe/pull/955)

# v1.64.2

## Enhancements

* Adds support for including no-code permissions to the `OrganizationPermissions` struct [#967](https://github.com/hashicorp/go-tfe/pull/967)

# v1.64.1

## Bug Fixes

* Fixes BETA feature regression in `Stacks` associated with decoding `StackVCSRepo` data by @brandonc [#964](https://github.com/hashicorp/go-tfe/pull/964)

# v1.64.0

* Adds support for creating different organization token types by @glennsarti [#943](https://github.com/hashicorp/go-tfe/pull/943)
* Adds more BETA support for `Stacks` resources, which is is EXPERIMENTAL, SUBJECT TO CHANGE, and may not be available to all users by @DanielMSchmidt [#963](https://github.com/hashicorp/go-tfe/pull/963)

# v1.63.0

## Enhancements

* Adds more BETA support for `Stacks` resources, which is is EXPERIMENTAL, SUBJECT TO CHANGE, and may not be available to all users by @brandonc [#957](https://github.com/hashicorp/go-tfe/pull/957) and @DanielMSchmidt [#960](https://github.com/hashicorp/go-tfe/pull/960)

# v1.62.0

## Bug Fixes

* Fixed `RegistryNoCodeModules` methods `CreateWorkspace` and `UpdateWorkspace` to return a `Workspace` type. This resulted in a BREAKING CHANGE, yet the previous type was not properly decoded nor reflective of the actual API result by @paladin-devops [#954](https://github.com/hashicorp/go-tfe/pull/954)

## Enhancements

* Adds `AllowMemberTokenManagement` permission to `Team` by @juliannatetreault [#922](https://github.com/hashicorp/go-tfe/pull/922)

# v1.61.0

## Enhancements

* Adds support for creating no-code workspaces by @paladin-devops [#927](https://github.com/hashicorp/go-tfe/pull/927)
* Adds support for upgrading no-code workspaces by @paladin-devops [#935](https://github.com/hashicorp/go-tfe/pull/935)

# v1.60.0

## Enhancements

* Adds more BETA support for `Stacks` resources, which is EXPERIMENTAL, SUBJECT TO CHANGE, and may not be available to all users by @brandonc. [#934](https://github.com/hashicorp/go-tfe/pull/934)

# v1.59.0

## Features

* Adds support for the Run Tasks Integration API by @karvounis-form3 [#929](https://github.com/hashicorp/go-tfe/pull/929)

# v1.58.0

## Enhancements

* Adds BETA support for `Stacks` resources, which is EXPERIMENTAL, SUBJECT TO CHANGE, and may not be available to all users by @brandonc. [#920](https://github.com/hashicorp/go-tfe/pull/920)

# v1.57.0

## Enhancements

* Adds the `IsUnified` field to `Project`, `Organization` and `Team` by @roncodingenthusiast [#915](https://github.com/hashicorp/go-tfe/pull/915)
* Adds Workspace auto-destroy notification types to `NotificationTriggerType` by @notchairmk [#918](https://github.com/hashicorp/go-tfe/pull/918)
* Adds `CreatedAfter` and `CreatedBefore` Date Time filters to `AdminRunsListOptions` by @maed223 [#916](https://github.com/hashicorp/go-tfe/pull/916)

# v1.56.0

## Enhancements
* Adds `ManageAgentPools` permission to team `OrganizationAccess` by @emlanctot [#901](https://github.com/hashicorp/go-tfe/pull/901)

# v1.55.0

## Enhancements
* Adds the `CurrentRunStatus` filter to allow filtering workspaces by their current run status by @arybolovlev [#899](https://github.com/hashicorp/go-tfe/pull/899)

# v1.54.0

## Enhancements
* Adds the `AutoDestroyActivityDuration` field to `Workspace` by @notchairmk [#902](https://github.com/hashicorp/go-tfe/pull/902)

## Deprecations
* The `IsSiteAdmin` field on User has been deprecated. Use the `IsAdmin` field instead [#900](https://github.com/hashicorp/go-tfe/pull/900)

# v1.53.0

## Enhancements
* Adds `ManageTeams`, `ManageOrganizationAccess`, and `AccessSecretTeams` permissions to team `OrganizationAccess` by @juliannatetreault [#874](https://github.com/hashicorp/go-tfe/pull/874)
* Mocks are now generated using the go.uber.org/mock package [#897](https://github.com/hashicorp/go-tfe/pull/897)

# v1.52.0

## Enhancements
* Add `EnforcementLevel` to `Policy` create and update options. This will replace the deprecated `[]Enforce` method for specifying enforcement level. @JarrettSpiker [#895](https://github.com/hashicorp/go-tfe/pull/895)

## Deprecations
* The `Enforce` fields on `Policy`, `PolicyCreateOptions`, and `PolicyUpdateOptions` have been deprecated. Use the `EnforcementLevel` instead. @JarrettSpiker [#895](https://github.com/hashicorp/go-tfe/pull/895)

# v1.51.0

## Enhancements
* Adds `Teams` field to `OrganizationMembershipCreateOptions` to allow users to be added to teams at the same time they are invited to an organization. by @JarrettSpiker [#886](https://github.com/hashicorp/go-tfe/pull/886)
* `IsCloud()` returns true when TFP-AppName is "HCP Terraform" by @sebasslash [#891](https://github.com/hashicorp/go-tfe/pull/891)
* `OrganizationScoped` attribute for `OAuthClient` is now generally available by @netramali [#873](https://github.com/hashicorp/go-tfe/pull/873)

# v1.50.0

## Enhancements
* Adds Bitbucket Data Center as a new `ServiceProviderType` and ensures similar validation as Bitbucket Server by @zainq11 [#879](https://github.com/hashicorp/go-tfe/pull/879)
* Add `GlobalRunTasks` field to `Entitlements`. by @glennsarti [#865](https://github.com/hashicorp/go-tfe/pull/865)
* Add `Global` field to `RunTask`. by @glennsarti [#865](https://github.com/hashicorp/go-tfe/pull/865)
* Add `Stages` field to `WorkspaceRunTask`. by @glennsarti [#865](https://github.com/hashicorp/go-tfe/pull/865)
* Changing BETA `OrganizationScoped` attribute of `OAuthClient` to be a pointer for bug fix by @netramali [884](https://github.com/hashicorp/go-tfe/pull/884)
* Adds `Query` parameter to `VariableSetListOptions` to allow searching variable sets by name, by @JarrettSpiker[#877](https://github.com/hashicorp/go-tfe/pull/877)

## Deprecations
* The `Stage` field has been deprecated on `WorkspaceRunTask`. Instead, use `Stages`. by @glennsarti [#865](https://github.com/hashicorp/go-tfe/pull/865)

# v1.49.0

## Enhancements
* Adds `post_apply` to list of possible `stages` for Run Tasks by @glennsarti [#878](https://github.com/hashicorp/go-tfe/pull/878)

# v1.48.0

## Features
* For Terraform Enterprise users who have data retention policies defined on Organizations or Workspaces: A new DataRetentionPolicyChoice relation has been added to reflect that [data retention policies are polymorphic](https://developer.hashicorp.com/terraform/enterprise/api-docs/data-retention-policies#data-retention-policy-types). Organizations and workspaces may be related to a `DataRetentionPolicyDeleteOlder` or `DataRetentionPolicyDontDelete` record through the `DataRetentionPolicyChoice` struct. Data retention policies can be read using `ReadDataRetentionPolicyChoice`, and set or updated (including changing their type) using `SetDataRetentionPolicyDeleteOlder` or `SetDataRetentionPolicyDontDelete` by  @JarrettSpiker [#652](https://github.com/hashicorp/go-tfe/pull/844)

## Deprecations
* The `DataRetentionPolicy` type, and the `DataRetentionPolicy` relationship on `Organization` and `Workspace`s have been deprecated. The `DataRetentionPolicy` type is equivalent to the new `DataRetentionPolicyDeleteOlder`. The Data retention policy relationships on `Organization` and `Workspace`s are now [polymorphic](https://developer.hashicorp.com/terraform/enterprise/api-docs/data-retention-policies#data-retention-policy-types), and are represented by the `DataRetentionPolicyChoice` relationship. The existing `DataRetentionPolicy` relationship will continue to be populated when reading an `Organization` or `Workspace`, but it may be removed in a future release. @JarrettSpiker [#652](https://github.com/hashicorp/go-tfe/pull/844)
* The `SetDataRetentionPolicy` function on `Organizations` and `Workspaces` is now deprecated in favour of `SetDataRetentionPolicyDeleteOlder` or `SetDataRetentionPolicyDontDelete`. `SetDataRetentionPolicy` will only update the data retention policy when communicating with TFE versions v202311 and v202312. @JarrettSpiker [#652](https://github.com/hashicorp/go-tfe/pull/844)
* The `ReadDataRetentionPolicy` function on `Organizations` and `Workspaces` is now deprecated in favour of `ReadDataRetentionPolicyChoice`. `ReadDataRetentionPolicyChoice` may return the different multiple data retention policy types added in TFE 202401-1. `SetDataRetentionPolicy` will only update the data retention policy when communicating with TFE versions v202311 and v202312. @JarrettSpiker [#652](https://github.com/hashicorp/go-tfe/pull/844)

## Enhancements
* Adds `Variables` relationship field to `Workspace` by @arybolovlev [#872](https://github.com/hashicorp/go-tfe/pull/872)

# v1.47.1

## Bug fixes
* Change the error message for `ErrWorkspaceStillProcessing` to be the same error message returned by the API by @uturunku1 [#864](https://github.com/hashicorp/go-tfe/pull/864)

# v1.47.0

## Enhancements
* Adds BETA `description` attribute to `Project` by @netramali [#861](https://github.com/hashicorp/go-tfe/pull/861)
* Adds `Read` method to `TestVariables` by @aaabdelgany [#851](https://github.com/hashicorp/go-tfe/pull/851)

# v1.46.0

## Enhancements
* Adds `Query` field to `Project` and `Team` list options, to allow projects and teams to be searched by name by @JarrettSpiker [#849](https://github.com/hashicorp/go-tfe/pull/849)
* Adds `AgenPool` relation to `OAuthClient` create options to support for Private VCS by enabling creation of OAuth Client when AgentPoolID is set (as an optional param) @roleesinhaHC [#841](https://github.com/hashicorp/go-tfe/pull/841)
* Add `Sort` field to workspace list options @Maed223 [#859](https://github.com/hashicorp/go-tfe/pull/859)

# v1.45.0

## Enhancements
* Updates go-tfe client to export the instance name using `AppName()` @sebasslash [#848](https://github.com/hashicorp/go-tfe/pull/848)
* Add `DeleteByName` API endpoint to `RegistryModule` @laurenolivia [#847](https://github.com/hashicorp/go-tfe/pull/847)
* Update deprecated `RegistryModule` endpoints `DeleteProvider` and `DeleteVersion` with new API calls @laurenolivia [#847](https://github.com/hashicorp/go-tfe/pull/847)

# v1.44.0

## Enhancements
* Updates `Workspaces` to include an `AutoDestroyAt` attribute on create and update by @notchairmk and @ctrombley [#786](https://github.com/hashicorp/go-tfe/pull/786)
* Adds `AgentsEnabled` and `PolicyToolVersion` attributes to `PolicySet` by @mrinalirao [#752](https://github.com/hashicorp/go-tfe/pull/752)

# v1.43.0

## Features
* Adds `AggregatedCommitStatusEnabled` field to `Organization` by @mjyocca [#829](https://github.com/hashicorp/go-tfe/pull/829)

## Enhancements
* Adds `GlobalProviderSharing` field to `AdminOrganization` by @alex-ikse [#837](https://github.com/hashicorp/go-tfe/pull/837)

# v1.42.0

## Deprecations
* The `Sourceable` field has been deprecated on `RunTrigger`. Instead, use `SourceableChoice` to locate the non-empty field representing the actual sourceable value by @brandonc [#816](https://github.com/hashicorp/go-tfe/pull/816)

## Features
* Added `AdminOPAVersion` and `AdminSentinelVersion` Terraform Enterprise admin endpoints by @mrinalirao [#758](https://github.com/hashicorp/go-tfe/pull/758)

## Enhancements
* Adds `LockedBy` relationship field to `Workspace` by @brandonc [#816](https://github.com/hashicorp/go-tfe/pull/816)
* Adds `CreatedBy` relationship field to `TeamToken`, `UserToken`, and `OrganizationToken` by @brandonc [#816](https://github.com/hashicorp/go-tfe/pull/816)
* Added `Sentinel` field to `PolicyResult` by @stefan-kiss. [Issue#790](https://github.com/hashicorp/go-tfe/issues/790)

# v1.41.0

## Enhancements

* Allow managing workspace and organization data retention policies by @mwudka [#801](https://github.com/hashicorp/go-tfe/pull/817)

# v1.40.0

## Bug Fixes
* Removed unused field `AgentPoolID` from the Workspace model. (Callers should be using the `AgentPool` relation instead) by @brandonc [#815](https://github.com/hashicorp/go-tfe/pull/815)

## Enhancements
* Add organization scope field for oauth clients by @Netra2104 [#812](https://github.com/hashicorp/go-tfe/pull/812)
* Added BETA support for including `projects` relationship to oauth_client on create by @Netra2104 [#806](https://github.com/hashicorp/go-tfe/pull/806)
* Added BETA method `AddProjects` and `RemoveProjects` for attaching/detaching oauth_client to projects by Netra2104 [#806](https://github.com/hashicorp/go-tfe/pull/806)
* Adds a missing interface `WorkspaceResources` and the `List` method by @stefan-kiss [Issue#754](https://github.com/hashicorp/go-tfe/issues/754)

# v1.39.2

## Bug Fixes
* Fixes a dependency build failure for 32 bit linux architectures by @brandonc [#814](https://github.com/hashicorp/go-tfe/pull/814)

# v1.39.1

## Bug Fixes
* Fixes an issue where the request body is not preserved during certain retry scenarios by @sebasslash [#813](https://github.com/hashicorp/go-tfe/pull/813)

# v1.39.0

## Features
* New WorkspaceSettingOverwritesOptions field for allowing workspaces to defer some settings to a default from their organization or project by @SwiftEngineer [#762](https://github.com/hashicorp/go-tfe/pull/762)
* Added support for setting a default execution mode and agent pool at the organization level by @SwiftEngineer [#762](https://github.com/hashicorp/go-tfe/pull/762)
* Added validation when configuring registry module publishing by @hashimoon [#804](https://github.com/hashicorp/go-tfe/pull/804)
* Removed BETA labels for StateVersion Upload method, ConfigurationVersion `provisional` field, and `save-plan` runs by @brandonc [#800](https://github.com/hashicorp/go-tfe/pull/800)
* Allow soft deleting, restoring, and permanently deleting StateVersion and ConfigurationVersion backing data by @mwudka [#801](https://github.com/hashicorp/go-tfe/pull/801)
* Added the `AutoApplyRunTrigger` attribute to Workspaces by @nfagerlund [#798](https://github.com/hashicorp/go-tfe/pull/798)
* Removed BETA labels for `priority` attribute in variable sets by @Netra2104 [#796](https://github.com/hashicorp/go-tfe/pull/796)

# v1.38.0

## Features
* Added BETA support for including `priority` attribute to variable_set on create and update by @Netra2104 [#778](https://github.com/hashicorp/go-tfe/pull/778)

# v1.37.0

## Features
* Add the tags attribute to VCSRepo to be used with registry modules by @hashimoon [#793](https://github.com/hashicorp/go-tfe/pull/793)

# v1.36.0

## Features
* Added BETA support for private module registry test variables by @aaabdelgany [#787](https://github.com/hashicorp/go-tfe/pull/787)

## Bug Fixes
* Fix incorrect attribute type for `RegistryModule.VCSRepo.Tags` by @hashimoon [#789](https://github.com/hashicorp/go-tfe/pull/789)
* Fix nil dereference panic within `StateVersions` `upload` after not handling certain state version create errors by @brandonc [#792](https://github.com/hashicorp/go-tfe/pull/792)

# v1.35.0

## Features
* Added BETA support for private module registry tests by @hashimoon [#781](https://github.com/hashicorp/go-tfe/pull/781)

## Enhancements
* Removed beta flags for `PolicySetProjects` and `PolicySetWorkspaceExclusions` by @Netra2104 [#770](https://github.com/hashicorp/go-tfe/pull/770)

# v1.34.0

## Features
* Added support for the new Terraform Test Runs API by @liamcervante [#755](https://github.com/hashicorp/go-tfe/pull/755)

## Bug Fixes
* "project" was being rejected as an invalid `Include` option when listing workspaces by @brandonc [#765](https://github.com/hashicorp/go-tfe/pull/765)


# v1.33.0

## Enhancements
* Removed beta tags for TeamProjectAccess by @rberecka [#756](https://github.com/hashicorp/go-tfe/pull/756)
* Added BETA support for including `workspaceExclusions` relationship to policy_set on create by @Netra2104 [#757](https://github.com/hashicorp/go-tfe/pull/757)
* Added BETA method `AddWorkspaceExclusions` and `RemoveWorkspaceExclusions` for attaching/detaching workspace-exclusions to a policy-set by @hs26gill [#761](https://github.com/hashicorp/go-tfe/pull/761)

# v1.32.1

## Dependency Update
* Updated go-slug dependency to v0.12.1

# v1.32.0

## Enhancements
* Added BETA support for adding and updating custom permissions to `TeamProjectAccesses`. A `TeamProjectAccessType` of `"custom"` can set various permissions applied at
the project level to the project itself (`TeamProjectAccessProjectPermissionsOptions`) and all of the workspaces in a project (`TeamProjectAccessWorkspacePermissionsOptions`) by @rberecka [#745](https://github.com/hashicorp/go-tfe/pull/745)
* Added BETA field `Provisional` to `ConfigurationVersions` by @brandonc [#746](https://github.com/hashicorp/go-tfe/pull/746)


# v1.31.0

## Enhancements
* Added BETA support for including `projects` relationship and `projects-count` attribute to policy_set on create by @hs26gill [#737](https://github.com/hashicorp/go-tfe/pull/737)
* Added BETA method `AddProjects` and `RemoveProjects` for attaching/detaching policy set to projects by @Netra2104 [#735](https://github.com/hashicorp/go-tfe/pull/735)

# v1.30.0

## Enhancements
* Adds `SignatureSigningMethod` and `SignatureDigestMethod` fields in `AdminSAMLSetting` struct by @karvounis-form3 [#731](https://github.com/hashicorp/go-tfe/pull/731)
* Adds `Certificate`, `PrivateKey`, `TeamManagementEnabled`, `AuthnRequestsSigned`, `WantAssertionsSigned`, `SignatureSigningMethod`, `SignatureDigestMethod` fields in `AdminSAMLSettingsUpdateOptions` struct by @karvounis-form3 [#731](https://github.com/hashicorp/go-tfe/pull/731)

# v1.29.0

## Enhancements
* Adds `RunPreApplyCompleted` run status by @uk1288 [#727](https://github.com/hashicorp/go-tfe/pull/727)
* Added BETA support for saved plan runs, by @nfagerlund [#724](https://github.com/hashicorp/go-tfe/pull/724)
    * New `SavePlan` fields in `Run` and `RunCreateOptions`
    * New `RunPlannedAndSaved` `RunStatus` value
    * New `PlannedAndSavedAt` field in `RunStatusTimestamps`
    * New `RunOperationSavePlan` constant for run list filters

# v1.28.0

## Enhancements
* Update `Workspaces` to include associated `project` resource by @glennsarti [#714](https://github.com/hashicorp/go-tfe/pull/714)
* Adds BETA method `Upload` method to `StateVersions` and support for pending state versions by @brandonc [#717](https://github.com/hashicorp/go-tfe/pull/717)
* Adds support for the query parameter `q` to search `Organization Tags` by name by @sharathrnair87 [#720](https://github.com/hashicorp/go-tfe/pull/720)
* Added ContextWithResponseHeaderHook support to `IPRanges` by @brandonc [#717](https://github.com/hashicorp/go-tfe/pull/717)

## Bug Fixes
* `ConfigurationVersions`, `PolicySetVersions`, and `RegistryModules` `Upload` methods were sending API credentials to the specified upload URL, which was unnecessary by @brandonc [#717](https://github.com/hashicorp/go-tfe/pull/717)

# v1.27.0

## Enhancements
* Adds `RunPreApplyRunning` and `RunQueuingApply` run statuses by @uk1288 [#712](https://github.com/hashicorp/go-tfe/pull/712)

## Bug Fixes
* AgentPool `Update` is not able to remove all allowed workspaces from an agent pool. That operation is now handled by a separate `UpdateAllowedWorkspaces` method using `AgentPoolAllowedWorkspacesUpdateOptions` by @hs26gill [#701](https://github.com/hashicorp/go-tfe/pull/701)

# v1.26.0

## Enhancements

* Adds BETA fields `ResourceImports` count to both `Plan` and `Apply` types as well as `AllowConfigGeneration` to the `Run` struct type. These fields are not generally available and are subject to change in a future release.

# v1.25.1

## Bug Fixes
* Workspace safe delete conflict error when workspace is locked has been restored
to the original message using the error `ErrWorkspaceLockedCannotDelete` instead of
`ErrWorkspaceLocked`

# v1.25.0

## Enhancements
* Workspace safe delete 409 conflict errors associated with resources still being managed or being processed (indicating that you should try again later) are now the named errors  `ErrWorkspaceStillProcessing` and `ErrWorkspaceNotSafeToDelete` by @brandonc [#703](https://github.com/hashicorp/go-tfe/pull/703)

# v1.24.0

## Enhancements
* Adds support for a new variable field `version-id` by @arybolovlev [#697](https://github.com/hashicorp/go-tfe/pull/697)
* Adds `ExpiredAt` field to `OrganizationToken`, `TeamToken`, and `UserToken`. This enhancement will be available in TFE release, v202305-1. @JuliannaTetreault [#672](https://github.com/hashicorp/go-tfe/pull/672)
* Adds `ContextWithResponseHeaderHook` context for use with the ClientRequest Do method that allows callers to define a callback which receives raw http Response headers.  @apparentlymart [#689](https://github.com/hashicorp/go-tfe/pull/689)


# v1.23.0

## Features
* `ApplyToProjects` and `RemoveFromProjects` to `VariableSets` endpoints now generally available.
* `ListForProject` to `VariableSets` endpoints now generally available.

## Enhancements
* Adds `OrganizationScoped` and `AllowedWorkspaces` fields for creating workspace scoped agent pools and adds `AllowedWorkspacesName` for filtering agents pools associated with a given workspace by @hs26gill [#682](https://github.com/hashicorp/go-tfe/pull/682/files)

## Bug Fixes


# v1.22.0

## Beta API Changes
* The beta `no_code` field in `RegistryModuleCreateOptions` has been changed from `bool` to `*bool` and will be removed in a future version because a new, preferred method for managing no-code registry modules has been added in this release.

## Features
* Add beta endpoints `Create`, `Read`, `Update`, and `Delete` to manage no-code provisioning for a `RegistryModule`. This allows users to enable no-code provisioning for a registry module, and to configure the provisioning settings for that module version. This also allows users to disable no-code provisioning for a module version. @dsa0x [#669](https://github.com/hashicorp/go-tfe/pull/669)

# v1.21.0

## Features
* Add beta endpoints `ApplyToProjects`  and `RemoveFromProjects` to `VariableSets`.  Applying a variable set to a project will apply that variable set to all current and future workspaces in that project.
* Add beta endpoint `ListForProject` to `VariableSets` to list all variable sets applied to a project.
* Add endpoint `RunEvents` which lists events for a specific run by @glennsarti [#680](https://github.com/hashicorp/go-tfe/pull/680)

## Bug Fixes
* `VariableSets.Read` did not honor the Include values due to a syntax error in the struct tag of `VariableSetReadOptions` by @sgap [#678](https://github.com/hashicorp/go-tfe/pull/678)

## Enhancements
* Adds `ProjectID` filter to allow filtering of workspaces of a given project in an organization by @hs26gill [#671](https://github.com/hashicorp/go-tfe/pull/671)
* Adds `Name` filter to allow filtering of projects by @hs26gill [#668](https://github.com/hashicorp/go-tfe/pull/668/files)
* Adds `ManageMembership` permission to team `OrganizationAccess` by @JarrettSpiker [#652](https://github.com/hashicorp/go-tfe/pull/652)
* Adds `RotateKey` and `TrimKey` Admin endpoints by @mpminardi [#666](https://github.com/hashicorp/go-tfe/pull/666)
* Adds `Permissions` to `User` by @jeevanragula [#674](https://github.com/hashicorp/go-tfe/pull/674)
* Adds `IsEnterprise` and `IsCloud` boolean methods to the client by @sebasslash [#675](https://github.com/hashicorp/go-tfe/pull/675)

# v1.20.0

## Enhancements
* Update team project access to include additional project roles by @joekarl [#642](https://github.com/hashicorp/go-tfe/pull/642)

# v1.19.0

## Enhancements
* Removed Beta tags from `Project` features by @hs26gill [#637](https://github.com/hashicorp/go-tfe/pull/637)
* Add `Filter` and `Sort` fields to `AdminWorkspaceListOptions` to allow filtering and sorting of workspaces by @laurenolivia [#641](https://github.com/hashicorp/go-tfe/pull/641)
* Add support for `List` and `Read` Github app installation APIs by @roleesinhaHC [#655](https://github.com/hashicorp/go-tfe/pull/655)
* Add `GHAInstallationID` field to `VCSRepoOptions` and `VCSRepo` structs by @roleesinhaHC [#655](https://github.com/hashicorp/go-tfe/pull/655)

# v1.18.0

## Enhancements
* Adds `BaseURL` and `BaseRegistryURL` methods to `Client` to expose its configuration by @brandonc [#638](https://github.com/hashicorp/go-tfe/pull/638)
* Adds `ReadWorkspaces` and `ReadProjects` permissions to `Organizations` by @JuliannaTetreault [#614](https://github.com/hashicorp/go-tfe/pull/614)

# v1.17.0

## Enhancements
* Add Beta endpoint `TeamProjectAccesses` to manage Project Access for Teams by @hs26gill [#599](https://github.com/hashicorp/go-tfe/pull/599)
* Updates api doc links from terraform.io to developer.hashicorp domain by @uk1288 [#629](https://github.com/hashicorp/go-tfe/pull/629)
* Adds `UploadTarGzip()` method to `RegistryModules` and `ConfigurationVersions` interface by @sebasslash [#623](https://github.com/hashicorp/go-tfe/pull/623)
* Adds `ManageProjects` field to `OrganizationAccess` struct by @hs26gill [#633](https://github.com/hashicorp/go-tfe/pull/633)
* Adds agent-count to `AgentPools` endpoint. @evilensky [#611](https://github.com/hashicorp/go-tfe/pull/611)
* Adds `Links` to `Workspace`, (currently contains "self" and "self-html" paths) @brandonc [#622](https://github.com/hashicorp/go-tfe/pull/622)

# v1.16.0

## Bug Fixes

* Project names were being incorrectly validated as ID's @brandonc [#608](https://github.com/hashicorp/go-tfe/pull/608)

## Enhancements
* Adds `List()` method to `GPGKeys` interface by @sebasslash [#602](https://github.com/hashicorp/go-tfe/pull/602)
* Adds `ProviderBinaryUploaded` field to `RegistryPlatforms` struct by @sebasslash [#602](https://github.com/hashicorp/go-tfe/pull/602)

# v1.15.0

## Enhancements

* Add Beta `Projects` endpoint. The API is in not yet available to all users @hs26gill [#564](https://github.com/hashicorp/go-tfe/pull/564)

# v1.14.0

## Enhancements

* Adds Beta parameter `Overridable` for OPA `policy set` update API (`PolicySetUpdateOptions`) @mrinalirao [#594](https://github.com/hashicorp/go-tfe/pull/594)
* Adds new task stage status values representing `canceled`, `errored`, `unreachable` @mrinalirao [#594](https://github.com/hashicorp/go-tfe/pull/594)

# v1.13.0

## Bug Fixes

* Fixes `AuditTrail` pagination parameters (`CurrentPage`, `PreviousPage`, `NextPage`, `TotalPages`, `TotalCount`), which were not deserialized after reading from the List endpoint by @brandonc [#586](https://github.com/hashicorp/go-tfe/pull/586)

## Enhancements

* Add OPA support to the Policy Set APIs by @mrinalirao [#575](https://github.com/hashicorp/go-tfe/pull/575)
* Add OPA support to the Policy APIs by @mrinalirao [#579](https://github.com/hashicorp/go-tfe/pull/579)
* Add support for enabling no-code provisioning in an existing or new `RegistryModule` by @miguelhrocha [#562](https://github.com/hashicorp/go-tfe/pull/562)
* Add Policy Evaluation and Policy Set Outcome APIs by @mrinalirao [#583](https://github.com/hashicorp/go-tfe/pull/583)
* Add OPA support to Task Stage APIs by @mrinalirao [#584](https://github.com/hashicorp/go-tfe/pull/584)

# v1.12.0

## Enhancements

* Add `search[wildcard-name]` to `WorkspaceListOptions` by @laurenolivia [#569](https://github.com/hashicorp/go-tfe/pull/569)
* Add `NotificationTriggerAssessmentCheckFailed` notification trigger type by @rexredinger [#549](https://github.com/hashicorp/go-tfe/pull/549)
* Add `RemoteTFEVersion()` to the `Client` interface, which exposes the `X-TFE-Version` header set by a remote TFE instance by @sebasslash [#563](https://github.com/hashicorp/go-tfe/pull/563)
* Validate the module version as a version instead of an ID [#409](https://github.com/hashicorp/go-tfe/pull/409)
* Add `AllowForceDeleteWorkspaces` setting to `Organizations` by @JarrettSpiker [#539](https://github.com/hashicorp/go-tfe/pull/539)
* Add `SafeDelete` and `SafeDeleteID` APIs to `Workspaces` by @JarrettSpiker [#539](https://github.com/hashicorp/go-tfe/pull/539)
* Add `ForceExecute()` to `Runs` to allow force executing a run by @annawinkler [#570](https://github.com/hashicorp/go-tfe/pull/570)
* Pre-plan and Pre-Apply Run Tasks are now generally available (beta comments removed) by @glennsarti [#555](https://github.com/hashicorp/go-tfe/pull/555)

# v1.11.0

## Enhancements

* Add `Query` and `Status` fields to `OrganizationMembershipListOptions` to allow filtering memberships by status or username by @sebasslash [#550](https://github.com/hashicorp/go-tfe/pull/550)
* Add `ListForWorkspace` method to `VariableSets` interface to enable fetching variable sets associated with a workspace by @tstapler [#552](https://github.com/hashicorp/go-tfe/pull/552)
* Add `NotificationTriggerAssessmentDrifted` and `NotificationTriggerAssessmentFailed` notification trigger types by @lawliet89 [#542](https://github.com/hashicorp/go-tfe/pull/542)

## Bug Fixes
* Fix marshalling of run variables in `RunCreateOptions`. The `Variables` field type in `Run` struct has changed from `[]*RunVariable` to `[]*RunVariableAttr` by @Uk1288 [#531](https://github.com/hashicorp/go-tfe/pull/531)

# v1.10.0

## Enhancements

* Add `Query` param field to `OrganizationListOptions` to allow searching based on name or email by @laurenolivia [#529](https://github.com/hashicorp/go-tfe/pull/529)
* Add optional `AssessmentsEnforced` to organizations and `AssessmentsEnabled` to workspaces for managing the workspace and organization health assessment (drift detection) setting by @rexredinger [#462](https://github.com/hashicorp/go-tfe/pull/462)

## Bug Fixes
* Fixes null value returned in variable set relationship in `VariableSetVariable` by @sebasslash [#521](https://github.com/hashicorp/go-tfe/pull/521)

# v1.9.0

## Enhancements
* `RunListOptions` is generally available, and rename field (Name -> User) by @mjyocca [#472](https://github.com/hashicorp/go-tfe/pull/472)
* [Beta] Adds optional `JsonState` field to `StateVersionCreateOptions` by @megan07 [#514](https://github.com/hashicorp/go-tfe/pull/514)

## Bug Fixes
* Fixed invalid memory address error when using `TaskResults` field by @glennsarti [#517](https://github.com/hashicorp/go-tfe/pull/517)

# v1.8.0

## Enhancements

* Adds support for reading and listing Agents by @laurenolivia [#456](https://github.com/hashicorp/go-tfe/pull/456)
* It was previously logged that we added an `Include` param field to `PolicySetListOptions` to allow policy list to include related resource data such as workspaces, policies, newest_version, or current_version by @Uk1288 [#497](https://github.com/hashicorp/go-tfe/pull/497) in 1.7.0, but this was a mistake and the field is added in v1.8.0

# v1.7.0

## Enhancements

* Adds new run creation attributes: `allow-empty-apply`, `terraform-version`, `plan-only` by @sebasslash [#482](https://github.com/hashicorp/go-tfe/pull/447)
* Adds additional Task Stage and Run Statuses for Pre-plan run tasks by @glennsarti [#469](https://github.com/hashicorp/go-tfe/pull/469)
* Adds `stage` field to the create and update methods for Workspace Run Tasks by @glennsarti [#469](https://github.com/hashicorp/go-tfe/pull/469)
* Adds `ResourcesProcessed`, `StateVersion`, `TerraformVersion`, `Modules`, `Providers`, and `Resources` fields to the State Version struct by @laurenolivia [#484](https://github.com/hashicorp/go-tfe/pull/484)
* Add `Include` param field to `PolicySetListOptions` to allow policy list to include related resource data such as workspaces, policies, newest_version, or current_version by @Uk1288 [#497](https://github.com/hashicorp/go-tfe/pull/497)
* Allow FileTriggersEnabled to be set to false when Git tags are present by @mjyocca @hashimoon  [#468] (https://github.com/hashicorp/go-tfe/pull/468)

# v1.6.0

## Enhancements
* Remove beta messaging for Run Tasks by @glennsarti [#447](https://github.com/hashicorp/go-tfe/pull/447)
* Adds `Description` field to the `RunTask` object by @glennsarti [#447](https://github.com/hashicorp/go-tfe/pull/447)
* Add `Name` field to `OAuthClient` by @barrettclark [#466](https://github.com/hashicorp/go-tfe/pull/466)
* Add support for creating both public and private `RegistryModule` with no VCS connection by @Uk1288 [#460](https://github.com/hashicorp/go-tfe/pull/460)
* Add `ConfigurationSourceAdo` configuration source option by @mjyocca [#467](https://github.com/hashicorp/go-tfe/pull/467)
* [beta] state version outputs may now include a detailed-type attribute in a future API release by @brandonc [#479](https://github.com/hashicorp/go-tfe/pull/429)

# v1.5.0

## Enhancements
* [beta] Add support for triggering Workspace runs through matching Git tags [#434](https://github.com/hashicorp/go-tfe/pull/434)
* Add `Query` param field to `AgentPoolListOptions` to allow searching based on agent pool name, by @JarrettSpiker [#417](https://github.com/hashicorp/go-tfe/pull/417)
* Add organization scope and allowed workspaces field for scope agents by @Netra2104 [#453](https://github.com/hashicorp/go-tfe/pull/453)
* Adds `Namespace` and `RegistryName` fields to `RegistryModuleID` to allow reading of Public Registry Modules by @Uk1288 [#464](https://github.com/hashicorp/go-tfe/pull/464)

## Bug fixes
* Fixed JSON mapping for Configuration Versions failing to properly set the `speculative` property [#459](https://github.com/hashicorp/go-tfe/pull/459)

# v1.4.0

## Enhancements
* Adds `RetryServerErrors` field to the `Config` object by @sebasslash [#439](https://github.com/hashicorp/go-tfe/pull/439)
* Adds support for the GPG Keys API by @sebasslash [#429](https://github.com/hashicorp/go-tfe/pull/429)
* Adds support for new `WorkspaceLimit` Admin setting for organizations [#425](https://github.com/hashicorp/go-tfe/pull/425)
* Adds support for new `ExcludeTags` workspace list filter field by @Uk1288 [#438](https://github.com/hashicorp/go-tfe/pull/438)
* [beta] Adds additional filter fields to `RunListOptions` by @mjyocca [#424](https://github.com/hashicorp/go-tfe/pull/424)
* [beta] Renames the optional StateVersion field `ExtState` to `JSONStateOutputs` and changes the purpose and type by @annawinkler [#444](https://github.com/hashicorp/go-tfe/pull/444) and @brandoncroft [#452](https://github.com/hashicorp/go-tfe/pull/452)

# v1.3.0

## Enhancements
* Adds support for Microsoft Teams notification configuration by @JarrettSpiker [#398](https://github.com/hashicorp/go-tfe/pull/389)
* Add support for Audit Trail API by @sebasslash [#407](https://github.com/hashicorp/go-tfe/pull/407)
* Adds Private Registry Provider, Provider Version, and Provider Platform APIs support by @joekarl and @annawinkler [#313](https://github.com/hashicorp/go-tfe/pull/313)
* Adds List Registry Modules endpoint by @chroju [#385](https://github.com/hashicorp/go-tfe/pull/385)
* Adds `WebhookURL` field to `VCSRepo` struct by @kgns [#413](https://github.com/hashicorp/go-tfe/pull/413)
* Adds `Category` field to `VariableUpdateOptions` struct by @jtyr [#397](https://github.com/hashicorp/go-tfe/pull/397)
* Adds `TriggerPatterns` to `Workspace` by @matejrisek [#400](https://github.com/hashicorp/go-tfe/pull/400)
* [beta] Adds `ExtState` field to `StateVersionCreateOptions` by @brandonc [#416](https://github.com/hashicorp/go-tfe/pull/416)

# v1.2.0

## Enhancements
* Adds support for reading current state version outputs to StateVersionOutputs, which can be useful for reading outputs when users don't have the necessary permissions to read the entire state by @brandonc [#370](https://github.com/hashicorp/go-tfe/pull/370)
* Adds Variable Set methods for `ApplyToWorkspaces` and `RemoveFromWorkspaces` by @byronwolfman [#375](https://github.com/hashicorp/go-tfe/pull/375)
* Adds `Names` query param field to `TeamListOptions` by @sebasslash [#393](https://github.com/hashicorp/go-tfe/pull/393)
* Adds `Emails` query param field to `OrganizationMembershipListOptions` by @sebasslash [#393](https://github.com/hashicorp/go-tfe/pull/393)
* Adds Run Tasks API support by @glennsarti [#381](https://github.com/hashicorp/go-tfe/pull/381), [#382](https://github.com/hashicorp/go-tfe/pull/382) and [#383](https://github.com/hashicorp/go-tfe/pull/383)


## Bug fixes
* Fixes ignored comment when performing apply, discard, cancel, and force-cancel run actions [#388](https://github.com/hashicorp/go-tfe/pull/388)

# v1.1.0

## Enhancements

* Add Variable Set API support by @rexredinger [#305](https://github.com/hashicorp/go-tfe/pull/305)
* Add Comments API support by @alex-ikse [#355](https://github.com/hashicorp/go-tfe/pull/355)
* Add beta support for SSOTeamID to `Team`, `TeamCreateOptions`, `TeamUpdateOptions` by @xlgmokha [#364](https://github.com/hashicorp/go-tfe/pull/364)

# v1.0.0

## Breaking Changes
* Renamed methods named Generate to Create for `AgentTokens`, `OrganizationTokens`, `TeamTokens`, `UserTokens` by @sebasslash [#327](https://github.com/hashicorp/go-tfe/pull/327)
* Methods that express an action on a relationship have been prefixed with a verb, e.g `Current()` is now `ReadCurrent()` by @sebasslash [#327](https://github.com/hashicorp/go-tfe/pull/327)
* All list option structs are now pointers @uturunku1 [#309](https://github.com/hashicorp/go-tfe/pull/309)
* All errors have been refactored into constants in `errors.go` @uturunku1 [#310](https://github.com/hashicorp/go-tfe/pull/310)
* The `ID` field in Create/Update option structs has been renamed to `Type` in accordance with the JSON:API spec by @omarismail, @uturunku1 [#190](https://github.com/hashicorp/go-tfe/pull/190), [#323](https://github.com/hashicorp/go-tfe/pull/323), [#332](https://github.com/hashicorp/go-tfe/pull/332)
* Nested URL params (consisting of an organization, module and provider name) used to identify a `RegistryModule` have been refactored into a struct `RegistryModuleID` by @sebasslash [#337](https://github.com/hashicorp/go-tfe/pull/337)


## Enhancements
* Added missing include fields for `AdminRuns`, `AgentPools`, `ConfigurationVersions`, `OAuthClients`, `Organizations`, `PolicyChecks`, `PolicySets`, `Policies` and `RunTriggers` by @uturunku1 [#339](https://github.com/hashicorp/go-tfe/pull/339)
* Cleanup documentation and improve consistency by @uturunku1 [#331](https://github.com/hashicorp/go-tfe/pull/331)
* Add more linters to our CI pipeline by @sebasslash [#326](https://github.com/hashicorp/go-tfe/pull/326)
* Resolve `TFE_HOSTNAME` as fallback for `TFE_ADDRESS` by @sebasslash [#340](https://github.com/hashicorp/go-tfe/pull/326)
* Adds a `fetching` status to `RunStatus` and adds the `Archive` method to the ConfigurationVersions interface by @mpminardi [#338](https://github.com/hashicorp/go-tfe/pull/338)
* Added a `Download` method to the `ConfigurationVersions` interface by @tylerwolf [#358](https://github.com/hashicorp/go-tfe/pull/358)
* API Coverage documentation by @laurenolivia [#334](https://github.com/hashicorp/go-tfe/pull/334)

## Bug Fixes
* Fixed invalid memory address error when `AdminSMTPSettingsUpdateOptions.Auth` field is empty and accessed by @uturunku1 [#335](https://github.com/hashicorp/go-tfe/pull/335)
