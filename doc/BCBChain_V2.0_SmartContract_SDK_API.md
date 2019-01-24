# 智能合约开发SDK接口说明

**V2.0.1**

<div STYLE="page-break-after: always;"></div>

# 修订历史

| 版本&日期         | 修订内容&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; |
| ----------------- | ------------------------------------------------------------ |
| V2.0.1：2018-12-7 | 初稿。                                                       |

<div STYLE="page-break-after: always;"></div>
[TOC]

<script src="./github/ltview.js"></script>
<div STYLE="page-break-after: always;"></div>
# 1 文档概述

BCBChain SDK是专门为程序员们开发BCBChain上运行的智能合约设计的编程接口。本文档详细介绍BCBChain SDK提供的各接口及类型说明。



# 2 sdk

## 2.1 function

### 2.1 Require()

```
func Require(expr bool, errCode uint32, errInfo string)
```

用于表达式断定，要求满足 expr 为 true，如果 expr 为 false，将触发智能合约 panic 一个类型为 types.Error 的对象，其中错误码为 errCode，错误信息为 errInfo。如果 expr 为 true，智能合约将被允许继续往下执行。



### 2.2 RequireNotError()

```
func RequireNotError(err error, errCode uint32)
```

用于错误断定，要求满足 err 必须为空，如果 err 不为空，将触发智能合约 panic 一个类型为 types.Error 的对象，其中错误码为 errCode，错误信息为 err 对象中的描述信息。如果 err 为空，智能合约将被允许继续往下执行。



### 2.2 RequireOwner()

```
func RequireOwner(sdk ISmartContract)
```

用于进行权限断定，要求满足智能合约调用者的必须是合约的拥有者。如果不满足要求，将触发智能合约 panic 一个类型为 types.Error 的对象，其中错误码为 ```ErrNoAuthorization```，错误信息为具体错误原因。如果满足要求，智能合约将被允许继续往下执行。



### 2.2 RequireAddress()

```
func RequireAddress(sdk ISmartContract, addr types.Address)
```

用于地址格式断定，要求满足 传入的地址格式正确，如果地址格式不正确，将触发智能合约 panic 一个类型为 types.Error 的对象，其中错误码为 ```ErrInvalidAddress```，错误信息为具体错误原因。如果地址格式正确，智能合约将被允许继续往下执行。



## 2.2 interface

### 2.2.1 interface IBlock

参见本文档章节```sdk/IBlock```。



### 2.2.2 interface IContract

参见本文档章节```sdk/IContract```。



### 2.2.3 interface ITx

参见本文档章节```sdk/ITx```。



### 2.2.4 interface IMessage

参见本文档章节```sdk/IMessage```。



### 2.2.5 interface IAccount

参见本文档章节```sdk/IAccount```。



### 2.2.6 interface IToken

参见本文档章节```sdk/IToken```。



### 2.2.7 interface IHelper

参见本文档章节```sdk/IHelper```。



### 2.2.8 interface IAccountHelper

参见本文档章节```sdk/IAccountHelper```。



### 2.2.9 interface IBlockChainHelper

参见本文档章节```sdk/IBlockChainHelper```。



### 2.2.10 interface IBuildHelper

参见本文档章节```sdk/IBuildHelper```。



### 2.2.11 interface IContractHelper

参见本文档章节```sdk/IContractHelper```。



### 2.2.12 interface IReceiptHelper

参见本文档章节```sdk/IReceiptHelper```。



### 2.2.13 interface IGenesisHelper

参见本文档章节```sdk/IGenesisHelper```。



### 2.2.14 interface IStateHelper

参见本文档章节```sdk/IStateHelper```。



### 2.2.15 interface ITokenHelper

参见本文档章节```sdk/ITokenHelper```。



<div STYLE="page-break-after: always;"></div>

# 3 sdk/bn

程序包 sdk/bn 封装了一个处理大数的类 Number，进行加减乘除操作时不必考虑溢出的问题。

## 3.1 function

本章描述程序包 sdk/bn 提供的关于类 Number 简便构造函数。

### 3.1.1 N()

```
func N(x int64) Number
```

将 int64 类型的 x 转换成 Number 类型对象并返回。



### 3.1.2 N1()

```
func N1(b int64, d int64) Number
```

根据传入的 b 和 d，生成一个结果为 b * d 的 Number 类型对象并返回。



### 3.1.3 N2()

```
func N2(b int64, d1, d2 int64) Number
```

根据传入的 b，d1，d2，生成一个结果为 b * d1 * d2 的 Number 类型对象并返回。



### 3.1.4 NB()

```
func NB(x *big.Int) Number
```

将 big.Int 类型的 x 转换成 Number 类型对象并返回。



### 3.1.7 NBS()

```
func NBS(x []byte) Number
```

将按大端表示无符号整数的字节切片 x 转换成 Number 类型对象并返回。



### 3.1.8 NBytes()

```
func NBytes(x []byte) Number
```

将按大端表示无符号整数的字节切片 x 转换成 Number 类型对象并返回。



### 3.1.5 NString()

```
func NS(s string) Number
```

将按十进制字符串表示的大整数 s 转换成 Number 类型对象并返回，如果解析失败将返回0。



### 3.1.6 NStringHex()

```
func NS(s string) Number
```

将以 0x 或 0X 开头的十六进制字符串表示的无符号大整数 s 转换成 Number 类型对象并返回，如果解析失败将返回0。



### 3.1.9 NewNumber()

```
func NewNumber(x int64) Number
```

将 int64 类型的 x 转换成 Number 类型对象并返回。



### 3.1.11 NewNumberStringBase()

```
func NewNumberStringBase(s string, base int) Number
```

将字符串表示的大整数 s 转换成 Number 类型对象并返回，字符串按给定的基数 base 进行解析，如果解析失败将返回0。

基数 base 必须是 0 或者 2到MaxBase之间的整数。如果基数为0，字符串的前缀决定实际的转换基数：前缀为 "0x"、"0X" 表示十六进制；前缀 "0b"、"0B" 表示二进制；前缀 "0" 表示八进制；其它都自动采用十进制作为基数。

针对 <= 36 的基数，大写和小写字母表达相同的数，字母 'a' 到 'z' 和 'A' 到 'Z' 都表达数值 10 到 35。

针对 > 36 的基数，大写字母 'A' 到 'Z' 表达数值 36 到 61。



### 3.1.12 NewNumberBigInt()

```
func NewNumberBigInt(x *big.Int) Number
```

将 big.Int 类型的 x 转换成 Number 类型对象并返回。



### 3.1.13 NewNumberLong()

```
func NewNumberLong(b int64, d int64) Number
```

根据传入的 b 和 d，生成一个结果为 b * d 的 Number 类型对象并返回。



### 3.1.14 NewNumberLongLong()

```
func NewNumberLongLong(b int64, d1, d2 int64) Number
```

根据传入的 b，d1，d2，生成一个结果为 b * d1 * d2 的 Number 类型对象并返回。



<div STYLE="page-break-after: always;"></div>

## 3.2 class Number

本章描述类 Number 的成员函数。

### 3.13 String()

