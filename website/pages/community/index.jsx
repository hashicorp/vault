import VerticalTextBlockList from '@hashicorp/react-vertical-text-block-list'
import SectionHeader from '@hashicorp/react-section-header'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'

function CommunityPage() {
  return (
    <div id="community">
      <HashiHead is={Head} title="Community | Vault by HashiCorp" />
      <SectionHeader
        headline="Community"
        description="Vault is an open source project with a growing community. There are active, dedicated users willing to help you through various mediums."
        use_h1={true}
      />
      <VerticalTextBlockList
        data={[
          {
            header: 'Discussion List',
            body:
              '[Vault Community Forum](https://discuss.hashicorp.com/c/vault)'
          },
          {
            header: 'Bug Tracker',
            body:
              '[Issue tracker on GitHub](https://github.com/hashicorp/vault/issues) for reporting bugs. Use IRC or the mailing list for general help.'
          },
          {
            header: 'Training',
            body:
              '[Paid HashiCorp](https://www.hashicorp.com/training) training courses are available in a city near you. Private training courses are also available.'
          }
        ]}
      />
    </div>
  )
}

export default CommunityPage
