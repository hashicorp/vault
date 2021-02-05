import styles from './HCPCalloutSection.module.css'
import Button from '@hashicorp/react-button'

export default function HcpCalloutSection({
  id,
  header,
  title,
  description,
  chin,
  image,
  links,
}) {
  return (
    <div className={styles.hcpCalloutSection} id={id}>
      <div className={styles.header}>
        <h2>{header}</h2>
      </div>

      <div className={styles.content}>
        <div className={styles.info}>
          <h1>{title}</h1>
          <span className={styles.chin}>{chin}</span>
          <p className={styles.description}>{description}</p>
          <div className={styles.links}>
            {links.map((link, index) => {
              const variant = index === 0 ? 'primary' : 'tertiary'
              return (
                <div key={link.text}>
                  <Button
                    title={link.text}
                    linkType={link.type}
                    url={link.url}
                    theme={{ variant, brand: 'neutral', background: 'light' }}
                  />
                </div>
              )
            })}
          </div>
        </div>
        <img alt={title} src={image} />
      </div>
    </div>
  )
}
