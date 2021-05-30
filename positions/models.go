package positions

// OandaOpenPositionsResults see https://developer.oanda.com/rest-live-v20/position-ep/
type OandaOpenPositionsResults struct {
	Positions         []interface{} `json:"positions"`
	LastTransactionID string        `json:"lastTransactionID"`
}

// // AccountUnits ...
// type AccountUnits string
//
// // Position model. See https://developer.oanda.com/rest-live-v20/position-df/#Position
// type Position struct {
// 	Instrument              AccountUnits `json:"instrument"`
// 	PL                      AccountUnits `json:"pl"`
// 	UnrealizedPL            AccountUnits `json:"unrealizedPL"`
// 	MarginUsed              AccountUnits `json:"marginUsed"`
// 	ResettablePL            AccountUnits `json:"resettablePL"`
// 	Financing               AccountUnits `json:"financing"`
// 	Commission              AccountUnits `json:"commission"`
// 	DividendAdjustment      AccountUnits `json:"dividendAdjustment"`
// 	GuaranteedExecutionFees AccountUnits `json:"guaranteedExecutionFees"`
// 	Long                    PositionSide `json:"long"`
// 	Short                   PositionSide `json:"short"`
// }
//
// // PositionSide model. See https://developer.oanda.com/rest-live-v20/position-df/#PositionSide
// type PositionSide struct {
// 	Units      float64 `json:"units"`
// 	PriceValue float64 `json:"priceValue"`
// }
