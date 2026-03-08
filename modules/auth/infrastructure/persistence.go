package infrastructure

import (
	"context"
	"fmt"
	"sync"

	projectent "github.com/aggi-tech/aggipay/ent"
	permissionent "github.com/aggi-tech/aggipay/ent/permission"
	roleent "github.com/aggi-tech/aggipay/ent/role"
	"github.com/aggi-tech/aggipay/ent/user"
	"github.com/aggi-tech/aggipay/modules/auth/domain"
)

type EntRepository struct {
	client   *projectent.Client
	seedOnce sync.Once
	seedErr  error
}

func NewEntRepository(client *projectent.Client) *EntRepository {
	return &EntRepository{client: client}
}

func (r *EntRepository) Create(ctx context.Context, u domain.User) (*domain.User, error) {
	if err := r.ensureRBACSeeded(ctx); err != nil {
		return nil, err
	}

	roleIDs, err := r.resolveRoleIDs(ctx, u.Roles)
	if err != nil {
		return nil, err
	}

	builder := r.client.User.Create().
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		SetEmail(u.Email).
		SetBalance(0).
		SetActive(u.Active)

	if len(roleIDs) > 0 {
		builder.AddRoleIDs(roleIDs...)
	}
	if u.PhoneNumber != nil {
		builder.SetPhoneNumber(*u.PhoneNumber)
	}
	if u.PasswordHash != nil {
		builder.SetPasswordHash(*u.PasswordHash)
	}
	if u.OAuthProvider != nil {
		builder.SetOauthProvider(*u.OAuthProvider)
	}
	if u.OAuthSub != nil {
		builder.SetOauthSub(*u.OAuthSub)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, created.ID)
}

func (r *EntRepository) Update(ctx context.Context, id string, u domain.User) (*domain.User, error) {
	builder := r.client.User.UpdateOneID(id).
		SetFirstName(u.FirstName).
		SetLastName(u.LastName)

	if u.PhoneNumber != nil {
		builder.SetPhoneNumber(*u.PhoneNumber)
	}

	if _, err := builder.Save(ctx); err != nil {
		return nil, mapNotFound(err)
	}

	return r.FindByID(ctx, id)
}

func (r *EntRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.findOneWithAuthz(ctx, func(query *projectent.UserQuery) *projectent.UserQuery {
		return query.Where(user.IDEQ(id))
	})
}

func (r *EntRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findOneWithAuthz(ctx, func(query *projectent.UserQuery) *projectent.UserQuery {
		return query.Where(user.EmailEQ(email))
	})
}

func (r *EntRepository) FindByOAuthSub(ctx context.Context, provider, sub string) (*domain.User, error) {
	return r.findOneWithAuthz(ctx, func(query *projectent.UserQuery) *projectent.UserQuery {
		return query.Where(user.OauthProviderEQ(provider), user.OauthSubEQ(sub))
	})
}

func (r *EntRepository) FindOrCreateOAuthUser(ctx context.Context, u domain.User) (*domain.User, error) {
	if u.OAuthProvider == nil || u.OAuthSub == nil {
		return nil, domain.ErrUserNotFound
	}

	existing, err := r.FindByOAuthSub(ctx, *u.OAuthProvider, *u.OAuthSub)
	if err == nil {
		return existing, nil
	}

	byEmail, err := r.FindByEmail(ctx, u.Email)
	if err == nil {
		if _, err := r.client.User.UpdateOneID(byEmail.ID).
			SetOauthProvider(*u.OAuthProvider).
			SetOauthSub(*u.OAuthSub).
			Save(ctx); err != nil {
			return nil, err
		}
		return r.FindByID(ctx, byEmail.ID)
	}

	return r.Create(ctx, u)
}

func (r *EntRepository) List(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	if err := r.ensureRBACSeeded(ctx); err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	entities, err := r.queryUsersWithAuthz().Offset(offset).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.User, 0, len(entities))
	for _, entity := range entities {
		result = append(result, *toDomain(entity))
	}

	return result, nil
}

func (r *EntRepository) Count(ctx context.Context) (int, error) {
	return r.client.User.Query().Count(ctx)
}

func (r *EntRepository) SoftDelete(ctx context.Context, id string) error {
	err := r.client.User.UpdateOneID(id).SetActive(false).Exec(ctx)
	return mapNotFound(err)
}

func (r *EntRepository) findOneWithAuthz(
	ctx context.Context,
	apply func(query *projectent.UserQuery) *projectent.UserQuery,
) (*domain.User, error) {
	if err := r.ensureRBACSeeded(ctx); err != nil {
		return nil, err
	}

	entity, err := apply(r.queryUsersWithAuthz()).Only(ctx)
	if err != nil {
		return nil, mapNotFound(err)
	}

	return r.ensureUserRolesAndMap(ctx, entity)
}

func (r *EntRepository) queryUsersWithAuthz() *projectent.UserQuery {
	return r.client.User.Query().
		WithRoles(func(query *projectent.RoleQuery) {
			query.WithPermissions()
		})
}