```
func (x Number) String() string
```

将 x 转换为十进制字符串 string。x 未设定初始大数值，返回nil。



### 3.14 Value()

```
func (x Number) Value() *big.Int
```

获取 x 的 big.Int 的值。x 未设定初始大数值，返回nil。



### 3.15 CmpI()

```
func (x Number) CmpI(y int64) int
```

将 x 与 y 进行比较，返回 -1 代表 x < y，0 代表 x == y，+1 代表 x > y。x 或 y 未设定初始大数值将触发panic。



### 3.16 Cmp()

```
func (x Number) Cmp(y Number) int
```

将 x 与 y 进行比较，返回 -1 代表 x < y，0 代表 x == y，+1 代表 x > y。x 或 y 未设定初始大数值将触发panic。



### 3.17 IsZero()

```
func (x Number) IsZero() bool
```

判断 x 是否为 0。x 等于 0 返回 true。x 未设定初始大数值将触发panic。



### 3.18 IsPositive()

```
func (x Number) IsPositive() bool
```

判断 x 是否为正数。x 为正数返回 true。x 未设定初始大数值将触发panic。



### 3.19 IsNegative()

```
func (x Number) IsNegative() bool
```

判断 x 是否为负数。x 为负数返回 true。x 未设定初始大数值将触发panic。



### 3.20 IsEqualI()

```
func (x Number) IsEqualI(y int64) bool
```

将 x 与 y 进行比较，true 代表 x == y，false 代表 x != y。x 或 y 未设定初始大数值将触发panic。



### 3.21 IsEqual()

```
func (x Number) IsEqual(y Number) bool
```

将 x 与 y 进行比较，true 代表 x == y，false 代表 x != y。x 或 y 未设定初始大数值将触发panic。



### 3.22 IsGreaterThanI()

```
func (x Number) IsGreaterThanI(y int64) bool
```

将 x 与 y 进行比较，true 代表 x > y，false 代表 x <= y。x 或 y 未设定初始大数值将触发panic。



### 3.23 IsGreaterThan()

```
func (x Number) IsGreaterThan(y Number) bool
```

将 x 与 y 进行比较，true 代表 x > y，false 代表 x <= y。x 或 y 未设定初始大数值将触发panic。



### 3.24 IsLessThanI()

```
func (x Number) IsLessThanI(y int64) bool
```

将 x 与 y 进行比较，true 代表 x < y，false 代表 x >= y。x 或 y 未设定初始大数值将触发panic。



### 3.25 IsLessThan()

```
func (x Number) IsLessThan(y Number) bool
```

将 x 与 y 进行比较，true 代表 x < y，false 代表 x >= y。x 或 y 未设定初始大数值将触发panic。



### 3.26 IsGEI()

```
func (x Number) IsGEI()(y int64) bool
```

将 x 与 y 进行比较，true 代表 x >= y，false 代表 x < y。x 或 y 未设定初始大数值将触发panic。



### 3.27 IsGE()

```
func (x Number) IsGE(y Number) bool
```

将 x 与 y 进行比较，true 代表 x >= y，false 代表 x < y。x 或 y 未设定初始大数值将触发panic。



### 3.28 IsLEI()

```
func (x Number) IsLEI(y int64) bool
```

将 x 与 y 进行比较，true 代表 x <= y，false 代表 x > y。x 或 y 未设定初始大数值将触发panic。



### 3.29 IsLE()

```
func (x Number) IsLE(y Number) bool
```

将 x 与 y 进行比较，true 代表 x <= y，false 代表 x > y。x 或 y 未设定初始大数值将触发panic。



### 3.30 AddI()

```
func (x Number) AddI(y int64) Number
```

计算 x + y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.31 Add()

```
func (x Number) Add(y Number) Number
```

计算 x + y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.32 SubI()

```
func (x Number) SubI(y int64) Number
```

计算 x - y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.33 Sub()

```
func (x Number) Sub(y Number) Number
```

计算 x - y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.34 MulI()

```
func (x Number) MulI(y int64) Number
```

计算 x * y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.35 Mul()

```
func (x Number) Mul(y Number) Number
```

计算 x * y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.36 DivI()

```
func (x Number) DivI(y int64) Number
```

计算 x / y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.37 Div()

```
func (x Number) Div(y Number) Number
```

计算 x / y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.38 ModI()

```
func (x Number) ModI(y int64) Number
```

计算 x % y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.39 Mod()

```
func (x Number) Mod(y Number) Number
```

计算 x % y ，并返回计算结果。x 或 y 未设定初始大数值将触发panic。



### 3.40 Sq()

```
func (x Number) Sq() Number
```

计算 x ** 2 ，并返回计算结果。x 未设定初始大数值将触发panic。



### 3.41 Sqrt()

```
func (x Number) Sqrt() Number
```

计算 x 平方根，并返回计算结果。x 未设定初始大数值将触发panic。



### 3.44 Exp()

```
func (x Number) Exp(y Number) Number
```

计算 x ** y， 并且返回结果。x 或 y 未设定初始大数值将触发panic。



### 3.42 Lsh()

```
func (x Number) Lsh(n uint) Number
```

将 x 向左移 n 位，并返回移位结果。x 未设定初始大数值将触发panic。



### 3.42 Rsh()

```
func (x Number) Rsh(n uint) Number
```

将 x 向右移 n 位，并返回移位结果。x 未设定初始大数值将触发panic。



### 3.43 And()

```
func (x Number) And(y Number) Number
```

按位计算 x & y ，并返回计算结果（需要注意这里x、y有可能为负数，需要采用二进制补码形式进行按位与）。x 或 y 未设定初始大数值将触发panic。



### 3.43 Or()

```
func (x Number) Or(y Number) Number
```

按位计算 x | y ，并返回计算结果（需要注意这里x、y有可能为负数，需要采用二进制补码形式进行按位或）。x 或 y 未设定初始大数值将触发panic。



### 3.43 Xor()

```
func (x Number) Xor(y Number) Number
```

按位计算 x ^ y ，并返回计算结果（需要注意这里x、y有可能为负数，需要采用二进制补码形式进行按位异或）。x 或 y 未设定初始大数值将触发panic。



### 3.43 Not()

```
func (x Number) Not() Number
```

按位计算 ^ x ，并返回计算结果（需要注意这里x有可能为负数，需要采用二进制补码形式进行取反）。x 未设定初始大数值将触发panic。



### 3.46 Bytes()

```
func (x Number) Bytes() []byte

```

将 x 转换为大端顺序的字节切片（第一字节解为表示符号，负数为0xFF，非负数为0x00），并返回转换结果。例如：```380```将被转换为```0x00017C```；```-380```将被转换为```0xFF017C```。x 未设定初始大数值将触发panic。



### 3.45 SetBytes()

```
func (x *Number) SetBytes(buf []byte) Number
```

将 x 的值设置为一个大端顺序的字节切片（第一字节为0xFF表示是一个负数，其绝对值从第二字节开始编码，否则整个字节切片表示一个非负整数），并将 x 的值返回。例如：```0x00017C```和```0x017C```将被转换为```380```；```0xFF017C```将被转换为```-380```。



