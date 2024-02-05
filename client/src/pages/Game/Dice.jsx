import { useEffect, useState } from 'react'
import { motion } from 'framer-motion'

export default function Dice({ value, shouldFloat, onClick, isRolling }) {
  const dots = []
  for (let i = 0; i < value; i++) {
    dots.push(i)
  }

  const variants = {
    rolling: {
      rotate: [0, 400, -45],
      transition: {
        duration: 1,
      },
    },

    floating: {
      y: [-3, 3],
      transition: {
        duration: 0.5,
        repeat: Infinity,
        repeatType: 'mirror',
        ease: 'easeInOut',
      },
    },

    stale: {},
  }

  let variantName = ''
  if (isRolling) variantName = 'rolling'
  else if (shouldFloat) variantName = 'floating'
  else variantName = 'stale'

  return (
    <motion.div
      variants={variants}
      animate={variantName}
      className={` flex flex-col gap-3 `}
    >
      <button
        className={`
        bg-[lightgrey] w-[80px] h-[80px] p-3 gap-3 border rounded
        flex flex-wrap items-center justify-around`}
        onClick={onClick}
      >
        {dots.map((d) => {
          return (
            <div key={d} className='w-3 h-3  bg-[#525252] rounded-full'></div>
          )
        })}
      </button>
    </motion.div>
  )
}
