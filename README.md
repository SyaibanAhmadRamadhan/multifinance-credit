# MULTI FINANCE CREDIT

## infra
1. observability: using otel for standard observability tools
2. using zipkin for showing data tracing from otel collector
3. s3: using minio cluster with load balancer nginx for performance upload and get image and photo

## link
1. [zipkin](http://localhost:9411)
2. [minio console admin](http://localhost:9001)
    - password: minioadmin
    - username: minioadmin
3. [app api](http://localhost:3002)
4. basic auth
    - username: admin
    - password: admin
5. postman collection in folder `docs/multifinance-credit.postman_collection.json`


## how to run
1. first step, make sure you install `docker`, `docker compose`, `go`, and `nodejs`.
2. step 2, you can visit the [minio console admin](http://localhost:9001) and go to the access key tab. and you can generate access keys and secret keys. then put the access key and secret key in `env.json` in the minio object section
3. you can run `go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2` or `make install`
4. and after that you run `make run-docker`

## api spec gui
1. makesure you can install redocly preview.
2. run `make preview_open_api`
3. and follow [preview open api](http://localhost:8080)

## docs 
1. diagram database
![diagram database](docs/multi-finance-erd%20(1).png)