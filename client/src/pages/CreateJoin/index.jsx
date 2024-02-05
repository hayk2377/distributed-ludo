import { COLORS } from '../../utils'

export default function CreateJoin({onCreateLobby, onJoinLobby}) {
    const onCreate = ()=>{
        console.log('game created')
    }
    const onJoin = (code)=>{
        console.log('game joined', code)
    }

  return (
    <div className='max-w-[500px] mx-auto flex flex-col gap-3'>
      <h1>Create or Join a Game</h1>
      <Join onJoin={onJoinLobby}/>
      <p>or</p>
      <Create onCreate={onCreateLobby}/>
    </div>
  )
}

function Create({onCreate}) {
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
    <form className='w-full' onSubmit={(e)=>{
        e.preventDefault()
        onJoin(e.target.code.value)
    }}>
      <p>Game Code</p>
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
          Join
        </button>
      </div>
    </form>
  )
}
