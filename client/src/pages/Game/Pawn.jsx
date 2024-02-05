import { getColor } from '../../utils'
import { motion } from 'framer-motion'

export default function Pawn({ pawnNumber, shouldFloat, onClick, top, left }) {
  // const showNumber = true
  const showNumber = false

  const color = getColor({ pawnNumber })


  return (
    <>
      <motion.button
        variants={{
          floating: {
            x: ['-50%', '-50%'],
            y: ['-30%', '-70%'],
            transition: {
              duration: 0.5,
              repeat: Infinity,
              repeatType: 'mirror',
              ease: 'easeInOut',
            },
          },
          stale: {
            x: '-50%',
            y: '-50%',
          },
        }}
        animate={shouldFloat ? 'floating' : 'stale'}
        onClick={() => onClick(pawnNumber)}
        className={` absolute w-[17px] h-[17px] rounded-full  flex justify-center items-center `}
        style={{
          boxShadow: '0px 2px 2px 1px rgba(0,0,0,0.2)',
          background: color,
          transition: 'all 150ms',
          top: `${top}px`,
          left: `${left}px`,
        }}
      >
        {showNumber && <p>{pawnNumber}</p>}
      </motion.button>
    </>
  )
}
