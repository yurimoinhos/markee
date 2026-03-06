package infrastructure

import (
	"context"

	projectent "github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/ent/user"
	"github.com/aggi-tech/aggipay/internal/modules/auth/domain"
)

type EntRepository struct {
	client *projectent.Client
}

func NewEntRepository(client *projectent.Client) *EntRepository {
	return &EntRepository{client: client}
}

func (r *EntRepository) Create(ctx context.Context, u domain.User) (*domain.User, error) {
	builder := r.client.User.Create().
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		SetEmail(u.Email).
		SetBalance(0).
		SetActive(u.Active)

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

	return toDomain(created), nil
}

func (r *EntRepository) Update(ctx context.Context, id string, u domain.User) (*domain.User, error) {
	builder := r.client.User.UpdateOneID(id).
		SetFirstName(u.FirstName).
		SetLastName(u.LastName)

	if u.PhoneNumber != nil {
		builder.SetPhoneNumber(*u.PhoneNumber)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if projectent.IsNotFound(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toDomain(updated), nil
}

func (r *EntRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	entity, err := r.client.User.Get(ctx, id)
	if err != nil {
		if projectent.IsNotFound(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toDomain(entity), nil
}

func (r *EntRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	entity, err := r.client.User.Query().Where(user.EmailEQ(email)).First(ctx)
	if err != nil {
		if projectent.IsNotFound(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toDomain(entity), nil
}

func (r *EntRepository) FindByOAuthSub(ctx context.Context, provider, sub string) (*domain.User, error) {
	entity, err := r.client.User.Query().
		Where(user.OauthProviderEQ(provider), user.OauthSubEQ(sub)).
		First(ctx)
	if err != nil {
		if projectent.IsNotFound(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toDomain(entity), nil
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
		updated, err := r.client.User.UpdateOneID(byEmail.ID).
			SetOauthProvider(*u.OAuthProvider).
			SetOauthSub(*u.OAuthSub).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return toDomain(updated), nil
	}

	return r.Create(ctx, u)
}

func (r *EntRepository) List(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	offset := (page - 1) * pageSize
	entities, err := r.client.User.Query().Offset(offset).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.User, 0, len(entities))
	for _, entity := range entities {
		result = append(result, *toDomain(entity))
	}

	return result, nil
}

func (r *EntRepository) SoftDelete(ctx context.Context, id string) error {
	err := r.client.User.UpdateOneID(id).SetActive(false).Exec(ctx)
	if err != nil {
		if projectent.IsNotFound(err) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func toDomain(entity *projectent.User) *domain.User {
	return &domain.User{
		ID:            entity.ID,
		FirstName:     entity.FirstName,
		LastName:      entity.LastName,
		Email:         entity.Email,
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
