package main

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OrganizationRepository struct {
}

type Organization struct {
	ID               string
	Name             string
	ContactFirstname string
	ContactLastname  string
	ContactEmail     string
	Country          string
	Language         string
	SignupDate       time.Time
}

type Domain struct {
	DomainName     string
	OrganizationID string
	Active         bool
	VerifyToken    string
}

var organizationRepository *OrganizationRepository
var organizationRepositoryOnce sync.Once

func GetOrganizationRepository() *OrganizationRepository {
	organizationRepositoryOnce.Do(func() {
		organizationRepository = &OrganizationRepository{}
		_, err := GetDatabase().DB().Exec("CREATE TABLE IF NOT EXISTS organizations (" +
			"id uuid DEFAULT uuid_generate_v4(), " +
			"name VARCHAR NOT NULL, " +
			"PRIMARY KEY (id))")
		if err != nil {
			panic(err)
		}
		_, err = GetDatabase().DB().Exec("CREATE TABLE IF NOT EXISTS organizations_domains (" +
			"domain VARCHAR NOT NULL, " +
			"organization_id uuid NOT NULL, " +
			"PRIMARY KEY (domain))")
		if err != nil {
			panic(err)
		}
		_, err = GetDatabase().DB().Exec("CREATE INDEX IF NOT EXISTS idx_organizations_domains_organization_id ON organizations_domains(organization_id)")
		if err != nil {
			panic(err)
		}
	})
	return organizationRepository
}

func (r *OrganizationRepository) RunSchemaUpgrade(curVersion, targetVersion int) {
	if curVersion < 3 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE organizations " +
			"ADD COLUMN contact_firstname VARCHAR, " +
			"ADD COLUMN contact_lastname VARCHAR, " +
			"ADD COLUMN contact_email VARCHAR"); err != nil {
			panic(err)
		}
	}
	if curVersion < 4 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE organizations_domains " +
			"ADD COLUMN active boolean NOT NULL DEFAULT FALSE, " +
			"ADD COLUMN verify_token uuid"); err != nil {
			panic(err)
		}
	}
	if curVersion < 5 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE organizations_domains " +
			"DROP CONSTRAINT organizations_domains_pkey"); err != nil {
			panic(err)
		}
		if _, err := GetDatabase().DB().Exec("CREATE INDEX IF NOT EXISTS idx_organizations_domains_domain ON organizations_domains(domain)"); err != nil {
			panic(err)
		}
	}
	if curVersion < 6 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE organizations " +
			"ADD COLUMN country VARCHAR, " +
			"ADD COLUMN language VARCHAR"); err != nil {
			panic(err)
		}
		if _, err := GetDatabase().DB().Exec("UPDATE organizations SET country = 'DE', language = 'de'"); err != nil {
			panic(err)
		}
	}
	if curVersion < 8 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE organizations " +
			"ADD COLUMN signup_date TIMESTAMP NOT NULL DEFAULT '2021-03-28 16:00:00'"); err != nil {
			panic(err)
		}
		if _, err := GetDatabase().DB().Exec("ALTER TABLE organizations " +
			"ALTER COLUMN signup_date DROP DEFAULT"); err != nil {
			panic(err)
		}
	}
}

func (r *OrganizationRepository) Create(e *Organization) error {
	var id string
	err := GetDatabase().DB().QueryRow("INSERT INTO organizations "+
		"(name, contact_firstname, contact_lastname, contact_email, country, language, signup_date) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) "+
		"RETURNING id",
		e.Name, e.ContactFirstname, e.ContactLastname, e.ContactEmail, e.Country, e.Language, e.SignupDate).Scan(&id)
	if err != nil {
		return err
	}
	e.ID = id
	GetSettingsRepository().InitDefaultSettingsForOrg(e.ID)
	return nil
}

