package errors

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("email already exists")

	ErrWrongEmailOrPassword = errors.New("wrong email or password")

	ErrUserNotFound = errors.New("user not found")

	ErrNoPostsFound = errors.New("no posts found")

	ErrInvalidFields = errors.New("invalid fields")

	ErrNotEnoughQuota = errors.New("not enough quota")

	ErrNoInvoicesFound = errors.New("no invoices found")

	ErrInvoiceNotAwaitingPayment = errors.New("invoice not awaiting payment")

	ErrInvalidStatusUpdate = errors.New("invalid Status Update")

	ErrNotAuthorized = errors.New("not authorized")

	ErrInvalidReferralCode = errors.New("invalid referral code")

	ErrInvalidVoucher = errors.New("invalid voucher")

	ErrVoucherExpired = errors.New("voucher expired")
)
