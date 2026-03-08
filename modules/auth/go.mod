module github.com/aggi-tech/aggipay/modules/auth

go 1.25.0

require (
github.com/aggi-tech/aggipay/contracts v0.0.0
github.com/aggi-tech/aggipay/ent v0.0.0
github.com/aggi-tech/aggipay/platform v0.0.0
github.com/coreos/go-oidc/v3 v3.17.0
github.com/gin-gonic/gin v1.12.0
github.com/go-playground/validator/v10 v10.30.1
github.com/golang-jwt/jwt/v5 v5.3.1
golang.org/x/oauth2 v0.35.0
)

replace (
github.com/aggi-tech/aggipay/contracts => ../../contracts
github.com/aggi-tech/aggipay/ent => ../../ent
github.com/aggi-tech/aggipay/platform => ../../platform
)
