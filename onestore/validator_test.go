package onestore

import (
	"context"
	"fmt"
	"testing"
)

var client_id = "package"
var client_secret = "client_secret"

func init() {

}

func TestNew(t *testing.T) {
	t.Parallel()

	New(client_id, client_secret, "purchaseToken")
}

func TestResponse(t *testing.T) {
	client := New(client_id, client_secret, "purchaseToken")

	ctx := context.Background()
	resp, err := client.Verify(ctx, "package", "productID", "purchaseToken")
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(resp.PurchaseState)
	fmt.Println(resp.PurchaseId)
}

func TestHandleError(t *testing.T) {
	// Exception scenario
	client := New(client_id, client_secret, "SANDBOX")

	ctx := context.Background()
	_, err := client.Verify(ctx, "package", "productID", "purchaseToken")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestBadURL(t *testing.T) {
	client := New(client_id, client_secret, "SANBOX")

	ctx := context.Background()
	_, err := client.Verify(ctx, "package", "productID", "purchaseToken")
	if err != nil {
		t.Errorf("%s", err)
	}
}
