import Image from 'next/image'
import VisuallyHidden from '@reach/visually-hidden'
import PlayIcon from './play-icon'
import s from './style.module.css'

export default function Video({ thumbnail, heading, description, person }) {
  return (
    <figure className={s.video}>
      <div className={s.thumbnail}>
        <button>
          <VisuallyHidden>Play video</VisuallyHidden>
          <PlayIcon />
        </button>
        <Image src={thumbnail} layout="fill" objectFit="cover" />
      </div>
      <figcaption className={s.content}>
        <h3 className={s.heading}>{heading}</h3>
        <p className={s.description}>{description}</p>
        <div className={s.person}>
          <div className={s.personThumbnail}>
            <Image
              src={person.thumbnail}
              width={52}
              height={52}
              alt={`${person.name} avatar`}
            />
          </div>
          <div>
            <p className={s.personName}>{person.name}</p>
            <p className={s.personDescription}>{person.description}</p>
          </div>
        </div>
      </figcaption>
    </figure>
  )
}
