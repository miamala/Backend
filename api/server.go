package api

import "miamala/api/codegen"

func RegisterHandlers(router codegen.EchoRouter, si codegen.ServerInterface) {
	wrapper := codegen.ServerInterfaceWrapper{
		Handler: si,
	}
	router.GET("/transactions", wrapper.GetTransactions)
	router.POST("/transactions", wrapper.PostTransactions)
	router.DELETE("/transactions/:id", wrapper.DeleteTransactionsTransactionId)
	router.GET("/transactions/:id", wrapper.GetTransactionsTransactionId)
}
