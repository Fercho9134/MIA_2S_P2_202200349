// Login.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import Swal from 'sweetalert2';
import Navbar from './Navbar';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [userId, setUserId] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = {
      commands: [`login -user=${username} -pass=${password} -id=${userId}`]
    };

    try {
      const response = await axios.post('http://54.152.52.39:8080/analyze', data);
      const { message, response: isSuccess } = response.data;

      if (isSuccess) {
        Swal.fire({
          icon: 'success',
          title: 'Inicio de sesión exitoso',
          text: message,
        });
        localStorage.setItem('isLoggedIn', 'true');
        navigate('/console');
      } else {
        Swal.fire({
          icon: 'error',
          title: 'Error al iniciar sesión',
          text: message,
        });
      }
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Ocurrió un error inesperado',
        text: 'Intente de nuevo más tarde.',
      });
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6 flex flex-col">
      <Navbar />
      <div className="w-full max-w-4xl mx-auto bg-gray-800 rounded-lg shadow-lg p-8 mt-4">
        <h1 className="text-3xl font-bold mb-6 text-center">Login</h1>
        <form onSubmit={handleSubmit} className="flex flex-col space-y-6">
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="p-3 bg-gray-700 text-white border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-200"
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="p-3 bg-gray-700 text-white border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-200"
            required
          />
          <input
            type="text"
            placeholder="Partition ID"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            className="p-3 bg-gray-700 text-white border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-200"
            required
          />
          <button
            type="submit"
            className="bg-blue-600 text-white px-6 py-2 rounded-lg shadow hover:bg-blue-700 transition duration-200"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;
