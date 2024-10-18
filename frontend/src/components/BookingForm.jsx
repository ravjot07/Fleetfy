import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { LoadScript, Autocomplete } from '@react-google-maps/api';

function BookingForm() {
  const [pickupLocation, setPickupLocation] = useState({ address: '', lat: null, lng: null });
  const [dropoffLocation, setDropoffLocation] = useState({ address: '', lat: null, lng: null });
  const [vehicleType, setVehicleType] = useState('');
  const [estimatedCost, setEstimatedCost] = useState('');
  const [estimatedDistance, setEstimatedDistance] = useState(null);
  const [message, setMessage] = useState('');
  const [pickupAutocomplete, setPickupAutocomplete] = useState(null);
  const [dropoffAutocomplete, setDropoffAutocomplete] = useState(null);
  const navigate = useNavigate();

  const libraries = ['places'];

  // Check if the user is logged in and has the 'user' role
  useEffect(() => {
    const role = localStorage.getItem('role');
    const userID = localStorage.getItem('userID');
    
    if (!role || role !== 'user' || !userID) {
      setMessage('You are not authorized to create a booking. Please log in.');
      navigate('/login');
    }
  }, [navigate]);

  // Handle loading Google Autocomplete for pickup
  const handleLoadPickup = (autocomplete) => {
    setPickupAutocomplete(autocomplete);
  };

  const handlePickupPlaceChanged = () => {
    if (pickupAutocomplete !== null) {
      const place = pickupAutocomplete.getPlace();
      setPickupLocation({
        address: place.formatted_address,
        lat: place.geometry.location.lat(),
        lng: place.geometry.location.lng(),
      });
    }
  };

  // Handle loading Google Autocomplete for dropoff
  const handleLoadDropoff = (autocomplete) => {
    setDropoffAutocomplete(autocomplete);
  };

  const handleDropoffPlaceChanged = () => {
    if (dropoffAutocomplete !== null) {
      const place = dropoffAutocomplete.getPlace();
      setDropoffLocation({
        address: place.formatted_address,
        lat: place.geometry.location.lat(),
        lng: place.geometry.location.lng(),
      });
    }
  };

  // Calculate distance using Google Distance Matrix API
  useEffect(() => {
    if (
      pickupLocation.lat &&
      pickupLocation.lng &&
      dropoffLocation.lat &&
      dropoffLocation.lng
    ) {
      const service = new window.google.maps.DistanceMatrixService();
      service.getDistanceMatrix(
        {
          origins: [{ lat: pickupLocation.lat, lng: pickupLocation.lng }],
          destinations: [{ lat: dropoffLocation.lat, lng: dropoffLocation.lng }],
          travelMode: window.google.maps.TravelMode.DRIVING,
        },
        (response, status) => {
          if (status === 'OK') {
            const distanceInMeters = response.rows[0].elements[0].distance.value;
            const distanceInKm = distanceInMeters / 1000;
            setEstimatedDistance(distanceInKm);
          } else {
            console.error('Error calculating distance:', status);
          }
        }
      );
    }
  }, [pickupLocation, dropoffLocation]);

  // Calculate cost based on distance and vehicle type
  useEffect(() => {
    if (estimatedDistance && vehicleType) {
      let costPerKm;
      switch (vehicleType) {
        case 'small':
          costPerKm = 5;
          break;
        case 'medium':
          costPerKm = 8;
          break;
        case 'large':
          costPerKm = 12;
          break;
        default:
          costPerKm = 5;
      }
      const cost = estimatedDistance * costPerKm;
      setEstimatedCost(cost.toFixed(2));
    } else {
      setEstimatedCost(null);
    }
  }, [estimatedDistance, vehicleType]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const userID = localStorage.getItem('userID'); // Get userID from localStorage
    const role = localStorage.getItem('role'); // Get role from localStorage

    if (!userID || !role) {
      setMessage('You are not authorized to create a booking. Please log in.');
      navigate('/login');
      return;
    }

    const bookingData = {
      pickup_location: pickupLocation.address,
      pickup_lat: pickupLocation.lat,
      pickup_lng: pickupLocation.lng,
      dropoff_location: dropoffLocation.address,
      dropoff_lat: dropoffLocation.lat,
      dropoff_lng: dropoffLocation.lng,
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
    <div className="flex items-center justify-center h-full" style={{ backgroundImage: `url('/booking.webp')` }}>
      <LoadScript googleMapsApiKey="AIzaSyCwP2TQiDq4eQPuMZf1OdnsLufGx72LbGo" libraries={libraries}>
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow-lg w-96">
          <h2 className="text-2xl font-bold text-geay-600 mb-6">Create Booking</h2>

          {message && <p className="text-red-600 mb-4">{message}</p>}

          <div className="mb-4">
            <label className="block text-gray-700">Pickup Location</label>
            <Autocomplete onLoad={handleLoadPickup} onPlaceChanged={handlePickupPlaceChanged}>
              <input
                type="text"
                value={pickupLocation.address}
                onChange={(e) => setPickupLocation({ ...pickupLocation, address: e.target.value })}
                className="w-full p-2 border rounded mt-2"
                required
              />
            </Autocomplete>
          </div>

          <div className="mb-4">
            <label className="block text-gray-700">Dropoff Location</label>
            <Autocomplete onLoad={handleLoadDropoff} onPlaceChanged={handleDropoffPlaceChanged}>
              <input
                type="text"
                value={dropoffLocation.address}
                onChange={(e) => setDropoffLocation({ ...dropoffLocation, address: e.target.value })}
                className="w-full p-2 border rounded mt-2"
                required
              />
            </Autocomplete>
          </div>

          <div className="mb-4">
            <label className="block text-gray-700">Vehicle Type</label>
            <select
              value={vehicleType}
              onChange={(e) => setVehicleType(e.target.value)}
              className="w-full p-2 border rounded mt-2"
              required
            >
              <option value="" disabled>Select vehicle type</option>
              <option value="small">Small Vehicle</option>
              <option value="medium">Medium Vehicle</option>
              <option value="large">Large Vehicle</option>
            </select>
          </div>

          {estimatedCost && (
            <div className="mb-4">
              <label className="block text-gray-700">Estimated Cost</label>
              <p className="text-green-600">${estimatedCost}</p>
            </div>
          )}

          <button
            type="submit"
            className="w-full bg-gray-600 text-white p-2 rounded mt-4 hover:bg-gray-500 transition"
            disabled={!estimatedCost}
          >
            Submit
          </button>
        </form>
      </LoadScript>
    </div>
  );
}

export default BookingForm;
