import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (response.ok) {
        // Destructure role, userID, and username from the response data
        const { role, userID, username } = data;

        // Log the role and userID to the console
        console.log('UserID:', userID);
        console.log('Role:', role);

        // Save role and userID to localStorage
        localStorage.setItem('userID', userID);
        localStorage.setItem('role', role);

        // Redirect based on role
        if (role === 'user') {
          navigate('/user-dashboard');
        } else if (role === 'admin') {
          navigate('/admin-dashboard');
        } else if (role === 'driver') {
          navigate('/driver-dashboard');
        } else {
          setMessage('Login failed: Unrecognized role');
        }
      } else {
        setMessage('Invalid credentials');
      }
    } catch (error) {
      setMessage('An error occurred during login');
    }
  };

  return (
    <div className="login-container">
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Login</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default Login;
