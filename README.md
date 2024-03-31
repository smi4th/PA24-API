# PA24-API

## Description
\<Description needed>

## Gestion des flux

```mermaid

sequenceDiagram

    autonumber

    %% Declaration of elements

    actor client
    
    participant main
    %% participant gateway
    %% participant user
    %% participant authentification
    %% participant dataBase

    Note over main: manage input/output data

    %% Links between elements

    client->>main: HTTPS (request)

    main->>gateway: 

    alt login
        create participant user
        gateway->>user: credentialsValid()
        user->>+dataBase: SQL
        dataBase->>-user: response
        destroy user
        user->>gateway: response
    else any
        create participant authentification
        gateway->>authentification: isRequestValid()
        Note over gateway,authentification: It will test if the token is valid<br>but also if the client has the right to do this request
        authentification->>dataBase: SQL
        dataBase->>authentification: response
        destroy authentification
        authentification->>gateway: response
        create participant anyAPIRoute as any API route
        gateway-->anyAPIRoute: request
        anyAPIRoute->>dataBase: SQL
        dataBase->>anyAPIRoute: response
        destroy anyAPIRoute
        anyAPIRoute->>gateway: response
        gateway->>main: response
    end

    main->>client: HTTPS (response)

```

## Installation

To install the API on the disired server you need to run theses commands.

There is a `installation.sh` file that contains bash code to install and launch the API server with default configuration.

```bash
git clone https://github.com/smi4th/PA24-API.git
sudo bash PA24-API/installation.sh
```

## Configuration

You can find a `config.json` file. It contains some informations that you can modify to your liking. Like the databse location, the database credentials...

The database credentials **must** be stored in a secure way to ensure security breach.