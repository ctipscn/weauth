// Package admin 提供管理员后台应用服务。
// 管理员操作绕过普通用户的归属检查，直接管理所有用户和站点资源。
package admin

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/domain/site"
	"github.com/ctipscn/weauth/backend/internal/domain/user"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"gorm.io/gorm"
)

// AdminApp 管理员应用服务
type AdminApp struct {
	db       *gorm.DB
	userRepo user.Repository
	siteRepo site.Repository
	authRepo auth.Repository
	log      *infralogger.LogService
}

// NewAdminApp 创建管理员应用服务
func NewAdminApp(
	db *gorm.DB,
	userRepo user.Repository,
	siteRepo site.Repository,
	authRepo auth.Repository,
	log *infralogger.LogService,
) *AdminApp {
	return &AdminApp{
		db:       db,
		userRepo: userRepo,
		siteRepo: siteRepo,
		authRepo: authRepo,
		log:      log,
	}
}

// ListUsers 获取所有用户列表
func (app *AdminApp) ListUsers() ([]*user.User, error) {
	users, err := app.userRepo.FindAll()
	if err != nil {
		app.log.Errorf("[ADMIN] ListUsers error: %v", err)
		return nil, err
	}
	app.log.Infof("[ADMIN] Listed users (count=%d)", len(users))
	return users, nil
}

// ListSites 获取所有站点列表（包含已软删除的记录，管理员可见）
func (app *AdminApp) ListSites() ([]*site.Site, error) {
	sites, err := app.siteRepo.FindAllUnscoped()
	if err != nil {
		app.log.Errorf("[ADMIN] ListSites error: %v", err)
		return nil, err
	}
	app.log.Infof("[ADMIN] Listed sites (count=%d)", len(sites))
	return sites, nil
}

// GetSiteByID 根据 ID 获取站点（包含已软删除的记录，管理员用）
func (app *AdminApp) GetSiteByID(id uint) (*site.Site, error) {
	s, err := app.siteRepo.FindByIDUnscoped(id)
	if err != nil {
		if errors.Is(err, domain.ErrSiteNotFound) {
			app.log.Warnf("[ADMIN] GetSiteByID: site not found (siteID=%d)", id)
			return nil, ErrSiteNotFound
		}
		app.log.Errorf("[ADMIN] GetSiteByID error (siteID=%d): %v", id, err)
		return nil, err
	}
	return s, nil
}

// DeleteUser 删除用户。
// 当 cascadeSites 为 true 时，级联删除该用户下的所有站点及关联数据。
// 注：sites.user_id 存在 DB 级的外键约束 ON DELETE CASCADE，
// 因此硬删除用户会自动删除其名下站点。非级联模式使用软删除避免触发外键清理。
func (app *AdminApp) DeleteUser(id uint, cascadeSites bool) error {
	// 验证用户存在
	u, err := app.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			app.log.Warnf("[ADMIN] DeleteUser: user not found (userID=%d)", id)
			return ErrUserNotFound
		}
		app.log.Errorf("[ADMIN] DeleteUser error (userID=%d): %v", id, err)
		return err
	}

	return app.db.Transaction(func(tx *gorm.DB) error {
		// 删除认证凭证（两种模式均清理，防止残留登录能力）
		if err := tx.Where("user_id = ?", id).Delete(&models.Auth{}).Error; err != nil {
			app.log.Errorf("[ADMIN] DeleteUser: failed to delete auth (userID=%d): %v", id, err)
			return err
		}

		if cascadeSites {
			// 级联模式：硬删除用户前，先清理名下站点及关联数据（无外键约束，需业务层保证一致性）
			var userSites []struct{ ID uint }
			if err := tx.Model(&models.Site{}).Unscoped().Where("user_id = ?", id).Find(&userSites).Error; err != nil {
				app.log.Errorf("[ADMIN] DeleteUser: failed to query sites (userID=%d): %v", id, err)
				return err
			}
			for _, s := range userSites {
				if err := tx.Unscoped().Where("site_id = ?", s.ID).Delete(&models.Domain{}).Error; err != nil {
					app.log.Errorf("[ADMIN] DeleteUser: failed to delete domains (userID=%d, siteID=%d): %v", id, s.ID, err)
					return err
				}
				if err := tx.Unscoped().Where("site_id = ?", s.ID).Delete(&models.IpPolicy{}).Error; err != nil {
					app.log.Errorf("[ADMIN] DeleteUser: failed to delete IP policy (userID=%d, siteID=%d): %v", id, s.ID, err)
					return err
				}
				if err := tx.Unscoped().Where("site_id = ?", s.ID).Delete(&models.IpRecord{}).Error; err != nil {
					app.log.Errorf("[ADMIN] DeleteUser: failed to delete IP records (userID=%d, siteID=%d): %v", id, s.ID, err)
					return err
				}
				if err := tx.Unscoped().Where("site_id = ?", s.ID).Delete(&models.Challenge{}).Error; err != nil {
					app.log.Errorf("[ADMIN] DeleteUser: failed to delete challenges (userID=%d, siteID=%d): %v", id, s.ID, err)
					return err
				}
				if err := tx.Unscoped().Where("site_id = ?", s.ID).Delete(&models.VerificationLog{}).Error; err != nil {
					app.log.Errorf("[ADMIN] DeleteUser: failed to delete verification logs (userID=%d, siteID=%d): %v", id, s.ID, err)
					return err
				}
			}
			// 删除用户下的所有站点
			if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Site{}).Error; err != nil {
				app.log.Errorf("[ADMIN] DeleteUser: failed to delete sites (userID=%d): %v", id, err)
				return err
			}
			// 硬删除用户
			if err := tx.Unscoped().Delete(&models.User{}, id).Error; err != nil {
				app.log.Errorf("[ADMIN] DeleteUser: failed to hard-delete user (userID=%d): %v", id, err)
				return err
			}
		} else {
			// 非级联模式：软删除用户（设置 deleted_at），不触发 FK CASCADE
			result := tx.Delete(&models.User{}, id)
			if result.Error != nil {
				app.log.Errorf("[ADMIN] DeleteUser: failed to soft-delete user (userID=%d): %v", id, err)
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrUserNotFound
			}
		}

		app.log.Infof("[ADMIN] User deleted (userID=%d, username=%s, cascadeSites=%v)", id, u.Username, cascadeSites)
		return nil
	})
}

