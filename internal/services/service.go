package services

import (
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/repositories"
)

type Services struct {
	Auth             IAuthService
	User             IUserService
	Post             IPostService
	History          IHistoryService
	Subscription     ISubscriptionService
	UserSubscription IUserSubscriptionService
	Invoice          IInvoiceService
	Gift             IGiftService
	Voucher          IVoucherService
}

func New(r *repositories.Repositories) *Services {
	return &Services{
		Auth: NewAuthService(&AuthServiceConfig{
			userRepository:      r.Users,
			userTokenRepository: r.UserTokens,
			hasher:              helpers.NewHasher(),
		}),
		User: NewUserService(&UserServiceConfig{userRepository: r.Users}),
		Post: NewPostService(&PostServiceConfig{postRepository: r.Posts, historyRepository: r.Histories}),
		History: NewHistoryService(&HistoryServiceConfig{
			historyRepository: r.Histories,
		}),
		Subscription: NewSubscriptionService(&SubscriptionServiceConfig{subscriptionRepository: r.Subscriptions}),
		UserSubscription: NewUserSubscriptionService(&UserSubscriptionServiceConfig{
			subscriptionRepository:     r.Subscriptions,
			userSubscriptionRepository: r.UserSubscriptions,
		}),
		Invoice: NewInvoiceService(&InvoiceServiceConfig{
			invoiceRepository:          r.Invoices,
			subscriptionRepository:     r.Subscriptions,
			userSubscriptionRepository: r.UserSubscriptions,
			giftRepository:             r.Gifts,
			userGiftRepository:         r.UserGifts,
			userRepository:             r.Users,
			voucherRepository:          r.Vouchers,
			userVoucherRepository:      r.UserVouchers,
			userSpendingRepository:     r.UserSpendings,
		}),
		Gift: NewGiftService(&GiftServiceConfig{
			giftRepository: r.Gifts,
		}),
		Voucher: NewVoucherService(&VoucherServiceConfig{
			voucherRepository:     r.Vouchers,
			userVoucherRepository: r.UserVouchers,
		}),
	}
}
