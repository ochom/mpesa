package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/arrays"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/controllers/c2b"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
	"gorm.io/gorm"
)

// HandleGetShortcodes ...
func HandleListAccounts(ctx fiber.Ctx) error {
	accounts := sql.FindAll[models.Account]()
	mapped := arrays.Map(accounts, func(account *models.Account) map[string]any {
		return map[string]any{
			"id":           account.ID,
			"short_code":   account.ShortCode,
			"name":         account.Name,
			"type":         account.Type,
			"created_at":   account.CreatedAt,
			"updated_at":   account.UpdatedAt,
			"consumer_key": account.ConsumerKey,
			"pass_key":     account.PassKey,
		}
	})

	return ctx.JSON(mapped)
}

// HandleSearchAccounts ...
func HandleSearchAccounts(ctx fiber.Ctx) error {
	id, shortCode, _type := ctx.Query("id"), ctx.Query("short_code"), ctx.Query("type")
	accounts := sql.FindAll[models.Account](func(d *gorm.DB) *gorm.DB {
		return d.Where("id = ? OR short_code = ? OR type = ?", id, shortCode, _type)
	})

	mapped := arrays.Map(accounts, func(account *models.Account) map[string]any {
		return map[string]any{
			"id":           account.ID,
			"short_code":   account.ShortCode,
			"name":         account.Name,
			"type":         account.Type,
			"created_at":   account.CreatedAt,
			"updated_at":   account.UpdatedAt,
			"consumer_key": account.ConsumerKey,
			"pass_key":     account.PassKey,
		}
	})

	return ctx.JSON(mapped)
}

// HandleCreateAccount ...
func HandleCreateAccount(ctx fiber.Ctx) error {
	req, err := parseDataValidate[domain.CreateAccountRequest](ctx)
	if err != nil {
		return err
	}

	count := sql.Count[models.Account](func(d *gorm.DB) *gorm.DB {
		return d.Where("short_code = ? AND type = ?", req.ShortCode, req.Type)
	})

	if count > 0 {
		return ctx.JSON(fiber.Map{"message": "account already exists"})
	}

	account := models.NewAccount(req.Type, req.ShortCode, req.Name, req.PassKey, req.ConsumerKey, req.ConsumerSecrete)
	account.ValidationUrl = req.ValidationUrl
	account.ConfirmationUrl = req.ConfirmationUrl
	account.InitiatorName = req.InitiatorName
	account.InitiatorPassword = req.InitiatorPassword
	account.Certificate = req.Certificate

	if err := sql.Create(account); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleUpdateAccount ...
func HandleUpdateAccount(ctx fiber.Ctx) error {
	account, err := sql.FindOneById[models.Account](ctx.Params("id"))
	if err != nil {
		return err
	}

	req, err := parseData[domain.CreateAccountRequest](ctx)
	if err != nil {
		return err
	}

	if req.Name != "" {
		account.Name = req.Name
	}

	if req.PassKey != "" {
		account.PassKey = req.PassKey
	}

	if req.ConsumerKey != "" {
		account.ConsumerKey = req.ConsumerKey
	}

	if req.ConsumerSecrete != "" {
		account.ConsumerSecrete = req.ConsumerSecrete
	}

	if req.ValidationUrl != "" {
		account.ValidationUrl = req.ValidationUrl
	}

	if req.ConfirmationUrl != "" {
		account.ConfirmationUrl = req.ConfirmationUrl
	}

	if req.InitiatorName != "" {
		account.InitiatorName = req.InitiatorName
	}

	if req.InitiatorPassword != "" {
		account.InitiatorPassword = req.InitiatorPassword
	}

	if req.Certificate != "" {
		account.Certificate = req.Certificate
	}

	if err := sql.Update(account); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"message": "success"})
}

// HandleC2BRegisterUrls ...
func HandleC2BRegisterUrls(ctx fiber.Ctx) error {
	req, err := parseData[map[string]string](ctx)
	if err != nil {
		return err
	}

	go c2b.RegisterUrls(req)
	return ctx.JSON(fiber.Map{"message": "success"})
}
