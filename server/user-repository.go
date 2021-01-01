package main

import (
	"errors"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
}

type User struct {
	ID             string
	OrganizationID string
	Email          string
	AtlassianID    NullString
	HashedPassword NullString
	AuthProviderID NullString
	OrgAdmin       bool
	SuperAdmin     bool
}

var userRepository *UserRepository
var userRepositoryOnce sync.Once

func GetUserRepository() *UserRepository {
	userRepositoryOnce.Do(func() {
		userRepository = &UserRepository{}
		_, err := GetDatabase().DB().Exec("CREATE TABLE IF NOT EXISTS users (" +
			"id uuid DEFAULT uuid_generate_v4(), " +
			"organization_id uuid NOT NULL, " +
			"email VARCHAR NOT NULL, " +
			"org_admin boolean NOT NULL DEFAULT FALSE, " +
			"super_admin boolean NOT NULL DEFAULT FALSE, " +
			"PRIMARY KEY (id))")
		if err != nil {
			panic(err)
		}
		_, err = GetDatabase().DB().Exec("CREATE UNIQUE INDEX IF NOT EXISTS users_email ON users(email)")
		if err != nil {
			panic(err)
		}
	})
	return userRepository
}

func (r *UserRepository) RunSchemaUpgrade(curVersion, targetVersion int) {
	if curVersion < 1 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE users " +
			"ADD COLUMN password VARCHAR, " +
			"ADD COLUMN auth_provider_id uuid"); err != nil {
			panic(err)
		}
	}
	if curVersion < 2 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE users " +
			"ALTER COLUMN id SET DEFAULT uuid_generate_v4()"); err != nil {
			panic(err)
		}
	}
	if curVersion < 7 {
		if _, err := GetDatabase().DB().Exec("ALTER TABLE users " +
			"ADD COLUMN atlassian_id VARCHAR"); err != nil {
			panic(err)
		}
		if _, err := GetDatabase().DB().Exec("CREATE INDEX IF NOT EXISTS users_atlassian_id ON users(atlassian_id)"); err != nil {
			panic(err)
		}
	}
}

