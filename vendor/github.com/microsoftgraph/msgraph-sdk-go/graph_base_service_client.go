package msgraphsdkgo

import (
    i25911dc319edd61cbac496af7eab5ef20b6069a42515e22ec6a9bc97bf598488 "github.com/microsoft/kiota-serialization-json-go"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i4bcdc892e61ac17e2afc10b5e2b536b29f4fd6c1ad30f4a5a68df47495db3347 "github.com/microsoft/kiota-serialization-form-go"
    i56887720f41ac882814261620b1c8459c4a992a0207af547c4453dd39fabc426 "github.com/microsoft/kiota-serialization-multipart-go"
    i7294a22093d408fdca300f11b81a887d89c47b764af06c8b803e2323973fdb83 "github.com/microsoft/kiota-serialization-text-go"
    i0013a2fa3da6391bf410919641ab77468fe95e3d38c64ffd840561d0ad959a62 "github.com/microsoftgraph/msgraph-sdk-go/applicationswithuniquename"
    i009f47bbce65ccdb7303730eed71e6bab3ae2f8e4e918bc9e94341d28624af97 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
    i07d47a144340607d6d6dbd93575e531530e4f1cc6091c947ea0766f7951ffd34 "github.com/microsoftgraph/msgraph-sdk-go/shares"
    i0906e75d8a44bf92212e084e1d2f62d03887dcec6a5c8535e92ccc04c1e5fdec "github.com/microsoftgraph/msgraph-sdk-go/solutions"
    i185698f71f6301975f0627ee999e6e91920d8fa9c00bdef3487b9f349e2df04e "github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
    i1a1369b1521a8ac4885166fd68eae4247248a891006fea464d2eea2a271b2cdb "github.com/microsoftgraph/msgraph-sdk-go/permissiongrants"
    i1b75be7b5675627960b4672ab148be21ff379d5cbc0e62f6bc5b97d54464f8b5 "github.com/microsoftgraph/msgraph-sdk-go/teamstemplates"
    i1be0f1b1da466bc62355d411ef490acbd8dc0ec5ca4d3448c7eb73e5caffafc3 "github.com/microsoftgraph/msgraph-sdk-go/education"
    i1d6652ecc686b20c37a9a3448b26db8187e284e1a4017cab8876b02b97557436 "github.com/microsoftgraph/msgraph-sdk-go/grouplifecyclepolicies"
    i1dc06c4b7f499cb445a6c55e466abd6d7466bb35a2683c675909db23c57898e7 "github.com/microsoftgraph/msgraph-sdk-go/authenticationmethodconfigurations"
    i20b08d3949f1191430a14a315e0758a1f131dc59bbdc93e654f1dd447a6af14c "github.com/microsoftgraph/msgraph-sdk-go/auditlogs"
    i286f3babd79fe9ec3b0f52b6ed5910842c0adaeff02be1843d0e01c56d9ba6d9 "github.com/microsoftgraph/msgraph-sdk-go/search"
    i2a252d42835bdab6d88bf938595da6cf029001f9ca970d6f599cecf0ca27f8e5 "github.com/microsoftgraph/msgraph-sdk-go/directoryroletemplates"
    i2e58fb6bd805debbded9a12526d5d812045740846d64ca75c71f1976d60e7f38 "github.com/microsoftgraph/msgraph-sdk-go/groupswithuniquename"
    i32d45c1243c349600fbe53b2f9641bb59857a3326037587cbe4e347b46ad207e "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance"
    i35d7bbcc8f7e8b8e9525ea0ee5b3c51c3a1a58f9ed512b727d181bfcd08eb032 "github.com/microsoftgraph/msgraph-sdk-go/security"
    i3e9b5129e2bb8b32b0374f7afe2536be6674d73df6c41d7c529f5a5432c4e0aa "github.com/microsoftgraph/msgraph-sdk-go/agreementacceptances"
    i4794c103c0d044c27a3ca3af0a0e498e93a9863420c1a4e7a29ef37590053c7b "github.com/microsoftgraph/msgraph-sdk-go/groupsettings"
    i49531dee6233ce59a870e4c19c3d202ad6358927c073b9d54fc443416439637d "github.com/microsoftgraph/msgraph-sdk-go/deviceswithdeviceid"
    i4a624e38d68c2a9fc4db1ea915bcaffde116f967f58ec2c99e2ea8bbff3690e1 "github.com/microsoftgraph/msgraph-sdk-go/schemaextensions"
    i4ac7f0a844871066493521918f268cafe2a25c71c28a98221ea3f22d5153090f "github.com/microsoftgraph/msgraph-sdk-go/policies"
    i4c91eeb51f03f9d59a342065f7c6ee027ad1fe84ada6b1946b8162c5ae146cfb "github.com/microsoftgraph/msgraph-sdk-go/devices"
    i51b9802eedc1a25686534d117657be902df58c07e90ac6ea84501100998084d9 "github.com/microsoftgraph/msgraph-sdk-go/communications"
    i5310ba7d4cfddbf5de4c1be94a30f9ca8c747c30a87e76587ce88d1cbfff01b4 "github.com/microsoftgraph/msgraph-sdk-go/applicationtemplates"
    i535d6c02ba98f73ff3a8c1c12a035ba5de51606f93aa2c0babdfed56fe505550 "github.com/microsoftgraph/msgraph-sdk-go/certificatebasedauthconfiguration"
    i58857a108d6e260e56ef0dd7e783668388f113eb436006780703ac59f0abb3b1 "github.com/microsoftgraph/msgraph-sdk-go/privacy"
    i62c2771f3f3a1e5e085aedcde54473e9f043cc57b9ce4dd88980a77aca7a5a10 "github.com/microsoftgraph/msgraph-sdk-go/identityproviders"
    i638650494f9db477daff56d31ff923f5c100f72df0257ed7fa5c222cb1a77a94 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement"
    i663c30678b300c2c4b619c4964b4326e471e4da61a44d7c39f752349da7a468e "github.com/microsoftgraph/msgraph-sdk-go/identityprotection"
    i6bf2d83eea06710580ad0d54b886ac4e14cbab0d1d84937f340f02b99f8f5738 "github.com/microsoftgraph/msgraph-sdk-go/reports"
    i738daeb889f22c1e163aee5a37a094b55b1d815dc76d4802d64e4e1b2e44206c "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement"
    i79097987e7b906ad07243b816ec81a9897e342c7f4eaff9e5dd8a8fcce18e841 "github.com/microsoftgraph/msgraph-sdk-go/functions"
    i79ca23a9ac0659e1330dd29e049fe157787d5af6695ead2ff8263396db68d027 "github.com/microsoftgraph/msgraph-sdk-go/identity"
    i7c9d1b36ac198368c1d8bed014b43e2a518b170ee45bf02c8bbe64544a50539a "github.com/microsoftgraph/msgraph-sdk-go/admin"
    i7d140130aac6882792a019b5ebe51fe8d69dfd63ec213c2e3cd98282ce2d0428 "github.com/microsoftgraph/msgraph-sdk-go/appcatalogs"
    i80d5f91f6f8d9dc3428331303d1837675adde9653ceda73f120faa5f0545ac4b "github.com/microsoftgraph/msgraph-sdk-go/tenantrelationships"
    i86cada4d4a5f2f8a9d1e7a85eacd70a661ea7b20d2737008c0719e95b5be3e16 "github.com/microsoftgraph/msgraph-sdk-go/oauth2permissiongrants"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    i8a94e224d4b14a30028778cb54ce1696f016a1e14d193c4593c6569d5c945b82 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
    i8b6ec7ad760ac5f470c53875acab4b511d51f95fd1aef26c4a386f184a390eac "github.com/microsoftgraph/msgraph-sdk-go/storage"
    i93194122344a685a2f9264205dc6d89a5ba39afdcea57fd0ade8f54b6f137c02 "github.com/microsoftgraph/msgraph-sdk-go/applications"
    i9429d7aae2f5c1dabbecc9411e8ad2b733d29338bc0c0436eeccc94605c461b7 "github.com/microsoftgraph/msgraph-sdk-go/print"
    i957076b10ba162b23efec7b94dd26b84c6475d285449c1cbc9c5b85910d36a12 "github.com/microsoftgraph/msgraph-sdk-go/domains"
    ia3e0f7c2d21d5c73ecb8a7552177d0fe444ae0522290dd1c4b5559e449b118af "github.com/microsoftgraph/msgraph-sdk-go/places"
    ia4b736f581ceef30e9ef8cebd9a6c2b932628e087982ff3dd2c9a0f1a920a918 "github.com/microsoftgraph/msgraph-sdk-go/compliance"
    ia6e876e3ed2d92c29c13dbc8c37513bc38d0d5f05ab9321e43a25ff336912a2d "github.com/microsoftgraph/msgraph-sdk-go/groups"
    iaca6694a878291d0e4021155b406c19d3080cdfc382b456e43c71264d4d9e519 "github.com/microsoftgraph/msgraph-sdk-go/domaindnsrecords"
    ib14d748b564c787931c10f1c7ba6856eeddea29a5b9e5c5c27eb1224ff65e5c4 "github.com/microsoftgraph/msgraph-sdk-go/directory"
    ib3217193884e00033cb8182cac52178dfa3b20ce9c4eb48e37a6217882d956ae "github.com/microsoftgraph/msgraph-sdk-go/external"
    ib33fc5e9889e020c0c572578957f59819123a589c61fd7f3eb37eb7958b525ee "github.com/microsoftgraph/msgraph-sdk-go/datapolicyoperations"
    ib68fa8e66bda853b3a33c491e8a66ca665897dab129192b2c97289266c4a1415 "github.com/microsoftgraph/msgraph-sdk-go/informationprotection"
    ib908319c645932a2c2abf7ce1571c02dfa73f84c9a76e6641ac843c4991c2f48 "github.com/microsoftgraph/msgraph-sdk-go/employeeexperience"
    ibaef614e7692eebc6aaa8080b8ac29169fdf539f24925bc1de4465a3fcdac177 "github.com/microsoftgraph/msgraph-sdk-go/chats"
    ic5e701d75e87f15ce153687b00984a314f7eeea8cfdc77cd9ad648e5ccbc7fbd "github.com/microsoftgraph/msgraph-sdk-go/invitations"
    ic949a0bb5066d68760e8502a7f9db83f571d9e01e38fad4aadf7268188e52df0 "github.com/microsoftgraph/msgraph-sdk-go/organization"
    icabdee72951e77325f237b36d388a199c87e65f67652b6bb85723aba847d7e83 "github.com/microsoftgraph/msgraph-sdk-go/connections"
    icb01e23b9e2000c995fe3ccbcce50c962c8a8bcf9515fda13f2213d1263e0e4f "github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
    ice10f31b9db59ba91184d2b882172edb754f885050cf0830aa2b7c8ff880556b "github.com/microsoftgraph/msgraph-sdk-go/scopedrolememberships"
    id007bc768abbff1131aab64890cdcd0411159a946e9df27140c5f7cf8f249647 "github.com/microsoftgraph/msgraph-sdk-go/subscribedskus"
    id2ac823944414906187dbe4e6ca3b5e46886b9db738d2c1c27de6df8b1bebd61 "github.com/microsoftgraph/msgraph-sdk-go/groupsettingtemplates"
    id4615a956cb1e7edabf8f5a4bc131d1ceca9a13d0f79ae0e122997452a9a0a4e "github.com/microsoftgraph/msgraph-sdk-go/directoryroles"
    id81f15a01b3ceaefa8b1b55f4ee944912f2179aafc4d873f0a2eaf0853eeccd0 "github.com/microsoftgraph/msgraph-sdk-go/authenticationmethodspolicy"
    idb79e5240e0d84d911f5352f7868cce57a8a6d96a758044a11abb7f38d0ba995 "github.com/microsoftgraph/msgraph-sdk-go/directoryroleswithroletemplateid"
    idb8230b65f4a369c23b4d9b41ebe568c657c92f8f77fe36d16d64528b3a317a3 "github.com/microsoftgraph/msgraph-sdk-go/subscriptions"
    ie05ac24b652f7d895cca374316c093c4ca40dd2df0f1518c465233d6432b1ef9 "github.com/microsoftgraph/msgraph-sdk-go/teamwork"
    ie3631868038c44f490dbc03525ac7249d0523c29cc45cbb25b2aebcf470d6c0c "github.com/microsoftgraph/msgraph-sdk-go/contracts"
    ie66b913c1bc1c536bc8db5d185910e9318f621374e016f95e36e9d59b7127f63 "github.com/microsoftgraph/msgraph-sdk-go/planner"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
    ieaa2790c8b7fa361674e69e4a385e279c8c641adf79d86e5b0ca566591a507e8 "github.com/microsoftgraph/msgraph-sdk-go/agreements"
    iefc72d8a17962d4db125c50866617eaa15d662c6e3fb13735d477380dcc0dbe3 "github.com/microsoftgraph/msgraph-sdk-go/drives"
    if39bc788926a05e976b265ecfc616408ca12af399df9ce3a2bb348fe89708057 "github.com/microsoftgraph/msgraph-sdk-go/teams"
    if51cca2652371587dbc02e65260e291435a6a8f7f2ffb419f26c3b9d2a033f57 "github.com/microsoftgraph/msgraph-sdk-go/contacts"
    if5372351befdb652f617b1ee71fbf092fa8dd2a161ba9c021bc265628b6ea82b "github.com/microsoftgraph/msgraph-sdk-go/sites"
    if5555fa41b6637688bcf8c25c62a041258f4dc6eacb38ad42d91c66f222ee182 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement"
    if6ffd1464db2d9c22e351b03e4c00ebd24a5353cd70ffb7f56cfad1c3ceec329 "github.com/microsoftgraph/msgraph-sdk-go/users"
    ifd912bc64ceed11eb9b85cc55c2e7c7a17f682cfe222749139d43f75cf28642a "github.com/microsoftgraph/msgraph-sdk-go/filteroperators"
)

