// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driverutil

import (
	"os"
	"strings"
)

const AwsLambdaPrefix = "AWS_Lambda_"

const (
	// FaaS environment variable names

	// EnvVarAWSExecutionEnv is the AWS Execution environment variable.
	EnvVarAWSExecutionEnv = "AWS_EXECUTION_ENV"
	// EnvVarAWSLambdaRuntimeAPI is the AWS Lambda runtime API variable.
	EnvVarAWSLambdaRuntimeAPI = "AWS_LAMBDA_RUNTIME_API"
	// EnvVarFunctionsWorkerRuntime is the functions worker runtime variable.
	EnvVarFunctionsWorkerRuntime = "FUNCTIONS_WORKER_RUNTIME"
	// EnvVarKService is the K Service variable.
	EnvVarKService = "K_SERVICE"
	// EnvVarFunctionName is the function name variable.
	EnvVarFunctionName = "FUNCTION_NAME"
	// EnvVarVercel is the Vercel variable.
	EnvVarVercel = "VERCEL"
	// EnvVarK8s is the K8s variable.
	EnvVarK8s = "KUBERNETES_SERVICE_HOST"
)

const (
	// FaaS environment variable names

	// EnvVarAWSRegion is the AWS region variable.
	EnvVarAWSRegion = "AWS_REGION"
	// EnvVarAWSLambdaFunctionMemorySize is the AWS Lambda function memory size variable.
	EnvVarAWSLambdaFunctionMemorySize = "AWS_LAMBDA_FUNCTION_MEMORY_SIZE"
	// EnvVarFunctionMemoryMB is the function memory in megabytes variable.
	EnvVarFunctionMemoryMB = "FUNCTION_MEMORY_MB"
	// EnvVarFunctionTimeoutSec is the function timeout in seconds variable.
	EnvVarFunctionTimeoutSec = "FUNCTION_TIMEOUT_SEC"
	// EnvVarFunctionRegion is the function region variable.
	EnvVarFunctionRegion = "FUNCTION_REGION"
	// EnvVarVercelRegion is the Vercel region variable.
	EnvVarVercelRegion = "VERCEL_REGION"
)

const (
	// FaaS environment names used by the client

	// EnvNameAWSLambda is the AWS Lambda environment name.
	EnvNameAWSLambda = "aws.lambda"
	// EnvNameAzureFunc is the Azure Function environment name.
	EnvNameAzureFunc = "azure.func"
	// EnvNameGCPFunc is the Google Cloud Function environment name.
	EnvNameGCPFunc = "gcp.func"
	// EnvNameVercel is the Vercel environment name.
	EnvNameVercel = "vercel"
)

// GetFaasEnvName parses the FaaS environment variable name and returns the
// corresponding name used by the client. If none of the variables or variables
// for multiple names are populated the client.env value MUST be entirely
// omitted. When variables for multiple "client.env.name" values are present,
// "vercel" takes precedence over "aws.lambda"; any other combination MUST cause
// "client.env" to be entirely omitted.
func GetFaasEnvName() string {
	envVars := []string{
		EnvVarAWSExecutionEnv,
		EnvVarAWSLambdaRuntimeAPI,
		EnvVarFunctionsWorkerRuntime,
		EnvVarKService,
		EnvVarFunctionName,
		EnvVarVercel,
	}

	// If none of the variables are populated the client.env value MUST be
	// entirely omitted.
	names := make(map[string]struct{})

	for _, envVar := range envVars {
		val := os.Getenv(envVar)
		if val == "" {
			continue
		}

		var name string

		switch envVar {
		case EnvVarAWSExecutionEnv:
			if !strings.HasPrefix(val, AwsLambdaPrefix) {
				continue
			}

			name = EnvNameAWSLambda
		case EnvVarAWSLambdaRuntimeAPI:
			name = EnvNameAWSLambda
		case EnvVarFunctionsWorkerRuntime:
			name = EnvNameAzureFunc
		case EnvVarKService, EnvVarFunctionName:
			name = EnvNameGCPFunc
		case EnvVarVercel:
			// "vercel" takes precedence over "aws.lambda".
			delete(names, EnvNameAWSLambda)

			name = EnvNameVercel
		}

		names[name] = struct{}{}
		if len(names) > 1 {
			// If multiple names are populated the client.env value
			// MUST be entirely omitted.
			names = nil

			break
		}
	}

	for name := range names {
		return name
	}

	return ""
}
