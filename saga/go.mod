module github.com/aggi-tech/aggipay/saga

go 1.25.0

require (
	github.com/aggi-tech/aggipay/contracts v0.0.0
	github.com/aggi-tech/aggipay/ent v0.0.0
	github.com/aggi-tech/aggipay/platform v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.11.2
)

replace (
	github.com/aggi-tech/aggipay/contracts => ../contracts
	github.com/aggi-tech/aggipay/ent => ../ent
	github.com/aggi-tech/aggipay/platform => ../platform
)
