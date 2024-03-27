import { useState } from 'react'
import  {BrowserRouter, Route,Routes}  from 'react-router-dom';
import CreateRoom from './components/CreateRoom'
import Room from './components/Room'

function App() {
  return (
     <div className='App'>
      <BrowserRouter>
      
        <Routes>
        <Route exact path="/" element={<CreateRoom/>}></Route>
        <Route path="/room/:roomid" element={<Room/>}></Route>
        </Routes>
      </BrowserRouter>
     </div>
    
  )
}

export default App