### 3.46 MarshalJSON()

```
func (x Number) MarshalJSON() (data []byte, err error)
```

实现标准的 JSON 序列化接口。将 x 转成简化版的 JSON 字符串，例如字符串```380```。x 未设定初始大数值将触发panic。



### 3.46 UnmarshalJSON()

```
func (x *Number) UnmarshalJSON(data []byte) error
```

实现标准的 JSON 反序列化接口。将 x 的值设为输入的 JSON 字符串对应的大数。支持简化版的 JSON 字符串（例如：字符串```380```）与结构版的 JSON 字符串（例如：字符串```{"v":380}```）。



<div STYLE="page-break-after: always;"></div>

# 4 sdk/common

```
type HexBytes []byte
```

此类型主要用来使用十六进制进行 json 的编码



## 4.1 array

### 4.1.1 Arr

```
func Arr(items ...interface{}) []interface{}
```

返回 items。



## 4.2 bytes

###  4.2.1 Marshal

```
func (bz HexBytes) Marshal() ([]byte, error)
```

返回 bz 的 []byte。



### 4.2.2 Unmarshal

```
func (bz *HexBytes) Unmarshal(data []byte) error
```

将 data 转换成 HexBytes 类型。



### 4.2.3 MarshalJSON

```
func (bz HexBytes) MarshalJSON() ([]byte, error)
```

将 bz 转换成 JSON 格式后的 []byte，并返回结果。



### 4.2.4 UnmarshalJSON

```
func (bz *HexBytes) UnmarshalJSON(data []byte) error
```

将 JSON 格式的 data 转换成 []byte，data 应以英文双引号开始和结束。



### 4.2.5 Bytes

```
func (bz HexBytes) Bytes() []byte
```

直接返回 bz。



### 4.2.6 String

```
func (bz HexBytes) String() string
```

将 bz 转换成全部大写的字符串。



### 4.2.7 Format

```
func (bz HexBytes) Format(s fmt.State, verb rune)
```

根据指定的格式，生成对应的 fmt.State。



## 4.3 byte slice

### 4.3.1 Fingerprint

```
func Fingerprint(slice []byte) []byte
```

返回 slice 的前六个字节，如果 slice 长度小于六，返回结果中其余为0。



### 4.3.2 IsZeros

```
func IsZeros(slice []byte) bool
```

判断 slice 是否全部为字节 0，是返回true。



### 4.3.3 RightPadBytes

```
func RightPadBytes(slice []byte, l int) []byte
```

对 slice 向右填充字节0，最终长度为 l，如果 l < len(slice)，返回 slice。



### 4.3.4 LeftPadBytes

```
func LeftPadBytes(slice []byte, l int) []byte
```

对 slice 向左填充字节0，最终长度为 l，如果 l < len(slice)，返回 slice。



### 4.3.5 TrimmedString

```
func TrimmedString(b []byte) string
```

返回将b前端所有字节 0 都去掉的子切片



### 4.3.6 PrefixEndBytes

```
func PrefixEndBytes(prefix []byte) []byte
```

将 prefix 最后一个字节加一并返回结果。



## 4.4 KVPair

收据结构，在 sdk/types 中收据类型引用。



### 4.4.1 KVPair

```
type KVPair struct {
   Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
   Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}
```

收据的基本结构。



## 4.5 math

sdk 中的 math 接口，主要包括比较大小等方法。



### 4.5.1 MaxInt8

```
func MaxInt8(a, b int8) int8
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.2 MaxUint8

```
func MaxUint8(a, b uint8) uint8
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.3 MaxInt16

```
func MaxInt16(a, b int16) int16
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.4 MaxUint16

```
func MaxUint16(a, b uint16) uint16
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.5 MaxInt32

```
func MaxInt32(a, b int32) int32
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.6 MaxUint32

```
func MaxUint32(a, b uint32) uint32
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.7 MaxInt64

```
func MaxInt64(a, b int64) int64
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.8 MaxUint64

```
func MaxUint64(a, b uint64) uint64
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.9 MaxInt

```
func MaxInt(a, b int) int
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.10 MaxUint

```
func MaxUint(a, b uint) uint
```

返回 a 和 b 中较大的值，如果相等返回 b。



### 4.5.11 MinInt8

```
func MinInt8(a, b int8) int8
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.12 MinUint8

```
func MinUint8(a, b uint8) uint8
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.13 MinInt16

```
func MinInt16(a, b int16) int16
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.14 MinUint16

```
func MinUint16(a, b uint16) uint16
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.15 MinInt32

```
func MinInt32(a, b int32) int32
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.16 MinUint32

```
func MinUint32(a, b uint32) uint32
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.17 MinInt64

```
func MinInt64(a, b int64) int64
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.18 MinUint64

```
func MinUint64(a, b uint64) uint64
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.19 MinInt

```
func MinInt(a, b int) int
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.20 MinUint

```
func MinUint(a, b uint) uint
```

返回 a 和 b 中较小的值，如果相等返回 b。



### 4.5.21 ExpUint64

```
func ExpUint64(a, b uint64) uint64
```

计算 a**b 并且返回结果。



## 4.6 nil

判断是否为 nil 或者为空。



### 4.6.1 IsTypedNil

```
func IsTypedNil(o interface{}) bool
```

判断输入的 o 是否为 nil，是 nil 返回 true。如果输入的是 chan、func、map、ptr 或者 slice，判断其是否为 nil，如果是，返回 true。如果是其他类型，返回 false。



### 4.6.2 IsEmpty

```
func IsEmpty(o interface{}) bool
```

判断输入的 o 是否为空，是空返回 true。如果输入的是 array、chan、map、slice 或者 string，为空值返回 true。如果输入的是其他类型，返回 false。



# 5 sdk/jsoniter

用来对 json 格式的数据进行 marshal 与 unmarshal。



## 5.1 Marshal

```
func Marshal(_v interface{}) ([]byte, error)
```

返回 _v 的 json 编码。



## 5.2 Unmarshal

```
func Unmarshal(_bz []byte, _v interface{}) (interface{}, error)
```

解析 json 编码的数据 _bz 并且将结果存入 _v 指向的值。



# 6 sdk/crypto

sdk 的加密接口，包括 ed25519 和 sha3算法。用来验证签名等。



## 6.1 ed25519

### 6.1.1 VerifySign

```
func VerifySign(pubkey, data, sign []byte) bool
```

验证签名，成功返回 true。pubkey 为公钥，data 为签名的数据， sign 为签名。



## 6.2 sha3

### 6.2.1 Sum224

```
func Sum224(datas ...[]byte) []byte
```

使用 sha3-224 算法计算 datas 的散列值并返回。



### 6.2.2 Sum256

```
func Sum256(datas ...[]byte) []byte
```

使用 sha3-256 算法计算 datas 的散列值并返回。



### 6.2.3 Sum384

```
func Sum384(datas ...[]byte) []byte
```

使用 sha-384 算法计算 datas 的散列值并返回。



### 6.2.4 Sum512

```
func Sum512(datas ...[]byte) []byte
```

