# wallet-service

make build

To run the containers:
make up

Steps to write the private key in the vault
docker exec -it sh-id-platform-test-vault sh

Then you will get a command shell, run:
vault write iden3/import/pbkey key_type=ethereum private_key=

(you need to get the private key from Metamask)

Write the vault token in the config.toml file, once the vault is initialized the token can be found in infrastructure/local/.vault/data/init.out or in the logs of the vault container.

Make sure that your database is properly configured (step 1) and run make db/migrate command. 

Run ./bin/platform command to start the issuer. Browse to http://localhost:3001 (or the port configured in ServerPort config entry) This will show you the api documentation.

To generate the api:
make api

