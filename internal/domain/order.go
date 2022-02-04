package domain

import (
	"bytes"
	"text/template"
	"time"
)

type OrderRequest struct {
	Symbol     string
	Side       string
	Size       float64
	LimitPrice float64
}

type (
	OrderResponse struct {
		Result     string     `json:"result"`
		SendStatus SendStatus `json:"sendStatus"`
		ServerTime time.Time  `json:"serverTime"`
		Error      string     `json:"error"`
	}

	SendStatus struct {
		OrderID      string       `json:"order_id"`
		Status       string       `json:"status"`
		ReceivedTime time.Time    `json:"receivedTime"`
		OrderEvents  []OrderEvent `json:"orderEvents"`
	}

	OrderEvent struct {
		Type            string  `json:"type"`
		ReducedQuantity float64 `json:"reducedQuantity"` // for place and edit orders
		Order           Order   `json:"order"`           // except edit order
		UID             string  `json:"uid"`             // for reject and cancel orders
		// For reject order
		Reason string `json:"reason"`
		// For execution order
		Amount               float64 `json:"amount"`
		Price                float64 `json:"price"`
		ExecutionID          string  `json:"executionId"`
		TakerReducedQuantity float64 `json:"takerReducedQuantity"`
		OrderPriorExecution  Order   `json:"orderPriorExecution"`
	}

	Order struct {
		OrderID     string    `json:"orderId"`
		ReducedOnly bool      `json:"reduceOnly"`
		Symbol      string    `json:"symbol"`
		Quantity    float64   `json:"quantity"`
		Side        string    `json:"side"`
		Filled      float64   `json:"filled"`
		Type        string    `json:"type"`
		Timestamp   time.Time `json:"timestamp"`
	}
)

var tgMsgTemplate *template.Template

func InitTemplate() (err error) {
	str := `
Result: {{ .Result }}
{{- if .Error }}
Error: {{ .Error }}
{{- else }}
{{ with .SendStatus }}
OrderID: {{ .OrderID }}
Status: {{ .Status }}
{{ range .OrderEvents}}
Type: {{ .Type }}
{{ if .Reason }}
Reason: {{ .Reason}}
{{ else }}{{ with .OrderPriorExecution }}
Symbol: {{ .Symbol }}
Side: {{ .Side }}{{- end }}
Amount: {{ .Amount }}
Price: {{ .Price }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

ServerTime: {{ .ServerTime }}
`
	tgMsgTemplate, err = template.New("tg").Parse(str)
	return
}

func (response OrderResponse) ToText() (string, error) {
	var text bytes.Buffer
	err := tgMsgTemplate.Execute(&text, response)

	return text.String(), err
}
