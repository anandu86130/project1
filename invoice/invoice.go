package invoice

import (
	"bytes"
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func Invoicedownload(c *gin.Context) {
	userid := c.GetUint("userid")
	var user model.UserModel
	if err := database.DB.Where("user_id=?", userid).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
		return
	}
	ID := c.Param("ID")
	if ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give order id"})
		return
	}
	var orderitems []model.Orderitems
	if err := database.DB.Preload("Product").Preload("Order").Where("order_id", ID).First(&orderitems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
		return
	}

	var order model.Order
	if err := database.DB.Preload("Address").Where("ID = ?", ID).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
		return
	}

	for _, status := range orderitems {
		if status.Orderstatus != "delivered" {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "The order is not delivered"})
			return
		}
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice to "+user.Name)

	pdf.Ln(10)
	pdf.Cell(40, 10, "Address:")
	pdf.Ln(10)
	pdf.Cell(40, 10, order.Address.Address)
	pdf.Ln(10)
	pdf.Cell(40, 10, order.Address.City+", "+order.Address.State+", "+order.Address.Country+" - "+order.Address.Pincode)
	pdf.Ln(20)

	pdf.Ln(10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(40, 10, "Order Item", "1", 0, "", true, 0, "")
	pdf.CellFormat(80, 10, "Description", "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 10, "Amount", "1", 0, "", true, 0, "")
	pdf.Ln(10)

	for _, items := range orderitems {
		quantityStr := strconv.Itoa(int(items.Quantity))
		amountStr := strconv.Itoa(int(items.Subtotal))
		pdf.CellFormat(40, 10, items.Product.Product_name, "1", 0, "", false, 0, "")
		pdf.CellFormat(80, 10, items.Product.Description, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, quantityStr, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, amountStr, "1", 0, "", false, 0, "")
		pdf.Ln(10)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")

	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
