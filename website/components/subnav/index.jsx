import Subnav from '@hashicorp/react-subnav'
import classNames from 'classnames'
import { useRouter } from 'next/router'
import s from './style.module.css'

export default function ProductSubnav({ menuItems }) {
  const router = useRouter()

  return (
    <Subnav
      className={classNames('g-product-subnav', s.subnav)}
      hideGithubStars={true}
      titleLink={{
        text: 'HashiCorp Vault',
        url: '/',
      }}
      ctaLinks={[
        {
          text: 'GitHub',
          url: 'https://www.github.com/hashicorp/vault',
        },
        {
          text: 'Try Cloud',
          url: 'https://portal.cloud.hashicorp.com/sign-up?utm_source=vault_io&utm_content=top_nav_vault',
        },
        {
          text: 'Download',
          url: '/downloads',
          theme: {
            brand: 'vault',
          },
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
