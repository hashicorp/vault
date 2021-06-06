import Link from 'next/link'

export default function Footer({ openConsentManager }) {
  return (
    <footer className="g-footer">
      <div className="g-container">
        <div className="left">
          <Link href="/docs">
            <a>Docs</a>
          </Link>
          <Link href="/api">
            <a>API</a>
          </Link>
          <a href="https://learn.hashicorp.com/vault">Learn</a>
          <Link href="/community">
            <a>Community</a>
          </Link>
          <a href="https://hashicorp.com/privacy">Privacy</a>
          <Link href="/security">
            <a>Security</a>
          </Link>
          <Link href="/files/press-kit.zip">
            <a>Press Kit</a>
          </Link>
          <a onClick={openConsentManager}>Consent Manager</a>
        </div>
      </div>
    </footer>
  )
}
