# Infermedica V3 API (unofficial)

Go interface to the infermedica REST API

## Description

This is a Go interface to the [Infermedica REST API](https://developer.infermedica.com/docs/api), I am just updating [mydoc-chat’s](https://github.com/mydoc-chat/infermedica) repository for my undergraduate thesis, so I don’t know if I’ll include all endpoints in this project.

## Current working

- [ ] Conditions
- [X] Diagnosis
- [ ] Explain
- [ ] Info
- [ ] Labtests
- [ ] Lookup
- [ ] Parse
- [ ] Riskfactors
- [ ] Search
- [ ] Suggest
- [ ] Symptoms
- [X] Triage

## Installation

```go get github.com/guiarnaldo/infermedica-v3```

## Usage examples

### Fetching symptoms
```go
app := infermedica.NewApp("appid", "appkey", "model", "source")
symptoms, err := app.Symptoms()
if err != nil {
    fmt.Printf("Could not fetch symptoms: %v", err)
}
fmt.Printf("All Symptoms: %v", symptoms)
```