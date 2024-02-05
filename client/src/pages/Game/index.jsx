import { useEffect, useRef, useState } from 'react'
import { COLORS, getPlayer } from '../../utils'
import Diamond from './Diamond'
import Dice from './Dice'
import PlayerName from './PlayerName'
import Tile, { getTileProperties } from './Tile'
import { toast } from 'react-toastify'
import Pawn from './Pawn'
import { AiOutlineLoading3Quarters } from 'react-icons/ai'
import { motion } from 'framer-motion'

export default function Game({
  gameState,
  onPawn,
  onRollDice,
  onGameEnd,
  playerId,
}) {
  //Not Game state related
  const tileProperties = getTileProperties()
  const boardRef = useRef(null)

  const [tileXYPositions, setTileXYPostions] = useState(null)
  useEffect(() => {
    const updateXYPositions = () => {
      if (!boardRef.current) {
        setTileXYPostions(null)
        return
      }

      const tileContianers = [
        ...boardRef.current.querySelectorAll('.tile-container'),
      ]
      const tileXYPositions = tileContianers.map((tile) => {
        const {
          x: tileX,
          y: tileY,
          width,
          height,
        } = tile.getBoundingClientRect()
        const { x: originX, y: originY } =
          tileContianers[0].getBoundingClientRect()

        const x = tileX - originX + width / 2
        const y = tileY - originY + height / 2
        return [x, y]
      })

      setTileXYPostions(tileXYPositions)
    }

    window.addEventListener('resize', updateXYPositions)
    updateXYPositions()
    return () => window.removeEventListener('resize', updateXYPositions)
  }, [boardRef.current === null])

  //Game state related
  const {
    gameCode,
    players,
    pawnPositions,
    turn,
    dice,
    hasDiced,
    movablePawns,
    recovering,
    winner,
    onlineStatuses,
    pawnAnimationPaths,
    offlineTimeout,
  } = gameState
  const playerNames = players.map((player) => player?.name)
  const myPlayerNumber = players.findIndex((player) => player?.id === playerId)
  const isMyTurn = myPlayerNumber === turn

  //Animation related but not directly from game state
  const [isAnimating, setIsAnimating] = useState(false)
  const [animatedPawnPositions, setAnimatedPawnPositions] =
    useState(pawnPositions)
  const [diceIsRolling, setDiceIsRolling] = useState(false)

  const [offlineTimeStamps, setOfflineTimeStamps] = useState(
    onlineStatuses.map(() => null)
  )

  //Announce winner if there is
  useEffect(() => {
    if (winner !== -1) {
      if (winner === myPlayerNumber) toast.success('You won the game')
      else toast.success(`${playerNames[winnerPlayerNumber]} won the game`)
      onGameEnd()
    }
  }, [winner, myPlayerNumber, playerNames])

  //Animate pawn movement using animation paths
  useEffect(() => {
    //Stop interactivity
    setIsAnimating(true)

    let longestPath = 0
    for (const path of pawnAnimationPaths) {
      if (!path) continue
      if (path.length > longestPath) longestPath = path.length
    }

    //Animate until another pawnAnimationPaths replaces this one
    let index = 0
    //No animation pawns stay put
    const newPawnPositions = [...pawnPositions]

    let interval
    const animateFrame = () => {
      console.log('animating pawns')
      for (
        let pawnNumber = 0;
        pawnNumber < pawnAnimationPaths.length;
        pawnNumber++
      ) {
        const path = pawnAnimationPaths[pawnNumber]
        if (!path) continue
        const position = path[index]
        //Update if still has path, leave as is if end path
        if (position || position == 0) {
          newPawnPositions[pawnNumber] = position
        }
      }

      if (index >= longestPath) {
        //Set true positions and allow interaction at last frame
        setAnimatedPawnPositions([...pawnPositions])
        clearInterval(interval)
        setIsAnimating(false)
      } else {
        //Progress animation, set temporary positions frame
        index++
        setAnimatedPawnPositions([...newPawnPositions])
      }
    }

    //Start animating in intervals
    interval = setInterval(animateFrame, 200)
    animateFrame()

    return () => {
      clearInterval(interval)
      setIsAnimating(false)
    }
  }, [JSON.stringify(pawnAnimationPaths, pawnPositions)])

  //Animate dice
  useEffect(() => {
    //works even if same dice back to back bc hasDiced changes
    if (!hasDiced) return

    setDiceIsRolling(true)
    setIsAnimating(true)
    const timeOut = setTimeout(() => {
      setIsAnimating(false)
      setDiceIsRolling(false)
    }, 1000)

    return () => {
      clearTimeout(timeOut)
      setIsAnimating(false)
      setDiceIsRolling(false)
    }
  }, [dice, hasDiced])

  //Animate offline timers
  useEffect(() => {
    console.log('animating offline timers')
    const newOfflineTimeStamps = [...offlineTimeStamps]

    for (let i = 0; i < onlineStatuses.length; i++) {
      //has already offline timed out or didn't event exist
      if (players[i] === null) continue

      //has time stamp means player became offline
      const hasTimeStamp = offlineTimeStamps[i] !== null
      const isOnline = onlineStatuses[i] === true

      const toOffline = !isOnline && !hasTimeStamp
      const toOnline = isOnline && hasTimeStamp

      //Announce and update
      if (toOffline) {
        newOfflineTimeStamps[i] = Date.now()
        toast.error(`${playerNames[i]} is offline`)
      }
      if (toOnline) {
        newOfflineTimeStamps[i] = null
        toast.success(`${playerNames[i]} is back online`)
      }
    }

    //reupdate the timestamps every 0.5 second to animate the count down
    setOfflineTimeStamps(newOfflineTimeStamps)
    const interval = setInterval(() => {
      console.log('updating countdown')
      setOfflineTimeStamps([...newOfflineTimeStamps])
    }, 500)

    //update time for 10 seconds
    const timeout = setTimeout(() => {
      console.log('clearing offline timer interval')
      clearInterval(interval)
    }, offlineTimeout * 1000 + 2000)

    return () => {
      clearInterval(interval)
      clearTimeout(timeout)
    }
  }, [JSON.stringify(onlineStatuses)])

  //Animate offline kick outs
  useEffect(() => {
    for (let playerNumber = 0; playerNumber < players.length; playerNumber++) {
      const wentOffline = offlineTimeStamps[playerNumber] !== null
      const notInGame = players[playerNumber] === null
      if (wentOffline && notInGame) {
        toast.error(
          `Player ${
            playerNumber + 1
          } was kicked out for being offline for too long`
        )
        //so that it doesn't get called again
        offlineTimeStamps[playerNumber] = null
        setOfflineTimeStamps([...offlineTimeStamps])
      }
    }
  }, [JSON.stringify(players)])

  //User Interactive Actions
  const movePawn = (pawnNumber) => {
    if (isAnimating) return
    onPawn(pawnNumber)
  }

  const rollDice = () => {
    if (isAnimating) return
    onRollDice()
  }

  const playerNamesComponents = playerNames.map((name, playerNumber) => (
    <PlayerName
      key={playerNumber}
      isWinner={winner === playerNumber}
      playerNumber={playerNumber}
      name={name ?? null}
      isTurn={playerNumber === turn}
      offlineTimeStamp={offlineTimeStamps[playerNumber]}
      countDownDuration={offlineTimeout}
    />
  ))

  return (
    <div className=' container flex flex-col items-center justify-center gap-6 min-h-screen'>
      <p>Code {gameCode}</p>
      {recovering && (
        <p className='flex items-center gap-3'>
          <motion.span
            animate={{ rotate: 360 }}
            transition={{ repeat: Infinity, duration: 1, ease: 'linear' }}
          >
            <AiOutlineLoading3Quarters />
          </motion.span>
          Trying to recover game after server crash, please be patient!
        </p>
      )}
      <div className='flex justify-between w-[410px] md:w-[710px]'>
        {playerNamesComponents[0]}
        {playerNamesComponents[1]}
      </div>
      <div
        ref={boardRef}
        className='sm relative tile-border 
        w-[410px] h-[410px] md:w-[710px] md:h-[710px]'
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(15, 1fr)',
          gridTemplateRows: 'repeat(15, 1fr)',
        }}
      >
        {animatedPawnPositions.map((position, pawnNumber) => {
          if (position == -1) return
          const shouldFloat = movablePawns.includes(pawnNumber) && !isAnimating
          const onClick = movePawn
          const xyPosition = tileXYPositions?.[position] || [0, 0]
          return (
            <Pawn
              key={pawnNumber}
              pawnNumber={pawnNumber}
              shouldFloat={shouldFloat}
              onClick={onClick}
              top={xyPosition[1]}
              left={xyPosition[0]}
            />
          )
        })}

        {tileProperties.map((property, position) => (
          <div className='tile-container' key={position}>
            <Tile {...property} position={position} />
          </div>
        ))}

        <Diamond />
      </div>
      <div className='flex justify-between w-[410px] md:w-[710px]'>
        {playerNamesComponents[3]}
        {playerNamesComponents[2]}
      </div>

      <Dice
        value={dice}
        shouldFloat={!hasDiced && isMyTurn && !isAnimating}
        isRolling={diceIsRolling}
        onClick={rollDice}
      />
    </div>
  )
}