使用 sha-512 算法计算 datas 的散列值并返回。



# 7 sdk/types

types 中包括 error 错误码和一些其他类型。



## 7.1 Error

Error 包括 ErrorCode（uint32） 和 ErrorDesc（string） 两个字段，分别是错误码以及错误信息描述。



### 7.1.1 Error

```
func (smcerror *Error) Error() string
```

返回错误信息。



### 7.1.2 错误码以及错误信息

```
CodeOK = 200                      ErrorDesc = ""
ErrStubDefined = 51001            ErrorDesc = ""
ErrAddSupplyNotEnabled = 52002    ErrorDesc = "Add supply is not enabled"
ErrBurnNotEnabled = 52003         ErrorDesc = "Burn supply is not enabled"
ErrInvalidAddress = 52004         ErrorDesc = "Invalid address"
ErrNoAuthorization = 53005        ErrorDesc = "No authorization to execute contract"
ErrInvalidParameter = 53006       ErrorDesc = "Invalid parameter"
ErrInsufficientBalance = 53007    ErrorDesc = "Insufficient balance"
ErrUserDefined = 55008            ErrorDesc = "Contract logic error"
```



## 7.2 types

sdk 常用类型。



### 7.2.1 Address

```
type Address = string
```

Account 或者 Contract 的地址。



### 7.2.2 HexBytes

```
type HexBytes = common.HexBytes
```

HexBytes 用来作为原始字节。



### 7.2.3 Hash

```
type Hash = common.HexBytes
```

SHA3 等算法的 hash。



### 7.2.4 Pubkey

```
type PubKey = common.HexBytes
```

用作公钥类型。



### 7.2.5 KVPair

```
type KVPair = common.KVPair
```

用作收据类型。



# 8 sdk/IBlock

智能合约调用时，IBlock在两种不同场景下对应的区块信息是有区别的。

> 在checkTx阶段，当前智能合约的调用还没有取得共识，此时IBlock为区块链上最后一个区块的信息（可能是一个空区块，总之这个区块与当前所执行的智能合约调用没有任何关系）。

> 在deliverTx阶段，此次调用智能合约的交易已经在区块链上达成共识，并写入到最新一个区块中，此时IBlock为区块链上最新一个实时区块的信息（不可能是一个空区块，这个区块中必然包含一笔本次智能合约调用的交易信息）。




## 8.1 ChainID

```
func (*IBlock) ChainID() string
```
获取当前区块链的链标识ID。



## 8.2 BlockHash

```
func (*IBlock) BlockHash() Hash
```
获取当前区块的区块哈希。



## 8.3 Height

```
func (*IBlock) Height() int64
```
获取当前区块的高度。



## 8.4 Time/Now

```
func (*IBlock) Time() int64
func (*IBlock) Now() Number
```
获取当前区块生成时的时间戳信息，结果为相对于1970-1-1 00:00:00过去的秒数。其中 Now() 是将当前区块生成时的时间戳信息以别名（now）的形式返回。



## 8.5 NumTxs

```
func (*IBlock) NumTxs() int32
```
获取当前区块纳入的交易笔数。



## 8.6 DataHash

```
func (*IBlock) DataHash() Hash
```
获取当前区块的数据哈希。



## 8.7 ProposerAddress

```
func (*IBlock) ProposerAddress() Address
```
获取当前区块的提案者的地址。



## 8.8 RewardAddress

```
func (*IBlock) RewardAddress() Address
```
获取接收当前区块的出块奖励的地址。



## 8.9 RandomNumber

```
func (*IBlock) RandomNumber() HexBytes
```
获取当前区块的区块随机数。



## 8.10 Version

```
func (*IBlock) Version() string
```
获取当前区块提案者的软件版本。



## 8.11 LastBlockHash

```
func (*IBlock) LastBlockHash() Hash
```
获取上一区块的区块哈希。



## 8.12 LastCommitHash

```
func (*IBlock) LastCommitHash() Hash
```
获取上一区块的确认哈希。



## 8.13 LastAppHash

```
func (*IBlock) LastAppHash() Hash
```
获取上一区块的应用层哈希。



## 8.14 LastFee

```
func (*IBlock) LastFee() int64
```
获取上一区块的手续费（单位：cong）。



# 9 sdk/IContract

智能合约调用时，IContract指向当前被调用的智能合约的详细信息。



## 9.1 Address

```
func (*IContract) Address() Address
```
获取当前智能合约的合约地址，合约地址的计算与合约名称、版本及拥有者地址相关，合约升级会导致智能合约的合约地址发生变化。



## 9.2 Account
```
func (*IContract) Account() Address
```
获取当前智能合约的账户地址，合约的账户地址可用于接收转给合约的资金，账户地址的计算只与合约名称相关，合约升级后合约地址发生变化但不会影响合约的账户地址，合约的账户地址没有私钥与之对应，合约账户上的资金只能通过合约进行操纵。



## 9.3 Owner

```
func (*IContract) Owner() Address
```
获取当前智能合约的拥有者的外部账户地址（外部账户地址有私钥与之对应）。




## 9.4 Name

```
func (*IContract) Name() string
```
获取当前智能合约的名称。




## 9.5 Version

```
func (*IContract) Version() string
```
获取当前智能合约的版本号。



## 9.6 CodeHash

```
func (*IContract) CodeHash() Hash
```
获取当前智能合约的代码所对应的哈希。



## 9.7 EffectHeigt

```
func (*IContract) EffectHeight() int64
```
获取当前智能合约生效时的区块高度，当最新区块高度达不到这个高度时对这个合约进行的所有调用都会失败（在checkTx阶段就会失败）。




## 9.8 LoseHeight
```
func (*IContract) LoseHeight() int64
```
获取当前智能合约失效时的区块高度，0表示没有失效。




## 9.9 KeyPrefix

```
func (*IContract) KeyPrefix() string
```
获取当前智能合约代码所能访问的状态数据库KEY值的前缀（格式为```/xxx```），用于进行合约之间数据的隔离保护。



## 9.10 Methods

```
func (*IContract) Methods() []Method
```
获取当前智能合约的所有方法详细信息的列表。Method结构定义如下：

  ```
  type Method struct {
  	MethodId 	string		// 方法ID，方法ID的计算与方法原型相关
    Gas 		int64		// 方法在调用时消耗的燃料
    ProtoType 	string		// 方法原型
  }
  ```




## 9.11 Token

```
func (*IContract) Token() Address
```
获取当前智能合约注册的代币地址，如果合约没有注册过代币，返回空地址。



## 9.12 SetOwner

```
func (*IContract) SetOwner( owner Address) Error
```
设置智能合约新的拥有者。owner 为智能合约新的拥有者的外部账户地址。返回 Error 的 ErrorCode 等于 200 表示成功。



## 9.13 Initialized

```
func (c *Contract) Initialized() bool
```

返回合约是否已初始化。



## 9.14 Interfaces

```
func (c *Contract) Interfaces() []std.Method
```

获取当前智能合约的所有接口的详细信息。



## 9.15 SetToken

```
func (c *Contract) SetToken(tokenAddr Address)
```

