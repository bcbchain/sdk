package helper

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl"
	"blockchain/smcsdk/sdkimpl/object"
	"bytes"
	"errors"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/go-crypto"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// BlockChainHelper block chain helper information
type BlockChainHelper struct {
	smc sdk.ISmartContract
}

var _ sdk.IBlockChainHelper = (*BlockChainHelper)(nil)
var _ sdkimpl.IAcquireSMC = (*BlockChainHelper)(nil)

// SMC get smartContract object
func (bh *BlockChainHelper) SMC() sdk.ISmartContract { return bh.smc }

// SetSMC set smartContract object
func (bh *BlockChainHelper) SetSMC(smc sdk.ISmartContract) { bh.smc = smc }

// CalcAccountFromPubKey calculate account address from pubKey
func (bh *BlockChainHelper) CalcAccountFromPubKey(pubKey types.PubKey) types.Address {
	sdk.Require(pubKey != nil && len(pubKey) == 32,
		types.ErrInvalidParameter, "invalid pubKey")

	pk := crypto.PubKeyEd25519FromBytes(pubKey)
	return pk.Address(bh.smc.Helper().GenesisHelper().ChainID())
}

// CalcAccountFromName calculate account address from name
func (bh *BlockChainHelper) CalcAccountFromName(name string, orgID string) types.Address {
	return bh.CalcContractAddress(name, "", orgID)
}

// nolint unhandled
// CalcContractAddress calculate contract address from name、version and owner
func (bh *BlockChainHelper) CalcContractAddress(name string, version string, orgID string) types.Address {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte(bh.smc.Helper().GenesisHelper().ChainID()))
	hasherSHA3256.Write([]byte(name))
	hasherSHA3256.Write([]byte(version))
	hasherSHA3256.Write([]byte(orgID))
	sha := hasherSHA3256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha) // does not error
	rpd := hasherRIPEMD160.Sum(nil)

	hasher := ripemd160.New()
	hasher.Write(rpd)
	md := hasher.Sum(nil)

	addr := make([]byte, 0, 0)
	addr = append(addr, rpd...)
	addr = append(addr, md[:4]...)

	return bh.smc.Helper().GenesisHelper().ChainID() + base58.Encode(addr)
}

// nolint unhandled
// CalcOrgID calculate organization ID
func (bh *BlockChainHelper) CalcOrgID(name string) string {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte(name))
	sha := hasherSHA3256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha) // does not error
	rpd := hasherRIPEMD160.Sum(nil)

	hasher := ripemd160.New()
	hasher.Write(rpd)
	md := hasher.Sum(nil)

	addr := make([]byte, 0, 0)
	addr = append(addr, rpd...)
	addr = append(addr, md[:4]...)

	return "org" + base58.Encode(addr)
}

// GetBlock get block data with height
func (bh *BlockChainHelper) GetBlock(height int64) sdk.IBlock {
	if height <= 0 {
		return nil
	}

	transID := bh.smc.(*sdkimpl.SmartContract).LlState().TransID()
	v := sdkimpl.GetBlockFunc(transID, height)

	block := object.NewBlock(
		bh.smc,
		v.ChainID,
		v.Version,
		v.BlockHash,
		v.DataHash,
		v.Height,
		v.Time,
		v.NumTxs,
		v.ProposerAddress,
		v.RewardAddress,
		v.RandomNumber,
		v.LastBlockHash,
		v.LastCommitHash,
		v.LastAppHash,
		v.LastFee,
	)

	return block
}

// nolint unhandled
// CheckAddress check address and return result
func (bh *BlockChainHelper) CheckAddress(addr types.Address) error {
	chainID := bh.smc.Block().ChainID()
	if strings.HasPrefix(addr, chainID) == false {
		return errors.New("Address chainID is error! ")
	}

	base58Addr := strings.Replace(addr, chainID, "", 1)
	addrData := base58.Decode(base58Addr)
	addrLen := len(addrData)
	if addrLen < 4 {
		return errors.New("Base58Addr parse error! ")
	}

	r160 := ripemd160.New()
	r160.Write(addrData[:addrLen-4])
	md := r160.Sum(nil)

	if bytes.Compare(md[:4], addrData[addrLen-4:]) != 0 {
		return errors.New("Address checksum is error! ")
	}
	return nil
}

// GetCurrentBlock get current block data
func GetCurrentBlock(smc sdk.ISmartContract) sdk.IBlock {
	// get current block if height is 0
	v := sdkimpl.GetBlockFunc(smc.(*sdkimpl.SmartContract).LlState().TransID(), 0)

	block := object.NewBlock(
		smc,
		v.ChainID,
		v.Version,
		v.BlockHash,
		v.DataHash,
		v.Height,
		v.Time,
		v.NumTxs,
		v.ProposerAddress,
		v.RewardAddress,
		v.RandomNumber,
		v.LastBlockHash,
		v.LastCommitHash,
		v.LastAppHash,
		v.LastFee,
	)

	return block
}
