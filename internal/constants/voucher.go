package constants

type VoucherStatus int

const (
	AVAILABLE VoucherStatus = iota + 1
	PENDING
	USED
)

func (e VoucherStatus) String() string {
	switch e {
	case AVAILABLE:
		return "AVAILABLE"
	case PENDING:
		return "PENDING"
	case USED:
		return "USED"
	default:
		return ""
	}
}
