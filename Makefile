.PHONY: build zip clean test lint \
        local-up local-down local-init-db local-api \
        local-db-tables local-db-scan local-db-clear

BIN_DIR          = bin
BINARY           = $(BIN_DIR)/bootstrap
GOOS             = linux
GOARCH           = arm64
LOCAL_DOCKER_NET = oracle-local
LOCAL_DDB_URL    = http://localhost:8000
LOCAL_REGION     = eu-west-3
TABLE_NAME       = arcana-oracle-card-cache

# ── Build ────────────────────────────────────────────────────────────────────

build:
	@mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go -C src build -o ../$(BINARY) ./cmd/oracle/

# SAM local needs a zip (Terraform config references bin/bootstrap.zip)
zip: build
	zip -j $(BINARY).zip $(BINARY)

clean:
	rm -rf $(BIN_DIR)

test:
	go -C src test ./...

lint:
	go -C src run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

# ── Local dev ─────────────────────────────────────────────────────────────────

# Start DynamoDB Local and create the cache table
local-up:
	docker compose up -d
	@echo "Waiting for DynamoDB Local..." && sleep 2
	@$(MAKE) local-init-db

local-down:
	docker compose down

# Idempotent — safe to run multiple times
local-init-db:
	AWS_ACCESS_KEY_ID=local AWS_SECRET_ACCESS_KEY=local \
	  aws --endpoint-url $(LOCAL_DDB_URL) --region $(LOCAL_REGION) \
	    dynamodb create-table \
	      --table-name $(TABLE_NAME) \
	      --attribute-definitions AttributeName=pk,AttributeType=S \
	      --key-schema AttributeName=pk,KeyType=HASH \
	      --billing-mode PAY_PER_REQUEST \
	    2>/dev/null || true
	@echo "DynamoDB table ready."

local-db-tables:
	AWS_ACCESS_KEY_ID=local AWS_SECRET_ACCESS_KEY=local \
	  aws --endpoint-url $(LOCAL_DDB_URL) --region $(LOCAL_REGION) \
	    dynamodb list-tables

local-db-scan:
	AWS_ACCESS_KEY_ID=local AWS_SECRET_ACCESS_KEY=local \
	  aws --endpoint-url $(LOCAL_DDB_URL) --region $(LOCAL_REGION) \
	    dynamodb scan --table-name $(TABLE_NAME)

local-db-clear:
	AWS_ACCESS_KEY_ID=local AWS_SECRET_ACCESS_KEY=local \
	  aws --endpoint-url $(LOCAL_DDB_URL) --region $(LOCAL_REGION) \
	    dynamodb delete-table --table-name $(TABLE_NAME) 2>/dev/null || true
	@$(MAKE) local-init-db

# Build, then start API Gateway + Lambda locally via SAM (Terraform hook).
# Requires: Docker, AWS SAM CLI, and env.local.json (see env.local.json.example).
# SAM reads infra/*.tf to discover routes and the Lambda; no SAM template needed.
local-api: zip
	cd infra && TF_CLI_ARGS_plan="-var-file=local.tfvars" sam local start-api \
	  --hook-name terraform \
	  --beta-features \
	  --warm-containers EAGER \
	  --port 3000
