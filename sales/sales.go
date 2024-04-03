package sales

import (
	"bytes"
	"fmt"
	"net/http"
	"project1/database"
	"project1/model"
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
	// case "monthly":
	// 	monthlySalesReport(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid sales report type"})
	}
}
func dailySalesReport(c *gin.Context) {
	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.Add(24 * time.Hour)

	var order []model.Orderitems
	if err := database.DB.Preload("Order").Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch orders"})
		return
	}
	var salesData []gin.H
	var total int
	for _, orders := range order {
		salesData = append(salesData, gin.H{
			"OrderID":     orders.Order.ID,
			"OrderDate":   orders.Order.Orderdate,
			"Payment":     orders.Order.Paymentmethod,
			"OrderStatus": orders.Orderstatus,
			"TotalAmount": orders.Order.Totalamount,
		})
		total += int(orders.Order.Totalamount)
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"SalesReport": salesData,
	// 	"Grandtotal":  total,
	// })
	generatePDF(c, salesData, total)
}

func weeklySalesReport(c *gin.Context) {
	startDate := time.Now().Truncate(24 * 7 * time.Hour)
	endDate := startDate.Add(7 * 24 * time.Hour)

	var order []model.Orderitems
	if err := database.DB.Preload("Order").Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch orders"})
		return
	}
	var salesData []gin.H
	var total int
	for _, orders := range order {
		salesData = append(salesData, gin.H{
			"OrderID":     orders.Order.ID,
			"OrderDate":   orders.Order.Orderdate,
			"Payment":     orders.Order.Paymentmethod,
			"OrderStatus": orders.Orderstatus,
			"TotalAmount": orders.Order.Totalamount,
		})
		total += int(orders.Order.Totalamount)

		// c.JSON(http.StatusOK, gin.H{
		// 	"SalesReport": salesData,
		// 	"GrandTotal":  total,
		// })
	}
	generatePDF(c, salesData, total)
}

func generatePDF(c *gin.Context, salesData []gin.H, grandTotal int) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(10)

	for _, data := range salesData {
		for key, value := range data {
			strValue := fmt.Sprintf("%v", value)
			pdf.CellFormat(40, 10, key+":", "", 0, "", false, 0, "")
			pdf.CellFormat(40, 10, strValue, "", 1, "", false, 0, "")
		}
		pdf.Ln(5)
	}

	strGrandTotal := fmt.Sprintf("%d", grandTotal)
	pdf.CellFormat(40, 10, "Grand Total:", "", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, strGrandTotal, "", 1, "", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
