import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextAndContent from '@hashicorp/react-text-and-content'
import BeforeAfterDiagram from '../../../components/before-after-diagram'
import UseCaseCtaSection from '../../../components/use-case-cta-section'

export default function SecretsManagmentUseCase() {
  return (
    <div id="use-cases" className="g-section-block page-wrap">
      <section className="g-container">
        {/* Header / Buttons */}

        <SectionHeader
          headline="Secrets Management in Low Trust Networks"
          description="Centrally store, access, and deploy secrets across applications, systems, and infrastructure"
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
              url:
                'https://www.datocms-assets.com/2885/1539885048-secrets-managementchallenge.svg',
              format: 'svg'
            }}
            beforeHeadline="The Challenge"
            beforeContent="Secrets for applications and systems need to be centralized and static IP-based solutions don't scale in dynamic environments with frequently changing applications and machines"
            afterImage={{
              url:
                'https://www.datocms-assets.com/2885/1539885054-secrets-managementsolution.svg',
              format: 'svg'
            }}
            afterHeadline="The Solution"
            afterContent="Vault centrally manages and enforces access to secrets and systems based on trusted sources of application and user identity"
          />
        </div>
      </section>

      {/* Case study slider */}

      <section className="g-section-block theme-black-background-white-text">
        <div className="g-container">
          <div className="g-case-study-slider">
            <div className="case-study-container">
              <div className="slider-container">
                <div className="slider-frame single">
                  <div className="case-study">
                    <div className="feature-image">
                      <a href="https://www.hashicorp.com/resources/adobe-100-trillion-transactions-hashicorp-vault">
                        <picture>
                          <source
                            type="image/webp"
                            srcSet="https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=156.25&amp;q=80&amp;w=250 250w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=312.5&amp;q=80&amp;w=500 500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=468.75&amp;q=80&amp;w=750 750w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=625&amp;q=80&amp;w=1000 1000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=937.5&amp;q=80&amp;w=1500 1500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=1250&amp;q=80&amp;w=2000 2000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=1562.5&amp;q=80&amp;w=2500 2500w"
                            sizes="100vw"
                          />
                          <img
                            src="https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=312.5&amp;q=80&amp;w=500"
                            srcSet="https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=156.25&amp;q=80&amp;w=250 250w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=312.5&amp;q=80&amp;w=500 500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=468.75&amp;q=80&amp;w=750 750w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=625&amp;q=80&amp;w=1000 1000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=937.5&amp;q=80&amp;w=1500 1500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=1250&amp;q=80&amp;w=2000 2000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=1562.5&amp;q=80&amp;w=2500 2500w"
                            sizes="100vw"
                            alt="Dan McTeer at HashiConf 2017"
                          />
                        </picture>
                      </a>
                    </div>
                    <div className="feature-content">
                      <div className="single-logo">
                        <img
                          src="https://www.datocms-assets.com/2885/1539889072-1524097013-adobe-white-1.svg"
                          alt="Adobe logo"
                        />
                      </div>
                      <h3>
                        Using Vault to Protect Adobe&apos;s Secrets and User
                        Data Across Clouds and Datacenters
                      </h3>
                      <p>
                        Securing secrets and application data is a complex task
                        for globally distributed organizations. For Adobe,
                        managing secrets for over 20 products across 100,000
                        hosts, four regions, and trillions of transactions
                        annually requires a different approach altogether.
                      </p>
                      <a
                        className="g-btn primary-hashicorp-light"
                        href="https://www.hashicorp.com/resources/adobe-100-trillion-transactions-hashicorp-vault"
                      >
                        Read Case Study
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features / Text and content */}
      <section className="g-container">
        <SectionHeader headline="Secret Management Features" />

        <TextAndContent
          data={{
            text: `### Dynamic Secrets

Dynamically create, revoke, and rotate secrets programmatically`,
            content: {
              __typename: 'SbcImageRecord',
              image: {
                url:
                  'https://www.datocms-assets.com/2885/1538684923-dynamic-secrets-screenshot.jpg',
                format: 'jpg',
                size: 71545,
                width: 668,
                height: 504
              }
            }
          }}
        />

        <TextAndContent
          data={{
            reverseDirection: true,
            text: `### Secret Storage

Encrypt data while at rest, in the storage backend of your choice`,
            content: {
              __typename: 'SbcCodeBlockRecord',
              chrome: true,
              code: `$ cat vault.config
storage "consul" {
    address = "127.0.0.1:8500"
    path    = "vault"
}
listener "tcp" {
    address = "127.0.0.1:8200"
}
telemetry {
    statsite_address = "127.0.0.1:8125"
    disable_hostname = true
}`
            }
          }}
        />

        <div className="g-text-and-content">
          <div className="text">
            <div>
              <h3>Identity Plugins</h3>
              <p>
                Improve the extensibility of Vault with pluggable identity
                backends
              </p>
            </div>
          </div>
          <div className="content logo-grid">
            <ul className="g-logo-grid large">
              {[
                'https://www.datocms-assets.com/2885/1506540090-color.svg',
                'https://www.datocms-assets.com/2885/1506540114-color.svg',
                'https://www.datocms-assets.com/2885/1506527176-color.svg',
                'https://www.datocms-assets.com/2885/1508434209-consul_primarylogo_fullcolor.svg',
                'https://www.datocms-assets.com/2885/1510033601-aws_logo_rgb_fullcolor.svg',
                'https://www.datocms-assets.com/2885/1506540175-color.svg',
                'https://www.datocms-assets.com/2885/1539818112-postgresql.svg',
                'https://www.datocms-assets.com/2885/1539817686-microsoft-sql-server.svg'
              ].map(logo => (
                <li key={logo}>
                  <img src={logo} alt="company logo" />
                </li>
              ))}
            </ul>
          </div>
        </div>

        <TextAndContent
          data={{
            reverseDirection: true,
            text: `### Detailed Audit Logs

Detailed audit log of all client interaction (authentication, token creation, secret access & revocation)`,
            content: {
              __typename: 'SbcCodeBlockRecord',
              chrome: true,
              code: `$ cat audit.log | jq {
    "time": "2018-08-27T13:17:11.609621226Z",
    "type": "response",
    "auth": {
        "client_token": "hmac-sha256:5c40f1e051ea75b83230a5bf16574090f697dfa22a78e437f12c1c9d226f45a5",
        "accessor": "hmac-sha256:f254a2d442f172f0b761c9fd028f599ad91861ed16ac3a1e8d96771fd920e862",
        "display_name": "token",
        "metadata": null,
        "entity_id": ""
    }
}`
            }
          }}
        />

        <TextAndContent
          data={{
            text: `### Leasing & Revoking Secrets

Manage authorization and create time-based tokens for automatic revocation or manual revocation`,
            content: {
              __typename: 'SbcCodeBlockRecord',
              chrome: true,
              code: `$ vault read database/creds/readonly
Key             Value
---             -----
lease_id        database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
lease_duration  1h0m0s
lease_renewable true
password        A1a-w2xv2zsq4r5ru940
username        v-token-readonly-48rt0t36sxp4wy81x8x1-1515627434
[...]
$ vault renew database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
Key             Value
---             -----
lease_id        database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
lease_duration  1h0m0s
lease_renewable true
$ vault lease revoke database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809`
            }
          }}
        />
      </section>

      <UseCaseCtaSection />
    </div>
  )
}
