package sales

import (
	"bytes"
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func Salesreport(c *gin.Context) {
	sales := c.Request.FormValue("salesreport")
	search := strings.ToLower(sales)
	switch search {
	case "daily":
		dailySalesReport(c)
	case "weekly":
		weeklySalesReport(c)
	case "monthly":
		monthlySalesReport(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid sales report type"})
	}
}
func dailySalesReport(c *gin.Context) {
	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.Add(24 * time.Hour)

	var order []model.Orderitems
	if err := database.DB.Preload("Order").Preload("Product").Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch orders"})
		return
	}
	generatepdf(c, order)
}

func weeklySalesReport(c *gin.Context) {
	startDate := time.Now().Truncate(24 * 7 * time.Hour)
	endDate := startDate.Add(7 * 24 * time.Hour)

	var order []model.Orderitems
	if err := database.DB.Preload("Order").Preload("Product").Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch orders"})
		return
	}
	generatepdf(c, order)
}

func monthlySalesReport(c *gin.Context) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var order []model.Orderitems
	if err := database.DB.Preload("Order").Preload("Product").Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch order"})
		return
	}
	generatepdf(c, order)
}
func generatepdf(c *gin.Context, order []model.Orderitems) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(10)

	pdf.Cell(40, 10, "Order Information:")
	pdf.Ln(10)

	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(20, 10, "Order ID", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Order Date", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Product Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Order Status", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Amount", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	var totalAmount float64
	for _, item := range order {
		if item.Orderstatus == "delivered" {
			pdf.CellFormat(20, 10, strconv.Itoa(int(item.OrderID)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(30, 10, item.Order.Orderdate.Format("2006-01-02"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(40, 10, item.Product.Product_name, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 10, item.Orderstatus, "1", 0, "", false, 0, "")
			pdf.CellFormat(30, 10, strconv.Itoa(int(item.Quantity)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(30, 10, strconv.FormatFloat(float64(item.Subtotal), 'f', 2, 64), "1", 0, "R", true, 0, "")
			pdf.Ln(10)
			totalAmount += float64(item.Subtotal)
		}
	}

	pdf.CellFormat(100, 10, "Grand Total", "1", 0, "C", false, 0, "")
	pdf.CellFormat(90, 10, strconv.FormatFloat(totalAmount, 'f', 2, 64), "1", 0, "R", true, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")

	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
