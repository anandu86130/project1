package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func Paymenthandler(orderId string, amount int) (string, error) {

	client := razorpay.NewClient(os.Getenv("RAZORPAY_ID"), os.Getenv("RAZORPAY_SECRET"))

	data := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  orderId,
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", errors.New("payment not initiated")
	}
	razorId, _ := body["id"].(string)
	return razorId, nil
}

func Paymentconfirmation(c *gin.Context) {
	var response struct {
		OrderID   string `json:"order_id"`
		PaymentID string `json:"payment_id"`
		Signature string `json:"signature"`
	}
	if err := c.BindJSON(&response); err != nil {
		fmt.Println("Error", err)
		return
	}
	err := Razorpaymentverification(response.OrderID, response.PaymentID, response.Signature)
	if err != nil {
		fmt.Println("Error", err)
		paymentfailed := model.Paymentdetails{
			PaymentId:     response.PaymentID,
			Paymentstatus: "failed",
		}
		database.DB.Where("order_id=?", response.OrderID).Updates(paymentfailed)
		return
	} else {
		fmt.Println("Payment done.")
	}
	payment := model.Paymentdetails{
		PaymentId:     response.PaymentID,
		Paymentstatus: "Success",
	}
	database.DB.Where("order_id=?", response.OrderID).Updates(&payment)
	c.JSON(http.StatusOK, gin.H{"Message": "Payment recieved successfullly"})
}

func Razorpaymentverification(OrderID, PaymentID, Signature string) error {
	signature := Signature
	secret := "MtKBBcRl8jkG7VzlGI2b3eoj"
	data := OrderID + "|" + PaymentID
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("INVALID SIGNATURE")
	} else {
		return nil
	}

}
