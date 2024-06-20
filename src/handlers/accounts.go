package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/sql"
	"github.com/ochom/mpesa/src/controllers/c2b"
	"github.com/ochom/mpesa/src/domain"
	"github.com/ochom/mpesa/src/models"
	"gorm.io/gorm"
)

// HandleGetShortcodes ...
func HandleListShortCodes(ctx fiber.Ctx) error {
	accounts := sql.FindAll[models.Account]()
	// mapped := arrays.Map(accounts, func(account *models.Account) map[string]any {
	// 	return map[string]any{
	// 		"id":         account.ID,
	// 		"short_code": account.ShortCode,
	// 		"type":       account.Type,
	// 	}
	// })

	return ctx.JSON(accounts)
}

// HandleCreateShortCode ...
func HandleCreateShortCode(ctx fiber.Ctx) error {
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

	account := models.NewAccount(req.Type, req.ShortCode, req.PassKey, req.ConsumerKey, req.ConsumerSecrete)
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

// HandleUpdateShortCode ...
func HandleUpdateShortCode(ctx fiber.Ctx) error {
	account, err := sql.FindOneById[models.Account](ctx.Params("id"))
	if err != nil {
		return err
	}

	req, err := parseData[domain.CreateAccountRequest](ctx)
	if err != nil {
		return err
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
