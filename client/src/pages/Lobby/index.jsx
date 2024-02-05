import { getColor } from '../../utils'
import { COLORS } from '../../utils'
import { AiOutlineLoading3Quarters } from 'react-icons/ai'
import { motion } from 'framer-motion'

export default function Lobby({ lobbyState, onStartGame, playerId }) {
  //players have id and name
  const { gameCode, players } = lobbyState
  const missingPlayers = []
  for (let i = 0; i < 4; i++) {
    if (!players[i]) missingPlayers.push({ id: i + 'missing', name: '' })
  }

  const hasEnoughPlayers = players.length >= 4
  const isOwner = players[0]?.id === playerId

  return (
    <div className='max-w-[500px] mx-auto flex flex-col gap-3'>
      <h1>Lobby</h1>
      <p>
        Code: <span className='font-mono border p-1'>{gameCode}</span>{' '}
      </p>
      {!hasEnoughPlayers && <p>Waiting for players to join ...</p>}

      <div className='flex flex-col gap-3'>
        {players.map((player, playerNumber) => (
          <Player
            key={player.id}
            name={player.name}
            id={player.id}
            color={getColor({ playerNumber })}
          />
        ))}

        {missingPlayers.map((player, playerNumber) => (
          <Player
            key={player.id}
            name={player.name}
            id={player.id}
            color={COLORS.grey}
            isMissing={true}
          />
        ))}
      </div>
      {isOwner && (
        <button
          className='p-3 border rounded text-white'
          style={{ background: COLORS.green }}
          onClick={onStartGame}
        >
          Start
        </button>
      )}
    </div>
  )
}

function Player({ color, name, isMissing = false }) {
  return (
    <div className='flex items-center gap-3'>
      <div
        className={`w-5 h-5 rounded-full`}
        style={{ background: color }}
      ></div>
      <p className='flex gap-3 items-center'>
        {name}
        {isMissing && (
          <motion.span
            animate={{ rotate: 360 }}
            transition={{ repeat: Infinity, duration: 1, ease: 'linear' }}
          >
            <AiOutlineLoading3Quarters />
          </motion.span>
        )}
      </p>
    </div>
  )
}
