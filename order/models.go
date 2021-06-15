package order

// Type ...
type Type string

// TimeInForce see http://developer.oanda.com/rest-live-v20/order-df/#TimeInForce
type TimeInForce string

// PositionFill see http://developer.oanda.com/rest-live-v20/order-df/#OrderPositionFill
type PositionFill string

const (
	// MARKET Order Type
	MARKET = "MARKET"

	/**
	 * Time In Force Constants
	 **/

	// GTC Good Until Cancelled
	GTC = "GTC"

	// GTD Good Until Date
	GTD = "GTD"

	// GFD Good For Day
	GFD = "GFD"

	// FOK Immediately Fill or Kill
	FOK = "FOK"

	// IOC Immediately Partially Fill or Kill
	IOC = "IOC"

	/**
	 * Position Fill Constants
	 **/

	// OPEN_ONLY ...
	OPEN_ONLY = "OPEN_ONLY"

	// REDUCE_FIRST ...
	REDUCE_FIRST = "REDUCE_FIRST"

	// REDUCE_ONLY ...
	REDUCE_ONLY = "REDUCE_ONLY"

	// DEFAULT ...
	DEFAULT = "DEFAULT"
)

// Order ...
type Order struct {
	Type          Type
	Units         int64
	Instrument    string
	TimeInForce   TimeInForce
	StopLoss      float64
	TakeProfit    float64
	ClientID      string
	ClientTag     string
	ClientComment string
}

type oandaOrderPayload struct {
	Order *oandaOrder `json:"order"`
}

// Payload ...
type oandaOrder struct {
	Type                   Type                          `json:"type"`
	Units                  string                        `json:"units"`
	Instrument             string                        `json:"instrument"`
	TimeInForce            TimeInForce                   `json:"timeInForce"`
	PriceBound             string                        `json:"priceBound,omitempty"`
	PositionFill           PositionFill                  `json:"positionFill"`
	TakeProfitOnFill       *oandaTakeProfitDetails       `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill         *oandaStopLossDetails         `json:"stopLossOnFill,omitempty"`
	TrailingStopLossOnFill *oandaTrailingStopLossDetails `json:"trailingStopLossOnFill,omitempty"`
	ClientExtensions       *oandaClientExtensions        `json:"clientExtensions,omitempty"`
}

// TakeProfitDetails ...
type oandaTakeProfitDetails struct {
	Price            string                 `json:"price"`
	TimeInForce      TimeInForce            `json:"timeInForce"`
	GtdTime          string                 `json:"gtdTime,omitempty"`
	ClientExtensions *oandaClientExtensions `json:"clientExtensions,omitempty"`
}

// StopLossDetails ...
type oandaStopLossDetails struct {
	Price            string                 `json:"price"`
	TimeInForce      TimeInForce            `json:"timeInForce"`
	GtdTime          string                 `json:"gtdTime,omitempty"`
	ClientExtensions *oandaClientExtensions `json:"clientExtensions,omitempty"`
}

// TrailingStopLossDetails ...
type oandaTrailingStopLossDetails struct {
	Distance         string                 `json:"distance"`
	TimeInForce      TimeInForce            `json:"timeInForce"`
	GtdTime          string                 `json:"gtdTime,omitempty"`
	ClientExtensions *oandaClientExtensions `json:"clientExtensions,omitempty"`
}

type oandaClientExtensions struct {
	ID      string `json:"id"`
	Tag     string `json:"tag"`
	Comment string `json:"comment,omitempty"`
}
