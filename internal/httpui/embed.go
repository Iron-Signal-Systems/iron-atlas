package httpui

import "embed"

//go:embed templates/*.html static/*
var webFiles embed.FS
