package constants

type GiftStatus int

const (
	PROCESSED GiftStatus = iota + 1
	COMPLETED
	CANCELLED
)

func (e GiftStatus) String() string {
	switch e {
	case PROCESSED:
		return "PROCESSED"
	case COMPLETED:
		return "COMPLETED"
	case CANCELLED:
		return "CANCELLED"
	default:
		return ""
	}
}
