import React from 'react';

function Home() {
  return (
    <div className="flex flex-col items-center justify-center h-full">
      <h1 className="text-4xl font-bold text-blue-600">Welcome to the Logistics Platform</h1>
      <p className="text-lg mt-4 text-gray-700">Book transportation services, track vehicles, and more!</p>
      <div className="mt-6">
        <a href="/login" className="bg-blue-600 text-white px-4 py-2 rounded mx-2">Login</a>
        <a href="/register" className="bg-blue-600 text-white px-4 py-2 rounded mx-2">Register</a>
      </div>
    </div>
  );
}

export default Home;
