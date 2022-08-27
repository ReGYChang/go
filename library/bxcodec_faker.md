- [Introduction](#introduction)
- [Installation](#installation)
- [Examples](#examples)
  - [Example Without Tag](#example-without-tag)
  - [Example With Tags](#example-with-tags)
  - [Example With Tags Length Bounds](#example-with-tags-length-bounds)
- [Conclusion](#conclusion)
- [Reference](#reference)

# Introduction

在日常開發過程中通常需要大批量的假資料用來測試程式, Go 開發中通常使用 `bxcodec/faker` 來批量生成測試資料

`bxcodec/faker` 功能齊全, 支持多種語言(包含中文), 足夠滿足日常開發所需, 其支援的資料型別如下:

- int, int8, int16, int32 & int64
- []int, []int8, []int16, []int32 & []int64
- bool & []bool
- string & []string
- float32, float64, []float32 &[]float64
- time.Time & []time.Time
- Nested Struct Field

>❗️ 生成資料的結構體字段必須為 public, 否則會觸發 panic

為生成資料而提供的方法如下:

```go
import (
	"fmt"

	"github.com/bxcodec/faker/v4"
)

// Single fake function can be used for retrieving particular values.
func Example_singleFakeData() {

	// Address
	fmt.Println(faker.Latitude())  // => 81.12195
	fmt.Println(faker.Longitude()) // => -84.38158

	// Datetime
	fmt.Println(faker.UnixTime())   // => 1197930901
	fmt.Println(faker.Date())       // => 1982-02-27
	fmt.Println(faker.TimeString()) // => 03:10:25
	fmt.Println(faker.MonthName())  // => February
	fmt.Println(faker.YearString()) // => 1994
	fmt.Println(faker.DayOfWeek())  // => Sunday
	fmt.Println(faker.DayOfMonth()) // => 20
	fmt.Println(faker.Timestamp())  // => 1973-06-21 14:50:46
	fmt.Println(faker.Century())    // => IV
	fmt.Println(faker.Timezone())   // => Asia/Jakarta
	fmt.Println(faker.Timeperiod()) // => PM

	// Internet
	fmt.Println(faker.Email())      // => mJBJtbv@OSAaT.com
	fmt.Println(faker.MacAddress()) // => cd:65:e1:d4:76:c6
	fmt.Println(faker.DomainName()) // => FWZcaRE.org
	fmt.Println(faker.URL())        // => https://www.oEuqqAY.org/QgqfOhd
	fmt.Println(faker.Username())   // => lVxELHS
	fmt.Println(faker.IPv4())       // => 99.23.42.63
	fmt.Println(faker.IPv6())       // => 975c:fb2c:2133:fbdd:beda:282e:1e0a:ec7d
	fmt.Println(faker.Password())   // => dfJdyHGuVkHBgnHLQQgpINApynzexnRpgIKBpiIjpTP

	// Words and Sentences
	fmt.Println(faker.Word())      // => nesciunt
	fmt.Println(faker.Sentence())  // => Consequatur perferendis voluptatem accusantium.
	fmt.Println(faker.Paragraph()) // => Aut consequatur sit perferendis accusantium voluptatem. Accusantium perferendis consequatur voluptatem sit aut. Aut sit accusantium consequatur voluptatem perferendis. Perferendis voluptatem aut accusantium consequatur sit.

	// Payment
	fmt.Println(faker.CCType())             // => American Express
	fmt.Println(faker.CCNumber())           // => 373641309057568
	fmt.Println(faker.Currency())           // => USD
	fmt.Println(faker.AmountWithCurrency()) // => USD 49257.100

	// Person
	fmt.Println(faker.TitleMale())       // => Mr.
	fmt.Println(faker.TitleFemale())     // => Mrs.
	fmt.Println(faker.FirstName())       // => Whitney
	fmt.Println(faker.FirstNameMale())   // => Kenny
	fmt.Println(faker.FirstNameFemale()) // => Jana
	fmt.Println(faker.LastName())        // => Rohan
	fmt.Println(faker.Name())            // => Mrs. Casandra Kiehn

	// Phone
	fmt.Println(faker.Phonenumber())         // -> 201-886-0269
	fmt.Println(faker.TollFreePhoneNumber()) // => (777) 831-964572
	fmt.Println(faker.E164PhoneNumber())     // => +724891571063

	//  UUID
	fmt.Println(faker.UUIDHyphenated()) // => 8f8e4463-9560-4a38-9b0c-ef24481e4e27
	fmt.Println(faker.UUIDDigit())      // => 90ea6479fd0e4940af741f0a87596b73

	fmt.Println(faker.Word())

	faker.ResetUnique() // Forget all generated unique values
}
```

# Installation

```go
go get github.com/bxcodec/faker/v3
```

# Examples

## Example Without Tag

```go
import (
	"fmt"

	"github.com/bxcodec/faker/v4"
)

// SomeStruct ...
type SomeStruct struct {
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	String   string
	Bool     bool
	SString  []string
	SInt     []int
	SInt8    []int8
	SInt16   []int16
	SInt32   []int32
	SInt64   []int64
	SFloat32 []float32
	SFloat64 []float64
	SBool    []bool
	Struct   AStruct
}

// AStruct ...
type AStruct struct {
	Number        int64
	Height        int64
	AnotherStruct BStruct
}

// BStruct ...
type BStruct struct {
	Image string
}

// You also can use faker to generate your structs data randomly without any tag.
// And it will fill the data based on its data-type.
func Example_withoutTag() {
	a := SomeStruct{}
	err := faker.FakeData(&a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", a)
	/*
		Result:
		{
		    Int:5231564546548329
		    Int8:52
		    Int16:8
		    Int32:2046951236
		    Int64:1486554863682234423
		    String:ELIQhSfhHmWxyzRPOTNblEQsp
		    Bool:false
		    SString:[bzYwplRGUAXPwatnpVMWSYjep zmeuJVGHHgmIsuyWmLJnDmbTI FqtejCwoDyMBWatoykIzorCJZ]
		    SInt:[11661230973193 626062851427 12674621422454 5566279673347]
		    SInt8:[12 2 58 22 11 66 5 88]
		    SInt16:[29295225 8411281 69902706328]
		    SInt32:[60525685140 2733640366211 278043484637 5167734561481]
		    SInt64:[81684520429188374184 9917955420365482658170 996818723778286568 163646873275501565]
		    SFloat32:[0.556428 0.7692596 0.6779895 0.29171365 0.95445055]
		    SFloat64:[0.44829454895586585 0.5495675898536803 0.6584538253883265]
		    SBool:[true false true false true true false]
		    Struct:{
		        Number:1
		        Height:26
		        AnotherStruct:{
		            Image:RmmaaHkAkrWHldVbBNuQSKlRb
		        }
		    }
		}
	*/
}
```

## Example With Tags

```go
import (
	"fmt"

	"github.com/bxcodec/faker/v4"
)

// SomeStructWithTags ...
type SomeStructWithTags struct {
	Latitude           float32 `faker:"lat"`
	Longitude          float32 `faker:"long"`
	CreditCardNumber   string  `faker:"cc_number"`
	CreditCardType     string  `faker:"cc_type"`
	Email              string  `faker:"email"`
	DomainName         string  `faker:"domain_name"`
	IPV4               string  `faker:"ipv4"`
	IPV6               string  `faker:"ipv6"`
	Password           string  `faker:"password"`
	Jwt                string  `faker:"jwt"`
	PhoneNumber        string  `faker:"phone_number"`
	MacAddress         string  `faker:"mac_address"`
	URL                string  `faker:"url"`
	UserName           string  `faker:"username"`
	TollFreeNumber     string  `faker:"toll_free_number"`
	E164PhoneNumber    string  `faker:"e_164_phone_number"`
	TitleMale          string  `faker:"title_male"`
	TitleFemale        string  `faker:"title_female"`
	FirstName          string  `faker:"first_name"`
	FirstNameMale      string  `faker:"first_name_male"`
	FirstNameFemale    string  `faker:"first_name_female"`
	LastName           string  `faker:"last_name"`
	Name               string  `faker:"name"`
	UnixTime           int64   `faker:"unix_time"`
	Date               string  `faker:"date"`
	Time               string  `faker:"time"`
	MonthName          string  `faker:"month_name"`
	Year               string  `faker:"year"`
	DayOfWeek          string  `faker:"day_of_week"`
	DayOfMonth         string  `faker:"day_of_month"`
	Timestamp          string  `faker:"timestamp"`
	Century            string  `faker:"century"`
	TimeZone           string  `faker:"timezone"`
	TimePeriod         string  `faker:"time_period"`
	Word               string  `faker:"word"`
	Sentence           string  `faker:"sentence"`
	Paragraph          string  `faker:"paragraph"`
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"amount"`
	AmountWithCurrency string  `faker:"amount_with_currency"`
	UUIDHypenated      string  `faker:"uuid_hyphenated"`
	UUID               string  `faker:"uuid_digit"`
	Skip               string  `faker:"-"`
	PaymentMethod      string  `faker:"oneof: cc, paypal, check, money order"` // oneof will randomly pick one of the comma-separated values supplied in the tag
	AccountID          int     `faker:"oneof: 15, 27, 61"`                     // use commas to separate the values for now. Future support for other separator characters may be added
	Price32            float32 `faker:"oneof: 4.95, 9.99, 31997.97"`
	Price64            float64 `faker:"oneof: 47463.9463525, 993747.95662529, 11131997.978767990"`
	NumS64             int64   `faker:"oneof: 1, 2"`
	NumS32             int32   `faker:"oneof: -3, 4"`
	NumS16             int16   `faker:"oneof: -5, 6"`
	NumS8              int8    `faker:"oneof: 7, -8"`
	NumU64             uint64  `faker:"oneof: 9, 10"`
	NumU32             uint32  `faker:"oneof: 11, 12"`
	NumU16             uint16  `faker:"oneof: 13, 14"`
	NumU8              uint8   `faker:"oneof: 15, 16"`
	NumU               uint    `faker:"oneof: 17, 18"`
	PtrNumU            *uint   `faker:"oneof: 19, 20"`
}

func Example_withTags() {

	a := SomeStructWithTags{}
	err := faker.FakeData(&a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", a)
	/*
		Result:
		{
			Latitude: 81.12195
			Longitude: -84.38158
			CreditCardType: American Express
			CreditCardNumber: 373641309057568
			Email: mJBJtbv@OSAaT.ru
			DomainName: FWZcaRE.ru,
			IPV4: 99.23.42.63
			IPV6: 975c:fb2c:2133:fbdd:beda:282e:1e0a:ec7d
			Password: dfJdyHGuVkHBgnHLQQgpINApynzexnRpgIKBpiIjpTPOmNyMFb
			Jwt: HDMNSOKhEIYkPIuHcVjfCtHlKkaqLGrUEqjKVkgR.HDMNSOKhEIYkPIuHcVjfCtHlKkaqLGrUEqjKVkgR.HDMNSOKhEIYkPIuHcVjfCtHlKkaqLGrUEqjKVkgR
			PhoneNumber: 792-153-4861
			MacAddress: cd:65:e1:d4:76:c6
			URL: https://www.oEuqqAY.org/QgqfOhd
			UserName: lVxELHS
			TollFreeNumber: (777) 831-964572
			E164PhoneNumber: +724891571063
			TitleMale: Mr.
			TitleFemale: Queen
			FirstName: Whitney
			FirstNameMale: Kenny
			FirstNameFemale: Jana
			LastName: Rohan
			Name: Miss Casandra Kiehn
			UnixTime: 1197930901
			Date: 1982-02-27
			Time: 03:10:25
			MonthName: February
			Year: 1996
			DayOfWeek: Sunday
			DayOfMonth: 20
			Timestamp: 1973-06-21 14:50:46
			Century: IV
			TimeZone: Canada/Eastern
			TimePeriod: AM
			Word: nesciunt
			Sentence: Consequatur perferendis aut sit voluptatem accusantium.
			Paragraph: Aut consequatur sit perferendis accusantium voluptatem. Accusantium perferendis consequatur voluptatem sit aut. Aut sit accusantium consequatur voluptatem perferendis. Perferendis voluptatem aut accusantium consequatur sit.
			Currency: IRR,
			Amount: 88.990000,
			AmountWithCurrency: XBB 49257.100000,
			UUIDHypenated: 8f8e4463-9560-4a38-9b0c-ef24481e4e27,
			UUID: 90ea6479fd0e4940af741f0a87596b73,
			PaymentMethod: paypal,
			AccountID: 61,
			Price32: 4.95,
			Price64: 993747.95662529
			NumS64:	1
			NumS32:	-3
			NumS16:	5
			NumS8:	-8
			NumU64:	9
			NumU32:	11
			NumU16:	13
			NumU8:	15
			NumU:	17
			PtrNumU: 19
			Skip:
		}
	*/

}
```

## Example With Tags Length Bounds

```go
import (
	"fmt"

	"github.com/bxcodec/faker/v4"
	"github.com/bxcodec/faker/v4/pkg/options"
)

// SomeStructWithLengthAndBoundary ...
type SomeStructWithLengthAndBoundary struct {
	Inta  int   `faker:"boundary_start=5, boundary_end=10"`
	Int8  int8  `faker:"boundary_start=100, boundary_end=1000"`
	Int16 int16 `faker:"boundary_start=123, boundary_end=1123"`
	Int32 int32 `faker:"boundary_start=-10, boundary_end=8123"`
	Int64 int64 `faker:"boundary_start=31, boundary_end=88"`

	UInta  uint   `faker:"boundary_start=35, boundary_end=152"`
	UInt8  uint8  `faker:"boundary_start=5, boundary_end=1425"`
	UInt16 uint16 `faker:"boundary_start=245, boundary_end=2125"`
	UInt32 uint32 `faker:"boundary_start=0, boundary_end=40"`
	UInt64 uint64 `faker:"boundary_start=14, boundary_end=50"`

	Float32 float32 `faker:"boundary_start=12.65, boundary_end=184.05"`
	Float64 float64 `faker:"boundary_start=1.256, boundary_end=3.4"`

	ASString []string          `faker:"len=50"`
	SString  string            `faker:"len=25"`
	MSString map[string]string `faker:"len=30"`
	MIint    map[int]int       `faker:"boundary_start=5, boundary_end=10"`
}

// You can set length for your random strings also set boundary for your integers.
func Example_withTagsLengthAndBoundary() {
	a := SomeStructWithLengthAndBoundary{}
	_ = faker.FakeData(&a, options.WithRandomMapAndSliceMaxSize(20)) // Random generated map or array size wont exceed 20...
	fmt.Printf("%+v", a)
	// Result:
	/*
	   {
	       Inta:7
	       Int8:-102
	       Int16:556
	       Int32:113
	       Int64:70
	       UInta:78
	       UInt8:54
	       UInt16:1797
	       UInt32:8
	       UInt64:34
	       Float32:60.999058
	       Float64:2.590148738554016
	       ASString:[
	           geHYIpEoQhQdijFooVEAOyvtTwJOofbQPJdbHvEEdjueZaKIgI
	           WVJBBtmrrVccyIydAiLSkMwWbFzFMEotEXsyUXqcmBTVORlkJK
	           xYiRTRSZRuGDcMWYoPALVMZgIXoTQtmdGXQfbISKJiavLspuBV
	           qsoiYlyRbXLDAMoIdQhgMriODYWCTEYepmjaldWLLjkulDGuQN
	           GQXUlqNkVjPKodMebPIeoZZlxfhbQJOjHSRjUTrcgBFPeDZIxn
	           MEeRkkLceDsqKLEJFEjJxHtYrYxQMxYcuRHEDGYSPbELDQLSsj
	           tWUIACjQWeiUhbboGuuEQIhUJCRSBzVImpYwOlFbsjCRmxboZW
	           ZDaAUZEgFKbMJoKIMpymTreeZGWLXNCfVzaEyWNdbkaZOcsfst
	           uwlsZBMlEknIBsALpXRaplZWVtXTKzsWglRVBpmfsQfqraiEYA
	           AXszbzsOzYPYeXHXRwoPmoPoBxopdFFvWMBTPCxESTepRpjlnB
	           kTuOPHlUrSzUQRmZMYplWbyoBbWzQYCiydyzurOduhjuyiGrCE
	           FZbeLMbelIeCMnixknIARZRbwALObGXADQqianJbkiEAqqpdnK
	           TiQrZbnkvxEciyKXlliUDOGVdpMoAsHSalFbLcYyXxNFLAhqjy
	           KlbjbloxkWKSqvUfJQPpFLoddWgeABfYUoaAnylKmEHwxgNsnO
	           ]
	       SString:VVcaPSFrOPYlEkpVyTRbSZneB
	       MSString:map[
	           ueFBFTTmqDwrXDoXAYTRhQRmLXhudA:AhQewvZfrlytbAROzGjpXUmNQzIoGl
	           fZwrCsFfZwqMsDJXOUYIacflFIeyFU:VMufFCRRHTtuFthOrRAMbzbKVJHnvJ
	           rHDQTyZqZVSPLwTtZfNSwKWrgmRghL:lRSXNHkhUyjDuBgoAfrQwOcHYilqRB
	           BvCpQJMHzKXKbOoAnTXkLCNxKshwWr:tiNFrXAXUtdywkyygWBrEVrmAcAepD
	           uWWKgHKTkUgAZiopAIUmgVWrkrceVy:GuuDNTUiaBtOKwWrMoZDiyaOPxywnq
	           HohMjOdMDkAqimKPTgdjUorydpKkly:whAjmraukcZczskqycoJELlMJTghca
	           umEgMBGUvBptdKImKsoWXMGJJoRbgT:tPpgHgLEyHmDOocOiSgTbXQHVduLxP
	           SRQLHjBXCXKvbLIktdKeLwMnIFOmbi:IJBpLyTcraOxOUtwSKTisjElpulkTL
	           dbnDeJZLqMXQGjbTSNxPSlfDHGCghU:JWrymovFwNWbIQBxPpQmlgJsgpXcui
	           roraKNGnBXnrJlsxTnFgxHyZeTXdAC:XIcLWqUAQAbfkRrgfjrTVxZCvRJXyl
	           TrvxqVVjXAboYDPvUglSJQrltPjzLx:nBhWdfNPybnNnCyQlSshWKOnwUMQzL
	           dTHhWJWMwfVvKpIKTFCaoBJgKmnfbD:ixjNHsvSkRkFiNLpgUzIKPsheqhCeY
	           lWyBrtfcGWiNbSTJZJXwOPvVngZZMk:kvlYeGgwguVtiafGKjHWsYWewbaXte
	           bigsYNfVcNMGtnzgaqEjeRRlIcUdbR:hYOnJupEOvblTTEYzZYPuTVmvTmiit
	           ]
	       MIint:map[7:7 5:7 8:8 9:5 6:5]
	   }
	*/
}
```

# Conclusion

`bxcodec/faker` 基本能滿足常見需求, 但也有一些需求無法滿足, 如其不支持以下資料型別:

- map[interface{}]interface{}
- map[any_type]interface{}
- map[interface{}]any_type

其不完全支持自定義型別, 最穩妥的使用方式為不使用任何自定義型別以避免 panic

# Reference

- [https://github.com/bxcodec/faker](https://github.com/bxcodec/faker)