package automation

type RunResult struct {
	RemindersSent          int `json:"reminders_sent"`
	ContractsNearExpiry    int `json:"contracts_near_expiry"`
	RecurringChargesMade   int `json:"recurring_charges_made"`
	ProjectsSuspended      int `json:"projects_suspended"`
	MonthlyReportGenerated int `json:"monthly_report_generated"`
}