为合约设置 token，tokenAddr 为 token 的地址。



## 9.16 OrgID

```
func (c *Contract) OrgID() string
```

获取当前合约所属的 org ID。



## 9.17 SetInitiated

```
func (c *Contract) SetInitiated()
```

设置当前合约状态为已初始化。



# 10 sdk/ITx

智能合约调用时，ITx指向当前调用智能合约的外部交易（区块链层面所指的交易，即由用户发起并签名的一次通讯，内部可以包含对多个智能合约的调用消息）的详细信息。



## 10.1 Note

```
func (*ITx) Note() string
```
获取当前交易当中传入的备注信息。



## 10.2 GasLimit

```
func (*ITx) GasLimit() int64
```

获取当前交易当中传入的最大燃料限制数量。



## 10.3 GasLeft

```
func (*ITx) GasLeft() int64
```

获取当前交易当中传入的燃料限制数量在被扣除需要消耗的燃料之后的剩余燃料数量。



## 10.4 Sender

```
func (*ITx) Sender() IAccount
```

获取当前交易发起方的账户信息。



# 11 sdk/IMessage

智能合约调用时，IMessage是针对某个智能合约方法的一次调用，在同一笔交易当中可以集成多个级联的智能合约方法调用消息。




## 11.1 To

```
func (*IMessage) To() Address
```

获取当前消息调用的智能合约地址。



## 11.2 MethodId

```
func (*IMessage) MethodId() uint32
```

获取当前消息调用的智能合约方法ID。



## 11.3 Data

```
func (*IMessage) Data() HexBytes
```

获取当前消息调用的参数数据字段的原始信息（参数表的RLP编码格式）。当发生跨智能合约调用时，在被调用合约内部，Data()获得的数据为空。



## 11.4 GasPrice

```
func (*IMessage) GasPrice() int64
```

获取当前消息调用的燃料价格（单位：cong）。



## 11.5 Sender

```
func (*IMessage) Sender() IAccount
```

获取当前消息调用发起方的账户信息。当发生跨智能合约调用时，在被调用合约内部，Sender()获得的为发起合约的账户地址对应的账户信息。



## 11.6 Origin

```
func (*IMessage) Origin() []Address
```

获取消息完整的调用链。在不是进行跨合约调用时，Origin为空。当发生跨智能合约调用时，Origin()用来表达调用的合约链，在被调用合约内部，Origin()获得的地址列表中最后一个为本次调用发起合约的合约地址。



## 11.7 InputReceipts

```
func (*IMessage) InputReceipts() []KVPair
```

获取级联消息中前一个消息输出的收据作为本次消息调用的输入，本函数获取这些输出的收据列表。



## 11.8 OutputReceipts

```
func (m *Message) OutputReceipts() []KVPair
```

获取级联消息中的输出收据。



## 11.9 FillOutputReceipts

```
func (m *Message) FillOutputReceipts(receipt KVPair)
```

向级联消息的输出收据填写收据信息。



## 11.10 AppendOutput

```
func (m *Message) AppendOutput(message IMessage)
```

向级联消息的输出收据添加收据。



## 11.11 GetTransferToMe

```
func (*IMessage) GetTransferToMe(tokenName string) (*receipt.Transfer, Error)
```

获取级联消息中前一个消息输出的收据作为本次消息调用的输入，本函数从这些输出的收据中直接解出第一个向本合约账户地址进行标准转账的收据（同时判定代币是否正确）。标准转账收据的结构定义如下：

  ```
  type Transfer struct {
      Token    Address	   `json:"token"`	// 通证或代币地址
      From     Address	   `json:"from"`	// 资金转出地址
      To       Address	   `json:"to"`		// 资金转入地址
      Value    Number	   `json:"value"`	// 转账金额（单位：cong）
  }
  ```

tokenName 为指定代币名称（为空表示创世通证）。返回标准转账对象和 Error，ErrorCode 等于 200 表示成功。



# 12 sdk/IAccount



## 12.1 Address

```
func (*IAccount) Address() Address
```

获取账户对象的账户地址。



## 12.2 PubKey

```
func (*IAccount) PubKey() PubKey
```

获取账户对象的公钥数据（可能为空）。



## 12.3 Balance

```
func (*IAccount) Balance() Number
```

在代币合约中，获取账户对象的在本合约对应的代币子账户的资金余额。在非代币合约中，返回0，如果账户地址之前没有创世通证或代币对应的子账户信息，返回余额为0。



## 12.4 BalanceOfToken

```
func (*IAccount) BalanceOfToken(token Address) Number
```

按代币地址获取账户对象指定代币子账户的资金余额。token 为指定代币地址。返回指定代币的资金余额（单位：cong），如果指定的代币地址不是一个代币，直接返回余额为0，如果账户地址之前没有指定代币对应的子账户信息，返回余额为0。



## 12.5 BalanceOfName

```
func (*IAccount) BalanceOfName(name string) Number
```

按代币名称获取账户对象指定代币子账户的资金余额。name 为指定代币名称。返回指定代币的资金余额（单位：cong），如果指定的代币名称找不到一个代币，直接返回余额为0，如果账户地址之前没有指定代币对应的子账户信息，返回余额为0。



## 12.6 BalanceOfSymbol

```
func (*IAccount) BalanceOfSymbol(symbol string) Number
```

按代币符号获取账户对象指定代币子账户的资金余额。symbol 为指定代币符号。返回指定代币的资金余额（单位：cong），如果指定的代币符号找不到一个代币，直接返回余额为0如果账户地址之前没有指定代币对应的子账户信息，返回余额为0。



## 12.7 Transfer

```
func (*IAccount) Transfer(to Address, value Number) Error
```

在代币合约中，向指定的账户地址的指定代币子账户转入资金，资金从消息发送者的指定代币子账户转出。在非代币合约中，直接返回错误。to 为接收转入资金的账户地址。value 为转入的资金数额（单位：cong）。ErrorCode 等于 200 表示执行成功。



## 12.8 TransferByToken 

```
func (*IAccount) TransferByToken(token Address, to Address, value Number) Error
```

向指定的账户地址的某种代币（指定代币地址）的子账户转入资金，资金从消息发送者的对应代币子账户转出。token 为代币地址（也可以是基础通证）。to 是接收转入资金的账户地址。value 是转入的资金数额（单位：cong）。ErrorCode等于 200 表示执行成功。




## 12.9 TransferByName

```
func (*IAccount) TransferByName(name string, to Address, value Number) Error
```

向指定的账户地址的某种代币（指定代币名称）的子账户转入资金，资金从消息发送者的对应代币子账户转出。name 为代币名称（也可以是基础通证的名称）。to 为接收转入资金的账户地址。value 为转入的资金数额（单位：cong）。ErrorCode等于 200 表示执行成功。




## 12.10 TransferBySymbol

```
func (*IAccount) TransferBySymbol(symbol string, to Address, value Number) Error
```

向指定的账户地址的某种代币（指定代币符号）的子账户转入资金，资金从消息发送者的对应代币子账户转出。symbol 为代币符号（也可以是基础通证的符号）。to 为接收转入资金的账户地址。value 为转入的资金数额（单位：cong）。ErrorCode等于 200 表示执行成功。



