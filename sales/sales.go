package sales

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
			"OrderID":          orders.Order.ID,
			"OrderDate":        orders.Order.Orderdate,
			"PaymentMethod":    orders.Order.Paymentmethod,
			"OrderStatus":      orders.Orderstatus,
			"TotalSalesAmount": orders.Order.Totalamount,
		})
		total += int(orders.Order.Totalamount)
	}
	c.JSON(http.StatusOK, gin.H{
		"SalesReport": salesData,
		"Grandtotal":  total,
	})
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
			"OrderID":          orders.Order.ID,
			"OrderDate":        orders.Order.Orderdate,
			"PaymentMethod":    orders.Order.Paymentmethod,
			"OrderStatus":      orders.Orderstatus,
			// "PaymentStatus":    orders.Order.Paymentdetails.Paymentstatus,
			"TotalSalesAmount": orders.Order.Totalamount,
		})
		total += int(orders.Order.Totalamount)

		c.JSON(http.StatusOK, gin.H{
			"SalesReport": salesData,
			"GrandTotal":  total,
		})
	}
}

// pdf := gofpdf.New("P", "mm", "A4", "")
// pdf.AddPage()
// pdf.SetFont("Arial", "B", 16)
// pdf.Cell(40, 10, "Sales Report")startDate
// pdf.Ln(10)

// pdf.CellFormat(190, 10, "Period: "+.Format("2006-01-02")+" to "+endDate.Format("2006-01-02"), "", 1, "L", false, 0, "")

// for _, data := range salesData {
// 	pdf.CellFormat(190, 10, data.Date.String(), "", 1, "L", false, 0, "")
// 	pdf.CellFormat(95, 10, "Total Sales: "+fmt.Sprintf("%.2f", data.TotalSales), "", 0, "L", false, 0, "")
// 	pdf.CellFormat(95, 10, "Total Orders: "+strconv.Itoa(data.TotalOrders), "", 1, "L", false, 0, "")
// }

// // Save the PDF to a buffer
// var buf bytes.Buffer
// pdf.Output(&buf)

// // Set response headers for PDF
// c.Header("Content-Type", "application/pdf")
// c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")

// // Write the PDF buffer to the response
// c.Data(http.StatusOK, "application/pdf", buf.Bytes())
