package auditloghandler

import (
	auditlogsvc "BackendCoursyclopedia/service/auditlogservice"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IAuditLogHandler interface {
	GetAuditLogs(c *fiber.Ctx) error
}

type AuditLogHandler struct {
	AuditLogService auditlogsvc.IAuditLogService // Use the interface type
}

func NewAuditLogHandler(auditlogservice auditlogsvc.IAuditLogService) IAuditLogHandler {
	return &AuditLogHandler{
		AuditLogService: auditlogservice,
	}
}

func (h *AuditLogHandler) withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (h *AuditLogHandler) GetAuditLogs(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	auditLogs, err := h.AuditLogService.GetAllAuditLogs(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Audit logs retrieved successfully",
		"data":    auditLogs,
	})
}
