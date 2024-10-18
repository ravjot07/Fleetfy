import React, { useState, useEffect } from 'react';
import Analytics from '../components/Analytics';

const AdminDashboard = () => {
  const [bookings, setBookings] = useState([]); // List of all bookings
  const [activeBookings, setActiveBookings] = useState([]); // Active bookings for each driver
  const [message, setMessage] = useState('');

  // Fetch all bookings
  const fetchAllBookings = async () => {
    const role = localStorage.getItem('role');
    try {
      const response = await fetch('http://localhost:8080/admin/bookings', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Role': role,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch bookings. Status code: ${response.status}`);
      }

      const data = await response.json();
      setBookings(data || []); // Ensure bookings is an array
    } catch (error) {
      console.error('Error fetching bookings:', error);
      setMessage(`An error occurred while fetching bookings: ${error.message}`);
    }
  };

  // Fetch active bookings list for each driver
  const fetchActiveBookings = async () => {
    const role = localStorage.getItem('role');
    try {
      const response = await fetch('http://localhost:8080/admin/drivers/active-bookings', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Role': role,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch active bookings. Status code: ${response.status}`);
      }

      const data = await response.json();
      setActiveBookings(data || []); // Ensure activeBookings is an array
    } catch (error) {
      console.error('Error fetching active bookings:', error);
      setMessage(`An error occurred while fetching active bookings: ${error.message}`);
    }
  };

  // Mark a booking as complete
  const markBookingComplete = async (bookingId) => {
    const role = localStorage.getItem('role');
    try {
      const response = await fetch(`http://localhost:8080/admin/bookings/${bookingId}/complete`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Role': role,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to mark booking as complete. Status code: ${response.status}`);
      }

      setMessage(`Booking ${bookingId} marked as complete successfully!`);
      // Refresh bookings after marking one as complete
      fetchAllBookings();
    } catch (error) {
      console.error('Error marking booking as complete:', error);
      setMessage(`An error occurred while marking the booking as complete: ${error.message}`);
    }
  };

  // Fetch all bookings and active bookings when the component mounts
  useEffect(() => {
    fetchAllBookings();
    fetchActiveBookings();
  }, []);

  return (
    <div className="container mx-auto p-4" >
      
      <h1 className="text-3xl font-bold mb-6">Admin Dashboard</h1>
      <Analytics />
      {message && <p className="text-red-500 mb-4">{message}</p>}

      <h2 className="text-2xl mb-4">All Bookings</h2>
      {bookings.length === 0 ? (
        <p>No bookings to display.</p>
      ) : (
        <table className="table-auto w-full border-collapse">
          <thead>
            <tr>
              <th className="border px-4 py-2">Booking ID</th>
              <th className="border px-4 py-2">Pickup Location</th>
              <th className="border px-4 py-2">Dropoff Location</th>
              <th className="border px-4 py-2">Vehicle Type</th>
              <th className="border px-4 py-2">Driver ID</th>
              <th className="border px-4 py-2">Status</th>
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
                <td className="border px-4 py-2">{booking.driver_id || 'Unassigned'}</td>
                <td className="border px-4 py-2">{booking.status}</td>
                <td className="border px-4 py-2">
                  {booking.status !== 'complete' ? (
                    <button
                      className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-700 transition"
                      onClick={() => markBookingComplete(booking.id)}
                    >
                      Mark as Complete
                    </button>
                  ) : (
                    <span className="text-gray-500">Completed</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <h2 className="text-2xl mt-6 mb-4">Active Bookings per Driver</h2>
      {activeBookings.length === 0 ? (
        <p>No active bookings to display.</p>
      ) : (
        <table className="table-auto w-full border-collapse">
          <thead>
            <tr>
              <th className="border px-4 py-2">Driver ID</th>
              <th className="border px-4 py-2">Active Booking IDs</th>
            </tr>
          </thead>
          <tbody>
            {activeBookings.map((driver) => (
              <tr key={driver.driver_id}>
                <td className="border px-4 py-2">{driver.driver_id}</td>
                <td className="border px-4 py-2">
                  {Array.isArray(driver.active_bookings) && driver.active_bookings.length > 0 ? (
                    driver.active_bookings.map((bookingId) => (
                      <span key={bookingId} className="inline-block mr-2">{bookingId}</span>
                    ))
                  ) : (
                    <span>No active bookings</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default AdminDashboard;
