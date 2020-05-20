import Link from 'next/link'
import { useEffect } from 'react'

function FourOhFour() {
  useEffect(() => {
    /* eslint-disable no-undef */
    if (
      typeof globalThis?.analytics?.track === 'function' &&
      typeof globalThis?.document?.referrer === 'string' &&
      typeof globalThis?.location?.href === 'string'
    )
      globalThis.analytics.track({
        event: '404 Response',
        action: globalThis.location.href,
        label: globalThis.document.referrer
      })
    /* eslint-enable no-undef */
  }, [])

  return (
    <header id="p-404">
      <h1>Page Not Found</h1>
      <p>
        We&apos;re sorry but we can&apos;t find the page you&apos;re looking
        for.
      </p>
      <p>
        <Link href="/">
          <a>Back to Home</a>
        </Link>
      </p>
    </header>
  )
}

export default FourOhFour
