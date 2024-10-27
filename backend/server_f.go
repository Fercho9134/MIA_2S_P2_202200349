package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proyecto1/Analyzer"
	"proyecto1/DiskManagement"
	"syscall"
	"time"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/prueba", Analyzer.ImprimirHandler)
	mux.HandleFunc("/analyze", Analyzer.AnalyzeHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: enableCors(mux),
	}

	// Canal para capturar señales del sistema
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Goroutine para iniciar el servidor
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en ListenAndServe: %v", err)
		}
	}()

	log.Println("Servidor iniciado en :8080 :b")

	// Espera la señal de interrupción
	<-stop

	log.Println("Apagando servidor...")

	// Llamada a la función de limpieza antes de apagar el servidor
	DiskManagement.Clean()

	// Contexto con timeout para darle tiempo al servidor de cerrar correctamente
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}
