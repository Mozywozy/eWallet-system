package transactions

import (
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransactionHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var request struct {
		Amount          float64       `json:"amount"`
		TransactionType TransactionType `json:"transaction_type"`
		Reference       string        `json:"reference"`
		Description     string        `json:"description"`
		AdditionalInfo  AdditionalInfo `json:"additional_info"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	err := h.service.InitiateTransaction(userID, request.Amount, request.TransactionType, request.Reference, request.Description, request.AdditionalInfo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Transaksi berhasil dibuat"})
}

func (h *TransactionHandler) UpdateTransactionHandler(c *fiber.Ctx) error {
	var request struct {
		Reference string          `json:"reference"`
		Status    TransactionStatus `json:"status"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	err := h.service.UpdateTransaction(request.Reference, request.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Status transaksi berhasil diperbarui"})
}

func (h *TransactionHandler) GetTransactionHandler(c *fiber.Ctx) error {
	reference := c.Params("reference")

	transaction, err := h.service.GetTransactionByReference(reference)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Transaksi tidak ditemukan"})
	}

	return c.JSON(transaction)
}
