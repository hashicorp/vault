import Subnav from '@hashicorp/react-subnav'
import menuItems from 'data/subnav'
import { useRouter } from 'next/router'

export default function ProductSubnav() {
  const router = useRouter()

  return (
    <Subnav
      hideGithubStars={true}
      titleLink={{
        text: 'vault',
        url: '/',
      }}
      ctaLinks={[
        {
          text: 'GitHub',
          url: 'https://www.github.com/hashicorp/vault',
        },
        {
          text: 'Try Cloud',
          url:
            'https://cloud.hashicorp.com/?utm_source=vault_io&utm_content=top_nav',
        },
        {
          text: 'Download',
          url: '/downloads',
        },
      ]}
      currentPath={router.asPath}
      menuItems={menuItems}
      menuItemsAlign="right"
      constrainWidth
    />
  )
}
