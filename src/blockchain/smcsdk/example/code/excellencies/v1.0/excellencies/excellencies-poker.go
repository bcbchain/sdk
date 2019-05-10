package excellencies

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/forx"
	"bytes"
	"encoding/binary"
	"sort"
)

type PokerCard struct {
	Color  int64 `json:"color"`  //扑克牌花色
	Number int64 `json:"number"` //扑克牌数字
}

type PokerCards []PokerCard

type Dealer struct {
	Random []byte  //随机数
	Pokers []int64 //牌对应的序号
}

func (cards PokerCards) Len() int {
	return len(cards)
}

func (cards PokerCards) Less(i, j int) bool {
	return cards[i].Number > cards[j].Number

}

func (cards PokerCards) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
}

type GamerInfo struct {
	Cards     PokerCards `json:"cards"`     //扑克牌
	Points    int64      `json:"points"`    //游戏点数
	BetAmount bn.Number  `json:"betAmount"` //本轮在该玩家上的下注总额
	IsWin     bool       `json:"isWin"`     //输赢判断
}

func BytesToInt(b []byte) int {
	var bytesBuffer *bytes.Buffer
	if len(b) < 8 {
		bytesBuffer = bytes.NewBuffer(make([]byte, 8-len(b)))
		bytesBuffer.Write(b)
	} else {
		bytesBuffer = bytes.NewBuffer(b)
	}

	var tmp int64
	err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	if err != nil {
		panic(err)
	}
	return int(tmp)
}

func (card *Dealer) GetCard() (poker PokerCard) {
	if len(card.Pokers) == 0 {
		panic("The pokers is empty")
	}

	if len(card.Random) == 0 {
		panic("The Random is empty")
	}

	cardSite := BytesToInt(card.Random[:2]) % len(card.Pokers)
	poker.Number = card.Pokers[cardSite]%13 + 1
	poker.Color = card.Pokers[cardSite] / 13 % 4

	//删除使用过的随机数和牌
	card.Random = card.Random[2:]
	card.Pokers = append(card.Pokers[:cardSite], card.Pokers[cardSite+1:]...)
	return
}

func (card *Dealer) Init(random []byte) {
	if len(random) < 12 {
		panic("The random length is too small")
	}
	card.Random = random

	//1副牌 1*52:52
	forx.Range(0, 52, func(i int) bool {
		card.Pokers = append(card.Pokers, int64(i))
		return forx.Continue
	})

}

func (gamer *GamerInfo) GetExcellencyCount() int64 {
	var count int64
	forx.Range(0, len(gamer.Cards), func(i int) bool {
		if gamer.Cards[i].Number > 10 {
			count++
		}
		return forx.Continue
	})

	return count
}

func (gamer *GamerInfo) GetMaxCardNumber() int64 {
	var number int64
	forx.Range(0, len(gamer.Cards), func(i int) bool {
		if gamer.Cards[i].Number > number {
			number = gamer.Cards[i].Number
		}
		return forx.Continue
	})

	return number
}

func (gamer *GamerInfo) GetMaxCardColor() int64 {
	var color int64
	forx.Range(0, len(gamer.Cards), func(i int) bool {
		if gamer.Cards[i].Color > color {
			color = gamer.Cards[i].Color
		}
		return forx.Continue
	})

	return color
}

func (gamer *GamerInfo) GetCardsString() string {
	var CardsString = map[int64]string{1: "A", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "10", 11: "J", 12: "Q", 13: "K"}
	var cardsStr string
	sort.Sort(gamer.Cards)
	forx.Range(0, len(gamer.Cards), func(i int) bool {
		temp, _ := CardsString[gamer.Cards[i].Number]
		cardsStr = cardsStr + temp
		return forx.Continue
	})

	return cardsStr
}

func (gamer *GamerInfo) JudgeWin(player *GamerInfo) {
	if len(gamer.Cards) != 3 || len(player.Cards) != 3 {
		panic("gamer is error")
	}

	if gamer.Points > player.Points {
		player.IsWin = false
	} else if gamer.Points < player.Points {
		player.IsWin = true
	} else if gamer.GetExcellencyCount() > player.GetExcellencyCount() {
		player.IsWin = false
	} else if gamer.GetExcellencyCount() < player.GetExcellencyCount() {
		player.IsWin = true
	} else if gamer.GetMaxCardNumber() > player.GetMaxCardNumber() {
		player.IsWin = false
	} else if gamer.GetMaxCardNumber() < player.GetMaxCardNumber() {
		player.IsWin = true
	} else if gamer.GetMaxCardColor() > player.GetMaxCardColor() {
		player.IsWin = false
	} else {
		player.IsWin = true
	}
}

func (gamer *GamerInfo) AddCard(poker PokerCard) {
	gamer.Cards = append(gamer.Cards, poker)

	if len(gamer.Cards) == 3 {
		gamer.CalcPoints()
	}
}

func (gamer *GamerInfo) CalcPoints() {
	var CardsToPoints = map[string]int64{"QJJ": 10, "QQJ": 11, "KJJ": 12, "KQJ": 13, "KQQ": 14, "KKJ": 15, "KKQ": 16, "AAA": 17, "222": 18, "333": 19,
		"444": 20, "555": 21, "666": 22, "777": 23, "888": 24, "999": 25, "101010": 26, "JJJ": 27, "QQQ": 28, "KKK": 29}
	var point int64 = 0
	if len(gamer.Cards) != 3 {
		panic("cards count is error")
	}

	point, ok := CardsToPoints[gamer.GetCardsString()]
	if ok {
		gamer.Points = point
		return
	}

	forx.Range(0, len(gamer.Cards), func(i int) bool {
		if gamer.Cards[i].Number <= 9 {
			point = point + gamer.Cards[i].Number
		}
		return forx.Continue
	})

	gamer.Points = point % 10
}
