package myballot

import (
	"blockchain/smcsdk/sdk/types"
)

//Voter 这里声明了一个新的复合类型用于稍后的变量
//     它用来表示一个选民
type Voter struct {
	weight   uint          // 计票的权重
	voted    bool          // 若为真，代表该人已投票
	delegate types.Address // 被委托人
	vote     uint          // 投票提案的索引
}

//Proposal 提案的类型
type Proposal struct {
	name      string // 简称（最长32个字节）
	voteCount uint   // 得票数
}
