package balance

import (
	"github.com/gofiber/fiber/v2"
)

type BalanceHandler struct {
	service BalanceService
}

func NewBalanceHandler(service BalanceService) *BalanceHandler {
	return &BalanceHandler{service: service}
}

func (h *BalanceHandler) GetBalanceHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) 
	balance, err := h.service.GetUserBalance(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"balance": balance})
}

func (h *BalanceHandler) TopUpBalanceHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var request struct {
		Amount               float64 `json:"amount"`
		WalletTransactionType string  `json:"wallet_transaction_type"`
		Reference            string  `json:"reference"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if request.WalletTransactionType != "CREDIT" && request.WalletTransactionType != "DEBIT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "wallet_transaction_type harus CREDIT atau DEBIT"})
	}

	err := h.service.ProcessBalanceTransaction(userID, request.Amount, request.WalletTransactionType, request.Reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Transaksi berhasil"})
}
