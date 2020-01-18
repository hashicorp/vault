import React, { Component } from 'react'
import Link from 'next/link'
import LinkWrap from '@hashicorp/react-link-wrap'
import InlineSvg from '@hashicorp/react-inline-svg'
import caratIcon from './img/carat.svg?include'
import downloadIcon from './img/download.svg?include'
import githubIcon from './img/github.svg?include'

class ProductSubnav extends Component {
  constructor(props) {
    super(props)

    this.state = {
      fixNav: true,
      menuOpen: false
    }

    this.toggleMenu = this.toggleMenu.bind(this)
    this.openSubmenu = this.openSubmenu.bind(this)
    this.closeSubmenu = this.closeSubmenu.bind(this)
  }

  toggleMenu() {
    this.setState({ menuOpen: !this.state.menuOpen })
  }

  openSubmenu(e) {
    const dropdown = e.currentTarget.nextElementSibling
    dropdown && dropdown.classList.add('open')
  }

  closeSubmenu(e) {
    e.currentTarget.parentElement.classList.remove('open')
  }

  render() {
    const {
      rootUrl,
      buttonText,
      hideDownloadIcon,
      buttonExternal,
      links
    } = this.props
    const currentProduct = {
      slug: 'vault',
      siteUrl: 'https://www.vaultproject.io',
      downloadUrl: '/downloads',
      githubUrl: 'https://github.com/hashicorp/vault',
      colorLogo: {
        url: 'https://www.datocms-assets.com/2885/1509990185-vault-dark.svg',
        alt: 'Vault Logo'
      }
    }

    return (
      <div className={`g-product-subnav ${currentProduct.slug}`}>
        <div className={`wrapper${this.state.fixed ? ' fixed' : ''}`}>
          <div className="g-container">
            <div className="logo-wrapper">
              <div className="logo">
                <Link href="/">
                  <a
                    data-ga-product-subnav={`Product logo | ${currentProduct.slug}`}
                  >
                    <img
                      src={currentProduct.colorLogo.url}
                      alt={currentProduct.colorLogo.alt}
                    />
                    <h4>
                      {currentProduct.slug.charAt(0).toUpperCase() +
                        currentProduct.slug.slice(1)}
                    </h4>
                  </a>
                </Link>
              </div>
            </div>
            <div
              className={`toggle${currentProduct.downloadUrl ? ' center' : ''}${
                this.state.menuOpen ? ' open' : ''
              }`}
            >
              <span onClick={this.toggleMenu}>
                <InlineSvg src={caratIcon} />
              </span>
            </div>
            <div
              className={`links-wrapper${this.state.menuOpen ? ' open' : ''}`}
            >
              <div className="links g-type-body-small-strong">
                <ul>
                  {links.map(link => {
                    if (link.title === '|')
                      return <li className="divider" key="divider" />
                    return link.links ? (
                      <li className="dropdown" key={link.title}>
                        <span onClick={this.openSubmenu}>
                          {link.title}
                          <InlineSvg src={caratIcon} />
                        </span>
                        <ul>
                          <li className="close" onClick={this.closeSubmenu}>
                            <span>
                              <InlineSvg src={caratIcon} /> Back
                            </span>
                          </li>
                          <li className="name">{link.title}</li>
                          {link.links.map(sublink => (
                            <li key={sublink.title}>
                              <LinkWrap
                                Link={Link}
                                href={resolveLocalUrl(rootUrl, sublink.url)}
                                {...(sublink.external && {
                                  rel: 'noopener',
                                  target: '_blank'
                                })}
                                data-ga-product-subnav={`Subnav sublink | ${sublink.title}`}
                              >
                                {sublink.title}
                              </LinkWrap>
                            </li>
                          ))}
                        </ul>
                      </li>
                    ) : (
                      <li key={link.url}>
                        <LinkWrap
                          Link={Link}
                          href={resolveLocalUrl(rootUrl, link.url)}
                          onClick={this.toggleMenu}
                          {...(link.external && {
                            rel: 'noopener',
                            target: '_blank'
                          })}
                          data-ga-product-subnav={`Subnav link | ${link.title}`}
                        >
                          {link.title}
                        </LinkWrap>
                      </li>
                    )
                  })}
                </ul>
                <div className="ctas-mobile">
                  {currentProduct.githubUrl && (
                    <a className="github g-btn" href={currentProduct.githubUrl}>
                      <InlineSvg src={githubIcon} /> GitHub
                    </a>
                  )}
                  {currentProduct.downloadUrl && (
                    <LinkWrap
                      Link={Link}
                      className="download g-btn variant-primary brand-vault background-light"
                      href={resolveLocalUrl(
                        rootUrl,
                        currentProduct.downloadUrl
                      )}
                      data-ga-product-subnav={`Download mobile | ${currentProduct.slug}`}
                      {...(buttonExternal && {
                        rel: 'noopener',
                        target: '_blank'
                      })}
                    >
                      {!hideDownloadIcon && <InlineSvg src={downloadIcon} />}
                      <span>{buttonText || 'Download'}</span>
                    </LinkWrap>
                  )}
                </div>
              </div>
            </div>
            <div className="ctas">
              {currentProduct.githubUrl && (
                <a className="github" href={currentProduct.githubUrl}>
                  <InlineSvg src={githubIcon} />
                </a>
              )}
              {currentProduct.downloadUrl && (
                <LinkWrap
                  Link={Link}
                  className="download g-btn variant-primary brand-vault background-light"
                  href={resolveLocalUrl(rootUrl, currentProduct.downloadUrl)}
                  data-ga-product-subnav={`Download desktop | ${currentProduct.slug}`}
                  {...(buttonExternal && { rel: 'noopener', target: '_blank' })}
                >
                  {!hideDownloadIcon && <InlineSvg src={downloadIcon} />}
                  <span>{buttonText || 'Download'}</span>
                </LinkWrap>
              )}
            </div>
          </div>
        </div>
      </div>
    )
  }
}

// This function ensures that if the root url is the same as the link's url,
// we strip the root and resolve from a relative path, so as to not redirect
// from the local dev environment to production unexpectedly.
function resolveLocalUrl(rootUrl, url) {
  const urlRegex = new RegExp(`^(?:http(?:s)*://)?(?:www.)?${rootUrl}`)
  return url.replace(/\.html$/, '').replace(urlRegex, '')
}

export default ProductSubnav