func (r *OrganizationRepository) GetOneByDomain(domain string) (*Organization, error) {
	e := &Organization{}
	err := GetDatabase().DB().QueryRow("SELECT organizations.id, organizations.name, organizations.contact_firstname, organizations.contact_lastname, organizations.contact_email, organizations.country, organizations.language, organizations.signup_date "+
		"FROM organizations_domains "+
		"INNER JOIN organizations ON organizations.id = organizations_domains.organization_id "+
		"WHERE LOWER(organizations_domains.domain) = $1 AND organizations_domains.active = TRUE",
		strings.ToLower(domain)).Scan(&e.ID, &e.Name, &e.ContactFirstname, &e.ContactLastname, &e.ContactEmail, &e.Country, &e.Language, &e.SignupDate)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *OrganizationRepository) GetOne(id string) (*Organization, error) {
	e := &Organization{}
	err := GetDatabase().DB().QueryRow("SELECT id, name, contact_firstname, contact_lastname, contact_email, country, language, signup_date "+
		"FROM organizations "+
		"WHERE id = $1",
		id).Scan(&e.ID, &e.Name, &e.ContactFirstname, &e.ContactLastname, &e.ContactEmail, &e.Country, &e.Language, &e.SignupDate)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *OrganizationRepository) GetByEmail(email string) (*Organization, error) {
	e := &Organization{}
	err := GetDatabase().DB().QueryRow("SELECT id, name, contact_firstname, contact_lastname, contact_email, country, language, signup_date "+
		"FROM organizations "+
		"WHERE LOWER(contact_email) = $1",
		strings.ToLower(email)).Scan(&e.ID, &e.Name, &e.ContactFirstname, &e.ContactLastname, &e.ContactEmail, &e.Country, &e.Language, &e.SignupDate)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *OrganizationRepository) GetAll() ([]*Organization, error) {
	var result []*Organization
	rows, err := GetDatabase().DB().Query("SELECT id, name, contact_firstname, contact_lastname, contact_email, country, language, signup_date " +
		"FROM organizations ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &Organization{}
		err = rows.Scan(&e.ID, &e.Name, &e.ContactFirstname, &e.ContactLastname, &e.ContactEmail, &e.Country, &e.Language, &e.SignupDate)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *OrganizationRepository) GetAllIDs() ([]string, error) {
	var result []string
	rows, err := GetDatabase().DB().Query("SELECT id " +
		"FROM organizations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ID string
		err = rows.Scan(&ID)
		if err != nil {
			return nil, err
		}
		result = append(result, ID)
	}
	return result, nil
}

func (r *OrganizationRepository) Update(e *Organization) error {
	_, err := GetDatabase().DB().Exec("UPDATE organizations SET "+
		"name = $1, contact_firstname = $2, contact_lastname = $3, contact_email = $4, country = $5, language = $6, signup_date = $7 "+
		"WHERE id = $8",
		e.Name, e.ContactFirstname, e.ContactLastname, e.ContactEmail, e.Country, e.Language, e.SignupDate, e.ID)
	return err
}

func (r *OrganizationRepository) Delete(e *Organization) error {
	if err := GetAuthProviderRepository().DeleteAll(e.ID); err != nil {
		return err
	}
	if err := GetLocationRepository().DeleteAll(e.ID); err != nil {
		return err
	}
	if err := GetSettingsRepository().DeleteAll(e.ID); err != nil {
		return err
	}
	if err := GetUserRepository().DeleteAll(e.ID); err != nil {
		return err
	}
	_, err := GetDatabase().DB().Exec("DELETE FROM organizations_domains WHERE organization_id = $1", e.ID)
	if err != nil {
		return err
	}
	_, err = GetDatabase().DB().Exec("DELETE FROM organizations WHERE id = $1", e.ID)
	return err
}

func (r *OrganizationRepository) GetDomain(org *Organization, domain string) (*Domain, error) {
	e := &Domain{}
	err := GetDatabase().DB().QueryRow("SELECT domain, organization_id, active, verify_token "+
		"FROM organizations_domains "+
		"WHERE domain = LOWER($1) AND organization_id = $2",
		strings.ToLower(domain), org.ID).Scan(&e.DomainName, &e.OrganizationID, &e.Active, &e.VerifyToken)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *OrganizationRepository) AddDomain(e *Organization, domain string, active bool) error {
	verifyToken := uuid.New().String()
	_, err := GetDatabase().DB().Exec("INSERT INTO organizations_domains "+
		"(domain, organization_id, active, verify_token) "+
		"VALUES ($1, $2, $3, $4)",
		strings.ToLower(domain), e.ID, active, verifyToken)
	return err
}

func (r *OrganizationRepository) RemoveDomain(e *Organization, domain string) error {
	_, err := GetDatabase().DB().Exec("DELETE FROM organizations_domains "+
		"WHERE domain = LOWER($1) AND organization_id = $2",
		strings.ToLower(domain), e.ID)
	return err
}

func (r *OrganizationRepository) ActivateDomain(e *Organization, domain string) error {
	_, err := GetDatabase().DB().Exec("UPDATE organizations_domains "+
		"SET active = TRUE "+
		"WHERE domain = LOWER($1) AND organization_id = $2",
		strings.ToLower(domain), e.ID)
	return err
}

func (r *OrganizationRepository) GetDomains(e *Organization) ([]*Domain, error) {
	var result []*Domain
	rows, err := GetDatabase().DB().Query("SELECT domain, organization_id, active, verify_token "+
		"FROM organizations_domains "+
		"WHERE organization_id = $1"+
		"ORDER BY domain",
		e.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		domain := &Domain{}
		err = rows.Scan(&domain.DomainName, &domain.OrganizationID, &domain.Active, &domain.VerifyToken)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *OrganizationRepository) isValidEmailForOrg(email string, org *Organization) bool {
	mailParts := strings.Split(email, "@")
	if len(mailParts) != 2 {
		return false
	}
	domain := strings.ToLower(mailParts[1])
	domains, err := GetOrganizationRepository().GetDomains(org)
	if err != nil {
		return false
	}
	for _, curDomain := range domains {
		if strings.ToLower(curDomain.DomainName) == domain {
			if curDomain.Active {
				return true
			}
		}
	}
	return false

}