func (r *EntRepository) ensureUserRolesAndMap(ctx context.Context, entity *projectent.User) (*domain.User, error) {
	if len(entity.Edges.Roles) > 0 {
		return toDomain(entity), nil
	}

	if err := r.assignRolesToUser(ctx, entity.ID, []string{domain.DefaultRole}); err != nil {
		return nil, err
	}

	reloaded, err := r.queryUsersWithAuthz().Where(user.IDEQ(entity.ID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toDomain(reloaded), nil
}

func (r *EntRepository) ensureRBACSeeded(ctx context.Context) error {
	r.seedOnce.Do(func() {
		r.seedErr = r.seedRBAC(ctx)
	})
	return r.seedErr
}

// seedRBAC aplica seed transacional para manter consistência ACID:
// ou todas as roles/permissões são gravadas, ou nenhuma alteração é persistida.
func (r *EntRepository) seedRBAC(ctx context.Context) (err error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	permissionByKey := make(map[string]*projectent.Permission, len(domain.PermissionDescriptions()))
	for key, description := range domain.PermissionDescriptions() {
		permissionEntity, findErr := tx.Permission.Query().Where(permissionent.KeyEQ(key)).Only(ctx)
		if findErr != nil {
			if !projectent.IsNotFound(findErr) {
				err = findErr
				return err
			}
			create := tx.Permission.Create().SetKey(key)
			if description != "" {
				create.SetDescription(description)
			}
			permissionEntity, err = create.Save(ctx)
			if err != nil {
				return err
			}
		} else if description != "" {
			permissionEntity, err = tx.Permission.UpdateOneID(permissionEntity.ID).SetDescription(description).Save(ctx)
			if err != nil {
				return err
			}
		}
		permissionByKey[key] = permissionEntity
	}

	for roleName, permissionKeys := range domain.RolePermissions() {
		roleEntity, findErr := tx.Role.Query().Where(roleent.NameEQ(roleName)).Only(ctx)
		if findErr != nil {
			if !projectent.IsNotFound(findErr) {
				err = findErr
				return err
			}
			roleEntity, err = tx.Role.Create().SetName(roleName).Save(ctx)
			if err != nil {
				return err
			}
		}

		permissionIDs := make([]string, 0, len(permissionKeys))
		for _, permissionKey := range permissionKeys {
			permissionEntity, ok := permissionByKey[permissionKey]
			if !ok {
				return fmt.Errorf("permission %q não encontrada durante seed", permissionKey)
			}
			permissionIDs = append(permissionIDs, permissionEntity.ID)
		}

		updater := tx.Role.UpdateOneID(roleEntity.ID).ClearPermissions()
		if len(permissionIDs) > 0 {
			updater.AddPermissionIDs(permissionIDs...)
		}
		if err := updater.Exec(ctx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *EntRepository) assignRolesToUser(ctx context.Context, userID string, roles []string) error {
	roleIDs, err := r.resolveRoleIDs(ctx, roles)
	if err != nil {
		return err
	}

	updater := r.client.User.UpdateOneID(userID).ClearRoles()
	if len(roleIDs) > 0 {
		updater.AddRoleIDs(roleIDs...)
	}
	return updater.Exec(ctx)
}

func (r *EntRepository) resolveRoleIDs(ctx context.Context, roles []string) ([]string, error) {
	normalizedRoles := domain.NormalizeRoles(roles)
	roleIDs := make([]string, 0, len(normalizedRoles))
	for _, roleName := range normalizedRoles {
		roleEntity, err := r.client.Role.Query().Where(roleent.NameEQ(roleName)).Only(ctx)
		if err != nil {
			if projectent.IsNotFound(err) {
				return nil, fmt.Errorf("role %q não existe", roleName)
			}
			return nil, err
		}
		roleIDs = append(roleIDs, roleEntity.ID)
	}
	return roleIDs, nil
}

func toDomain(entity *projectent.User) *domain.User {
	roles := make([]string, 0, len(entity.Edges.Roles))
	permissionsByKey := make(map[string]struct{})

	for _, roleEntity := range entity.Edges.Roles {
		roles = append(roles, roleEntity.Name)
		for _, permissionEntity := range roleEntity.Edges.Permissions {
			permissionsByKey[permissionEntity.Key] = struct{}{}
		}
	}

	normalizedRoles := domain.NormalizeRoles(roles)
	permissions := domain.NormalizePermissions(mapsKeys(permissionsByKey))
	permissions = domain.EnsurePermissionsForRoles(normalizedRoles, permissions)

	return &domain.User{
		ID:            entity.ID,
		FirstName:     entity.FirstName,
		LastName:      entity.LastName,
		Email:         entity.Email,
		Role:          domain.PrimaryRole(normalizedRoles),
		Roles:         normalizedRoles,
		Permissions:   permissions,
		PhoneNumber:   entity.PhoneNumber,
		Balance:       entity.Balance,
		Active:        entity.Active,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		PasswordHash:  entity.PasswordHash,
		OAuthProvider: entity.OauthProvider,
		OAuthSub:      entity.OauthSub,
	}
}

func mapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if projectent.IsNotFound(err) {
		return domain.ErrUserNotFound
	}
	return err
}

func mapsKeys[K comparable](in map[K]struct{}) []K {
	out := make([]K, 0, len(in))
	for key := range in {
		out = append(out, key)
	}
	return out
}
