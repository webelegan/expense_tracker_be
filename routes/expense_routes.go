package routes

import (
	"net/http"
	"strconv"
	"time"

	"expense-tracker-backend/database"
	"expense-tracker-backend/models"

	"github.com/gin-gonic/gin"
)

// SetupExpenseRoutes configures all expense-related routes
func SetupExpenseRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		expenses := api.Group("/expenses")
		{
			expenses.GET("", getExpenses)
			expenses.POST("", createExpense)
			expenses.PUT("/:id", updateExpense)
			expenses.DELETE("/:id", deleteExpense)
		}
	}
}

// getExpenses handles GET /api/expenses
// Query parameters: startDate, endDate (format: YYYY-MM-DD)
func getExpenses(c *gin.Context) {
	var expenses []models.Expense
	query := database.DB

	// Filter by date range if provided
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use YYYY-MM-DD"})
			return
		}
		// Start of the day
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		query = query.Where("date >= ?", start)
	}

	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD"})
			return
		}
		// End of the day (23:59:59)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
		query = query.Where("date <= ?", end)
	}

	if err := query.Order("date DESC").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// createExpense handles POST /api/expenses
func createExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// updateExpense handles PUT /api/expenses/:id
func updateExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var expense models.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var updateData models.Expense
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	expense.Title = updateData.Title
	expense.Category = updateData.Category
	expense.Amount = updateData.Amount
	expense.Date = updateData.Date
	if updateData.Type != "" {
		expense.Type = updateData.Type
	}

	if err := database.DB.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// deleteExpense handles DELETE /api/expenses/:id
func deleteExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var expense models.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	if err := database.DB.Delete(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

