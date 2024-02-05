import Pawn from './Pawn'
import { getColor, COLORS } from '../../utils'

import { motion } from 'framer-motion'

export default function Tile({ color, hasBorder, position }) {
  // const showNumber = true
  const showNumber = false
  return (
    <motion.div
      animate={{ scale: [0.2, 1] }}
      transition={{ delay: position * 0.005 }}
      className={`
      w-full h-full
      flex justify-center items-center ${hasBorder ? 'tile-border' : ''} ]`}
      style={{ background: color }}
    >
      <span className={` opacity-50 ${showNumber ? '' : 'hidden'} `}>
        {position}
      </span>
    </motion.div>
  )
}

function rotateBy90(matrix) {
  const n = matrix.length

  // Create a new matrix with rotated dimensions
  const rotatedMatrix = new Array(n)
  for (let i = 0; i < n; i++) {
    rotatedMatrix[i] = new Array(n)
  }

  // Perform the rotation
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < n; j++) {
      rotatedMatrix[-(j - n + 1)][i] = matrix[i][j]
    }
  }

  return rotatedMatrix
}

function emptyMatrix(size) {
  const matrix = []
  for (let i = 0; i < size; i++) {
    const row = []
    for (let j = 0; j < size; j++) {
      row.push(undefined)
    }
    matrix.push(row)
  }
  return matrix
}

function flattenMatrix(matrix) {
  let list = []
  for (const row of matrix) {
    list = [...list, ...row]
  }
  return list
}

export function getTileProperties() {
  const gridSize = 15
  let tileProperties = emptyMatrix(gridSize)
  for (let playerNumber = 0; playerNumber < 4; playerNumber++) {
    tileProperties = rotateBy90(tileProperties)

    for (let y = 0; y < 8; y++) {
      for (let x = 0; x < 7; x++) {
        const isExit = x == 1 && y == 6
        const hasBorder = y >= 6 || x == 6 || isExit
        const horizontalColor =
          (y == 0 && x < 6) || (y == 5 && x < 6) || (y == 7 && x > 0)
        const verticalColor = (x == 0 && y < 6) || (x == 5 && y < 6)
        const hasColor = horizontalColor || verticalColor || isExit
        const color = hasColor ? getColor({ playerNumber }) : COLORS.white
        tileProperties[y][x] = { color, hasBorder }
      }
    }
  }
  tileProperties = rotateBy90(tileProperties)
  return flattenMatrix(tileProperties)
}
