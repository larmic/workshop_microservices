// Package archtest holds the repository-wide Guard-Gates that run as plain Go
// tests: architecture rules (import graph), OpenAPI <-> handler consistency and
// the colour palette shared by slides and dashboard.
//
// The package contains no production code. Run it on its own with
//
//	go test ./archtest/
//
// or a single gate with -run (TestArchitecture, TestOpenAPI, TestPalette).
package archtest
