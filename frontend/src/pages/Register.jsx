import React, { useState } from 'react';

function Register() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('');
  const [message, setMessage] = useState('');

  const handleRegister = async (e) => {
    e.preventDefault();
    
    const response = await fetch('http://localhost:8080/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password, role }),
    });

    const data = await response.json();
    
    if (response.ok) {
      setMessage(data.message);
    } else {
      setMessage(data.error || 'Error occurred');
    }
  };

  return (
    <div className="flex items-center justify-center h-full">
      <form onSubmit={handleRegister} className="bg-white p-6 rounded-lg shadow-lg">
        <h2 className="text-2xl font-bold text-blue-600 mb-6">Register</h2>
        <input 
          type="text" 
          placeholder="Username" 
          value={username} 
          onChange={(e) => setUsername(e.target.value)} 
          className="w-full p-2 border rounded mb-4"
          required
        />
        <input 
          type="password" 
          placeholder="Password" 
          value={password} 
          onChange={(e) => setPassword(e.target.value)} 
          className="w-full p-2 border rounded mb-4"
          required
        />
        <input 
          type="text" 
          placeholder="Role (e.g., admin, user)" 
          value={role} 
          onChange={(e) => setRole(e.target.value)} 
          className="w-full p-2 border rounded mb-4"
          required
        />
        <button type="submit" className="w-full bg-blue-600 text-white p-2 rounded">Register</button>
        {message && <p className="mt-4 text-red-600">{message}</p>}
      </form>
    </div>
  );
}

export default Register;