## 12.11 SetBalanceOfToken

```
func (a *Account) SetBalanceOfToken(tokenAddr Address, bal Number)
```

对指定的 token 设置余额。tokenAddr 为 token 的地址，bal 为余额。



# 13 sdk/IToken

token 的接口，包括查询 token 的信息以及设置 token owner 等。



## 13.1 Address

```
func (*IToken) Address() Address
```

获取代币对象的代币地址。



## 13.2 Owner

```
func (*IToken) Owner() Address
```

获取代币对象的拥有者的外部账户地址（外部账户地址有私钥与之对应）。



## 13.3 Name

```
func (*IToken) Name() string
```

获取代币的名称。



## 13.4 Symbol

```
func (*IToken) Symbol() string
```

获取代币的符号。



## 13.5 TotalSupply

```
func (*IToken) TotalSupply() Number
```

获取代币的总供应量（单位为cong）。



## 13.6 AddSupplyEnabled

```
func (*IToken) AddSupplyEnabled() bool
```

获取代币是否允许增发。



## 13.7 BurnEnabled

```
func (*IToken) BurnEnabled() bool
```

获取代币是否允许燃烧。



## 13.8 GasPrice

```
func (*IToken) GasPrice() int64
```

获取代币的燃料价格（单位为cong）。



## 13.9 SetOwner

```
func (*IToken) SetOwner(owner Address) Error
```

设置代币新的拥有者。owner 为代币新的拥有者的外部账户地址。ErrorCode等于 200 表示执行成功。



## 13.10 SetTotalSupply

```
func (*IToken) SetTotalSupply(totalSupply Number) Error
```
设置代币新的总供应量。totalSupply 为代币新的总供应量（单位为cong）。ErrorCode等于 200 表示执行成功。



## 13.11 SetGasPrice

```
func (*IToken) SetGasPrice(gasPrice Address) Error
```

设置智能合约新的拥有者。owner 为智能合约新的拥有者的外部账户地址。ErrorCode等于 200 表示执行成功。



# 14 sdk/IHelper

通过此接口可以获取其他 helper 接口。



## 14.1 AccountHelper

```
func (*IHelper) AccountHelper() IAccountHelper
```

获取账户Helper对象。



## 14.2 BlockchainHelper

```
func (*IHelper) BlockChainHelper() IBlockChainHelper
```

获取区块链Helper对象。



## 14.3 BuildHelper

```
func (*IHelper) BuildHelper() IBuildHelper
```
获取合约构建Helper对象。



## 14.4 ContractHelper

```
func (*IHelper) ContractHelper() IContractHelper
```

获取智能合约Helper对象。



## 14.5 ReceiptHelper

```
func (*IHelper) ReceiptHelper() IReceiptHelper
```

获取收据Helper对象。



## 14.6 GenesisHelper

```
func (*IHelper) GenesisHelper() IGenesisHelper
```

获取创世Helper对象。



## 14.7 StateHelper

```
func (*IHelper) StateHelper() IStateHelper
```

获取状态数据Helper对象。



## 14.8 TokenHelper

```
func (*IHelper) TokenHelper() ITokenHelper
```

获取通证/代币Helper对象。



# 15 sdk/IAccountHelper



## 15.1 AccountOf

```
func (*IAccountHelper) AccountOf(addr Address) IAccount
```

根据账户地址构造一个账户对象，用来对账户进行一些操作。addr 为账户地址。返回账户对象（可能为空）。如果输入的账户地址为空，则返回空。



## 15.2 AccountOfPubKey

```
func (*IAccountHelper) AccountOfPubKey(pubkey PubKey) IAccount
```

根据账户公钥构造一个账户对象，用来对账户进行一些操作。pubkey 为账户公钥。返回账户对象（可能为空）。如果输入的账户公钥为空，则返回空。



# 16 sdk/IBlockChainHelper

区块链信息的接口，包括计算账户地址和获取区块信息等。



## 16.1 CalcAccountFromPubkey

```
func (*IBlockchainHelper) CalcAccountFromPubkey(pubkey []byte) Address
```

根据公钥计算账户地址。pubkey 为公钥。返回计算得出的地址。如果输入的公钥长度不等于32字节，返回空地址。



## 16.2 CalcContractFromName

```
func (*IBlockchainHelper) CalcContractFromName(name string) Address
```

根据合约名称计算出合约的账户地址。name 为合约名称。返回计算得出的地址。如果输入的合约名称为空，则返回空地址。



## 16.3 CalcContractAddress

```
func (*IBlockchainHelper) CalcContractAddress(
                                name string, 
                                version string,
                                owner Address) Address
```

根据给定的合约参数计算合约地址。name 为合约名称。version 为合约版本。owner 为合约所有者的地址（注：合约被转移所有者以后合约地址不会改变）。返回计算得出的地址。



## 16.4 CheckAddress

```
func (*IBlockchainHelper) CheckAddress(addr Address)
```

校验地址的合法性。addr 为地址。



## 16.5 GetBlock

```
func (bh *BlockChainHelper) GetBlock(height int64) IBlock
```

根据高度读取区块信息。height 为区块高度。返回区块信息对象（可能为空）。如果输入的高度小于等于0，则返回空，如果输入的高度没有保存区块信息或区块数据不正常，则返回一个该高度的区块信息对象，其它参数全部为空值。



## 16.6 GetCurrentBlock

```
func GetCurrentBlock(smc ISmartContract) IBlock
```

获取当前区块的信息。



# 17 sdk/IBuildHelper

构建智能合约接口。



## 17.1 Build

```
func (*IBuildHelper) Build(meta ContractMeta) (BuildResult, Error)
```

构建合约。
结构定义如下：

  ```
  type ContractMeta struct {
  	Name         string		//合约名称
  	ContractAddr Address	//合约地址
  	OrgId        string		//合约所属组织ID
  	EffectHeight int64		//合约生效高度
  	LoseHeight   int64		//合约失效高度（固定为0）
  	CodeData     []byte		//合约代码压缩包数据
  	CodeHash     []byte		//合约代码的哈希
  	CodeDevSig   []byte		//合约开发者对合约代码压缩包的签名数据
  	CodeOrgSig   []byte		//组织对合约开发者的签名进行的签名数据
  }
  
  type Method struct {
  	MethodId 	string		// 方法ID，方法ID的计算与方法原型相关
    Gas 		int64		// 方法在调用时消耗的燃料
    ProtoType 	string		// 方法原型
  }
  
  type BuildResult struct {
  	Methods     []Method 	// 公开的方法列表
  	Interfaces  []Method 	// 公开的接口列表
  	OrgCodeHash []byte   	// 组织代码哈希（编译以后的程序名称）
  }
  ```
meta 为合约元数据。返回构建成功后合约的相关信息。ErrorCode等于 200 表示执行成功。



# 18 sdk/IContractHelper

智能合约接口，可以查询指定合约的具体信息。



## 18.1 ContractOfAddress

```
func (*IContractHelper) ContractOfAddress(addr Address) IContract
```

