# Infermedica V3 API (unofficial)

Go interface to the infermedica REST API. Original repository: [mydoc-chat’s](https://github.com/mydoc-chat/infermedica)

## Description

This is a Go interface to the [Infermedica REST API: Triage](https://developer.infermedica.com/documentation/api-triage/quickstart/).

## Installation

```go get github.com/guiarnaldo/infermedica-v3```

# Usage examples

## Get diagnosis using Parse NLP
```go
    app := infermedica.NewApp("appid", "appkey", "model", "source")

	age := infermedica.Age{
		Value: 21,
	}

	parseRes, err := app.Parse(ParseReq{
		Text:            "I have diabetes",
		Age:             age,
		Sex:             SexMale,
		CorrectSpelling: true,
	})
	if err != nil {
		// Error Handling
	}

	evidences, err := parseRes.ParseToEvidence()
	if err != nil {
		// Error Handling
	}

	diagnosis, err := app.Diagnosis(DiagnosisReq{
		Sex:       SexMale,
		Age:       age,
		Evidences: evidences,
	})
	if err != nil {
		// Error Handling
	}

	fmt.Println(diagnosis)
```