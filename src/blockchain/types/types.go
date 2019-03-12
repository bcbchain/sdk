package types

import (
	"time"

	"github.com/tendermint/go-crypto"

	"github.com/tendermint/tmlibs/common"
)

// Address 地址用string
type Address = string

// Hash uses for public key and others, SHA3-256
type Hash = common.HexBytes

// HashLen -- Size of Hash is 32 bytes
const HashLen = 32

//MaxSizeNote maximum length of note of transaction
const MaxSizeNote = 256

// PubKey uses for public key and others, PubKeyEd25519
type PubKey = common.HexBytes

// ModuleHealth 啥意思我也不知道，請加注釋 TODO
type ModuleHealth struct {
	Tm     time.Time
	Status int
}

// Health 健康情況？TODO
type Health struct {
	Tm        time.Time
	SubHealth map[string]ModuleHealth
}

// Query abci query
type Query struct {
	QueryKey string
}

// Transaction 定义交易数据结构
type Transaction struct {
	Nonce    uint64    `json:"nonce"`    // 交易发起者发起交易的计数值，从1开始，必须单调增长，增长步长为1。
	GasLimit int64     `json:"gasLimit"` // 交易发起者愿意为执行此次交易支付的GAS数量的最大值。
	Note     string    `json:"note"`     // UTF-8编码的备注信息，要求小于256个字符。
	Messages []Message `json:"messages"` // 交易消息，RLP编码格式。
}

// Message - a contract method with its params
type Message struct {
	Contract Address           `json:"contract"` //调用合约地址
	MethodID uint32            `json:"methodID"` //方法id
	Items    []common.HexBytes `json:"items"`    //调用参数
}

//Ed25519Sig 定义加密算法结构
type Ed25519Sig struct {
	SigType  string
	PubKey   crypto.PubKeyEd25519
	SigValue crypto.SignatureEd25519
}

// Receipt 定义收据
type Receipt struct {
	Name            string  `json:"name"`                      //收据名称：标准名称（trnsfer，...) 非标准名称（...）
	ContractAddress Address `json:"contractAddress,omitempty"` //事件发起方的合约地址
	ReceiptBytes    []byte  `json:"receiptBytes"`
	ReceiptHash     Hash    `json:"receiptHash"`
}

// RPCInvokeCallParam 合约调用参数
// todo 设计转账接收人的账户余额， 注：考虑多个message中分别有不同的接收者
type RPCInvokeCallParam struct {
	Sender          Address         `json:"sender"`          // 交易发送者
	Balances        []byte          `json:"balances"`        // 交易发送者账户余额
	SenderPublicKey PubKey          `json:"senderPublicKey"` // 公钥
	To              Address         `json:"to"`              // 转账交易接收者账户
	ToBalance       []byte          `json:"tobalance"`       // 转账交易接收者账户余额
	Tx              Transaction     `json:"tx,omitempty"`    // 交易注释
	GasLeft         int64           `json:"gasleft"`         // 剩余gaslimit
	Message         Message         `json:"message"`         // 交易消息
	Receipts        []common.KVPair `json:"receipts"`        // 上个交易输出的消息列表
}

// Response 合约调用的 Response
type Response struct {
	Code     uint32          `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Log      string          `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	Data     string          `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Info     string          `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasLimit int64           `protobuf:"varint,5,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed  int64           `protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Fee      int64           `protobuf:"varint,7,opt,name=fee,json=fee,proto3" json:"fee,omitempty"`
	Tags     []common.KVPair `protobuf:"bytes,8,rep,name=tags" json:"tags,omitempty"`
	TxHash   common.HexBytes `protobuf:"bytes,9,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	Height   int64           `protobuf:"varint,10,opt,name=height,proto3" json:"height,omitempty"`
}
