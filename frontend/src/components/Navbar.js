import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { FaSignInAlt, FaHdd, FaTerminal, FaUser, FaCaretDown } from 'react-icons/fa';

const Navbar = () => {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useState(() => {
    const LoggedIn = localStorage.getItem('isLoggedIn') === 'true';
    setIsLoggedIn(LoggedIn);
  } , []);

  const toggleDropdown = () => {
    setDropdownOpen((prev) => !prev);
  };

  const handleLogin = () => {
    window.location.href = '/login';
  };

  const handleLogout = () => {
    localStorage.setItem('isLoggedIn', 'false');
    setIsLoggedIn(false);
    setDropdownOpen(false);
  };

  return (
    <nav className="bg-gray-800 h-16 p-4 rounded-lg shadow-lg mb-6 flex items-center">
      <div className="flex justify-between w-full">
        <div className="flex space-x-4">
          <Link to="/login" className="text-white flex items-center hover:text-blue-500 transition-colors">
            <FaSignInAlt className="mr-2" /> Login
          </Link>
          <Link to="/disk-manager" className="text-white flex items-center hover:text-blue-500 transition-colors">
            <FaHdd className="mr-2" /> Disk Manager
          </Link>
          <Link to="/console" className="text-white flex items-center hover:text-blue-500 transition-colors">
            <FaTerminal className="mr-2" /> Console
          </Link>
        </div>

        <div className="relative">
          {isLoggedIn ? (
            <div className="flex items-center">
              <button
                className="flex items-center text-white focus:outline-none"
                onClick={toggleDropdown}
              >
                <FaUser className="mr-2" />
                User Menu
                <FaCaretDown className="ml-1" />
              </button>
              {dropdownOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-gray-700 rounded-md shadow-lg z-10">
                  <button
                    onClick={handleLogout}
                    className="block w-full text-left px-4 py-2 text-white hover:bg-gray-600 focus:outline-none"
                  >
                    Cerrar Sesión
                  </button>
                </div>
              )}
            </div>
          ) : (
            <button
              onClick={handleLogin}
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
            >
              Iniciar Sesión
            </button>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