根据合约地址构造一个合约对象并读取合约信息，用来对合约进行一些操作。addr 为合约地址。返回合约对象（可能为空）。如果输入的合约地址为空，则返回空。



## 18.2 ContractOfName

```
func (*IContractHelper) ContractOfName(name Address) IContract
```

根据合约名称构造一个合约对象并读取合约信息，用来对合约进行一些操作。name 为合约名称。返回合约对象（可能为空）。 如果输入的合约地址为空，则返回空。



## 18.3 ContractOfToken

```
func (*IContractHelper) ContractOfToken(tokenAddr Address) IContract
```

根据代币地址构造一个针对该代币的合约对象并读取最新的合约信息，用来对合约进行一些操作。tokenAddr 为代币地址。返回合约对象（可能为空）。如果输入的代币地址为空，则返回空。



## 18.4 UpdateContractsToken

```
func (ch *ContractHelper) UpdateContractsToken(tokenAddr Address) Error
```

更新当前智能合约的 token 并返回结果。



# 19 sdk/IReceiptHelper

收据接口，用来发送一个收据。



## 19.1 Emit

```
func (*IReceiptHelper) Emit(interface Interface{})
```

发送一个收据，SDK底层实现会自动将传入的收据对象进行序列化作为本次调用合约的输出数据集中的一员。interface 为收据对象。明面上函数没有返回值，在合约调用上下文环境中会将收据输出到输出数据集。



# 20 sdk/IGenesisHelper

创世接口，获取创世的信息。



## 20.1 ChainId

```
func (*IGenesisHelper) ChainId() string
```

获取创世时指定的区块链ID。



## 20.2 Contracts

```
func (*IGenesisHelper) Contracts() []IContract
```

获取创世时设定的智能合约列表。



## 20.3 Token

```
func (*IGenesisHelper) Token() IToken
```

获取创世时设定的基础通证。



# 21 sdk/IStateHelper

状态查询以及设置接口。



## 21.1 Check

```
func (*IStateHelper) Check(key string) bool
```

根据给定的KEY值，检测在智能合约被许可的范围内存在KEY值指定的数据。key 为KEY值。



## 21.2 McCheck

```
func (sh *StateHelper) McCheck(key string) bool
```

根据给定的KEY值，检测在智能合约被许可的范围内存在KEY值指定的数据，包括缓存与数据库。key 为KEY值。



## 21.3 Get

```
func (*IStateHelper) Get(key string, defaultData Interface{}) Interface{}
```

根据给定的KEY值，在智能合约被许可的范围内读取数据。key 为KEY值，defaultData 为Value对应存储对象类型的模板。



## 21.4 GetEx

```
func (*IStateHelper) GetEx(key string, defaultData Interface{}) Interface{}
```

根据给定的KEY值，在智能合约被许可的范围内读取数据。key 为KEY值，defaultData 为Value对应存储对象类型的模板。



## 21.5 GetXXX

```
func (*IStateHelper) GetInt(key string) int
func (*IStateHelper) GetInt8(key string) int8
func (*IStateHelper) GetInt16(key string) int16
func (*IStateHelper) GetInt32(key string) int32
func (*IStateHelper) GetInt64(key string) int64
func (*IStateHelper) GetUint(key string) uint
func (*IStateHelper) GetUint8(key string) uint8
func (*IStateHelper) GetUint16(key string) uint16
func (*IStateHelper) GetUint32(key string) uint32
func (*IStateHelper) GetUint64(key string) uint64
func (*IStateHelper) GetFloat32(key string) float32
func (*IStateHelper) GetFloat64(key string) float64
func (*IStateHelper) GetBool(key string) bool
func (*IStateHelper) GetString(key string) string

func (*IStateHelper) GetInts(key string) []int
func (*IStateHelper) GetInt8s(key string) []int8
func (*IStateHelper) GetInt16s(key string) []int16
func (*IStateHelper) GetInt32s(key string) []int32
func (*IStateHelper) GetInt64s(key string) []int64
func (*IStateHelper) GetUints(key string) []uint
func (*IStateHelper) GetUint8s(key string) []uint8
func (*IStateHelper) GetUint16s(key string) []uint16
func (*IStateHelper) GetUint32s(key string) []uint32
func (*IStateHelper) GetUint64s(key string) []uint64
func (*IStateHelper) GetFloat32s(key string) []float32
func (*IStateHelper) GetFloat64s(key string) []float64
func (*IStateHelper) GetBools(key string) []bool
func (*IStateHelper) GetStrings(key string) []string
```

根据给定的KEY值，在智能合约被许可的范围内读取数据。key 为KEY值，返回基础类型数据值。如果从状态数据库中读不到数据，直接返回基础类型对应的默认值；



## 21.6 Set

```
func (*IStateHelper) Set(key string, data Interface{})
```

将输入的数据保存到状态数据库智能合约被允许的KEY值下。key 为KEY值，data 为要保存的数据对象。



## 21.7 SetXXX

```
func (*IStateHelper) SetInt(key string, v int)
func (*IStateHelper) SetInt8(key string, v int8)
func (*IStateHelper) SetInt16(key string, v int16)
func (*IStateHelper) SetInt32(key string, v int32)
func (*IStateHelper) SetInt64(key string, v int64)
func (*IStateHelper) SetUint(key string, v uint)
func (*IStateHelper) SetUint8(key string, v uint8)
func (*IStateHelper) SetUint16(key string, v uint16)
func (*IStateHelper) SetUint32(key string, v uint32)
func (*IStateHelper) SetUint64(key string, v uint64)
func (*IStateHelper) SetFloat32(key string, v float32)
func (*IStateHelper) SetFloat64(key string, v float64)
func (*IStateHelper) SetBool(key string, v bool)
func (*IStateHelper) SetString(key string, v string)

func (*IStateHelper) SetInts( key string, v []int )
func (*IStateHelper) SetInt8s( key string, v []int8 )
func (*IStateHelper) SetInt16s( key string, v []int16 )
func (*IStateHelper) SetInt32s( key string, v []int32 )
func (*IStateHelper) SetInt64s( key string, v []int64 )
func (*IStateHelper) SetUints( key string, v []uint )
func (*IStateHelper) SetUint8s( key string, v []uint8 )
func (*IStateHelper) SetUint16s( key string, v []uint16 )
func (*IStateHelper) SetUint32s( key string, v []uint32 )
func (*IStateHelper) SetUint64s( key string, v []uint64 )
func (*IStateHelper) SetFloat32s( key string, v []float32 )
func (*IStateHelper) SetFloat64s( key string, v []float64 )
func (*IStateHelper) SetBools( key string, v []bool )
func (*IStateHelper) SetStrings( key string, v []string )
```

将输入的数据保存到状态数据库智能合约被允许的KEY值下。key 为KEY值， v 为基础类型数据值。



## 21.8 McGet

```
func (*IStateHelper) McGet(key string, defaultData Interface{}) Interface{}
```

根据给定的KEY值，在智能合约被许可的范围内读取数据，并将数据缓存在内存中，在后续智能合约的调用消息中可以直接从内存中读取，而不需要再次访问数据库。key 为KEY值，defaultData 为Value对应存储对象类型的模板。



