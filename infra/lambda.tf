resource "aws_lambda_function" "oracle" {
  function_name = "arcana-oracle"
  filename      = "../bin/bootstrap.zip"
  role          = aws_iam_role.oracle_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  environment {
    variables = {
      SCW_API_URL               = var.scw_api_url
      SCW_SECRET_KEY            = var.scw_secret_key
      DYNAMODB_TABLE            = aws_dynamodb_table.card_cache.name
      AWS_ENDPOINT_URL_DYNAMODB = var.dynamodb_endpoint
    }
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.oracle.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.oracle.execution_arn}/*/*"
}

resource "aws_apigatewayv2_api" "oracle" {
  name          = "arcana-oracle"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = [var.cloudfront_origin]
    allow_methods = ["POST"]
    allow_headers = ["content-type"]
    max_age       = 86400
  }
}

resource "aws_apigatewayv2_integration" "oracle" {
  api_id                 = aws_apigatewayv2_api.oracle.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.oracle.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "card" {
  api_id    = aws_apigatewayv2_api.oracle.id
  route_key = "POST /oracle/card"
  target    = "integrations/${aws_apigatewayv2_integration.oracle.id}"
}

resource "aws_apigatewayv2_route" "spread" {
  api_id    = aws_apigatewayv2_api.oracle.id
  route_key = "POST /oracle/spread"
  target    = "integrations/${aws_apigatewayv2_integration.oracle.id}"
}

resource "aws_apigatewayv2_route" "astro" {
  api_id    = aws_apigatewayv2_api.oracle.id
  route_key = "POST /oracle/astro"
  target    = "integrations/${aws_apigatewayv2_integration.oracle.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.oracle.id
  name        = "$default"
  auto_deploy = true

  # card  → 20 req/min ≈ 0.33 req/s
  route_settings {
    route_key              = aws_apigatewayv2_route.card.route_key
    throttling_rate_limit  = 0.33
    throttling_burst_limit = 5
  }

  # spread → 5 req/min ≈ 0.08 req/s
  route_settings {
    route_key              = aws_apigatewayv2_route.spread.route_key
    throttling_rate_limit  = 0.08
    throttling_burst_limit = 2
  }

  # astro  → 5 req/min ≈ 0.08 req/s
  route_settings {
    route_key              = aws_apigatewayv2_route.astro.route_key
    throttling_rate_limit  = 0.08
    throttling_burst_limit = 2
  }
}

output "oracle_url" {
  value = aws_apigatewayv2_stage.default.invoke_url
}
