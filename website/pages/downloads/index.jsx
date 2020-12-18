import { VERSION } from '../../data/version'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'
import ProductDownloader from '@hashicorp/react-product-downloader'
import styles from './style.module.css'
import logo from '@hashicorp/mktg-assets/dist/product/vault-logo/color.svg'

export default function DownloadsPage({ releases }) {
  return (
    <>
      <HashiHead is={Head} title={`Downloads | Vault by HashiCorp`} />

      <ProductDownloader
        releases={releases}
        packageManagers={[
          {
            label: 'Homebrew',
            commands: [
              'brew tap hashicorp/tap',
              'brew install hashicorp/tap/vault',
            ],
            os: 'darwin',
          },
          {
            label: 'Ubuntu/Debian',
            commands: [
              'curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -',
              'sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"',
              'sudo apt-get update && sudo apt-get install vault',
            ],
            os: 'linux',
          },
          {
            label: 'CentOS/RHEL',
            commands: [
              'sudo yum install -y yum-utils',
              'sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo',
              'sudo yum -y install vault',
            ],
            os: 'linux',
          },
          {
            label: 'Fedora',
            commands: [
              'sudo dnf install -y dnf-plugins-core',
              'sudo dnf config-manager --add-repo https://rpm.releases.hashicorp.com/fedora/hashicorp.repo',
              'sudo dnf -y install vault',
            ],
            os: 'linux',
          },
          {
            label: 'Amazon Linux',
            commands: [
              'sudo yum install -y yum-utils',
              'sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo',
              'sudo yum -y install vault',
            ],
            os: 'linux',
          },
        ]}
        productName="Vault"
        productId="vault"
        latestVersion={VERSION}
        getStartedDescription="Follow step-by-step tutorials to get hands on with Vault."
        getStartedLinks={[
          {
            label: 'Get Started with CLI',
            href:
              'https://learn.hashicorp.com/collections/vault/getting-started',
          },
          {
            label: 'Get Started with UI',
            href:
              'https://learn.hashicorp.com/collections/vault/getting-started-ui',
          },
          {
            label: 'Vault on Kubernetes',
            href: 'https://learn.hashicorp.com/collections/vault/kubernetes',
          },
        ]}
        logo={<img className={styles.logo} alt="Vault" src={logo} />}
        brand="vault"
        tutorialLink={{
          href: 'https://learn.hashicorp.com/vault',
          label: 'View Tutorials at HashiCorp Learn',
        }}
      />
    </>
  )
}

export async function getStaticProps() {
  return fetch(`https://releases.hashicorp.com/vault/index.json`, {
    headers: {
      'Cache-Control': 'no-cache',
    },
  })
    .then((res) => res.json())
    .then((result) => {
      return {
        props: {
          releases: result,
        },
      }
    })
    .catch(() => {
      throw new Error(
        `--------------------------------------------------------
        Unable to resolve version ${VERSION} on releases.hashicorp.com from link
        <https://releases.hashicorp.com/vault/${VERSION}/index.json>. Usually this
        means that the specified version has not yet been released. The downloads page
        version can only be updated after the new version has been released, to ensure
        that it works for all users.
        ----------------------------------------------------------`
      )
    })
}
