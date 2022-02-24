import * as React from 'react'
import { DialogOverlay, DialogContent, DialogOverlayProps } from '@reach/dialog'
import { AnimatePresence, motion } from 'framer-motion'
import s from './style.module.css'

export interface IoDialogProps extends DialogOverlayProps {
  label: string
}

export default function IoDialog({
  isOpen,
  onDismiss,
  children,
  label,
}: IoDialogProps): React.ReactElement {
  const AnimatedDialogOverlay = motion(DialogOverlay)
  return (
    <AnimatePresence>
      {isOpen && (
        <AnimatedDialogOverlay
          className={s.dialogOverlay}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onDismiss={onDismiss}
        >
          <div className={s.dialogWrapper}>
            <motion.div
              initial={{ y: 50 }}
              animate={{ y: 0 }}
              exit={{ y: 50 }}
              transition={{ min: 0, max: 100, bounceDamping: 8 }}
              style={{ width: '100%', maxWidth: 800 }}
            >
              <DialogContent className={s.dialogContent} aria-label={label}>
                <button onClick={onDismiss} className={s.dialogClose}>
                  Close
                </button>
                {children}
              </DialogContent>
            </motion.div>
          </div>
        </AnimatedDialogOverlay>
      )}
    </AnimatePresence>
  )
}