## 21.9 McGetEx

```
func (*IStateHelper) McGetEx(key string, defaultData Interface{}) Interface{}
```

根据给定的KEY值，在智能合约被许可的范围内读取数据，并将数据缓存在内存中，在后续智能合约的调用消息中可以直接从内存中读取，而不需要再次访问数据库。key 为KEY值，defaultData 为Value对应存储对象类型的模板。



## 21.10 McGetXXX

```
func (*IStateHelper) McGetInt(key string) int
func (*IStateHelper) McGetInt8(key string) int8
func (*IStateHelper) McGetInt16(key string) int16
func (*IStateHelper) McGetInt32(key string) int32
func (*IStateHelper) McGetInt64(key string) int64
func (*IStateHelper) McGetUint(key string) uint
func (*IStateHelper) McGetUint8(key string) uint8
func (*IStateHelper) McGetUint16(key string) uint16
func (*IStateHelper) McGetUint32(key string) uint32
func (*IStateHelper) McGetUint64(key string) uint64
func (*IStateHelper) McGetFloat32(key string) float32
func (*IStateHelper) McGetFloat64(key string) float64
func (*IStateHelper) McGetBool(key string) bool
func (*IStateHelper) McGetString(key string) string

func (*IStateHelper) McGetInts(key string) []int
func (*IStateHelper) McGetInt8s(key string) []int8
func (*IStateHelper) McGetInt16s(key string) []int16
func (*IStateHelper) McGetInt32s(key string) []int32
func (*IStateHelper) McGetInt64s(key string) []int64
func (*IStateHelper) McGetUints(key string) []uint
func (*IStateHelper) McGetUint8s(key string) []uint8
func (*IStateHelper) McGetUint16s(key string) []uint16
func (*IStateHelper) McGetUint32s(key string) []uint32
func (*IStateHelper) McGetUint64s(key string) []uint64
func (*IStateHelper) McGetFloat32s(key string) []float32
func (*IStateHelper) McGetFloat64s(key string) []float64
func (*IStateHelper) McGetBools(key string) []bool
func (*IStateHelper) McGetStrings(key string) []string
```

根据给定的KEY值，在智能合约被许可的范围内读取数据，并将数据缓存在内存中，在后续智能合约的调用消息中可以直接从内存中读取，而不需要再次访问数据库。key 为KEY值。返回基础类型数据值。



## 21.11 McSet

```
func (*IStateHelper) McSet(key string, interface Interface{})
```

将输入的数据保存到状态数据库智能合约被允许的KEY值下，同时更新内存缓存，在后续智能合约的调用消息中可以直接从内存中读取，而不需要再次访问数据库。key 为KEY值，interface 为要保存的数据对象。



## 21.12 McSetXXX

```
func (*IStateHelper) McSetInt(key string, v int)
func (*IStateHelper) McSetInt8(key string, v int8)
func (*IStateHelper) McSetInt16(key string, v int16)
func (*IStateHelper) McSetInt32(key string, v int32)
func (*IStateHelper) McSetInt64(key string, v int64)
func (*IStateHelper) McSetUint(key string, v uint)
func (*IStateHelper) McSetUint8(key string, v uint8)
func (*IStateHelper) McSetUint16(key string, v uint16)
func (*IStateHelper) McSetUint32(key string, v uint32)
func (*IStateHelper) McSetUint64(key string, v uint64)
func (*IStateHelper) McSetFloat32(key string, v float32)
func (*IStateHelper) McSetFloat64(key string, v float64)
func (*IStateHelper) McSetBool(key string, v bool)
func (*IStateHelper) McSetString(key string, v string)

func (*IStateHelper) McSetInts( key string, v []int )
func (*IStateHelper) McSetInt8s( key string, v []int8 )
func (*IStateHelper) McSetInt16s( key string, v []int16 )
func (*IStateHelper) McSetInt32s( key string, v []int32 )
func (*IStateHelper) McSetInt64s( key string, v []int64 )
func (*IStateHelper) McSetUints( key string, v []uint )
func (*IStateHelper) McSetUint8s( key string, v []uint8 )
func (*IStateHelper) McSetUint16s( key string, v []uint16 )
func (*IStateHelper) McSetUint32s( key string, v []uint32 )
func (*IStateHelper) McSetUint64s( key string, v []uint64 )
func (*IStateHelper) McSetFloat32s( key string, v []float32 )
func (*IStateHelper) McSetFloat64s( key string, v []float64 )
func (*IStateHelper) McSetBools( key string, v []bool )
func (*IStateHelper) McSetStrings( key string, v []string )
```

将输入的数据保存到状态数据库智能合约被允许的KEY值下，同时更新内存缓存，在后续智能合约的调用消息中可以直接从内存中读取，而不需要再次访问数据库。key 为KEY值，v 为基础类型数据值。



## 20.13 McClear

```
func (*IStateHelper) McClear(key string)
```

清除内存缓存中指定的数据。key 为KEY值。



# 22 sdk/ITokenHelper

查询 token 信息以及注册 token 接口。



## 22.1 RegisterToken

```
func (*ITokenHelper) RegisterToken(name string, 
                                   symbol string, 
                                   totalSupply Number,
                                   addSupplyEnabled bool,
                                   burnEnabled bool) (IToken, Error)
```

向区块链注册一个BRC20标准代币。name 为代币名称，symbol 为代币符号，totalSupply  为总的供应量，addSupplyEnabled 为是否允许增发，burnEnabled 为是否允许燃烧。ErrorCode等于 200 表示执行成功。



## 22.2 Token

```
func (*ITokenHelper) Token() IToken
```

在代币合约中，获取在本合约注册的代币的信息。返回创世通证或代币的信息。在非代币合约中，直接返回空。



## 22.3 TokenOfAddress

```
func (*ITokenHelper) TokenOfAddress(tokenAddr Address) IToken
```

按代币地址获取代币的信息。tokenAddr 为指定代币地址。返回指定代币的信息，如果指定的代币地址不是一个代币，直接返回空对象。



## 22.4 TokenOfName

```
func (*ITokenHelper) TokenOfName(name string) IToken
```

按代币名称获取代币的信息。name 为指定代币名称。返回指定代币的信息。如果指定的代币名称不是一个代币，直接返回空对象。



## 22.5 TokenOfSymbol

```
func (*ITokenHelper) TokenOfSymbol(symbol string) IToken
```

按代币符号获取代币的信息。symbol 为指定代币符号。返回指定代币的信息。如果指定的代币符号不是一个代币，直接返回空对象。



## 22.6 TokenOfContract

```
func (*ITokenHelper) TokenOfContract(contractAddr Address) IToken
```

按合约地址获取代币的信息。contractAddr 为指定合约地址。返回指定代币的信息。如果指定的合约地上没有注册一个代币，直接返回空对象。



## 22.7 BaseGasPrice

```
func (*ITokenHelper) BaseGasPrice() int64
```

读取基础燃料价格。返回基础燃料价格（单位为Cong）。
