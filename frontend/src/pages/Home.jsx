import React from 'react';


function Home() {
  return (
    <div className="flex flex-col items-center justify-center h-full bg-cover bg-center" style={{ backgroundImage: `url('/home.webp')` }}>
      <h1 className="text-6xl font-bold text-gray-1000">
        Welcome to the Logistics Platform
      </h1>
      <p className="text-3xl mt-4 text-gray-700">
        Book transportation services, track vehicles, and more!
      </p>
      <div className="mt-6">
        <a href="/login" className="text-xl bg-gray-800 text-white px-4 py-2 rounded mx-2">
          Login
        </a>
        <a href="/register" className="text-xl bg-gray-800 text-white px-4 py-2 rounded mx-2">
          Register
        </a>
      </div>
    </div>
  );
}

export default Home;
