import { Button } from "./components/ui/button"

function App() {
  return (
    <main className='h-screen w-full overflow-hidden text-slate-100'>
      <div className='h-full w-full flex flex-col bg-gray-600/30'>
        <div className='drag-zone cursor-move flex items-center justify-between bg-gray-600/20 pt-1 px-1'>
          <h1 className='text-2xl font-bold mb-2'>goDrawer</h1>
          <Button className='bg-orange-500 hover:bg-orange-300'>
            Settings
          </Button>
        </div>
        <div className='grow flex flex-col items-center justify-center bg-gray-800/50'>
          <p className='text-lg'>Drag files here to open them</p>
          <p className='text-lg'>Drag files here to open them</p>
          <p className='text-lg'>Drag files here to open them</p>
          <p className='text-lg'>Drag files here to open them</p>
        </div>
        <div className='flex items-center justify-center p-1'>
          <Button className='bg-sky-500 hover:bg-sky-300'>Add Drawer</Button>
        </div>
      </div>
    </main>
  )
}

export default App
