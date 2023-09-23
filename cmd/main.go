package main

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"miamala/api"
	"net/http"
	"strconv"
)

type Transaction struct {
	ID       int     `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Type     string  `json:"type"`
}

type Server struct {
	DB *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		DB: db,
	}
}

func (s *Server) GetTransactions(ctx echo.Context) error {
	// Query the database to retrieve all transactions
	rows, err := s.DB.Query("SELECT id, amount, category, type FROM transactions")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Iterate over the rows and populate a slice of transactions
	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Category,
			&transaction.Type,
		); err != nil {
			fmt.Println("Error scanning row:", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		transactions = append(transactions, transaction)
	}
	return ctx.JSON(http.StatusOK, transactions)
}

func (s *Server) PostTransactions(ctx echo.Context) error {
	// Parse the request body into a Transaction struct
	var newTransaction Transaction
	if err := ctx.Bind(&newTransaction); err != nil {
		fmt.Println("Error parsing request body:", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	spew.Dump("cheeckio", newTransaction.Amount, newTransaction.Category, newTransaction.Type)

	// Insert the new transaction into,  the database
	_, err := s.DB.Exec("INSERT INTO transactions (amount, category, type) VALUES (?,?,?)",
		newTransaction.Amount,
		newTransaction.Category,
		newTransaction.Type,
	)

	spew.Dump("ndio hii", err)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return ctx.JSON(http.StatusCreated, newTransaction)
}

func (s *Server) DeleteTransactionsTransactionId(ctx echo.Context, transactionId int) error {
	// Get the transaction ID from the URL parameter
	transactionID, err := strconv.Atoi(ctx.Param("transactionId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Transaction ID"})
	}

	// Delete the specific transaction by ID from the database
	_, err = s.DB.Exec("DELETE FROM transactions WHERE id = ?", transactionID)
	if err != nil {
		fmt.Println("Error deleting from database:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (s *Server) GetTransactionsTransactionId(ctx echo.Context, transactionId int) error {
	// Get the transaction ID from the URL parameter
	transactionID, err := strconv.Atoi(ctx.Param("transactionId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Transaction ID"})
	}

	// Query the database to retrieve the specific transaction by ID
	row := s.DB.QueryRow("SELECT id, amount, category, type FROM transactions WHERE id = ?", transactionID)

	// Populate a Transaction struct with the retrieved data
	var transaction Transaction
	if err := row.Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Category,
		&transaction.Type,
	); err != nil {
		fmt.Println("Error scanning row:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return ctx.JSON(http.StatusOK, transaction)
}

func (s *Server) PutTransactionsTransactionId(ctx echo.Context, transactionId int) error {
	// Get the transaction ID from the URL parameter
	transactionID, err := strconv.Atoi(ctx.Param("transactionId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Transaction ID"})
	}

	// Parse the request body into a Transaction struct
	var updatedTransaction Transaction
	if err := ctx.Bind(&updatedTransaction); err != nil {
		fmt.Println("Error parsing request body:", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	// Update the specific transaction by ID in the database
	_, err = s.DB.Exec("UPDATE transactions SET amount = ?, category = ?, type = ? WHERE id = ?",
		updatedTransaction.Amount,
		updatedTransaction.Category,
		updatedTransaction.Type,
		transactionID,
	)
	if err != nil {
		fmt.Println("Error updating database:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return ctx.JSON(http.StatusOK, Transaction{})
}

func main() {
	// Initialize SQLite database connection
	db, err := sql.Open("sqlite3", "./transactions.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create the 'transactions' table if it doesn't exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount REAL,
			category TEXT,
			type TEXT
		)
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Error creating 'transactions' table:", err)
		return
	}

	//Start Echo
	e := echo.New()

	// Create a new Server instance with the database connection
	server := NewServer(db)

	api.RegisterHandlers(e, server)

	// Start the server
	err = e.Start(":8088")
	if err != nil {
		panic(err)
	}
}
