import React, { useState, useEffect } from 'react';
import { Line, Bar, Doughnut } from 'react-chartjs-2';
import { Chart, CategoryScale, LinearScale, PointElement, LineElement, BarElement, ArcElement, Title, Tooltip, Legend } from 'chart.js';
import axios from 'axios';

// Register Chart.js components
Chart.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, ArcElement, Title, Tooltip, Legend);

const Analytics = () => {
  // Initializing state with empty arrays for labels and datasets
  const [bookingsData, setBookingsData] = useState({
    labels: [],
    datasets: []
  });

  const [vehicleStatusData, setVehicleStatusData] = useState({
    labels: [],
    datasets: []
  });

  const [driverPerformanceData, setDriverPerformanceData] = useState({
    labels: [],
    datasets: []
  });

  const [revenueData, setRevenueData] = useState({
    labels: [],
    datasets: []
  });

  const [bookingStatusData, setBookingStatusData] = useState({
    labels: [],
    datasets: []
  });

  // Fetch analytics data when the component mounts
  useEffect(() => {
    // Fetch Bookings Over Time data
    axios.get('http://localhost:8080/admin/analytics/bookings-over-time', {
      headers: { 'Role': 'admin' }
    })
    .then(response => {
      const dates = response.data.map(item => new Date(item.date).toLocaleDateString());
      const counts = response.data.map(item => item.count);

      setBookingsData({
        labels: dates,
        datasets: [{
          label: 'Bookings Over Time',
          data: counts,
          borderColor: 'rgba(75, 192, 192, 1)',
          backgroundColor: 'rgba(75, 192, 192, 0.2)',
          fill: true,
        }]
      });
    })
    .catch(error => {
      console.error('Error fetching bookings data:', error);
    });

    // Fetch Vehicle Status data
    axios.get('http://localhost:8080/admin/analytics/vehicle-status', {
      headers: { 'Role': 'admin' }
    })
    .then(response => {
      setVehicleStatusData({
        labels: ['Active', 'Idle'],
        datasets: [{
          data: [response.data.active, response.data.idle],
          backgroundColor: ['#36A2EB', '#FF6384'],
        }]
      });
    })
    .catch(error => {
      console.error('Error fetching vehicle status data:', error);
    });

    // Fetch Driver Performance data
    axios.get('http://localhost:8080/admin/analytics/driver-performance', {
      headers: { 'Role': 'admin' }
    })
    .then(response => {
      const drivers = response.data.map(item => item.driver_name);
      const deliveries = response.data.map(item => item.deliveries);

      setDriverPerformanceData({
        labels: drivers,
        datasets: [{
          label: 'Completed Deliveries',
          data: deliveries,
          backgroundColor: 'rgba(153, 102, 255, 0.2)',
          borderColor: 'rgba(153, 102, 255, 1)',
        }]
      });
    })
    .catch(error => {
      console.error('Error fetching driver performance data:', error);
    });

    // Fetch Revenue Over Time data
    axios.get('http://localhost:8080/admin/analytics/revenue-over-time', {
      headers: { 'Role': 'admin' }
    })
    .then(response => {
      const dates = response.data.map(item => new Date(item.date).toLocaleDateString());
      const revenues = response.data.map(item => item.revenue);

      setRevenueData({
        labels: dates,
        datasets: [{
          label: 'Revenue Over Time',
          data: revenues,
          borderColor: 'rgba(54, 162, 235, 1)',
          backgroundColor: 'rgba(54, 162, 235, 0.2)',
          fill: true,
        }]
      });
    })
    .catch(error => {
      console.error('Error fetching revenue data:', error);
    });

    // Fetch Booking Status Distribution data
    axios.get('http://localhost:8080/admin/analytics/booking-status-distribution', {
      headers: { 'Role': 'admin' }
    })
    .then(response => {
      const statuses = response.data.map(item => item.status);
      const counts = response.data.map(item => item.count);

      setBookingStatusData({
        labels: statuses,
        datasets: [{
          data: counts,
          backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56'],
        }]
      });
    })
    .catch(error => {
      console.error('Error fetching booking status distribution data:', error);
    });

  }, []); // Empty dependency array to run effect only once when component mounts

  return (
    <div>
      <h2>Admin Dashboard - Analytics</h2>

      {/* Line Chart: Bookings Over Time */}
      <div>
        <h3>Number of Bookings Over the Last 7 Days</h3>
        {bookingsData.datasets.length > 0 ? (
          <Line data={bookingsData} />
        ) : (
          <p>Loading data...</p>
        )}
      </div>

      {/* Doughnut Chart: Vehicle Availability Status */}
      {/* <div>
        <h3>Vehicle Availability Status</h3>
        {vehicleStatusData.datasets.length > 0 ? (
          <Doughnut data={vehicleStatusData} />
        ) : (
          <p>Loading data...</p>
        )}
      </div> */}

      {/* Bar Chart: Driver Performance */}
      {/* <div>
        <h3>Driver Performance</h3>
        {driverPerformanceData.datasets.length > 0 ? (
          <Bar data={driverPerformanceData} />
        ) : (
          <p>Loading data...</p>
        )}
      </div> */}

      {/* Line Chart: Revenue Over Time */}
      <div>
        <h3>Revenue Over Time</h3>
        {revenueData.datasets.length > 0 ? (
          <Line data={revenueData} />
        ) : (
          <p>Loading data...</p>
        )}
      </div>

      {/* Doughnut Chart: Booking Status Distribution */}
      <div>
        <h3>Booking Status Distribution</h3>
        {bookingStatusData.datasets.length > 0 ? (
          <Doughnut data={bookingStatusData} />
        ) : (
          <p>Loading data...</p>
        )}
      </div>
    </div>
  );
};
  
export default Analytics;
