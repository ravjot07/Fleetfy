import React from 'react';
import { useNavigate } from 'react-router-dom'; // Import useNavigate

function UserDashboard() {
  const navigate = useNavigate(); // Initialize useNavigate

  const handleCreateBooking = () => {
    navigate('/create-booking'); // Redirect to the create booking page
  };

  return (
    <div className="flex flex-col items-center justify-center h-full">
      <h1 className="text-4xl font-bold text-blue-600 mb-8">Welcome to the User Dashboard</h1>

      {/* Add Create Booking button */}
      <button
        onClick={handleCreateBooking}
        className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-500 transition"
      >
        Create Booking
      </button>
    </div>
  );
}

export default UserDashboard;
