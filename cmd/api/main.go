package main

import (
	"log"

	"github.com/pedropiedade7/mundo-invest-test/internal/app"
)

func main() {
	log.Println("Iniciando a aplicação Mundo Invest...")

	api, err := app.New()
	if err != nil {
		log.Fatalf("Falha ao inicializar a aplicação: %v", err)
	}
	defer api.Close()

	log.Println("Servidor Gin rodando em http://localhost:8080")
	if err := api.Run(":8080"); err != nil {
		log.Fatalf("Erro ao rodar o servidor: %v", err)
	}
}
