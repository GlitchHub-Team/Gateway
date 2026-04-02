# Gateway

MVP repository per lo sviluppo Gateway

## Dipendenze
- NATS


# Funzionalità
- Analisi statica: `golangci-lint run`
- Formattazione(automatica on save): gofumpt -w nomefile.go
- Delve Debugger: premere F5 per avviare il debugger e per mettere i breakpoint basta mettere i pallini rossi alla sinistra del codice

# Organizzazione cartelle
Le *cartelle* e di conseguenza i *package* sono organizzati per *Component* come suggerito nell'ultimo capitolo di "Clean Architecture" di Robert C. Martin.

# Leggere Code Coverage
Seguire i seguenti passaggi per leggere la code coverage:
1. `chmod +x testCoverage.sh`
2. `./testCoverage.sh`
3. Aprire il file `coverage.html` con un browser per visualizzare la code coverage in modo interattivo.

# Come creare Tenant - nsc edition
1. `nsc add account <tenant-name>`, convenzione "tenant-<tenant-id>"
2. `nsc add import --account <tenant-name> \
  --src-account <application-core-account-id> \
  --remote-subject "sensor.<tenant-id>.>" \
  --name "sensor_service_import" \
  --service
`
1. `nsc describe account <tenant-name>`, per controllare l'account creato
2. `nsc push -u nats://localhost:4222 --ca-cert ca.pem`, ovviamente la ca.pem deve essere nella cartella di esecuzione del comando                               

# Come creare Gateway - nsc edition
1. `nsc add user -a <tenant-name> -n <gateway-name> 
        --allow-pub "sensor.<tenant-id>.<gateway-id>.>"
        --allow-sub "\$JS.API.>,_INBOX.>"
        --public-key <gateway-public-key>`, 
    convenzione "gateway-<gateway-id>"
2. `nsc describe user -a <tenant-name> -n <gateway-name> -R`, per recuperare il JWT del gateway