import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextAndContent from '@hashicorp/react-text-and-content'
import BeforeAfterDiagram from '../../../components/before-after-diagram'
import UseCaseCtaSection from '../../../components/use-case-cta-section'

export default function DataEncryptionUseCase() {
  return (
    <div id="use-cases" className="g-section-block page-wrap">
      {/* Header / Buttons */}
      <section className="g-container">
        <SectionHeader
          headline="Leverage Trusted Identities in Low Trust Networks"
          description="Authenticate and access different clouds, systems, and endpoints using trusted identities"
          useH1={true}
        />

        <div className="button-container">
          <Button title="Download" url="/downloads" />
          <Button title="Get Started" url="/intro" theme="dark-outline" />
        </div>
      </section>

      {/* Before/After Diagram */}
      <section>
        <div className="g-container">
          <BeforeAfterDiagram
            beforeImage={{
              url: require('./img/challenge.png'),
              format: 'png',
            }}
            beforeHeadline="The Challenge"
            beforeContent="With the proliferation of different clouds, services, and systems all with their own identity providers, organizations need a way to manage identity sprawl"
            afterImage={{
              url: require('./img/solution.png'),
              format: 'png',
            }}
            afterHeadline="The Solution"
            afterContent="Vault merges identities across providers and uses a unified ACL system to broker access to systems and secrets"
          />
        </div>
      </section>

      {/* Features / Text and content */}
      <section className="g-container">
        <SectionHeader headline="Identity-based Access Features" />

        <div className="g-text-and-content">
          <div className="text">
            <div>
              <h3 id="secure-plugins">Identity Plugins</h3>
              <p>
                Improve the extensibility of Vault with pluggable identity
                backends
              </p>
            </div>
          </div>
          <div className="content logo-grid">
            <ul className="g-logo-grid large">
              <li key="aws">
                <img src={require('./img/logos/aws.png')} alt="company logo" />
              </li>
              <li key="azure">
                <img
                  src={require('./img/logos/azure.png')}
                  alt="company logo"
                />
              </li>
              <li key="gcp">
                <img src={require('./img/logos/gcp.png')} alt="company logo" />
              </li>
              <li key="kubernetes">
                <img
                  src={require('./img/logos/kubernetes.png')}
                  alt="company logo"
                />
              </li>
              <li key="nomad">
                <img
                  src={require('./img/logos/nomad.png')}
                  alt="company logo"
                />
              </li>
              <li key="okta">
                <img src={require('./img/logos/okta.png')} alt="company logo" />
              </li>
              <li key="pivotalcf">
                <img
                  src={require('./img/logos/pivotalcf.png')}
                  alt="company logo"
                />
              </li>
              <li key="ssh">
                <img src={require('./img/logos/ssh.png')} alt="company logo" />
              </li>
            </ul>
          </div>
        </div>

        <TextAndContent
          data={{
            reverseDirection: true,
            text: `### Entities

Integrated identities across platforms and using this information for policy and access control decisions.`,
            content: {
              __typename: 'SbcImageRecord',
              image: {
                url: require('./img/screenshot-entities.png'),
                format: 'png',
              },
            },
          }}
        />

        <TextAndContent
          data={{
            text: `### Control Groups

Require multiple Identity Entities or members of Identity Groups to authorize an requested action.`,
            content: {
              __typename: 'SbcImageRecord',
              image: {
                url: require('./img/screenshot-control-groups.png'),
                format: 'png',
              },
            },
          }}
        />

        <TextAndContent
          data={{
            text: `### ACL Templates and Policy Control

Create and manage policies that authorize access control throughout your infrastructure and organization`,
            content: {
              __typename: 'SbcCodeBlockRecord',
              chrome: true,
              code: `# User template (user-tmpl.hcl)
# Grant permissions on user specific path

path "user-kv/data/{{identity.entity.name}}/*" {
  capabilities = [ "create", "update", "read", "delete", "list" ]
}

# For Web UI usage
path "user-kv/metadata" {
  capabilities = ["list"]
}

# Group template (group-tmpl.hcl)
# Grant permissions on the group specific path
# The region is specified in the group metadata
path "group-kv/data/education/{{identity.groups.names.education.metadata.region}}/*" {
  capabilities = [ "create", "update", "read", "delete", "list" ]
}

# Group member can update the group information
path "identity/group/id/{{identity.groups.names.education.id}}" {
  capabilities = [ "update", "read" ]
}

# For Web UI usage
path "group-kv/metadata" {
  capabilities = ["list"]
}

path "identity/group/id" {
  capabilities = [ "list" ]
}`,
            },
          }}
        />

        <TextAndContent
          data={{
            reverseDirection: true,
            text: `### Identity Groups

Group trusted identities into logical groups for group-based access control.`,
            content: {
              __typename: 'SbcImageRecord',
              image: {
                url: require('./img/screenshot-identity-groups.png'),
                format: 'png',
              },
            },
          }}
        />

        <TextAndContent
          data={{
            text: `### Multi-factor Authentication

Enforce MFA workflows when accessing a secret or a secret path`,
            content: {
              __typename: 'SbcCodeBlockRecord',
              chrome: true,
              code: `$ curl \
--header "X-Vault-Token: ..." \\
--header "X-Vault-MFA:my_totp:695452" \\
http://127.0.0.1:8200/v1/secret/foo`,
            },
          }}
        />
      </section>

      <UseCaseCtaSection />
    </div>
  )
}
