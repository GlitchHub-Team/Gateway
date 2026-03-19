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