package postgres

import (
	"context"
	"errors"
	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetAll(ctx context.Context) (models.Accounts, error) {
	var accounts models.Accounts
	if err := r.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Account, error) {
	var account models.Account

	if err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Account{}, errors.New("account not found")
		}
		return models.Account{}, err
	}
	return account, nil
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *AccountRepositoryPostgres) Update(ctx context.Context, id uuid.UUID, updates models.AccountUpdate) error {
	updateData := map[string]interface{}{"updatedAt": time.Now()}

	if updates.Balance != nil {
		updateData["balance"] = *updates.Balance
	}

	result := r.db.WithContext(ctx).Model(&models.Account{}).Where("id = ?", id).Updates(updateData)
	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return result.Error
}
func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Account{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return result.Error
}
