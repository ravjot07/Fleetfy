import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const DriverDashboard = () => {
  const [bookings, setBookings] = useState([]);  // Initialize as an empty array
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  // Fetch pending bookings from the backend
  const fetchPendingBookings = async () => {
    const driverID = localStorage.getItem('userID'); // Get Driver ID from localStorage
    const role = localStorage.getItem('role');       // Get Role from localStorage

    // If driverID or role is not present, navigate to login
    if (!driverID || role !== 'driver') {
      setMessage('Unauthorized access. Please log in as a driver.');
      navigate('/login');
      return;
    }

    try {
      console.log('Sending Request Headers:', {
        'Content-Type': 'application/json',
        'Driver-ID': driverID,
        'Role': role,
      });

      const response = await fetch('http://localhost:8080/driver/bookings/pending', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Driver-ID': driverID,
          'Role': role,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch bookings. Status code: ${response.status}`);
      }

      const data = await response.json();
      setBookings(data || []);  // Ensure bookings is an array, even if the data is null or undefined
    } catch (error) {
      console.error('Error fetching bookings:', error);
      setMessage(`An error occurred while fetching bookings: ${error.message}`);
    }
  };

  // Function to accept a booking
  const acceptBooking = async (bookingId) => {
    const driverID = localStorage.getItem('userID'); // Get Driver ID from localStorage
    const role = localStorage.getItem('role');

    if (!driverID || role !== 'driver') {
      setMessage('Unauthorized access. Please log in as a driver.');
      navigate('/login');
      return;
    }

    try {
      console.log(`Accepting booking ${bookingId} for driver ${driverID}`);

      const response = await fetch(`http://localhost:8080/driver/bookings/${bookingId}/accept`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Driver-ID': driverID,  // Send Driver-ID from localStorage
          'Role': role,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('Unauthorized. Please check your credentials.');
        }
        throw new Error(`Failed to accept booking. Status code: ${response.status}`);
      }

      setMessage(`Booking ${bookingId} accepted successfully!`);
      // After accepting a booking, refresh the bookings list
      fetchPendingBookings();
    } catch (error) {
      console.error('Error accepting booking:', error);
      setMessage(`An error occurred while accepting the booking: ${error.message}`);
    }
  };

  // Fetch pending bookings when the component mounts
  useEffect(() => {
    fetchPendingBookings();
  }, []);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Driver Dashboard</h1>

      {message && <p className="text-red-500 mb-4">{message}</p>}

      {bookings.length === 0 ? (
        <p>No pending bookings to display.</p>
      ) : (
        <table className="table-auto w-full border-collapse">
          <thead>
            <tr>
              <th className="border px-4 py-2">Booking ID</th>
              <th className="border px-4 py-2">Pickup Location</th>
              <th className="border px-4 py-2">Dropoff Location</th>
              <th className="border px-4 py-2">Vehicle Type</th>
              <th className="border px-4 py-2">Estimated Cost</th>
              <th className="border px-4 py-2">Action</th>
            </tr>
          </thead>
          <tbody>
            {bookings.map((booking) => (
              <tr key={booking.id}>
                <td className="border px-4 py-2">{booking.id}</td>
                <td className="border px-4 py-2">{booking.pickup_location}</td>
                <td className="border px-4 py-2">{booking.dropoff_location}</td>
                <td className="border px-4 py-2">{booking.vehicle_type}</td>
                <td className="border px-4 py-2">${booking.estimated_cost.toFixed(2)}</td>
                <td className="border px-4 py-2">
                  <button
                    className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
                    onClick={() => acceptBooking(booking.id)}
                  >
                    Accept
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default DriverDashboard;
