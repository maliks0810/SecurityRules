package handlers

import (
	"securityrules/security-rules/internal/app/models"
	"securityrules/security-rules/internal/utils/snowflake"

	"github.com/gofiber/fiber/v2"
)

func GetInformation(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("Welcome to Go microservices using Fiber")
}

func GetMukesh(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("Mukesh work faster")
}

func GetSecurityExceptions(ctx *fiber.Ctx) error {
	aladdinID := ctx.Query("aladdin_id")
	if aladdinID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "aladdin_id query parameter is required"})
	}

	rows, err := snowflake.DB.Query(
		"SELECT SECURITY_EXCEPTION_ID, RULE_ID, ALADDIN_ID, RUN_DATE, RUN_START_TIME, "+
			"RESULT_TYPE_ID, EXCEPTION_SOURCE_ID, EXCEPTION_STATUS_ID, SEVERITY_TYPE_ID, "+
			"PROCESS_TYPE_ID, CATEGORY_TYPE_ID, ASSIGN_TO, ASSIGN_TO_DATE, ASSIGNED_BY, "+
			"RESOLVE_DATE, ISSUE_DESCRIPTION, SOURCE_SYSTEM, CREATED_DATE, CREATED_BY, "+
			"MODIFIED_DATE, MODIFIED_BY FROM SECURITY_EXCEPTION WHERE ALADDIN_ID = ?", aladdinID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to query security exceptions"})
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
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to scan security exception"})
		}
		exceptions = append(exceptions, e)
	}

	if exceptions == nil {
		exceptions = []models.SecurityException{}
	}
	return ctx.Status(fiber.StatusOK).JSON(exceptions)
}
