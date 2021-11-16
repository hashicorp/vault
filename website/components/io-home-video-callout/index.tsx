import * as React from 'react'
import Image from 'next/image'
import VisuallyHidden from '@reach/visually-hidden'
import Dialog from 'components/dialog'
import PlayIcon from './play-icon'
import s from './style.module.css'

export default function IoHomeVideoCallout({
  thumbnail,
  heading,
  description,
  person,
}) {
  const [showDialog, setShowDialog] = React.useState(false)
  const showVideo = () => setShowDialog(true)
  const hideVideo = () => setShowDialog(false)
  return (
    <>
      <figure className={s.videoCallout}>
        <button className={s.thumbnail} onClick={showVideo}>
          <VisuallyHidden>Play video</VisuallyHidden>
          <PlayIcon />
          <Image src={thumbnail} layout="fill" objectFit="cover" />
        </button>
        <figcaption className={s.content}>
          <h3 className={s.heading}>{heading}</h3>
          <p className={s.description}>{description}</p>
          {person && (
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
          )}
        </figcaption>
      </figure>
      <Dialog
        isOpen={showDialog}
        onDismiss={hideVideo}
        label={`${heading} video}`}
      >
        <div className={s.video}></div>
      </Dialog>
    </>
  )
}
