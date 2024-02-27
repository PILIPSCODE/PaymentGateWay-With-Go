package productcontrollers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var s snap.Client

type CustDetail struct {
	Fname, LName, Email, Phone, City string
	Item_Detail                      []ItemDetails
	Ongkir, Diskon                   int
}
type ItemDetails struct {
	ID, Name        string
	Price, Quantity int32
}

var CustDetails CustDetail

func setupGlobalMidtransConfig() {
	midtrans.ServerKey = os.Getenv("SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

}

func GenerateSnapReq() *snap.Request {
	var items []midtrans.ItemDetails
	var GrossAmt int
	fmt.Println()
	for _, itemData := range CustDetails.Item_Detail {
		item := midtrans.ItemDetails{
			ID:    itemData.ID,
			Name:  itemData.Name,
			Price: int64(itemData.Price),
			Qty:   int32(itemData.Quantity),
		}
		items = append(items, item)
		GrossAmt += int(itemData.Price) * int(itemData.Quantity)
	}
	GrossAmt += CustDetails.Ongkir
	GrossAmt += CustDetails.Diskon

	Shipping := midtrans.ItemDetails{
		ID:    "Shipping",
		Name:  "Ongkir",
		Price: int64(CustDetails.Ongkir),
		Qty:   1,
	}
	items = append(items, Shipping)

	Diskon := midtrans.ItemDetails{
		ID:    "Diskon",
		Name:  "Diskon",
		Price: int64(CustDetails.Diskon),
		Qty:   1,
	}
	items = append(items, Diskon)

	custAddress := &midtrans.CustomerAddress{
		FName: CustDetails.Fname,
		LName: CustDetails.LName,
		Phone: CustDetails.Phone,
		City:  CustDetails.City,
	}

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-GO-ID-" + Random(),
			GrossAmt: int64(GrossAmt),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    CustDetails.Fname,
			LName:    CustDetails.LName,
			Email:    CustDetails.Email,
			Phone:    CustDetails.Phone,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items:           &items,
	}
	return snapReq
}

func createTransactionWithGlobalConfig(c *gin.Context) {
	res, err := snap.CreateTransaction(GenerateSnapReq())
	if err != nil {
		fmt.Println("Snap Request Error", err.GetMessage())
	}

	c.JSON(http.StatusOK, gin.H{"Res": res, "Cust": CustDetails})

}

func Initial(c *gin.Context) {
	if err := c.ShouldBindJSON(&CustDetails); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	setupGlobalMidtransConfig()
	createTransactionWithGlobalConfig(c)

}

func Random() string {
	time.Sleep(500 * time.Millisecond)
	return strconv.FormatInt(time.Now().Unix(), 10)
}
