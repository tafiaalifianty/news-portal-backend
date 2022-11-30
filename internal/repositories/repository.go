package repositories

import "gorm.io/gorm"

type Repositories struct {
	Users             IUserRepository
	Posts             IPostRepository
	Histories         IHistoryRepository
	Subscriptions     ISubscriptionRepository
	UserSubscriptions IUserSubscriptionRepository
	UserTokens        IUserTokenRepository
	Invoices          IInvoiceRepository
	Gifts             IGiftRepository
	UserGifts         IUserGiftRepository
	Vouchers          IVoucherRepository
	UserVouchers      IUserVoucherRepository
	UserSpendings     IUserSpendingRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:         NewUserRepository(&UserRepositoryConfig{db: db}),
		Posts:         NewPostRepository(&PostRepositoryConfig{db: db}),
		Histories:     NewHistoryRepository(&HistoryRepositoryConfig{db: db}),
		Subscriptions: NewSubscriptionRepository(&SubscriptionRepositoryConfig{db: db}),
		UserSubscriptions: NewUserSubscriptionRepository(&UserSubscriptionRepositoryConfig{
			db: db,
		}),
		UserTokens: NewUserTokenRepository(&UserTokenRepositoryConfig{
			db: db,
		}),
		Invoices: NewInvoiceRepository(&InvoiceRepositoryConfig{
			db: db,
		}),
		Gifts: NewGiftRepository(&GiftRepositoryConfig{
			db: db,
		}),
		UserGifts: NewUserGiftRepository(&UserGiftRepositoryConfig{
			db: db,
		}),
		Vouchers: NewVoucherRepository(&VoucherRepositoryConfig{
			db: db,
		}),
		UserVouchers: NewUserVoucherRepository(&UserVoucherRepositoryConfig{
			db: db,
		}),
		UserSpendings: NewUserSpendingRepository(&UserSpendingRepositoryConfig{
			db: db,
		}),
	}
}
