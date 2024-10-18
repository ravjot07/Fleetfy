import React from 'react';
import { useNavigate } from 'react-router-dom';

function UserDashboard() {
  const navigate = useNavigate();

  const handleCreateBooking = () => {
    navigate('/create-booking'); // Redirect to the create booking page
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100 p-6" style={{ backgroundImage: `url('/user.webp')` }}>
      <div className="bg-white shadow-lg rounded-lg p-8 max-w-md w-full">
        <h1 className="text-4xl font-bold text-gray-800 text-center mb-6">
          Welcome to the User Dashboard
        </h1>
        <p className="text-gray-600 text-center mb-8">
          Easily book transportation for your goods by creating a new booking.
        </p>
        
        {/* Create Booking Button */}
        <button
          onClick={handleCreateBooking}
          className="w-full bg-gray-600 text-white font-semibold px-6 py-3 rounded-md shadow hover:bg-blue-500 transition duration-200 ease-in-out"
        >
          Create Booking
        </button>
      </div>
    </div>
  );
}

export default UserDashboard;
