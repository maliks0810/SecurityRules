package models

import "time"

type SecurityException struct {
	SecurityExceptionID int        `json:"security_exception_id"`
	RuleID              *int       `json:"rule_id"`
	AladdinID           *string    `json:"aladdin_id"`
	RunDate             *time.Time `json:"run_date"`
	RunStart            *time.Time `json:"run_start"`
	ResultTypeID        *int       `json:"result_type_id"`
	ExceptionSourceID   *int       `json:"exception_source_id"`
	ExceptionStatusID   *int       `json:"exception_status_id"`
	SeverityTypeID      *int       `json:"severity_type_id"`
	ProcessTypeID       *int       `json:"process_type_id"`
	CategoryTypeID      *int       `json:"category_type_id"`
	AssignTo            *string    `json:"assign_to"`
	AssignToDate        *int       `json:"assign_to_date"`
	AssignedBy          *string    `json:"assigned_by"`
	ResolveDate         *time.Time `json:"resolve_date"`
	IssueDescription    *string    `json:"issue_description"`
	CreatedDate         *time.Time `json:"created_date"`
	CreatedBy           *string    `json:"created_by"`
	ModifiedDate        *time.Time `json:"modified_date"`
	ModifiedBy          *string    `json:"modified_by"`
}
