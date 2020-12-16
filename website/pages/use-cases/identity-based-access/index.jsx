import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextSplitWithLogoGrid from '@hashicorp/react-text-split-with-logo-grid'
import TextSplitWithCode from '@hashicorp/react-text-split-with-code'
import TextSplitWithImage from '@hashicorp/react-text-split-with-image'
import BeforeAfterDiagram from 'components/before-after-diagram'
import UseCaseCtaSection from 'components/use-case-cta-section'

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
      <section className="no-spacing">
        <div className="g-grid-container">
          <SectionHeader headline="Identity-based Access Features" />
        </div>

        <TextSplitWithLogoGrid
          textSplit={{
            heading: 'Identity Plugins',
            content:
              'Improve the extensibility of Vault with pluggable identity backends',
          }}
          logoGrid={[
            'aws',
            'microsoft-azure',
            'google',
            'kubernetes',
            { url: require('./img/logos/nomad.png?url') },
            { url: require('./img/logos/okta.png?url') },
            { url: require('./img/logos/pivotalcf.png?url') },
            { url: require('./img/logos/ssh.png?url') },
            {
              url:
                'https://www.datocms-assets.com/2885/1608143270-ellipsis.png',
            },
          ]}
        />

        <TextSplitWithImage
          textSplit={{
            heading: 'Entities',
            content:
              'Integrated identities across platforms and using this information for policy and access control decisions.',
            textSide: 'right',
          }}
          image={{
            url: require('./img/screenshot-entities.png?url'),
          }}
        />

        <TextSplitWithImage
          textSplit={{
            heading: 'Control Groups',
            content:
              'Require multiple Identity Entities or members of Identity Groups to authorize an requested action.',
          }}
          image={{
            url: require('./img/screenshot-control-groups.png?url'),
          }}
        />

        <TextSplitWithCode
          textSplit={{
            heading: 'ACL Templates and Policy Control',
            content:
              'Create and manage policies that authorize access control throughout your infrastructure and organization',
          }}
          codeBlock={{
            options: { showWindowBar: true },
            code:
              '# User template (user-tmpl.hcl)\n# Grant permissions on user specific path\npath "user-kv/data/{{identity.entity.name}}/*" {\n  capabilities = [ "create", "update", "read", "delete", "list" ]\n}\n# For Web UI usage\npath "user-kv/metadata" {\n  capabilities = ["list"]\n}\n# Group template (group-tmpl.hcl)\n# Grant permissions on the group specific path\n# The region is specified in the group metadata\npath "group-kv/data/education/{{identity.groups.names.education.metadata.region}}/*" {\n  capabilities = [ "create", "update", "read", "delete", "list" ]\n}\n# Group member can update the group information\npath "identity/group/id/{{identity.groups.names.education.id}}" {\n  capabilities = [ "update", "read" ]\n}\n# For Web UI usage\npath "group-kv/metadata" {\n  capabilities = ["list"]\n}\npath "identity/group/id" {\n  capabilities = [ "list" ]\n}\n',
          }}
        />

        <TextSplitWithImage
          textSplit={{
            heading: 'Identity Groups',
            content:
              'Group trusted identities into logical groups for group-based access control.',
            textSide: 'right',
          }}
          image={{
            url: require('./img/screenshot-identity-groups.png?url'),
          }}
        />

        <TextSplitWithCode
          textSplit={{
            heading: 'Multi-factor Authentication',
            content:
              'Enforce MFA workflows when accessing a secret or a secret path',
          }}
          codeBlock={{
            options: { showWindowBar: true },
            code:
              '$ curl --header "X-Vault-Token: ..." \\\n--header "X-Vault-MFA:my_totp:695452" \\\nhttp://127.0.0.1:8200/v1/secret/foo',
          }}
        />
      </section>

      <UseCaseCtaSection />
    </div>
  )
}
