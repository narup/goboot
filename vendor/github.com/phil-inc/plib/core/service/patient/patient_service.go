package patient

import (
	"context"

	"github.com/phil-inc/plib/core/data"
)

// GetDefaultAddress get patient default address
func GetDefaultAddress(ctx context.Context, addressID string) *data.Address {
	session := data.Session()
	defer session.Close()

	addr := new(data.Address)
	err := session.FindByID(addressID, addr)
	if err != nil {
		return nil
	}
	return addr
}
