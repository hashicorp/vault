import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import s from './style.module.css'

interface IoHomeFeatureProps {
  link: string
  // thumnail: string
  heading: string
  description: string
}

export default function IoHomeFeature({
  link,
  // thumbnail,
  heading,
  description,
}: IoHomeFeatureProps) {
  return (
    <a href={link} className={s.feature}>
      <div className={s.featureMedia}>
        {/* <Image src={thumbnail} width={400} height={200} alt="" /> */}
      </div>
      <div className={s.featureContent}>
        <h3 className={s.featureHeading}>{heading}</h3>
        <p className={s.featureDescription}>{description}</p>
        <span className={s.featureCta} aria-hidden={true}>
          Learn more{' '}
          <span>
            <IconArrowRight16 />
          </span>
        </span>
      </div>
    </a>
  )
}
