import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function BookingForm() {
  const [pickupLocation, setPickupLocation] = useState('');
  const [dropoffLocation, setDropoffLocation] = useState('');
  const [vehicleType, setVehicleType] = useState('');
  const [estimatedCost, setEstimatedCost] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  // Check if the user is logged in and has the 'user' role
  useEffect(() => {
    const role = localStorage.getItem('role');
    const userID = localStorage.getItem('userID');
    
    // If no userID or role is found, or if the role is not 'user', redirect to login
    if (!role || role !== 'user' || !userID) {
      setMessage('You are not authorized to create a booking. Please log in.');
      navigate('/login');
    }
  }, [navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const userID = localStorage.getItem('userID'); // Get userID from localStorage
    const role = localStorage.getItem('role'); // Get role from localStorage

    // Check if userID and role are available
    if (!userID || !role) {
      setMessage('You are not authorized to create a booking. Please log in.');
      navigate('/login');
      return;
    }

    const bookingData = {
      pickup_location: pickupLocation,
      dropoff_location: dropoffLocation,
      vehicle_type: vehicleType,
      estimated_cost: parseFloat(estimatedCost),
    };

    try {
      const response = await fetch('http://localhost:8080/user/bookings', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'User-ID': userID,  // Send the User-ID from localStorage
          'Role': role,        // Send the Role from localStorage
        },
        body: JSON.stringify(bookingData),
      });

      const data = await response.json();

      if (response.ok) {
        setMessage(`Booking created successfully! Booking ID: ${data.booking_id}`);
      } else {
        setMessage(`Error creating booking: ${data.error || 'Unknown error'}`);
      }
    } catch (error) {
      setMessage('An error occurred while creating the booking');
    }
  };

  return (
    <div className="flex items-center justify-center h-full">
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow-lg w-96">
        <h2 className="text-2xl font-bold text-blue-600 mb-6">Create Booking</h2>

        {message && <p className="text-red-600 mb-4">{message}</p>}

        <div className="mb-4">
          <label className="block text-gray-700">Pickup Location</label>
          <input
            type="text"
            value={pickupLocation}
            onChange={(e) => setPickupLocation(e.target.value)}
            className="w-full p-2 border rounded mt-2"
            required
          />
        </div>

        <div className="mb-4">
          <label className="block text-gray-700">Dropoff Location</label>
          <input
            type="text"
            value={dropoffLocation}
            onChange={(e) => setDropoffLocation(e.target.value)}
            className="w-full p-2 border rounded mt-2"
            required
          />
        </div>

        <div className="mb-4">
          <label className="block text-gray-700">Vehicle Type</label>
          <input
            type="text"
            value={vehicleType}
            onChange={(e) => setVehicleType(e.target.value)}
            className="w-full p-2 border rounded mt-2"
            required
          />
        </div>

        <div className="mb-4">
          <label className="block text-gray-700">Estimated Cost</label>
          <input
            type="number"
            value={estimatedCost}
            onChange={(e) => setEstimatedCost(e.target.value)}
            className="w-full p-2 border rounded mt-2"
            required
          />
        </div>

        <button
          type="submit"
          className="w-full bg-blue-600 text-white p-2 rounded mt-4 hover:bg-blue-500 transition"
        >
          Submit
        </button>
      </form>
    </div>
  );
}

export default BookingForm;
