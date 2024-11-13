package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

func Razopay(razopaykey, razopaysecret string) (string, error) {
	client := razorpay.NewClient(razopaykey, razopaysecret)

	data := map[string]interface{}{
		"amount":   600,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}
	idFromResponse, _ := body["id"].(string)
	return idFromResponse, nil
}

func VerifyPayment(orderID, paymentID, providedSignature, razorpaySecret string) bool {
    // Concatenate orderID and paymentID with "|" as Razorpay requires
    data := orderID + "|" + paymentID

	fmt.Println("secret = ",razorpaySecret)
    // Create HMAC SHA-256 hash using the Razorpay secret key
    h := hmac.New(sha256.New, []byte(razorpaySecret))
    h.Write([]byte(data)) // Write the concatenated data to HMAC

    // Convert the generated HMAC to a hexadecimal string
    generatedSignature := hex.EncodeToString(h.Sum(nil))

    // Compare the generated signature with the provided signature
    isValid := generatedSignature == providedSignature
    fmt.Println("Order Id = ", orderID)
    fmt.Println("PaymentId = ", paymentID)
    fmt.Println("Generated Signature:", generatedSignature)
    fmt.Println("Provided Signature:", providedSignature)
    fmt.Println("Signature Match:", isValid)

    return isValid
}