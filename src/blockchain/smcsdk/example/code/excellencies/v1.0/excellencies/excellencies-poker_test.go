package excellencies

//
//import (
//	"blockchain/smcsdk/sdk/bn"
//	"testing"
//)
//
//func TestGamerInfo_JudgeWin(t *testing.T) {
//	var tests = []struct {
//		banker   GamerInfo
//		player   GamerInfo
//		expected bool
//	}{
//		//庄家K Q 10(黑桃)玩家Q J 10(红桃) 公牌 K > Q 庄家赢
//		{GamerInfo{[]PokerCard{{3, 13}, {3, 12}, {3, 10}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 12}, {2, 11}, {2, 10}}, 0, bn.N(1E10), false},
//			false},
//		//庄家K 2 10(黑桃)玩家KKK(大三公) 庄家输
//		{GamerInfo{[]PokerCard{{3, 13}, {3, 2}, {3, 4}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 13}, {1, 13}, {0, 13}}, 0, bn.N(1E10), false},
//			true},
//		//庄家kkk(大三公) 玩家QQQ(大三公) KKK > QQQ 庄家赢
//		{GamerInfo{[]PokerCard{{3, 13}, {2, 13}, {1, 13}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{3, 12}, {2, 12}, {1, 12}}, 0, bn.N(1E10), false},
//			false},
//		//庄家K 3 2 (5点 黑桃) 玩家K 2 3 (5点红桃) 庄家赢
//		{GamerInfo{[]PokerCard{{3, 13}, {3, 3}, {3, 2}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 13}, {2, 2}, {2, 3}}, 0, bn.N(1E10), false},
//			false},
//		//庄家 J K 3 玩家8 K 5 （均3点） 庄家两张公牌 庄家赢
//		{GamerInfo{[]PokerCard{{3, 11}, {3, 13}, {3, 3}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{3, 8}, {3, 12}, {3, 5}}, 0, bn.N(1E10), false},
//			false},
//		//庄家K24 玩家JJJ （玩家大三公）庄家输
//		{GamerInfo{[]PokerCard{{3, 13}, {3, 2}, {3, 4}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 13}, {1, 13}, {0, 13}}, 0, bn.N(1E10), false},
//			true},
//		//庄家 K（红） 2 4 玩家 K (黑) 2 4 （同点 黑 > 红）庄家输
//		{GamerInfo{[]PokerCard{{2, 13}, {3, 2}, {3, 4}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{3, 13}, {1, 2}, {0, 4}}, 0, bn.N(1E10), false},
//			true},
//		//庄家K 2 4 玩家K Q J （玩家大三公）庄家输
//		{GamerInfo{[]PokerCard{{3, 13}, {3, 2}, {3, 4}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 13}, {1, 12}, {0, 11}}, 0, bn.N(1E10), false},
//			true},
//		//庄家K 2 4 玩家K Q 6 （同点 玩家两张公牌）庄家输
//		{GamerInfo{[]PokerCard{{3, 13}, {3, 2}, {3, 4}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 13}, {1, 12}, {0, 6}}, 0, bn.N(1E10), false},
//			true},
//		//庄家K K K 玩家6 6 6 （庄家大三公 > 玩家小三公）庄家赢
//		{GamerInfo{[]PokerCard{{3, 13}, {2, 13}, {1, 13}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 6}, {1, 6}, {0, 6}}, 0, bn.N(1E10), false},
//			false},
//		//庄家K Q K 玩家6 6 6 （庄家混三公 < 玩家小三公）庄家输
//		{GamerInfo{[]PokerCard{{3, 13}, {2, 12}, {1, 13}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 6}, {1, 6}, {0, 6}}, 0, bn.N(1E10), false},
//			true},
//		//庄家K 2 2 玩家K 3 3 （庄家4点 < 玩家6点）庄家输
//		{GamerInfo{[]PokerCard{{3, 13}, {2, 2}, {1, 2}}, 0, bn.N(1E10), false},
//			GamerInfo{[]PokerCard{{2, 13}, {1, 3}, {0, 3}}, 0, bn.N(1E10), false},
//			true},
//	}
//
//	for _, test := range tests {
//		test.banker.CalcPoints()
//		test.player.CalcPoints()
//		test.banker.JudgeWin(&test.player)
//		if test.expected == test.player.IsWin {
//			t.Log("测试成功", test.expected, test.player.IsWin, test.banker.Points, test.player.Points, test.banker.Cards[0].Color, test.player.Cards[0].Color)
//		} else {
//			t.Error("测试失败", test.expected, test.player.IsWin, test.banker.Points, test.player.Points, test.banker.Cards[0].Color, test.player.Cards[0].Color)
//		}
//	}
//}
