package finance

type Dashboard struct {
	MonthlyRevenueCents uint64 `json:"monthly_revenue_cents"`
	MRRCents            uint64 `json:"mrr_cents"`
	PendingPayments     int    `json:"pending_payments"`
	PaymentsReceived    int    `json:"payments_received"`
	ActiveContracts     int    `json:"active_contracts"`
	ExpiringContracts   int    `json:"expiring_contracts"`
}

type CashFlowPoint struct {
	Date         string `json:"date"`
	InCents      uint64 `json:"in_cents"`
	PendingCents uint64 `json:"pending_cents"`
}

type DefaultMetrics struct {
	OverdueCharges     int     `json:"overdue_charges"`
	DefaultRatePercent float64 `json:"default_rate_percent"`
	GrowthPercent      float64 `json:"growth_percent"`
}
