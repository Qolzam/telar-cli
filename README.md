# Telar Social on-click setup

### Setup Options

- [OpenFaaS Cloud Community Cluster]
  + Ask for github username.
  + Check some ingredients
    - Git already is installed.
    - Check openfaas-cli.
    - Check whether user name exist in openfaas cloud CUSTOMERS file.
    - Check each repository is forked (https://github.com/red-gold/telar-web, https://github.com/red-gold/ts-ui, https://github.com/red-gold/ts-serverless)
    - Already added the OpenFaaS GitHub App (https://github.com/apps/openfaas-cloud-community-cluster) and added the 3 repositories.
    - Ask if you applied the github app of openfaas.
    - You alread installed github SSH.
  * Create an internal trust secret (PAYLOAD_SECRET=$(head -c 12 /dev/urandom | shasum| cut -d’ ‘ -f1))
  + Check Firebase Storage (Service Account file and bucket name)
  + Create a cluster and database on Mongodb cloud (your_mongodb_password, mongo_user, mongo_database )
  + Enable Google reCAPTCHA (recaptcha_site_key, your_recaptcha_key)
  + Enable Github OAuth app (your_github_client_secret)
  * Generate a key/pair
    - openssl ecparam -genkey -name prime256v1 -noout -out key
    - openssl ec -in key -pubout -out key.pub
  + User Manegement
    - Enable email verification (email, password)
    - Enable Admin Account (username,password)
  + Check websocket connection.
  * git add . && git commit -sm ‘Deploy Telar Social.’ && git push
- AWS EKS
- Google Kubernetes Engine
- Azure Kubernetes Service
- Bare-metal Kubernetes
- Self-hosted OpenFaaS Cloud
- Bare-metal K3S
- Local Kind
- Local Minikube

## Author

- [Amirhossein Movahedi](https://amir.red-gold.tech)
## License

MIT