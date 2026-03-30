package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/handler"
	clienthdl "ecommerce-api/internal/handler/cliente"
	pedidohdl "ecommerce-api/internal/handler/pedido"
	produtohdl "ecommerce-api/internal/handler/produto"
	clienterepo "ecommerce-api/internal/repository/cliente"
	pedidorepo "ecommerce-api/internal/repository/pedido"
	produtorepo "ecommerce-api/internal/repository/produto"
	clientesvc "ecommerce-api/internal/service/cliente"
	pedidosvc "ecommerce-api/internal/service/pedido"
	produtosvc "ecommerce-api/internal/service/produto"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./ecommerce.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&entity.Cliente{},
		&entity.Produto{},
		&entity.Pedido{},
		&entity.ItemPedido{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	clienteRepo := clienterepo.New(db)
	produtoRepo := produtorepo.New(db)
	pedidoRepo := pedidorepo.New(db)

	clienteSvc := clientesvc.New(&clienteRepo)
	produtoSvc := produtosvc.New(&produtoRepo)
	pedidoSvc := pedidosvc.New(&pedidoRepo)

	clienteHandler := clienthdl.New(clienteSvc)
	produtoHandler := produtohdl.New(produtoSvc)
	pedidoHandler := pedidohdl.New(pedidoSvc)

	r := handler.SetupRouter(clienteHandler, produtoHandler, pedidoHandler, apiKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
