import { motion } from 'framer-motion'
import { FaCrown } from 'react-icons/fa6'
import { AiOutlineLoading3Quarters } from 'react-icons/ai'

export default function PlayerName({
  isWinner,
  playerNumber,
  name,
  isTurn,
  offlineTimeStamp,
  countDownDuration,
}) {
  console.log('player name called')
  //online statuses
  const isOnline = offlineTimeStamp === null
  const inGame = name !== null
  const isCountingDown = inGame && !isOnline
  let timeRemaining = null
  if (isCountingDown) {
    const milliseconds =
      countDownDuration * 1000 - (Date.now() - offlineTimeStamp)
    timeRemaining = Math.floor(milliseconds / 1000)
  }

  //animation
  const variants = {
    floating: {
      y: [2, -2],
      transition: {
        duration: 0.5,
        repeat: Infinity,
        repeatType: 'mirror',
        ease: 'easeInOut',
      },
    },
    stale: {},
  }

  const currentVariant = isTurn ? 'floating' : 'stale'
  return (
    <div
      style={{
        gridColumn: 'span-7',
      }}
    >
      <motion.p
        variants={variants}
        animate={currentVariant}
        className='flex gap-3 items-center'
      >
        {isWinner ? <FaCrown /> : ''}

        {isOnline && <span> {name}</span>}

        {isCountingDown && (
          <>
            <motion.span
              animate={{ rotate: 360 }}
              transition={{ repeat: Infinity, duration: 1, ease: 'linear' }}
            >
              <AiOutlineLoading3Quarters />
            </motion.span>
            Waiting ... {timeRemaining}s
          </>
        )}
      </motion.p>
    </div>
  )
}
