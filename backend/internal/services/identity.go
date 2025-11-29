package services

import (
	"database/sql"
	"fmt"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
)

type IdentityService struct {
	db *sqlx.DB
}

func NewIdentityService(db *sqlx.DB) *IdentityService {
	return &IdentityService{db: db}
}

// FindOrCreateCustomer finds an existing customer or creates a new one
// based on provided identifiers
func (s *IdentityService) FindOrCreateCustomer(tenantID int64, identifiers []models.CustomerIdentifier) (*models.Customer, error) {
	// Try to find existing customer by identifiers
	customer, err := s.FindCustomerByIdentifiers(tenantID, identifiers)
	if err == nil && customer != nil {
		return customer, nil
	}

	// Create new customer
	customer = &models.Customer{
		TenantID: tenantID,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert customer
	err = tx.QueryRowx(
		`INSERT INTO customers (tenant_id) VALUES ($1) RETURNING id, created_at, updated_at`,
		tenantID,
	).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Insert identifiers
	for i := range identifiers {
		identifiers[i].CustomerID = customer.ID
		_, err = tx.Exec(
			`INSERT INTO customer_identifiers (customer_id, type, value, source_system, is_primary)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (customer_id, type, value) DO NOTHING`,
			identifiers[i].CustomerID,
			identifiers[i].Type,
			identifiers[i].Value,
			identifiers[i].SourceSystem,
			identifiers[i].IsPrimary,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert identifier: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return customer, nil
}

// FindCustomerByIdentifiers finds a customer by matching identifiers
func (s *IdentityService) FindCustomerByIdentifiers(tenantID int64, identifiers []models.CustomerIdentifier) (*models.Customer, error) {
	if len(identifiers) == 0 {
		return nil, fmt.Errorf("no identifiers provided")
	}

	// Build query to find customer by any identifier
	query := `
		SELECT DISTINCT c.id, c.tenant_id, c.created_at, c.updated_at
		FROM customers c
		INNER JOIN customer_identifiers ci ON c.id = ci.customer_id
		WHERE c.tenant_id = $1 AND (
	`

	args := []interface{}{tenantID}
	argPos := 2

	for i, ident := range identifiers {
		if i > 0 {
			query += " OR "
		}
		query += fmt.Sprintf("(ci.type = $%d AND ci.value = $%d)", argPos, argPos+1)
		args = append(args, ident.Type, ident.Value)
		argPos += 2
	}
	query += ") LIMIT 1"

	var customer models.Customer
	err := s.db.Get(&customer, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	return &customer, nil
}

// GetCustomerJourney returns the complete journey for a customer
func (s *IdentityService) GetCustomerJourney(customerID int64, from, to *string) (*models.CustomerJourney, error) {
	journey := &models.CustomerJourney{
		CustomerID: customerID,
		Identifiers: []models.CustomerIdentifier{},
		Interactions: []models.InteractionWithChannel{},
		ConversionEvents: []models.ConversionEventWithSource{},
	}

	// Get customer identifiers
	err := s.db.Select(&journey.Identifiers,
		`SELECT id, customer_id, type, value, source_system, is_primary, created_at
		 FROM customer_identifiers
		 WHERE customer_id = $1
		 ORDER BY is_primary DESC, created_at ASC`,
		customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get identifiers: %w", err)
	}

	// Build interactions query
	interactionsQuery := `
		SELECT i.*, c.name as channel_name
		FROM interactions i
		INNER JOIN channels c ON i.channel_id = c.id
		WHERE i.customer_id = $1
	`
	args := []interface{}{customerID}
	argPos := 2

	if from != nil {
		interactionsQuery += fmt.Sprintf(" AND i.started_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}
	if to != nil {
		interactionsQuery += fmt.Sprintf(" AND i.started_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}
	interactionsQuery += " ORDER BY i.started_at ASC"

	err = s.db.Select(&journey.Interactions, interactionsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get interactions: %w", err)
	}

	// Get participants for each interaction
	for i := range journey.Interactions {
		var participants []models.InteractionParticipant
		err = s.db.Select(&participants,
			`SELECT ip.*
			 FROM interaction_participants ip
			 WHERE ip.interaction_id = $1`,
			journey.Interactions[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get participants: %w", err)
		}
		journey.Interactions[i].Participants = participants
	}

	// Build conversion events query
	conversionsQuery := `
		SELECT ce.*, es.name as event_source_name, cur.code as currency_code
		FROM conversion_events ce
		INNER JOIN event_sources es ON ce.event_source_id = es.id
		INNER JOIN currencies cur ON ce.currency_id = cur.id
		WHERE ce.customer_id = $1
	`
	args = []interface{}{customerID}
	argPos = 2

	if from != nil {
		conversionsQuery += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}
	if to != nil {
		conversionsQuery += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}
	conversionsQuery += " ORDER BY ce.occurred_at ASC"

	err = s.db.Select(&journey.ConversionEvents, conversionsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversion events: %w", err)
	}

	return journey, nil
}


