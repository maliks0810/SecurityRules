package handlers

import (
	"SecurityRules/database"
	"SecurityRules/internal/app/models"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SecurityRuleHandler struct{}

func NewSecurityRuleHandler() *SecurityRuleHandler {
	return &SecurityRuleHandler{}
}

func (h *SecurityRuleHandler) GetAll(c *fiber.Ctx) error {
	rows, err := database.Pool.Query(context.Background(),
		"SELECT id, name, description, severity, category, enabled, created_at, updated_at FROM security_rules ORDER BY id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to query rules"})
	}
	defer rows.Close()

	var rules []models.SecurityRule
	for rows.Next() {
		var r models.SecurityRule
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.Severity, &r.Category, &r.Enabled, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to scan rule"})
		}
		rules = append(rules, r)
	}

	if rules == nil {
		rules = []models.SecurityRule{}
	}
	return c.JSON(rules)
}

func (h *SecurityRuleHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var r models.SecurityRule
	err = database.Pool.QueryRow(context.Background(),
		"SELECT id, name, description, severity, category, enabled, created_at, updated_at FROM security_rules WHERE id = $1", id).
		Scan(&r.ID, &r.Name, &r.Description, &r.Severity, &r.Category, &r.Enabled, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "rule not found"})
	}

	return c.JSON(r)
}

func (h *SecurityRuleHandler) Create(c *fiber.Ctx) error {
	var req models.CreateRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Name == "" || req.Severity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name and severity are required"})
	}

	if req.Severity != "low" && req.Severity != "medium" && req.Severity != "high" && req.Severity != "critical" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "severity must be low, medium, high, or critical"})
	}

	var r models.SecurityRule
	err := database.Pool.QueryRow(context.Background(),
		`INSERT INTO security_rules (name, description, severity, category, enabled)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, description, severity, category, enabled, created_at, updated_at`,
		req.Name, req.Description, req.Severity, req.Category, req.Enabled).
		Scan(&r.ID, &r.Name, &r.Description, &r.Severity, &r.Category, &r.Enabled, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create rule"})
	}

	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h *SecurityRuleHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req models.UpdateRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Severity != nil {
		s := *req.Severity
		if s != "low" && s != "medium" && s != "high" && s != "critical" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "severity must be low, medium, high, or critical"})
		}
	}

	var r models.SecurityRule
	err = database.Pool.QueryRow(context.Background(),
		`UPDATE security_rules SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			severity = COALESCE($3, severity),
			category = COALESCE($4, category),
			enabled = COALESCE($5, enabled),
			updated_at = $6
		 WHERE id = $7
		 RETURNING id, name, description, severity, category, enabled, created_at, updated_at`,
		req.Name, req.Description, req.Severity, req.Category, req.Enabled, time.Now(), id).
		Scan(&r.ID, &r.Name, &r.Description, &r.Severity, &r.Category, &r.Enabled, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "rule not found"})
	}

	return c.JSON(r)
}

func (h *SecurityRuleHandler) GetSecurityException(c *fiber.Ctx) error {
	aladdinID := c.Query("AladdinId")
	if aladdinID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "AladdinId query parameter is required"})
	}

	rows, err := database.Pool.Query(context.Background(),
		"SELECT * FROM GET_SECURITY_EXCEPTION($1)", aladdinID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to query security exceptions"})
	}
	defer rows.Close()

	var exceptions []models.SecurityException
	for rows.Next() {
		var e models.SecurityException
		if err := rows.Scan(
			&e.SecurityExceptionID, &e.RuleID, &e.AladdinID,
			&e.RunDate, &e.RunStartTime, &e.ResultTypeID,
			&e.ExceptionSourceID, &e.ExceptionStatusID, &e.SeverityTypeID,
			&e.ProcessTypeID, &e.CategoryTypeID, &e.AssignTo,
			&e.AssignToDate, &e.AssignedBy, &e.ResolveDate,
			&e.IssueDescription, &e.SourceSystem,
			&e.CreatedDate, &e.CreatedBy, &e.ModifiedDate, &e.ModifiedBy,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to scan security exception"})
		}
		exceptions = append(exceptions, e)
	}

	if exceptions == nil {
		exceptions = []models.SecurityException{}
	}
	return c.JSON(exceptions)
}

func (h *SecurityRuleHandler) InsertSecurityException(c *fiber.Ctx) error {
	var req models.InsertSecurityExceptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	_, err := database.Pool.Exec(context.Background(),
		`CALL INSERT_SECURITY_EXCEPTION($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`,
		req.RuleID, req.AladdinID, req.RunDate, req.RunStartTime,
		req.ResultTypeID, req.ExceptionSourceID, req.ExceptionStatusID,
		req.SeverityTypeID, req.ProcessTypeID, req.CategoryTypeID,
		req.AssignTo, req.AssignToDate, req.AssignedBy, req.ResolveDate,
		req.IssueDescription, req.SourceSystem,
		req.CreatedDate, req.CreatedBy, req.ModifiedDate, req.ModifiedBy,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to insert security exception: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "security exception created successfully"})
}

func (h *SecurityRuleHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := database.Pool.Exec(context.Background(),
		"DELETE FROM security_rules WHERE id = $1", id)
	if err != nil || result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "rule not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
