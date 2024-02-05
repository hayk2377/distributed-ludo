import { COLORS } from '../../utils'

export default function CreateJoin({ onCreateLobby, onJoinLobby, name }) {
  return (
    <div className='max-w-[500px] mx-auto flex flex-col gap-6'>
      <div className='flex flex-col gap-3'>
        <h1 className='text-3xl'>Welcome {name}</h1>
        <p>Want to play? </p>
        <div className='flex flex-col gap-3'></div>

        <Join onJoin={onJoinLobby} />
        <p>or</p>
        <Create onCreate={onCreateLobby} />
      </div>
    </div>
  )
}

function Create({ onCreate }) {
  return (
    <button
      style={{ background: COLORS.green }}
      onClick={onCreate}
      className=' text-white p-3 border rounded'
    >
      Create New Game
    </button>
  )
}

function Join({ onJoin }) {
  return (
    <form
      className='w-full'
      onSubmit={(e) => {
        e.preventDefault()
        onJoin(e.target.code.value)
      }}
    >
      <div className='w-full flex'>
        <input
          type='text'
          name='code'
          defaultValue={1234}
          className='p-3 flex-1 border'
        />
        <button
          className='p-3 border rounded'
          style={{ borderColor: COLORS.green }}
        >
          Join Game
        </button>
      </div>
    </form>
  )
}
