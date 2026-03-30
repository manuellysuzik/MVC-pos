package handler

import (
	clienthdl "ecommerce-api/internal/handler/cliente"
	pedidohdl "ecommerce-api/internal/handler/pedido"
	produtohdl "ecommerce-api/internal/handler/produto"
	"ecommerce-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	clienteHandler *clienthdl.ClienteHandler,
	produtoHandler *produtohdl.ProdutoHandler,
	pedidoHandler *pedidohdl.PedidoHandler,
	apiKey string,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	api.Use(middleware.APIKeyMiddleware(apiKey))
	{
		clientes := api.Group("/clientes")
		{
			clientes.POST("", clienteHandler.Create)
			clientes.GET("", clienteHandler.FindAll)
			clientes.GET("/count", clienteHandler.Count)
			clientes.GET("/search", clienteHandler.FindByName)
			clientes.GET("/:id", clienteHandler.FindByID)
			clientes.PUT("/:id", clienteHandler.Update)
			clientes.DELETE("/:id", clienteHandler.Delete)
		}

		produtos := api.Group("/produtos")
		{
			produtos.POST("", produtoHandler.Create)
			produtos.GET("", produtoHandler.FindAll)
			produtos.GET("/count", produtoHandler.Count)
			produtos.GET("/search", produtoHandler.FindByName)
			produtos.GET("/:id", produtoHandler.FindByID)
			produtos.PUT("/:id", produtoHandler.Update)
			produtos.DELETE("/:id", produtoHandler.Delete)
		}

		pedidos := api.Group("/pedidos")
		{
			pedidos.POST("", pedidoHandler.Create)
			pedidos.GET("", pedidoHandler.FindAll)
			pedidos.GET("/count", pedidoHandler.Count)
			pedidos.GET("/search", pedidoHandler.FindByStatus)
			pedidos.GET("/cliente/:id", pedidoHandler.FindByClienteID)
			pedidos.GET("/:id", pedidoHandler.FindByID)
			pedidos.PUT("/:id", pedidoHandler.Update)
			pedidos.DELETE("/:id", pedidoHandler.Delete)
		}
	}

	return r
}
