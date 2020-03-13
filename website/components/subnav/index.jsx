import Subnav from '@hashicorp/react-subnav'
import menuItems from '../../data/subnav'
import { useRouter } from 'next/router'

export default function ProductSubnav() {
  const router = useRouter()

  return (
    <Subnav
      titleLink={{
        text: 'vault',
        url: '/'
      }}
      ctaLinks={[
        {
          text: 'GitHub',
          url: 'https://www.github.com/hashicorp/vault'
        },
        {
          text: 'Download',
          url: '/downloads'
        }
      ]}
      currentPath={router.pathname}
      menuItems={menuItems}
      menuItemsAlign="right"
    />
  )
}
