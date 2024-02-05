export const COLORS = {
  red: '#eb1c22',
  blue: '#27abfd',
  green: '#02a049',
  yellow: '#ffdd15',
  white: 'transparent',
  grey:'grey',
  lightGrey:'lightgrey'
}

export function getPlayer(pawnNumber) {
  return Math.floor(pawnNumber / 4)
}

export function getColor({ playerNumber, pawnNumber }) {
  const colorList = [COLORS.red, COLORS.green, COLORS.yellow, COLORS.blue]

  if (playerNumber !== undefined) return colorList[playerNumber]
  else return colorList[Math.floor(pawnNumber / 4)]
}

const redHomes = [32, 33, 47, 48]
const blueHomes = [167, 168, 182, 183]
const greenHomes = [41, 42, 56, 57]
const yellowHomes = [176, 177, 191, 192]
export const initialPositions = [
  ...redHomes,
  ...greenHomes,
  ...yellowHomes,
  ...blueHomes,
]
let str = ""
initialPositions.forEach((i) => {
  str += i + ", "
})
console.log(str)

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
        const hasBorder = y >= 7 || x == 6 || isExit
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