// GraphBaseServiceClient the main entry point of the SDK, exposes the configuration and the fluent API.
type GraphBaseServiceClient struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// Admin provides operations to manage the admin singleton.
// returns a *AdminRequestBuilder when successful
func (m *GraphBaseServiceClient) Admin()(*i7c9d1b36ac198368c1d8bed014b43e2a518b170ee45bf02c8bbe64544a50539a.AdminRequestBuilder) {
    return i7c9d1b36ac198368c1d8bed014b43e2a518b170ee45bf02c8bbe64544a50539a.NewAdminRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AgreementAcceptances provides operations to manage the collection of agreementAcceptance entities.
// returns a *AgreementAcceptancesRequestBuilder when successful
func (m *GraphBaseServiceClient) AgreementAcceptances()(*i3e9b5129e2bb8b32b0374f7afe2536be6674d73df6c41d7c529f5a5432c4e0aa.AgreementAcceptancesRequestBuilder) {
    return i3e9b5129e2bb8b32b0374f7afe2536be6674d73df6c41d7c529f5a5432c4e0aa.NewAgreementAcceptancesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Agreements provides operations to manage the collection of agreement entities.
// returns a *AgreementsRequestBuilder when successful
func (m *GraphBaseServiceClient) Agreements()(*ieaa2790c8b7fa361674e69e4a385e279c8c641adf79d86e5b0ca566591a507e8.AgreementsRequestBuilder) {
    return ieaa2790c8b7fa361674e69e4a385e279c8c641adf79d86e5b0ca566591a507e8.NewAgreementsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AppCatalogs provides operations to manage the appCatalogs singleton.
// returns a *AppCatalogsRequestBuilder when successful
func (m *GraphBaseServiceClient) AppCatalogs()(*i7d140130aac6882792a019b5ebe51fe8d69dfd63ec213c2e3cd98282ce2d0428.AppCatalogsRequestBuilder) {
    return i7d140130aac6882792a019b5ebe51fe8d69dfd63ec213c2e3cd98282ce2d0428.NewAppCatalogsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Applications provides operations to manage the collection of application entities.
// returns a *ApplicationsRequestBuilder when successful
func (m *GraphBaseServiceClient) Applications()(*i93194122344a685a2f9264205dc6d89a5ba39afdcea57fd0ade8f54b6f137c02.ApplicationsRequestBuilder) {
    return i93194122344a685a2f9264205dc6d89a5ba39afdcea57fd0ade8f54b6f137c02.NewApplicationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ApplicationsWithAppId provides operations to manage the collection of application entities.
// returns a *ApplicationsWithAppIdRequestBuilder when successful
func (m *GraphBaseServiceClient) ApplicationsWithAppId(appId *string)(*icb01e23b9e2000c995fe3ccbcce50c962c8a8bcf9515fda13f2213d1263e0e4f.ApplicationsWithAppIdRequestBuilder) {
    return icb01e23b9e2000c995fe3ccbcce50c962c8a8bcf9515fda13f2213d1263e0e4f.NewApplicationsWithAppIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, appId)
}
// ApplicationsWithUniqueName provides operations to manage the collection of application entities.
// returns a *ApplicationsWithUniqueNameRequestBuilder when successful
func (m *GraphBaseServiceClient) ApplicationsWithUniqueName(uniqueName *string)(*i0013a2fa3da6391bf410919641ab77468fe95e3d38c64ffd840561d0ad959a62.ApplicationsWithUniqueNameRequestBuilder) {
    return i0013a2fa3da6391bf410919641ab77468fe95e3d38c64ffd840561d0ad959a62.NewApplicationsWithUniqueNameRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, uniqueName)
}
// ApplicationTemplates provides operations to manage the collection of applicationTemplate entities.
// returns a *ApplicationTemplatesRequestBuilder when successful
func (m *GraphBaseServiceClient) ApplicationTemplates()(*i5310ba7d4cfddbf5de4c1be94a30f9ca8c747c30a87e76587ce88d1cbfff01b4.ApplicationTemplatesRequestBuilder) {
    return i5310ba7d4cfddbf5de4c1be94a30f9ca8c747c30a87e76587ce88d1cbfff01b4.NewApplicationTemplatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuditLogs provides operations to manage the auditLogRoot singleton.
// returns a *AuditLogsRequestBuilder when successful
func (m *GraphBaseServiceClient) AuditLogs()(*i20b08d3949f1191430a14a315e0758a1f131dc59bbdc93e654f1dd447a6af14c.AuditLogsRequestBuilder) {
    return i20b08d3949f1191430a14a315e0758a1f131dc59bbdc93e654f1dd447a6af14c.NewAuditLogsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthenticationMethodConfigurations provides operations to manage the collection of authenticationMethodConfiguration entities.
// returns a *AuthenticationMethodConfigurationsRequestBuilder when successful
func (m *GraphBaseServiceClient) AuthenticationMethodConfigurations()(*i1dc06c4b7f499cb445a6c55e466abd6d7466bb35a2683c675909db23c57898e7.AuthenticationMethodConfigurationsRequestBuilder) {
    return i1dc06c4b7f499cb445a6c55e466abd6d7466bb35a2683c675909db23c57898e7.NewAuthenticationMethodConfigurationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthenticationMethodsPolicy provides operations to manage the authenticationMethodsPolicy singleton.
// returns a *AuthenticationMethodsPolicyRequestBuilder when successful
func (m *GraphBaseServiceClient) AuthenticationMethodsPolicy()(*id81f15a01b3ceaefa8b1b55f4ee944912f2179aafc4d873f0a2eaf0853eeccd0.AuthenticationMethodsPolicyRequestBuilder) {
    return id81f15a01b3ceaefa8b1b55f4ee944912f2179aafc4d873f0a2eaf0853eeccd0.NewAuthenticationMethodsPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CertificateBasedAuthConfiguration provides operations to manage the collection of certificateBasedAuthConfiguration entities.
// returns a *CertificateBasedAuthConfigurationRequestBuilder when successful
func (m *GraphBaseServiceClient) CertificateBasedAuthConfiguration()(*i535d6c02ba98f73ff3a8c1c12a035ba5de51606f93aa2c0babdfed56fe505550.CertificateBasedAuthConfigurationRequestBuilder) {
    return i535d6c02ba98f73ff3a8c1c12a035ba5de51606f93aa2c0babdfed56fe505550.NewCertificateBasedAuthConfigurationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Chats provides operations to manage the collection of chat entities.
// returns a *ChatsRequestBuilder when successful
func (m *GraphBaseServiceClient) Chats()(*ibaef614e7692eebc6aaa8080b8ac29169fdf539f24925bc1de4465a3fcdac177.ChatsRequestBuilder) {
    return ibaef614e7692eebc6aaa8080b8ac29169fdf539f24925bc1de4465a3fcdac177.NewChatsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Communications provides operations to manage the cloudCommunications singleton.
// returns a *CommunicationsRequestBuilder when successful
func (m *GraphBaseServiceClient) Communications()(*i51b9802eedc1a25686534d117657be902df58c07e90ac6ea84501100998084d9.CommunicationsRequestBuilder) {
    return i51b9802eedc1a25686534d117657be902df58c07e90ac6ea84501100998084d9.NewCommunicationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Compliance provides operations to manage the compliance singleton.
// returns a *ComplianceRequestBuilder when successful
func (m *GraphBaseServiceClient) Compliance()(*ia4b736f581ceef30e9ef8cebd9a6c2b932628e087982ff3dd2c9a0f1a920a918.ComplianceRequestBuilder) {
    return ia4b736f581ceef30e9ef8cebd9a6c2b932628e087982ff3dd2c9a0f1a920a918.NewComplianceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Connections provides operations to manage the collection of externalConnection entities.
// returns a *ConnectionsRequestBuilder when successful
func (m *GraphBaseServiceClient) Connections()(*icabdee72951e77325f237b36d388a199c87e65f67652b6bb85723aba847d7e83.ConnectionsRequestBuilder) {
    return icabdee72951e77325f237b36d388a199c87e65f67652b6bb85723aba847d7e83.NewConnectionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewGraphBaseServiceClient instantiates a new GraphBaseServiceClient and sets the default values.
func NewGraphBaseServiceClient(requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactory)(*GraphBaseServiceClient) {
    m := &GraphBaseServiceClient{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}", map[string]string{}),
    }
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultSerializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriterFactory { return i25911dc319edd61cbac496af7eab5ef20b6069a42515e22ec6a9bc97bf598488.NewJsonSerializationWriterFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultSerializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriterFactory { return i7294a22093d408fdca300f11b81a887d89c47b764af06c8b803e2323973fdb83.NewTextSerializationWriterFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultSerializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriterFactory { return i4bcdc892e61ac17e2afc10b5e2b536b29f4fd6c1ad30f4a5a68df47495db3347.NewFormSerializationWriterFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultSerializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriterFactory { return i56887720f41ac882814261620b1c8459c4a992a0207af547c4453dd39fabc426.NewMultipartSerializationWriterFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultDeserializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNodeFactory { return i25911dc319edd61cbac496af7eab5ef20b6069a42515e22ec6a9bc97bf598488.NewJsonParseNodeFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultDeserializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNodeFactory { return i7294a22093d408fdca300f11b81a887d89c47b764af06c8b803e2323973fdb83.NewTextParseNodeFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultDeserializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNodeFactory { return i4bcdc892e61ac17e2afc10b5e2b536b29f4fd6c1ad30f4a5a68df47495db3347.NewFormParseNodeFactory() })
    if m.BaseRequestBuilder.RequestAdapter.GetBaseUrl() == "" {
        m.BaseRequestBuilder.RequestAdapter.SetBaseUrl("https://graph.microsoft.com/v1.0")
    }
    m.BaseRequestBuilder.PathParameters["baseurl"] = m.BaseRequestBuilder.RequestAdapter.GetBaseUrl()
    m.BaseRequestBuilder.RequestAdapter.EnableBackingStore(backingStore);
    return m
}
// Contacts provides operations to manage the collection of orgContact entities.
// returns a *ContactsRequestBuilder when successful
func (m *GraphBaseServiceClient) Contacts()(*if51cca2652371587dbc02e65260e291435a6a8f7f2ffb419f26c3b9d2a033f57.ContactsRequestBuilder) {
    return if51cca2652371587dbc02e65260e291435a6a8f7f2ffb419f26c3b9d2a033f57.NewContactsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Contracts provides operations to manage the collection of contract entities.
// returns a *ContractsRequestBuilder when successful
func (m *GraphBaseServiceClient) Contracts()(*ie3631868038c44f490dbc03525ac7249d0523c29cc45cbb25b2aebcf470d6c0c.ContractsRequestBuilder) {
    return ie3631868038c44f490dbc03525ac7249d0523c29cc45cbb25b2aebcf470d6c0c.NewContractsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DataPolicyOperations provides operations to manage the collection of dataPolicyOperation entities.
// returns a *DataPolicyOperationsRequestBuilder when successful
func (m *GraphBaseServiceClient) DataPolicyOperations()(*ib33fc5e9889e020c0c572578957f59819123a589c61fd7f3eb37eb7958b525ee.DataPolicyOperationsRequestBuilder) {
    return ib33fc5e9889e020c0c572578957f59819123a589c61fd7f3eb37eb7958b525ee.NewDataPolicyOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceAppManagement provides operations to manage the deviceAppManagement singleton.
// returns a *DeviceAppManagementRequestBuilder when successful
func (m *GraphBaseServiceClient) DeviceAppManagement()(*i638650494f9db477daff56d31ff923f5c100f72df0257ed7fa5c222cb1a77a94.DeviceAppManagementRequestBuilder) {
    return i638650494f9db477daff56d31ff923f5c100f72df0257ed7fa5c222cb1a77a94.NewDeviceAppManagementRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceManagement provides operations to manage the deviceManagement singleton.
// returns a *DeviceManagementRequestBuilder when successful
func (m *GraphBaseServiceClient) DeviceManagement()(*i738daeb889f22c1e163aee5a37a094b55b1d815dc76d4802d64e4e1b2e44206c.DeviceManagementRequestBuilder) {
    return i738daeb889f22c1e163aee5a37a094b55b1d815dc76d4802d64e4e1b2e44206c.NewDeviceManagementRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Devices provides operations to manage the collection of device entities.
// returns a *DevicesRequestBuilder when successful
func (m *GraphBaseServiceClient) Devices()(*i4c91eeb51f03f9d59a342065f7c6ee027ad1fe84ada6b1946b8162c5ae146cfb.DevicesRequestBuilder) {
    return i4c91eeb51f03f9d59a342065f7c6ee027ad1fe84ada6b1946b8162c5ae146cfb.NewDevicesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DevicesWithDeviceId provides operations to manage the collection of device entities.
// returns a *DevicesWithDeviceIdRequestBuilder when successful
func (m *GraphBaseServiceClient) DevicesWithDeviceId(deviceId *string)(*i49531dee6233ce59a870e4c19c3d202ad6358927c073b9d54fc443416439637d.DevicesWithDeviceIdRequestBuilder) {
    return i49531dee6233ce59a870e4c19c3d202ad6358927c073b9d54fc443416439637d.NewDevicesWithDeviceIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deviceId)
}
// Directory provides operations to manage the directory singleton.
// returns a *DirectoryRequestBuilder when successful
func (m *GraphBaseServiceClient) Directory()(*ib14d748b564c787931c10f1c7ba6856eeddea29a5b9e5c5c27eb1224ff65e5c4.DirectoryRequestBuilder) {
    return ib14d748b564c787931c10f1c7ba6856eeddea29a5b9e5c5c27eb1224ff65e5c4.NewDirectoryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DirectoryObjects provides operations to manage the collection of directoryObject entities.
// returns a *DirectoryObjectsRequestBuilder when successful
func (m *GraphBaseServiceClient) DirectoryObjects()(*i185698f71f6301975f0627ee999e6e91920d8fa9c00bdef3487b9f349e2df04e.DirectoryObjectsRequestBuilder) {
    return i185698f71f6301975f0627ee999e6e91920d8fa9c00bdef3487b9f349e2df04e.NewDirectoryObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DirectoryRoles provides operations to manage the collection of directoryRole entities.
// returns a *DirectoryRolesRequestBuilder when successful
func (m *GraphBaseServiceClient) DirectoryRoles()(*id4615a956cb1e7edabf8f5a4bc131d1ceca9a13d0f79ae0e122997452a9a0a4e.DirectoryRolesRequestBuilder) {
    return id4615a956cb1e7edabf8f5a4bc131d1ceca9a13d0f79ae0e122997452a9a0a4e.NewDirectoryRolesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DirectoryRolesWithRoleTemplateId provides operations to manage the collection of directoryRole entities.
// returns a *DirectoryRolesWithRoleTemplateIdRequestBuilder when successful
func (m *GraphBaseServiceClient) DirectoryRolesWithRoleTemplateId(roleTemplateId *string)(*idb79e5240e0d84d911f5352f7868cce57a8a6d96a758044a11abb7f38d0ba995.DirectoryRolesWithRoleTemplateIdRequestBuilder) {
    return idb79e5240e0d84d911f5352f7868cce57a8a6d96a758044a11abb7f38d0ba995.NewDirectoryRolesWithRoleTemplateIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, roleTemplateId)
}
// DirectoryRoleTemplates provides operations to manage the collection of directoryRoleTemplate entities.
// returns a *DirectoryRoleTemplatesRequestBuilder when successful
func (m *GraphBaseServiceClient) DirectoryRoleTemplates()(*i2a252d42835bdab6d88bf938595da6cf029001f9ca970d6f599cecf0ca27f8e5.DirectoryRoleTemplatesRequestBuilder) {
    return i2a252d42835bdab6d88bf938595da6cf029001f9ca970d6f599cecf0ca27f8e5.NewDirectoryRoleTemplatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DomainDnsRecords provides operations to manage the collection of domainDnsRecord entities.
// returns a *DomainDnsRecordsRequestBuilder when successful
func (m *GraphBaseServiceClient) DomainDnsRecords()(*iaca6694a878291d0e4021155b406c19d3080cdfc382b456e43c71264d4d9e519.DomainDnsRecordsRequestBuilder) {
    return iaca6694a878291d0e4021155b406c19d3080cdfc382b456e43c71264d4d9e519.NewDomainDnsRecordsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Domains provides operations to manage the collection of domain entities.
// returns a *DomainsRequestBuilder when successful
func (m *GraphBaseServiceClient) Domains()(*i957076b10ba162b23efec7b94dd26b84c6475d285449c1cbc9c5b85910d36a12.DomainsRequestBuilder) {
    return i957076b10ba162b23efec7b94dd26b84c6475d285449c1cbc9c5b85910d36a12.NewDomainsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drives provides operations to manage the collection of drive entities.
// returns a *DrivesRequestBuilder when successful
func (m *GraphBaseServiceClient) Drives()(*iefc72d8a17962d4db125c50866617eaa15d662c6e3fb13735d477380dcc0dbe3.DrivesRequestBuilder) {
    return iefc72d8a17962d4db125c50866617eaa15d662c6e3fb13735d477380dcc0dbe3.NewDrivesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Education provides operations to manage the educationRoot singleton.
// returns a *EducationRequestBuilder when successful
func (m *GraphBaseServiceClient) Education()(*i1be0f1b1da466bc62355d411ef490acbd8dc0ec5ca4d3448c7eb73e5caffafc3.EducationRequestBuilder) {
    return i1be0f1b1da466bc62355d411ef490acbd8dc0ec5ca4d3448c7eb73e5caffafc3.NewEducationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EmployeeExperience provides operations to manage the employeeExperience singleton.
// returns a *EmployeeExperienceRequestBuilder when successful
func (m *GraphBaseServiceClient) EmployeeExperience()(*ib908319c645932a2c2abf7ce1571c02dfa73f84c9a76e6641ac843c4991c2f48.EmployeeExperienceRequestBuilder) {
    return ib908319c645932a2c2abf7ce1571c02dfa73f84c9a76e6641ac843c4991c2f48.NewEmployeeExperienceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// External provides operations to manage the external singleton.
// returns a *ExternalRequestBuilder when successful
func (m *GraphBaseServiceClient) External()(*ib3217193884e00033cb8182cac52178dfa3b20ce9c4eb48e37a6217882d956ae.ExternalRequestBuilder) {
    return ib3217193884e00033cb8182cac52178dfa3b20ce9c4eb48e37a6217882d956ae.NewExternalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FilterOperators provides operations to manage the collection of filterOperatorSchema entities.
// returns a *FilterOperatorsRequestBuilder when successful
func (m *GraphBaseServiceClient) FilterOperators()(*ifd912bc64ceed11eb9b85cc55c2e7c7a17f682cfe222749139d43f75cf28642a.FilterOperatorsRequestBuilder) {
    return ifd912bc64ceed11eb9b85cc55c2e7c7a17f682cfe222749139d43f75cf28642a.NewFilterOperatorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Functions provides operations to manage the collection of attributeMappingFunctionSchema entities.
// returns a *FunctionsRequestBuilder when successful
func (m *GraphBaseServiceClient) Functions()(*i79097987e7b906ad07243b816ec81a9897e342c7f4eaff9e5dd8a8fcce18e841.FunctionsRequestBuilder) {
    return i79097987e7b906ad07243b816ec81a9897e342c7f4eaff9e5dd8a8fcce18e841.NewFunctionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GroupLifecyclePolicies provides operations to manage the collection of groupLifecyclePolicy entities.
// returns a *GroupLifecyclePoliciesRequestBuilder when successful
func (m *GraphBaseServiceClient) GroupLifecyclePolicies()(*i1d6652ecc686b20c37a9a3448b26db8187e284e1a4017cab8876b02b97557436.GroupLifecyclePoliciesRequestBuilder) {
    return i1d6652ecc686b20c37a9a3448b26db8187e284e1a4017cab8876b02b97557436.NewGroupLifecyclePoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Groups provides operations to manage the collection of group entities.
// returns a *GroupsRequestBuilder when successful
func (m *GraphBaseServiceClient) Groups()(*ia6e876e3ed2d92c29c13dbc8c37513bc38d0d5f05ab9321e43a25ff336912a2d.GroupsRequestBuilder) {
    return ia6e876e3ed2d92c29c13dbc8c37513bc38d0d5f05ab9321e43a25ff336912a2d.NewGroupsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GroupSettings provides operations to manage the collection of groupSetting entities.
// returns a *GroupSettingsRequestBuilder when successful
func (m *GraphBaseServiceClient) GroupSettings()(*i4794c103c0d044c27a3ca3af0a0e498e93a9863420c1a4e7a29ef37590053c7b.GroupSettingsRequestBuilder) {
    return i4794c103c0d044c27a3ca3af0a0e498e93a9863420c1a4e7a29ef37590053c7b.NewGroupSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GroupSettingTemplates provides operations to manage the collection of groupSettingTemplate entities.
// returns a *GroupSettingTemplatesRequestBuilder when successful
func (m *GraphBaseServiceClient) GroupSettingTemplates()(*id2ac823944414906187dbe4e6ca3b5e46886b9db738d2c1c27de6df8b1bebd61.GroupSettingTemplatesRequestBuilder) {
    return id2ac823944414906187dbe4e6ca3b5e46886b9db738d2c1c27de6df8b1bebd61.NewGroupSettingTemplatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GroupsWithUniqueName provides operations to manage the collection of group entities.
// returns a *GroupsWithUniqueNameRequestBuilder when successful
func (m *GraphBaseServiceClient) GroupsWithUniqueName(uniqueName *string)(*i2e58fb6bd805debbded9a12526d5d812045740846d64ca75c71f1976d60e7f38.GroupsWithUniqueNameRequestBuilder) {
    return i2e58fb6bd805debbded9a12526d5d812045740846d64ca75c71f1976d60e7f38.NewGroupsWithUniqueNameRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, uniqueName)
}
// Identity provides operations to manage the identityContainer singleton.
// returns a *IdentityRequestBuilder when successful
func (m *GraphBaseServiceClient) Identity()(*i79ca23a9ac0659e1330dd29e049fe157787d5af6695ead2ff8263396db68d027.IdentityRequestBuilder) {
    return i79ca23a9ac0659e1330dd29e049fe157787d5af6695ead2ff8263396db68d027.NewIdentityRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IdentityGovernance provides operations to manage the identityGovernance singleton.
// returns a *IdentityGovernanceRequestBuilder when successful
func (m *GraphBaseServiceClient) IdentityGovernance()(*i32d45c1243c349600fbe53b2f9641bb59857a3326037587cbe4e347b46ad207e.IdentityGovernanceRequestBuilder) {
    return i32d45c1243c349600fbe53b2f9641bb59857a3326037587cbe4e347b46ad207e.NewIdentityGovernanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IdentityProtection provides operations to manage the identityProtectionRoot singleton.
// returns a *IdentityProtectionRequestBuilder when successful
func (m *GraphBaseServiceClient) IdentityProtection()(*i663c30678b300c2c4b619c4964b4326e471e4da61a44d7c39f752349da7a468e.IdentityProtectionRequestBuilder) {
    return i663c30678b300c2c4b619c4964b4326e471e4da61a44d7c39f752349da7a468e.NewIdentityProtectionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IdentityProviders provides operations to manage the collection of identityProvider entities.
// returns a *IdentityProvidersRequestBuilder when successful
func (m *GraphBaseServiceClient) IdentityProviders()(*i62c2771f3f3a1e5e085aedcde54473e9f043cc57b9ce4dd88980a77aca7a5a10.IdentityProvidersRequestBuilder) {
    return i62c2771f3f3a1e5e085aedcde54473e9f043cc57b9ce4dd88980a77aca7a5a10.NewIdentityProvidersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// InformationProtection provides operations to manage the informationProtection singleton.
// returns a *InformationProtectionRequestBuilder when successful
func (m *GraphBaseServiceClient) InformationProtection()(*ib68fa8e66bda853b3a33c491e8a66ca665897dab129192b2c97289266c4a1415.InformationProtectionRequestBuilder) {
    return ib68fa8e66bda853b3a33c491e8a66ca665897dab129192b2c97289266c4a1415.NewInformationProtectionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Invitations provides operations to manage the collection of invitation entities.
// returns a *InvitationsRequestBuilder when successful
func (m *GraphBaseServiceClient) Invitations()(*ic5e701d75e87f15ce153687b00984a314f7eeea8cfdc77cd9ad648e5ccbc7fbd.InvitationsRequestBuilder) {
    return ic5e701d75e87f15ce153687b00984a314f7eeea8cfdc77cd9ad648e5ccbc7fbd.NewInvitationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Oauth2PermissionGrants provides operations to manage the collection of oAuth2PermissionGrant entities.
// returns a *Oauth2PermissionGrantsRequestBuilder when successful
func (m *GraphBaseServiceClient) Oauth2PermissionGrants()(*i86cada4d4a5f2f8a9d1e7a85eacd70a661ea7b20d2737008c0719e95b5be3e16.Oauth2PermissionGrantsRequestBuilder) {
    return i86cada4d4a5f2f8a9d1e7a85eacd70a661ea7b20d2737008c0719e95b5be3e16.NewOauth2PermissionGrantsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Organization provides operations to manage the collection of organization entities.
// returns a *OrganizationRequestBuilder when successful
func (m *GraphBaseServiceClient) Organization()(*ic949a0bb5066d68760e8502a7f9db83f571d9e01e38fad4aadf7268188e52df0.OrganizationRequestBuilder) {
    return ic949a0bb5066d68760e8502a7f9db83f571d9e01e38fad4aadf7268188e52df0.NewOrganizationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PermissionGrants provides operations to manage the collection of resourceSpecificPermissionGrant entities.
// returns a *PermissionGrantsRequestBuilder when successful
func (m *GraphBaseServiceClient) PermissionGrants()(*i1a1369b1521a8ac4885166fd68eae4247248a891006fea464d2eea2a271b2cdb.PermissionGrantsRequestBuilder) {
    return i1a1369b1521a8ac4885166fd68eae4247248a891006fea464d2eea2a271b2cdb.NewPermissionGrantsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Places the places property
// returns a *PlacesRequestBuilder when successful
func (m *GraphBaseServiceClient) Places()(*ia3e0f7c2d21d5c73ecb8a7552177d0fe444ae0522290dd1c4b5559e449b118af.PlacesRequestBuilder) {
    return ia3e0f7c2d21d5c73ecb8a7552177d0fe444ae0522290dd1c4b5559e449b118af.NewPlacesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Planner provides operations to manage the planner singleton.
// returns a *PlannerRequestBuilder when successful
func (m *GraphBaseServiceClient) Planner()(*ie66b913c1bc1c536bc8db5d185910e9318f621374e016f95e36e9d59b7127f63.PlannerRequestBuilder) {
    return ie66b913c1bc1c536bc8db5d185910e9318f621374e016f95e36e9d59b7127f63.NewPlannerRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Policies provides operations to manage the policyRoot singleton.
// returns a *PoliciesRequestBuilder when successful
func (m *GraphBaseServiceClient) Policies()(*i4ac7f0a844871066493521918f268cafe2a25c71c28a98221ea3f22d5153090f.PoliciesRequestBuilder) {
    return i4ac7f0a844871066493521918f268cafe2a25c71c28a98221ea3f22d5153090f.NewPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Print provides operations to manage the print singleton.
// returns a *PrintRequestBuilder when successful
func (m *GraphBaseServiceClient) Print()(*i9429d7aae2f5c1dabbecc9411e8ad2b733d29338bc0c0436eeccc94605c461b7.PrintRequestBuilder) {
    return i9429d7aae2f5c1dabbecc9411e8ad2b733d29338bc0c0436eeccc94605c461b7.NewPrintRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Privacy provides operations to manage the privacy singleton.
// returns a *PrivacyRequestBuilder when successful
func (m *GraphBaseServiceClient) Privacy()(*i58857a108d6e260e56ef0dd7e783668388f113eb436006780703ac59f0abb3b1.PrivacyRequestBuilder) {
    return i58857a108d6e260e56ef0dd7e783668388f113eb436006780703ac59f0abb3b1.NewPrivacyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Reports provides operations to manage the reportRoot singleton.
// returns a *ReportsRequestBuilder when successful
func (m *GraphBaseServiceClient) Reports()(*i6bf2d83eea06710580ad0d54b886ac4e14cbab0d1d84937f340f02b99f8f5738.ReportsRequestBuilder) {
    return i6bf2d83eea06710580ad0d54b886ac4e14cbab0d1d84937f340f02b99f8f5738.NewReportsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoleManagement provides operations to manage the roleManagement singleton.
// returns a *RoleManagementRequestBuilder when successful
func (m *GraphBaseServiceClient) RoleManagement()(*if5555fa41b6637688bcf8c25c62a041258f4dc6eacb38ad42d91c66f222ee182.RoleManagementRequestBuilder) {
    return if5555fa41b6637688bcf8c25c62a041258f4dc6eacb38ad42d91c66f222ee182.NewRoleManagementRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SchemaExtensions provides operations to manage the collection of schemaExtension entities.
// returns a *SchemaExtensionsRequestBuilder when successful
func (m *GraphBaseServiceClient) SchemaExtensions()(*i4a624e38d68c2a9fc4db1ea915bcaffde116f967f58ec2c99e2ea8bbff3690e1.SchemaExtensionsRequestBuilder) {
    return i4a624e38d68c2a9fc4db1ea915bcaffde116f967f58ec2c99e2ea8bbff3690e1.NewSchemaExtensionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ScopedRoleMemberships provides operations to manage the collection of scopedRoleMembership entities.
// returns a *ScopedRoleMembershipsRequestBuilder when successful
func (m *GraphBaseServiceClient) ScopedRoleMemberships()(*ice10f31b9db59ba91184d2b882172edb754f885050cf0830aa2b7c8ff880556b.ScopedRoleMembershipsRequestBuilder) {
    return ice10f31b9db59ba91184d2b882172edb754f885050cf0830aa2b7c8ff880556b.NewScopedRoleMembershipsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Search provides operations to manage the searchEntity singleton.
// returns a *SearchRequestBuilder when successful
func (m *GraphBaseServiceClient) Search()(*i286f3babd79fe9ec3b0f52b6ed5910842c0adaeff02be1843d0e01c56d9ba6d9.SearchRequestBuilder) {
    return i286f3babd79fe9ec3b0f52b6ed5910842c0adaeff02be1843d0e01c56d9ba6d9.NewSearchRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Security provides operations to manage the security singleton.
// returns a *SecurityRequestBuilder when successful
func (m *GraphBaseServiceClient) Security()(*i35d7bbcc8f7e8b8e9525ea0ee5b3c51c3a1a58f9ed512b727d181bfcd08eb032.SecurityRequestBuilder) {
    return i35d7bbcc8f7e8b8e9525ea0ee5b3c51c3a1a58f9ed512b727d181bfcd08eb032.NewSecurityRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServicePrincipals provides operations to manage the collection of servicePrincipal entities.
// returns a *ServicePrincipalsRequestBuilder when successful
func (m *GraphBaseServiceClient) ServicePrincipals()(*i009f47bbce65ccdb7303730eed71e6bab3ae2f8e4e918bc9e94341d28624af97.ServicePrincipalsRequestBuilder) {
    return i009f47bbce65ccdb7303730eed71e6bab3ae2f8e4e918bc9e94341d28624af97.NewServicePrincipalsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServicePrincipalsWithAppId provides operations to manage the collection of servicePrincipal entities.
// returns a *ServicePrincipalsWithAppIdRequestBuilder when successful
func (m *GraphBaseServiceClient) ServicePrincipalsWithAppId(appId *string)(*i8a94e224d4b14a30028778cb54ce1696f016a1e14d193c4593c6569d5c945b82.ServicePrincipalsWithAppIdRequestBuilder) {
    return i8a94e224d4b14a30028778cb54ce1696f016a1e14d193c4593c6569d5c945b82.NewServicePrincipalsWithAppIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, appId)
}
// Shares provides operations to manage the collection of sharedDriveItem entities.
// returns a *SharesRequestBuilder when successful
func (m *GraphBaseServiceClient) Shares()(*i07d47a144340607d6d6dbd93575e531530e4f1cc6091c947ea0766f7951ffd34.SharesRequestBuilder) {
    return i07d47a144340607d6d6dbd93575e531530e4f1cc6091c947ea0766f7951ffd34.NewSharesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sites provides operations to manage the collection of site entities.
// returns a *SitesRequestBuilder when successful
func (m *GraphBaseServiceClient) Sites()(*if5372351befdb652f617b1ee71fbf092fa8dd2a161ba9c021bc265628b6ea82b.SitesRequestBuilder) {
    return if5372351befdb652f617b1ee71fbf092fa8dd2a161ba9c021bc265628b6ea82b.NewSitesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Solutions provides operations to manage the solutionsRoot singleton.
// returns a *SolutionsRequestBuilder when successful
func (m *GraphBaseServiceClient) Solutions()(*i0906e75d8a44bf92212e084e1d2f62d03887dcec6a5c8535e92ccc04c1e5fdec.SolutionsRequestBuilder) {
    return i0906e75d8a44bf92212e084e1d2f62d03887dcec6a5c8535e92ccc04c1e5fdec.NewSolutionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Storage provides operations to manage the storage singleton.
// returns a *StorageRequestBuilder when successful
func (m *GraphBaseServiceClient) Storage()(*i8b6ec7ad760ac5f470c53875acab4b511d51f95fd1aef26c4a386f184a390eac.StorageRequestBuilder) {
    return i8b6ec7ad760ac5f470c53875acab4b511d51f95fd1aef26c4a386f184a390eac.NewStorageRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SubscribedSkus provides operations to manage the collection of subscribedSku entities.
// returns a *SubscribedSkusRequestBuilder when successful
func (m *GraphBaseServiceClient) SubscribedSkus()(*id007bc768abbff1131aab64890cdcd0411159a946e9df27140c5f7cf8f249647.SubscribedSkusRequestBuilder) {
    return id007bc768abbff1131aab64890cdcd0411159a946e9df27140c5f7cf8f249647.NewSubscribedSkusRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Subscriptions provides operations to manage the collection of subscription entities.
// returns a *SubscriptionsRequestBuilder when successful
func (m *GraphBaseServiceClient) Subscriptions()(*idb8230b65f4a369c23b4d9b41ebe568c657c92f8f77fe36d16d64528b3a317a3.SubscriptionsRequestBuilder) {
    return idb8230b65f4a369c23b4d9b41ebe568c657c92f8f77fe36d16d64528b3a317a3.NewSubscriptionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Teams provides operations to manage the collection of team entities.
// returns a *TeamsRequestBuilder when successful
func (m *GraphBaseServiceClient) Teams()(*if39bc788926a05e976b265ecfc616408ca12af399df9ce3a2bb348fe89708057.TeamsRequestBuilder) {
    return if39bc788926a05e976b265ecfc616408ca12af399df9ce3a2bb348fe89708057.NewTeamsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TeamsTemplates provides operations to manage the collection of teamsTemplate entities.
// returns a *TeamsTemplatesRequestBuilder when successful
func (m *GraphBaseServiceClient) TeamsTemplates()(*i1b75be7b5675627960b4672ab148be21ff379d5cbc0e62f6bc5b97d54464f8b5.TeamsTemplatesRequestBuilder) {
    return i1b75be7b5675627960b4672ab148be21ff379d5cbc0e62f6bc5b97d54464f8b5.NewTeamsTemplatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Teamwork provides operations to manage the teamwork singleton.
// returns a *TeamworkRequestBuilder when successful
func (m *GraphBaseServiceClient) Teamwork()(*ie05ac24b652f7d895cca374316c093c4ca40dd2df0f1518c465233d6432b1ef9.TeamworkRequestBuilder) {
    return ie05ac24b652f7d895cca374316c093c4ca40dd2df0f1518c465233d6432b1ef9.NewTeamworkRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TenantRelationships provides operations to manage the tenantRelationship singleton.
// returns a *TenantRelationshipsRequestBuilder when successful
func (m *GraphBaseServiceClient) TenantRelationships()(*i80d5f91f6f8d9dc3428331303d1837675adde9653ceda73f120faa5f0545ac4b.TenantRelationshipsRequestBuilder) {
    return i80d5f91f6f8d9dc3428331303d1837675adde9653ceda73f120faa5f0545ac4b.NewTenantRelationshipsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Users provides operations to manage the collection of user entities.
// returns a *UsersRequestBuilder when successful
func (m *GraphBaseServiceClient) Users()(*if6ffd1464db2d9c22e351b03e4c00ebd24a5353cd70ffb7f56cfad1c3ceec329.UsersRequestBuilder) {
    return if6ffd1464db2d9c22e351b03e4c00ebd24a5353cd70ffb7f56cfad1c3ceec329.NewUsersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
