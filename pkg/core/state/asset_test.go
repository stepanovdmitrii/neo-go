package state

import (
	"testing"

	"github.com/CityOfZion/neo-go/pkg/core/transaction"
	"github.com/CityOfZion/neo-go/pkg/crypto/keys"
	"github.com/CityOfZion/neo-go/pkg/internal/random"
	"github.com/CityOfZion/neo-go/pkg/io"
	"github.com/CityOfZion/neo-go/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeAssetState(t *testing.T) {
	asset := &Asset{
		ID:         random.Uint256(),
		AssetType:  transaction.Token,
		Name:       "super cool token",
		Amount:     util.Fixed8(1000000),
		Available:  util.Fixed8(100),
		Precision:  0,
		FeeMode:    feeMode,
		Owner:      keys.PublicKey{},
		Admin:      random.Uint160(),
		Issuer:     random.Uint160(),
		Expiration: 10,
		IsFrozen:   false,
	}

	buf := io.NewBufBinWriter()
	asset.EncodeBinary(buf.BinWriter)
	assert.Nil(t, buf.Err)
	assetDecode := &Asset{}
	r := io.NewBinReaderFromBuf(buf.Bytes())
	assetDecode.DecodeBinary(r)
	assert.Nil(t, r.Err)
	assert.Equal(t, asset, assetDecode)
}

func TestAssetState_GetName_NEO(t *testing.T) {
	asset := &Asset{AssetType: transaction.GoverningToken}
	assert.Equal(t, "NEO", asset.GetName())
}

func TestAssetState_GetName_NEOGas(t *testing.T) {
	asset := &Asset{AssetType: transaction.UtilityToken}
	assert.Equal(t, "NEOGas", asset.GetName())
}
