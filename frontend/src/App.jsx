import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import AdminDashboard from './pages/AdminDashboard';
import DriverDashboard from './pages/DriverDashboard';
import UserDashboard from './pages/UserDashboard';
import BookingForm from './components/BookingForm';
import Navbar from './components/Navbar';

function App() {
  return (
    <Router>
      <div className="bg-gray-100 h-screen">
        <Navbar />        
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/admin-dashboard" element={<AdminDashboard />} />
          <Route path="/driver-dashboard" element={<DriverDashboard />} />
          <Route path="/user-dashboard" element={<UserDashboard />} />
          <Route path="/create-booking" element={<BookingForm />} /> {/* Add route to booking form */}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
