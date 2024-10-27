// App.js
import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Console from './components/Console'; // Ajusta según la ruta de tu archivo
import Login from './components/Login';
import DiskManager from './components/DiskManager';


const App = () => {
  //Creamos por defecto en localStorage el estado de login
  if (!localStorage.getItem('isLoggedIn')) {
    localStorage.setItem('isLoggedIn', 'false');
  }

  //Creamos un arreglos de discos vacio si no existe
  if (!localStorage.getItem('disks')) {
    localStorage.setItem('disks', JSON.stringify([]));
  }

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/console" />} /> {/* Redirige a la consola por defecto */}
        <Route path="/console" element={<Console />} />
        <Route path="/login" element={<Login />} />
        <Route path="/disk-manager" element={<DiskManager />} />
      </Routes>
    </Router>
  );
};

export default App;