// DeleteSite 删除指定站点及所有关联数据（域名、IP 策略、IP 记录、挑战记录等）
// 直接硬删除，跳过预检查，支持已软删除的站点。
// 业务层主动清理关联数据，DB 层无外键约束。
func (app *AdminApp) DeleteSite(id uint) error {
	return app.db.Transaction(func(tx *gorm.DB) error {
		// 清理关联数据（无外键约束，需业务层保证一致性）
		if err := tx.Unscoped().Where("site_id = ?", id).Delete(&models.Domain{}).Error; err != nil {
			app.log.Errorf("[ADMIN] DeleteSite: failed to delete domains (siteID=%d): %v", id, err)
			return err
		}
		if err := tx.Unscoped().Where("site_id = ?", id).Delete(&models.IpPolicy{}).Error; err != nil {
			app.log.Errorf("[ADMIN] DeleteSite: failed to delete IP policy (siteID=%d): %v", id, err)
			return err
		}
		if err := tx.Unscoped().Where("site_id = ?", id).Delete(&models.IpRecord{}).Error; err != nil {
			app.log.Errorf("[ADMIN] DeleteSite: failed to delete IP records (siteID=%d): %v", id, err)
			return err
		}
		if err := tx.Unscoped().Where("site_id = ?", id).Delete(&models.Challenge{}).Error; err != nil {
			app.log.Errorf("[ADMIN] DeleteSite: failed to delete challenges (siteID=%d): %v", id, err)
			return err
		}
		if err := tx.Unscoped().Where("site_id = ?", id).Delete(&models.VerificationLog{}).Error; err != nil {
			app.log.Errorf("[ADMIN] DeleteSite: failed to delete verification logs (siteID=%d): %v", id, err)
			return err
		}

		// 删除站点本身
		result := tx.Unscoped().Delete(&models.Site{}, id)
		if result.Error != nil {
			app.log.Errorf("[ADMIN] DeleteSite: failed to delete site (siteID=%d): %v", id, result.Error)
			return result.Error
		}
		if result.RowsAffected == 0 {
			app.log.Warnf("[ADMIN] DeleteSite: site not found (siteID=%d)", id)
			return ErrSiteNotFound
		}

		app.log.Infof("[ADMIN] Site deleted (siteID=%d)", id)
		return nil
	})
}

// DeleteDomain 删除指定站点下的域名
func (app *AdminApp) DeleteDomain(siteID, domainID uint) error {
	// 验证站点存在
	if _, err := app.siteRepo.FindByID(siteID); err != nil {
		if errors.Is(err, domain.ErrSiteNotFound) {
			app.log.Warnf("[ADMIN] DeleteDomain: site not found (siteID=%d, domainID=%d)", siteID, domainID)
			return ErrSiteNotFound
		}
		app.log.Errorf("[ADMIN] DeleteDomain error (siteID=%d, domainID=%d): %v", siteID, domainID, err)
		return err
	}

	if err := app.siteRepo.DeleteDomain(siteID, domainID); err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			app.log.Warnf("[ADMIN] DeleteDomain: domain not found (siteID=%d, domainID=%d)", siteID, domainID)
			return ErrDomainNotFound
		}
		app.log.Errorf("[ADMIN] DeleteDomain error (siteID=%d, domainID=%d): %v", siteID, domainID, err)
		return err
	}

	app.log.Infof("[ADMIN] Domain deleted (siteID=%d, domainID=%d)", siteID, domainID)
	return nil
}