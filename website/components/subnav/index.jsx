import Subnav from '@hashicorp/react-subnav'
import menuItems from 'data/subnav'
import { useRouter } from 'next/router'

export default function ProductSubnav() {
  const router = useRouter()

  return (
    <Subnav
      className="g-product-subnav"
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
            'https://portal.cloud.hashicorp.com/sign-up?utm_source=vault_io&utm_content=top_nav_vault',
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
      matchOnBasePath
    />
  )
}
