import React, { useState, useRef, useEffect } from 'react';
import axios from 'axios';
import { FaUpload, FaPlay, FaBroom, FaTerminal, FaSpinner } from 'react-icons/fa';
import Navbar from './Navbar';

const Console = () => {
  const [command, setCommand] = useState('');
  const [responses, setResponses] = useState([]);
  const [loading, setLoading] = useState(false);  // Estado de carga
  const consoleEndRef = useRef(null);  // Referencia para el final de la consola

  // Función para desplazar la consola hacia abajo
  const scrollToBottom = () => {
    consoleEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  // Ejecutar el scroll cada vez que haya nuevas respuestas
  useEffect(() => {
    scrollToBottom();
  }, [responses]);


  // Función para ejecutar comandos línea por línea
  const handleSubmit = async (e) => {
    e.preventDefault();

    // Dividir el comando por líneas
    const commandsArray = command.split('\n').filter(line => line.trim() !== '');
    console.log('Commands:', commandsArray);
    setLoading(true);  // Iniciar la carga

    try {
      for (let cmd of commandsArray) {
        const data = { commands: [cmd] };  // Enviar un solo comando a la vez
        const response = await axios.post('http://54.152.52.39:8080/analyze', data);

        console.log('Response:', response.data);

        //Si hay mensajes de exito sobre el comando mkdisk, guardamos en el array disks de localStorage el path del disco
        //Ejemplo de respuesta {"command": "mkdisk", "message": "> Comando mkdisk con parámetros: -path=/home/fernando/discos/disco1.mia -size=15 -unit=k ejecutado exitosamente"}
        //El path no tiene una posicion fija, por lo que se busca con la expresion regular
        
        if (response.data[0].command === 'mkdisk' && response.data[0].message.includes('ejecutado exitosamente')){
          const regex = /-path=(".*?"|\S+)/;
          const path = response.data[0].message.match(regex)[1];
          const disks = JSON.parse(localStorage.getItem('disks'));
          disks.push(path);
          localStorage.setItem('disks', JSON.stringify(disks));
        }

        //Si es rmDisk, se elimina el disco del array disks de localStorage
        //Ejemplo de respuesta {"command": "rmdisk", "message": "> Comando rmdisk con parámetros: -path=/home/fernando/discos/disco1.mia ejecutado exitosamente"}

        if (response.data[0].command === 'rmdisk' && response.data[0].message.includes('ejecutado exitosamente')){
          const regex = /-path=(".*?"|\S+)/;
          const path = response.data[0].message.match(regex)[1];
          const disks = JSON.parse(localStorage.getItem('disks'));
          const index = disks.indexOf(path);
          if (index > -1) {
            disks.splice(index, 1);
          }
          localStorage.setItem('disks', JSON.stringify(disks));
        }

        // Añadir cada respuesta al estado, sin borrar las anteriores
        setResponses((prevResponses) => [...prevResponses, ...response.data]);

        // Simulación de retraso para efectos visuales (opcional)
        await new Promise((r) => setTimeout(r, 500));
      }
    } catch (error) {
      console.error('Error submitting command:', error);
      setResponses((prevResponses) => [...prevResponses, { command: 'Error', message: 'Error processing command.' }]);
    } finally {
      setLoading(false);  // Terminar la carga
    }
  };

  // Función para cargar archivos de comandos
  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();

    reader.onload = (event) => {
      const fileContent = event.target.result;
      setCommand(fileContent);  // Pegar contenido en el textarea
    };

    if (file) {
      reader.readAsText(file);
    }
  };

  // Función para limpiar la consola
  const clearConsole = () => {
    setResponses([]);  // Limpiar el estado de las respuestas
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6 flex flex-col">
      <Navbar  />
      <div className="w-full max-w-4xl mx-auto bg-gray-800 rounded-lg shadow-lg p-6 mt-4">
        <h1 className="text-3xl font-bold flex items-center mb-4">
          <FaTerminal className="mr-2 text-green-400" /> Command Console
        </h1>

        {/* Formulario para ingresar los comandos */}
        <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
          <textarea
            value={command}
            onChange={(e) => setCommand(e.target.value)}
            rows="10"
            className="w-full p-4 bg-gray-900 text-green-400 font-mono text-sm border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner transition-all duration-200"
            placeholder="Enter your commands here..."
          />

          {/* Botones: Upload, Execute y Clear Console */}
          <div className="flex space-x-4">
            <label className="flex items-center space-x-2 p-2 bg-blue-500 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer transition-transform transform hover:scale-105">
              <FaUpload className="text-white" />
              <span>Upload File</span>
              <input
                type="file"
                accept="*"
                onChange={handleFileUpload}
                className="hidden"
              />
            </label>

            <button
              type="submit"
              className="flex items-center space-x-2 p-2 bg-green-500 rounded-lg hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 transition-transform transform hover:scale-105"
              disabled={loading}
            >
              {loading ? (
                <>
                  <FaSpinner className="animate-spin text-white" />
                  <span>Executing...</span>
                </>
              ) : (
                <>
                  <FaPlay className="text-white" />
                  <span>Execute</span>
                </>
              )}
            </button>

            <button
              type="button"
              onClick={clearConsole}
              className="flex items-center space-x-2 p-2 bg-red-500 rounded-lg hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 transition-transform transform hover:scale-105"
            >
              <FaBroom className="text-white" />
              <span>Clear Console</span>
            </button>
          </div>
        </form>

        {/* Consola para mostrar los resultados */}
        <div className="mt-6 bg-gray-900 p-4 rounded-md border border-gray-700">
          <h2 className="text-lg font-semibold mb-2">Output:</h2>

          <div className="p-4 bg-black text-green-400 font-mono text-sm rounded h-48 overflow-y-auto whitespace-pre-wrap shadow-inner transition-opacity duration-500 ease-in-out">
            {responses.length > 0 ? (
              responses.map((response, index) => (
                <p key={index} className="mb-2">
                  {response.message}
                </p>
              ))
            ) : (
              <p className="text-gray-500">No output yet...</p>
            )}
            <div ref={consoleEndRef} />  {/* Referencia para el desplazamiento automático */}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Console;
