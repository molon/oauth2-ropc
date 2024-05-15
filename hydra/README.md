# hydra

```shell
git clone https://github.com/molon/hydra
cd ./hydra
go build -tags sqlite,json1,hsm .

export DSN="sqlite://$(pwd)/_db.sqlite?_fk=true"
./hydra migrate -c ./contrib/quickstart/5-min/hydra.yml sql -e --yes
./hydra serve -c ./contrib/quickstart/5-min/hydra.yml all --dev

##### client_credentials
client=$(./hydra create client \
    --endpoint http://127.0.0.1:4445/ \
    --format json \
    --grant-type client_credentials)
client_id=$(echo $client | jq -r '.client_id')
client_secret=$(echo $client | jq -r '.client_secret')
echo $client_id 
echo $client_secret

#### authorization_code
./hydra perform client-credentials \
  --endpoint http://127.0.0.1:4444/ \
  --client-id "$client_id" \
  --client-secret "$client_secret"

code_client=$(./hydra create client \
    --endpoint http://127.0.0.1:4445 \
    --grant-type authorization_code,refresh_token \
    --response-type code,id_token \
    --format json \
    --scope openid --scope offline \
    --redirect-uri http://127.0.0.1:5555/callback)

code_client_id=$(echo $code_client | jq -r '.client_id')
code_client_secret=$(echo $code_client | jq -r '.client_secret')
echo $code_client_id 
echo $code_client_secret

# 实际上这里会无法回调，因为文档里的是基于 docker-compose quickstart.yml 里的 oryd/hydra-login-consent-node 镜像来做的，而我们没启动这个
# 这里只是给一个例子
./hydra perform authorization-code \
    --client-id $code_client_id \
    --client-secret $code_client_secret \
    --endpoint http://127.0.0.1:4444/ \
    --port 5555 \
    --scope openid --scope offline

########## password
password_client=$(./hydra create client \
    --endpoint http://127.0.0.1:4445/ \
    --format json \
    --grant-type password,refresh_token \
    --scope openid --scope offline)
password_client_id=$(echo $password_client | jq -r '.client_id')
password_client_secret=$(echo $password_client | jq -r '.client_secret')
echo $password_client_id 
echo $password_client_secret

./hydra perform password \
    --client-id $password_client_id \
    --client-secret $password_client_secret \
    --endpoint http://127.0.0.1:4444/ \
    --scope openid --scope offline \
    --username fake-kratos-username \
    --password fake-kratos-password
```
