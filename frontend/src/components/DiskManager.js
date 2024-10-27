// DiskManager.js
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Navbar from './Navbar';
import Swal from 'sweetalert2';
import { FaHdd } from 'react-icons/fa';

const DiskManager = () => {
  const [disks, setDisks] = useState([]);
  const [selectedPartitions, setSelectedPartitions] = useState(null); // Para las particiones del disco seleccionado
  const [modalOpen, setModalOpen] = useState(false); // Controlar la ventana emergente

  useEffect(() => {
    const storedDisks = JSON.parse(localStorage.getItem('disks')) || [];
    setDisks(storedDisks);
  }, []);

  // Función para manejar la solicitud de particiones
  const handleDiskClick = async (diskPath) => {
    const data = { commands: [`listpartitions -path=${diskPath}`] };

    try {
      const response = await axios.post('http://54.152.52.39:8080/analyze', data);
      const partitions = response.data.response;

      if (partitions === false || partitions.length === 0) {
        Swal.fire({
          icon: 'info',
          title: 'Sin particiones',
          text: `El disco ${diskPath} no tiene particiones.`,
        });
      } else {
        setSelectedPartitions(partitions);
        setModalOpen(true); // Abrir la ventana emergente con las particiones
      }
    } catch (error) {
      console.error('Error fetching partitions:', error);
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: 'Hubo un error al obtener las particiones del disco.',
      });
    }
  };

  // Función para cerrar la ventana emergente
  const handleCloseModal = () => {
    setModalOpen(false);
    setSelectedPartitions(null); // Limpiar las particiones seleccionadas
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6 flex flex-col">
      <Navbar />
      <div className="w-full max-w-4xl mx-auto bg-gray-800 rounded-lg shadow-lg p-6 mt-4">
        <h1 className="text-3xl font-bold mb-6 text-center">Disk Manager</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {disks.length > 0 ? (
            disks.map((diskPath, index) => {
              const diskName = diskPath.split('/').pop(); // Obtener el nombre del disco
              return (
                <div
                  key={index}
                  className="bg-gray-700 rounded-lg p-4 flex flex-col items-center justify-center cursor-pointer hover:bg-gray-600 transition-colors"
                  onClick={() => handleDiskClick(diskPath)}
                >
                  <FaHdd className="text-6xl text-green-400 mb-4" />
                  <p className="text-lg">{diskName}</p>
                </div>
              );
            })
          ) : (
            <p className="text-gray-500">No hay discos disponibles.</p>
          )}
        </div>
      </div>

      {/* Ventana emergente para mostrar particiones */}
      {modalOpen && selectedPartitions && (
        <div className="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
          <div className="bg-gray-800 text-white p-6 rounded-lg shadow-lg w-full max-w-3xl">
            <h2 className="text-2xl font-bold mb-4">Particiones del disco</h2>
            <div className="space-y-4">
              {selectedPartitions.map((partition, index) => (
                <div key={index} className="bg-gray-700 p-4 rounded-lg shadow-md">
                  <p><strong>Nombre:</strong> {partition.name}</p>
                  <p><strong>Tipo:</strong> {partition.type}</p>
                  <p><strong>Tamaño:</strong> {partition.size} bytes</p>
                  <p><strong>Estado:</strong> {partition.status === '0' ? 'Desmontada' : 'Montada'}</p>
                  <p><strong>Inicio:</strong> {partition.start}</p>
                </div>
              ))}
            </div>
            <button
              onClick={handleCloseModal}
              className="mt-6 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition duration-200"
            >
              Cerrar
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default DiskManager;
