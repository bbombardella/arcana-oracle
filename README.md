<div align="center">
  <img src="assets/banner.svg" alt="Arcana Oracle — API Backend" width="860" />
</div>

<br />

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.24-20232a?style=flat-square&logo=go&logoColor=00add8)](https://go.dev)
[![AWS Lambda](https://img.shields.io/badge/AWS_Lambda-arm64-20232a?style=flat-square&logo=awslambda&logoColor=ff9900)](https://aws.amazon.com/lambda)
[![Terraform](https://img.shields.io/badge/Terraform-1.7+-20232a?style=flat-square&logo=terraform&logoColor=7b42bc)](https://www.terraform.io)
[![Scaleway AI](https://img.shields.io/badge/Scaleway_AI-Mistral-20232a?style=flat-square&logo=scaleway&logoColor=4f0599)](https://www.scaleway.com/en/generative-apis)
[![License: MIT](https://img.shields.io/badge/License-MIT-20232a?style=flat-square)](LICENSE)

</div>

---

## ✦ À propos

**Arcana Oracle** est le backend API de l'application [Arcana](../arcana). Il expose un service HTTP de génération de lectures de tarot et d'astrologie, propulsé par **Mistral Small 3.2 24B** via l'API Scaleway — en gardant la clé secrète côté serveur.

Déployé en tant que fonction AWS Lambda (Graviton2, arm64), derrière API Gateway HTTP API V2.

---

## ☽ Endpoints

| Endpoint | Limite | Description |
|---|---|---|
| `POST /oracle/card` | 20 req/min | Interprétation intime d'une carte, endroit ou renversée |
| `POST /oracle/spread` | 5 req/min | Vision globale d'un tirage multi-cartes |
| `POST /oracle/astro` | 5 req/min | Lecture combinée signe astrologique × carte tirée |

Langues supportées : `fr` (défaut), `en`.

---

## ⬡ Stack technique

| Technologie | Rôle |
|---|---|
| **Go** | Runtime Lambda, logique métier, validation |
| **AWS Lambda** `provided.al2023` arm64 | Exécution serverless sur Graviton2 |
| **API Gateway HTTP API V2** | Routage, CORS, throttling par route |
| **DynamoDB** | Cache des interprétations de cartes individuelles |
| **Terraform** | Infrastructure as code (Lambda, API GW, IAM, DynamoDB) |
| **Scaleway AI** | Inférence Mistral Small 3.2 24B |

---

## ✧ Architecture

```
Browser (arcana SPA)
  └─► POST https://<api-id>.execute-api.<region>.amazonaws.com/oracle/<type>
        body: { ...payload }

AWS API Gateway HTTP API
  └─► Go Lambda
        ├─ Validation du payload et de la carte
        ├─ [/card] Lookup DynamoDB → retourne le cache si trouvé
        ├─ Appel Scaleway (Mistral, génération de texte)
        ├─ [/card] Écriture asynchrone en cache DynamoDB
        └─ Réponse HTTP texte vers le navigateur
```

**Cache DynamoDB** — uniquement pour `/oracle/card`. Clé composite : `card#<id>#<reversed>#<lang>`.

---

## ✦ Structure

```
src/
├── cmd/oracle/
│   └── main.go                 # Entrée Lambda — routing, wiring des dépendances
└── internal/
    ├── cards/
    │   └── cards.go            # Map statique des 78 cartes (id → nom français)
    ├── types/
    │   ├── card.go             # CardRequest, CardInfo, PositionInfo
    │   ├── spread.go           # SpreadRequest
    │   └── astro.go            # AstroRequest, SignInfo
    ├── prompts/
    │   └── prompts.go          # SystemPrompt + BuildCardPrompt/BuildSpreadPrompt/BuildAstroPrompt
    ├── handler/
    │   ├── card.go             # CardHandler — cache → Scaleway → async write
    │   ├── spread.go           # SpreadHandler — Scaleway direct
    │   └── astro.go            # AstroHandler — Scaleway direct
    ├── cache/
    │   └── dynamodb.go         # Get/Set DynamoDB
    └── scaleway/
        └── client.go           # Client HTTP Scaleway SSE
infra/
├── main.tf                     # Provider AWS + backend S3
├── lambda.tf                   # Lambda, API Gateway, routes, throttling
├── iam.tf                      # Rôle d'exécution Lambda + policy DynamoDB
├── dynamodb.tf                 # Table de cache
└── variables.tf
```

---

## ✧ Mise en route locale

### Prérequis

- **Go** ≥ 1.24
- **Docker** + **Docker Compose**
- **AWS SAM CLI**
- **AWS CLI**
- **Terraform** ≥ 1.7

### Lancement

```bash
# 1. Démarrer DynamoDB Local et créer la table de cache
make local-up

# 2. Configurer les variables locales
cp infra/local.tfvars.example infra/local.tfvars
# Remplir SCW_API_URL, SCW_SECRET_KEY, CLOUDFRONT_ORIGIN dans local.tfvars

# 3. Compiler et démarrer l'API en local (port 3000)
make local-api
```

SAM lit les fichiers Terraform pour découvrir le Lambda et les routes — aucun template SAM séparé n'est nécessaire.

### Commandes disponibles

```bash
make build          # Cross-compile linux/arm64 → bin/bootstrap
make zip            # build + zip → bin/bootstrap.zip
make test           # go test ./...
make lint           # golangci-lint
make clean          # Supprime bin/

make local-up       # Lance DynamoDB Local + crée la table
make local-down     # Arrête Docker Compose
make local-api      # Lance API Gateway + Lambda en local (port 3000)

make local-db-tables  # Liste les tables DynamoDB locales
make local-db-scan    # Affiche le contenu du cache
make local-db-clear   # Vide la table de cache
```

---

## ☽ Oracle IA

Les lectures sont générées par **Mistral Small 3.2 24B** (`mistral-small-3.2-24b-instruct-2506`) via l'API Scaleway.

```
temperature : 0.92
max_tokens  : 300
```

L'oracle adopte la voix d'une sorcière — envoûtante, poétique, adulte. Trois types de lectures :

- **Carte** — interprétation intime (3-4 phrases), endroit ou renversée, avec position nommée dans le tirage
- **Tirage** — vision d'ensemble tissée entre les cartes (5-6 phrases) — Passé / Présent / Futur, Situation / Obstacle / Fondation…
- **Astrologie** — fusion signe zodiacal × carte (4-5 phrases)

---

## ✦ Déploiement

Le déploiement est automatisé via GitHub Actions sur push sur `main` :

1. `make test` — tests unitaires
2. `make zip` — compilation et packaging
3. `terraform apply` — déploiement de l'infrastructure et mise à jour du Lambda

Les secrets suivants doivent être configurés dans le dépôt GitHub :

| Secret | Description |
|---|---|
| `AWS_ROLE_ARN` | ARN du rôle IAM assumé via OIDC |
| `SCW_API_URL` | Endpoint Scaleway AI |
| `SCW_SECRET_KEY` | Clé secrète Scaleway |
| `CLOUDFRONT_ORIGIN` | URL CloudFront du SPA Arcana |

---

<div align="center">
<br />

*✦ &nbsp; Que les arcanes t'éclairent &nbsp; ✦*

<br />
</div>