func (r *UserRepository) Create(e *User) error {
	var id string
	err := GetDatabase().DB().QueryRow("INSERT INTO users "+
		"(organization_id, email, org_admin, super_admin, password, auth_provider_id, atlassian_id) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) "+
		"RETURNING id",
		e.OrganizationID, strings.ToLower(e.Email), e.OrgAdmin, e.SuperAdmin, CheckNullString(e.HashedPassword), CheckNullString(e.AuthProviderID), CheckNullString(e.AtlassianID)).Scan(&id)
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (r *UserRepository) GetOne(id string) (*User, error) {
	e := &User{}
	err := GetDatabase().DB().QueryRow("SELECT id, organization_id, email, org_admin, super_admin, password, auth_provider_id, atlassian_id "+
		"FROM users "+
		"WHERE id = $1",
		id).Scan(&e.ID, &e.OrganizationID, &e.Email, &e.OrgAdmin, &e.SuperAdmin, &e.HashedPassword, &e.AuthProviderID, &e.AtlassianID)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	e := &User{}
	err := GetDatabase().DB().QueryRow("SELECT id, organization_id, email, org_admin, super_admin, password, auth_provider_id, atlassian_id "+
		"FROM users "+
		"WHERE LOWER(email) = $1",
		strings.ToLower(email)).Scan(&e.ID, &e.OrganizationID, &e.Email, &e.OrgAdmin, &e.SuperAdmin, &e.HashedPassword, &e.AuthProviderID, &e.AtlassianID)
	if err != nil {
		return nil, err
	}
	return e, nil
}
func (r *UserRepository) GetByAtlassianID(atlassianID string) (*User, error) {
	e := &User{}
	err := GetDatabase().DB().QueryRow("SELECT id, organization_id, email, org_admin, super_admin, password, auth_provider_id, atlassian_id "+
		"FROM users "+
		"WHERE LOWER(atlassian_id) = $1",
		strings.ToLower(atlassianID)).Scan(&e.ID, &e.OrganizationID, &e.Email, &e.OrgAdmin, &e.SuperAdmin, &e.HashedPassword, &e.AuthProviderID, &e.AtlassianID)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserRepository) GetByKeyword(organizationID string, keyword string) ([]*User, error) {
	var result []*User
	rows, err := GetDatabase().DB().Query("SELECT id, organization_id, email, org_admin, super_admin, password, auth_provider_id, atlassian_id "+
		"FROM users "+
		"WHERE organization_id = $1 AND LOWER(email) LIKE '%' || $2 || '%' "+
		"ORDER BY email", organizationID, strings.ToLower(keyword))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &User{}
		err = rows.Scan(&e.ID, &e.OrganizationID, &e.Email, &e.OrgAdmin, &e.SuperAdmin, &e.HashedPassword, &e.AuthProviderID, &e.AtlassianID)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *UserRepository) GetAll(organizationID string, maxResults int, offset int) ([]*User, error) {
	var result []*User
	rows, err := GetDatabase().DB().Query("SELECT id, organization_id, email, org_admin, super_admin, password, auth_provider_id, atlassian_id "+
		"FROM users "+
		"WHERE organization_id = $1 "+
		"ORDER BY email "+
		"LIMIT $2 OFFSET $3", organizationID, maxResults, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &User{}
		err = rows.Scan(&e.ID, &e.OrganizationID, &e.Email, &e.OrgAdmin, &e.SuperAdmin, &e.HashedPassword, &e.AuthProviderID, &e.AtlassianID)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *UserRepository) Update(e *User) error {
	_, err := GetDatabase().DB().Exec("UPDATE users SET "+
		"organization_id = $1, "+
		"email = $2, "+
		"org_admin = $3, "+
		"super_admin = $4, "+
		"password = $5, "+
		"auth_provider_id = $6, "+
		"atlassian_id = $7 "+
		"WHERE id = $8",
		e.OrganizationID, strings.ToLower(e.Email), e.OrgAdmin, e.SuperAdmin, CheckNullString(e.HashedPassword), CheckNullString(e.AuthProviderID), CheckNullString(e.AtlassianID), e.ID)
	return err
}

func (r *UserRepository) Delete(e *User) error {
	_, err := GetDatabase().DB().Exec("DELETE FROM users WHERE id = $1", e.ID)
	return err
}

func (r *UserRepository) DeleteAll(organizationID string) error {
	_, err := GetDatabase().DB().Exec("DELETE FROM users WHERE organization_id = $1", organizationID)
	return err
}

func (r *UserRepository) GetCount(organizationID string) (int, error) {
	var res int
	err := GetDatabase().DB().QueryRow("SELECT COUNT(id) "+
		"FROM users "+
		"WHERE organization_id = $1",
		organizationID).Scan(&res)
	return res, err
}

func (r *UserRepository) GetHashedPassword(password string) string {
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pwHash)
}

func (r *UserRepository) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (r *UserRepository) mergeUsers(source, target *User) error {
	if source.OrganizationID != target.OrganizationID {
		return errors.New("Organization ID of source and target users don't match")
	}
	if _, err := GetDatabase().DB().Exec("UPDATE bookings SET user_id = $2 WHERE user_id = $1", source.ID, target.ID); err != nil {
		return err
	}
	if target.AtlassianID == "" {
		target.AtlassianID = source.AtlassianID
	}
	target.OrgAdmin = target.OrgAdmin || source.OrgAdmin
	target.SuperAdmin = target.SuperAdmin || source.SuperAdmin
	if err := r.Delete(source); err != nil {
		return err
	}
	if err := r.Update(target); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) canCreateUser(org *Organization) bool {
	maxUsers, _ := GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	curUsers, _ := GetUserRepository().GetCount(org.ID)
	return curUsers < maxUsers
}
