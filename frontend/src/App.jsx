import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import AdminDashboard from './pages/AdminDashboard';
import DriverDashboard from './pages/DriverDashboard';
import UserDashboard from './pages/UserDashboard';
import BookingForm from './components/BookingForm';

function App() {
  return (
    <Router>
      <div className="bg-gray-100 h-screen">
        <nav className="p-4 bg-blue-600 text-white">
          <ul className="flex justify-between">
            <li><Link to="/" className="text-lg">Home</Link></li>
            <div>
              <Link to="/login" className="mx-2">Login</Link>
              <Link to="/register" className="mx-2">Register</Link>
              <Link to="/create-booking" className="mx-2">Create Booking</Link> {/* Add link to booking form */}
            </div>
          </ul>
        </nav>

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
