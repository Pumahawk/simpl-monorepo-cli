# Esperienza dello sviluppatore

Il programmatore può trovarsi in questi stati:

- Utente deve fare l'onboarding (prima installazione)
- Utente deve gestire l'ambiente locale.

# Prima installazione

L'utente deve ancora fare l'onboarding.

I primi step consistono:

- Installare i requisiti minimi.
  - Docker
- Scaricare la cli dedicata.
- Avviare la cli.

# Utente deve gestire l'ambiente locale

L'utente ha gia fatto l'onboarding.
Ha gia installato tutti gli strumenti.
Deve poter gestire e monitorare l'ambiente locale.

# Attività di preparazione ambiente locale

Assumento la presenza dei requisiti minimi deve poterci essere una procedura automatica
che prepara l'ambiente locale per permettere lo sviluppatore di lavorare.

La procedura in installazione consiste:

- Scaricare il repository simpl-monorepo.
- Avviare il processo di installazione programmi dipendenze (mise install).
- Creazione del cluster.
- Installazione dipendenze interne al cluster (helm install).
- Compilazione degli applicativi o download degli applicativi data una versione specifica.
- Avvio dei servizi in modalità locale.
- Lettura dei log.

Esecuzione dei test avviene tramite un processo separato.

# Workflow con la CLI

- Download della cli tramite link condiviso.
- Doppio click sull'eseguibile appena scaricato.
- La cli controlla se sono rispettati i requisiti (docker-desktop).
- La cli verifica se sono presenti le dipendenze software. Se mancano le installa (mise).
- La cli controlla se il progetto è presente, altrimenti lo scarica.
- La cli parte mostrando la schermata principale.

# Schermata principale della CLI

La schermata principale mostra le informazioni piu importanti.

- Stato del cluster
  - Stato minikube
  - Stato port-forward
- Stato servizi su kubernetes
  - Database
  - PostgREST
  - Ejbca
  - Keycloak
  - RedPanda
  - Stato dei microservizi frontend.
- Stato dei microservizi backend.

Mostra un menu selezionabile di funzionalità

## Funzionalità selezionabili dalla schermata principale

- Cluster: Operazioni che riguardano minikube
  - Start/Restart
  - Stop
  - Destroy
- Services: Operazioni che riguardano i servizi.
  - List: Apertura dell'elenco dei servizi (Controllo specifico per servizio)
  - Backend: Operazioni che possono essere eseguete su tutti i servizi backend. (Stop all, restart all, start all).
  - Frontend: Operazioni che possono essere eseguete su tutti i servizi frontend. (Stop all, restart all, start all).
  - Update all: Fa l'upgrade di tutti i servizi allineando all'ultima versine (develop o main).
- Logs: Accesso a funzionalità log.

# Funzionalità CLI

- Analisi prerequisiti - Presenza di docker.
- Download dipendenze e praparazione sistema come download mise, installazione dipendenze (kubectl, helm, minikube).
- Recupero stato servizi kubectl.
- Controllo stato servizi kubectl (start, stop).
- Orchestrazione processi per controllo backend.
- Aggiornamento versione microservizi simpl (backend e frontend).
- Controllo loggin (estrazione log filtrando per microservizio o processo.)
