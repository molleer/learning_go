import React, { useEffect, useState } from 'react';
import logo from './logo.svg';
import './App.css';

const ws = new WebSocket("ws://localhost:8080/ws")

ws.onopen = () => {
  console.log("Websocket has opened")
  ws.send("Hello from frontend")
}

ws.onmessage = (msg) => {
  console.log(msg.data)
}

function App() {

  return (
    <div className="App">
     
    </div>
  );
}

export default App;
