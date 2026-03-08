package finance

import (
	"context"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	entcharge "github.com/aggi-tech/aggipay/ent/charge"
	entpaymentrecord "github.com/aggi-tech/aggipay/ent/paymentrecord"
	entservicecontract "github.com/aggi-tech/aggipay/ent/servicecontract"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Dashboard(ctx context.Context, ownerUserID string, now time.Time) (Dashboard, error) {
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	payments, err := r.client.PaymentRecord.Query().
		Where(
			entpaymentrecord.Status("paid"),
			entpaymentrecord.PaidAtGTE(monthStart),
			entpaymentrecord.PaidAtLT(monthEnd),
		).All(ctx)
	if err != nil {
		return Dashboard{}, err
	}
	var monthlyRevenue uint64
	for _, p := range payments {
		monthlyRevenue += p.AmountCents
	}

	contracts, err := r.client.ServiceContract.Query().
		Where(
			entservicecontract.OwnerUserID(ownerUserID),
			entservicecontract.Status("active"),
			entservicecontract.BillingType("monthly"),
		).All(ctx)
	if err != nil {
		return Dashboard{}, err
	}
	var mrr uint64
	for _, c := range contracts {
		mrr += c.AmountCents
	}

	pendingCount, err := r.client.Charge.Query().Where(
		entcharge.OwnerUserID(ownerUserID),
		entcharge.Status("pending"),
	).Count(ctx)
	if err != nil {
		return Dashboard{}, err
	}

	receivedCount, err := r.client.PaymentRecord.Query().Where(
		entpaymentrecord.Status("paid"),
	).Count(ctx)
	if err != nil {
		return Dashboard{}, err
	}

	activeContracts, err := r.client.ServiceContract.Query().Where(
		entservicecontract.OwnerUserID(ownerUserID),
		entservicecontract.Status("active"),
	).Count(ctx)
	if err != nil {
		return Dashboard{}, err
	}

	expiringCount, err := r.client.ServiceContract.Query().Where(
		entservicecontract.OwnerUserID(ownerUserID),
		entservicecontract.Status("active"),
		entservicecontract.EndDateNotNil(),
		entservicecontract.EndDateLTE(now.AddDate(0, 0, 30)),
	).Count(ctx)
	if err != nil {
		return Dashboard{}, err
	}

	return Dashboard{
		MonthlyRevenueCents: monthlyRevenue,
		MRRCents:            mrr,
		PendingPayments:     pendingCount,
		PaymentsReceived:    receivedCount,
		ActiveContracts:     activeContracts,
		ExpiringContracts:   expiringCount,
	}, nil
}

func (r *Repository) CashFlow(ctx context.Context, ownerUserID string, now time.Time) ([]CashFlowPoint, error) {
	start := now.AddDate(0, 0, -30)
	payments, err := r.client.PaymentRecord.Query().Where(
		entpaymentrecord.PaidAtGTE(start),
		entpaymentrecord.PaidAtNotNil(),
	).All(ctx)
	if err != nil {
		return nil, err
	}

	pendingCharges, err := r.client.Charge.Query().Where(
		entcharge.OwnerUserID(ownerUserID),
		entcharge.Status("pending"),
		entcharge.DueDateGTE(start),
	).All(ctx)
	if err != nil {
		return nil, err
	}

	inByDay := map[string]uint64{}
	for _, p := range payments {
		if p.PaidAt == nil {
			continue
		}
		key := p.PaidAt.Format("2006-01-02")
		inByDay[key] += p.AmountCents
	}

	pendingByDay := map[string]uint64{}
	for _, c := range pendingCharges {
		key := c.DueDate.Format("2006-01-02")
		pendingByDay[key] += c.AmountCents
	}

	points := make([]CashFlowPoint, 0, 31)
	for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
		k := d.Format("2006-01-02")
		points = append(points, CashFlowPoint{Date: k, InCents: inByDay[k], PendingCents: pendingByDay[k]})
	}
	return points, nil
}

func (r *Repository) Defaults(ctx context.Context, ownerUserID string, now time.Time) (DefaultMetrics, error) {
	overdueCount, err := r.client.Charge.Query().Where(
		entcharge.OwnerUserID(ownerUserID),
		entcharge.Status("overdue"),
	).Count(ctx)
	if err != nil {
		return DefaultMetrics{}, err
	}
	allCharges, err := r.client.Charge.Query().Where(entcharge.OwnerUserID(ownerUserID)).Count(ctx)
	if err != nil {
		return DefaultMetrics{}, err
	}

	defaultRate := 0.0
	if allCharges > 0 {
		defaultRate = float64(overdueCount) * 100 / float64(allCharges)
	}

	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	prevMonthStart := thisMonthStart.AddDate(0, -1, 0)

	thisMonthPayments, err := r.client.PaymentRecord.Query().Where(
		entpaymentrecord.PaidAtGTE(thisMonthStart),
		entpaymentrecord.PaidAtLT(thisMonthStart.AddDate(0, 1, 0)),
	).All(ctx)
	if err != nil {
		return DefaultMetrics{}, err
	}
	prevMonthPayments, err := r.client.PaymentRecord.Query().Where(
		entpaymentrecord.PaidAtGTE(prevMonthStart),
		entpaymentrecord.PaidAtLT(thisMonthStart),
	).All(ctx)
	if err != nil {
		return DefaultMetrics{}, err
	}

	var thisMonth uint64
	for _, p := range thisMonthPayments {
		thisMonth += p.AmountCents
	}
	var prevMonth uint64
	for _, p := range prevMonthPayments {
		prevMonth += p.AmountCents
	}

	growth := 0.0
	if prevMonth > 0 {
		growth = (float64(thisMonth-prevMonth) * 100) / float64(prevMonth)
	}

	return DefaultMetrics{OverdueCharges: overdueCount, DefaultRatePercent: defaultRate, GrowthPercent: growth}, nil
}
