package directsale

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

const (
	TOKENNAME     = "Diamond Coin"
	PROPORTION    = 1000 //proportion
	EXCHANGERATIO = 1000000000
	MAXLEVEL      = 11 //见点奖最大层数
)

type Settings struct {
	StarSettings []starReward `json:"starSetting"`
	Foster       int64        `json:"foster"`       // 抚育奖奖励比例（千分比）
	ShareReward  int64        `json:"shareReward"`  // 直推奖比例（千分比）
	PointReward  bn.Number    `json:"pointReward"`  // 见点奖
	SalerExpense bn.Number    `json:"salerExpense"` // 入会费用
}

type starReward struct {
	MaxIncome     bn.Number `json:"maxIncome"`   // 该星级对应的流水最大值
	DividendRatio int64     `json:"rewardRatio"` // 奖励比例 （千分比）
}

type Saler struct {
	ContractNames []string                `json:"contractNames"` // 运营的合约名称
	RefAddr       types.Address           `json:"refAddr"`       // 直推人地址
	RefCounts     int64                   `json:"refCounts"`     // 推荐人数
	Sons          []types.Address         `json:"sons"`          // 下级地址
	Parent        types.Address           `json:"parent"`        // 上级地址
	InvestRate    int64                   `json:"investRate"`    // 出资比例
	Accounts      map[string]SalerAccount `json:"accounts"`      // 会员账户信息
}

type SalerAccount struct {
	InvestAmount  bn.Number `json:"investAmount"` //总出资额(saler + owner)
	IncomeBalance bn.Number `json:"incomeAmount"` //本金余额(统计income接口中的转入的资金)
	LockedAmount  bn.Number `json:"lockedAmount"` //总锁定
	FlowAmount    bn.Number `json:"incomeAmount"` //累计流水(统计pay接口中的incomeAmount)
	TotalBalance  bn.Number `json:"totalBalance"` //总余额
}

type SalerApp struct {
	AppInfoMap map[string]AppStatByToken `json:"appInfoList"` // 合约对应锁定金额
}

type AppStatByToken struct {
	LockedAmount bn.Number `json:"lockedAmount"` // 锁定的钱
	FlowAmount   bn.Number `json:"flowAmount"`   // 累计流水金额(结算清零)
	PayCount     int64     `json:"payCount"`     // 支出次数统计
}

type Global struct {
	TotalLocked bn.Number `json:"totalLocked"` // 资金池的锁定金额
	PoolBalance bn.Number `json:"poolBalance"` // 资金池可用余额
}

type Settle struct {
	Address       types.Address        `json:"address"`       // 会员地址
	Ratio         map[string]int64     `json:"ratio"`         // 钱币兑换汇率
	StartLevel    int                  `json:"startLevel"`    // 星级
	FosterForward map[string]bn.Number `json:"fosterForward"` // 抚育奖总和
	SalerForward  map[string]bn.Number `json:"salerForward"`  // 会员奖总和
	DCAmount      bn.Number            `json:"dcAmount"`      // 换算成dc的总流水
}
