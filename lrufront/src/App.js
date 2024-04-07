import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');

  const handleGet = async () => {
    try {
      const response = await axios.get(`http://localhost:8585/get/${key}`);
      alert(`Value for key ${key}: ${response.data}`);
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  };

  const handleSet = async () => {
    try {
      await axios.post(`http://localhost:8585/set/${key}/${value}`);
      alert(`Key ${key} set with value ${value}`);
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  };

  return (
    <div className="container">
      <h1>LRU Cache</h1>
      <div className="input-container">
        <input
          type="text"
          placeholder="Enter Key"
          value={key}
          onChange={(e) => setKey(e.target.value)}
          className="input-field"
        />
        <input
          type="text"
          placeholder="Enter Value"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          className="input-field"
        />
        <button onClick={handleSet} className="button">Set value</button>
        <button onClick={handleGet} className="button">Get value</button>
      </div>
    </div>
  );
}

export default App;
