# p3-ugc-7-8

## How to run

Prepare two terminal. One one terminal, go to folder `github.com/nafisalfiani/p3-ugc-7-8/account-service`:

```shell
cd github.com/nafisalfiani/p3-ugc-7-8/account-service
```

And execute:

```shell
make run
```

And on the other, go to folder `"github.com/nafisalfiani/p3-ugc-7-8/api-gateway`:

```shell
cd "github.com/nafisalfiani/p3-ugc-7-8/api-gateway
```

And execute:

```shell
make run
```

Make sure that you already have swaggo installed. If not, you can use:

```shell
make swag-install
```

. Now that both applications are running, go to your web browser and go to:

```shell
localhost:8080/swagger/index.html
```
